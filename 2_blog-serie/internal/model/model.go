package model

import (
	"Golang_Programming_Journey/2_blog-serie/global"
	"Golang_Programming_Journey/2_blog-serie/pkg/otgorm"
	"Golang_Programming_Journey/2_blog-serie/pkg/setting"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"time"
)

const (
	STATE_OPEN  = 1
	STATE_CLOSE = 0
)

type Model struct {
	ID         uint32 `gorm:"primary_key" json:"id"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	CreatedOn  uint32 `json:"created_on"`
	ModifiedOn uint32 `json:"modified_on"`
	DeletedOn  uint32 `json:"deleted_on"`
	IsDel      uint8  `json:"is_del"`
}

func NewDBEngine(databaseSetting *setting.DatabaseSettings) (*gorm.DB, error) {
	dst := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
		databaseSetting.UserName,
		databaseSetting.Password,
		databaseSetting.Host,
		databaseSetting.DBName,
		databaseSetting.Charset,
		databaseSetting.ParseTime,
	)
	db, err := gorm.Open(mysql.Open(dst), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	if global.ServerSetting.RunMode == "debug" {
		//db = LogMode(true)
		//db.LogMode(true)
		db.Logger = db.Logger.LogMode(logger.Info)
	}
	//db.SingularTable(true)
	db.NamingStrategy = schema.NamingStrategy{
		TablePrefix:   "blog_",
		SingularTable: true,
	}

	//if err = db.AutoMigrate(&Tag{}); err != nil {
	//	return nil, err
	//}

	db.Callback().Create().Before("gorm:create").Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Before("gorm:update").Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Before("gorm:delete").Replace("gorm:delete", deleteCallback)
	//sqlDB（*sql.DB）是 Go 标准库 database/sql 中的数据库连接池对象，
	//主要负责管理数据库连接，如设置最大空闲连接数、最大打开连接数等，它并不具备 GORM 的回调注册功能。

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(databaseSetting.MaxIdleConns)
	sqlDB.SetMaxOpenConns(databaseSetting.MaxOpenConns)

	otgorm.AddGormCallbacks(db)

	return db, nil

}

func updateTimeStampForCreateCallback(db *gorm.DB) {

	if db.Statement.Schema != nil {
		timeNow := uint32(time.Now().Unix())
		if CreateFiled := db.Statement.Schema.LookUpField("CreatedOn"); CreateFiled != nil {
			if err := CreateFiled.Set(db.Statement.Context, db.Statement.ReflectValue, timeNow); err != nil {
				global.Logger.Errorf(db.Statement.Context, "updateTimeStampForCreateCallback Set err:%v", err)
			}
		}
		if ModifyFiled := db.Statement.Schema.LookUpField("ModifiedOn"); ModifyFiled != nil {
			if err := ModifyFiled.Set(db.Statement.Context, db.Statement.ReflectValue, timeNow); err != nil {
				global.Logger.Errorf(db.Statement.Context, "updateTimeStampForCreateCallback Set err:%v", err)
			}
		}
	}
}

func updateTimeStampForUpdateCallback(db *gorm.DB) {

	if _, ok := db.Statement.Settings.Load("gorm:update_column"); !ok {
		// 获取 ModifiedOn 字段
		if ModifyFiled := db.Statement.Schema.LookUpField("ModifiedOn"); ModifyFiled != nil {
			if err := ModifyFiled.Set(db.Statement.Context, db.Statement.ReflectValue, time.Now().Unix()); err != nil {
				global.Logger.Errorf(db.Statement.Context, "updateTimeStampForUpdateCallback Set err:%v", err)
			}
		}
	}
}

func deleteCallback(db *gorm.DB) {
	log.Println("deleteCallback is called")
	if db.Statement.Schema != nil {
		var extraOption string
		if str, ok := db.Statement.Settings.Load("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}
		deletedOnFiles := db.Statement.Schema.LookUpField("DeletedOn")
		iSDelField := db.Statement.Schema.LookUpField("IsDel")
		if !db.Statement.Unscoped && deletedOnFiles != nil && iSDelField != nil {
			//软删除 将是否删除 位置 置为已删除
			db.Statement.SQL.Reset()    // 清空 SQL 语句，防止拼接错误
			db.Statement.Build("WHERE") // 重新构建 WHERE 子句

			whereClause := db.Statement.SQL.String()
			if whereClause == "" {
				global.Logger.Errorf(db.Statement.Context, "警告: WHERE 条件为空，避免误软删整张表")
				return
			}
			log.Println(db.Statement.Vars)

			timeNow := uint32(time.Now().Unix())

			log.Println(addExtraSpaceIfExist(whereClause))
			log.Println(addExtraSpaceIfExist(extraOption))
			sql := fmt.Sprintf(
				"UPDATE %v SET %v=%v, %v=%v %v %v",
				db.Statement.Table,
				deletedOnFiles.DBName,
				timeNow, //表示软删除时间
				iSDelField.DBName,
				1, //表示软删除标记位标志
				addExtraSpaceIfExist(whereClause),
				addExtraSpaceIfExist(extraOption),
			)

			//这个很奇怪，会优先自动把原来的delete中的两个参数填充到占位符，所以需要在构造语句时就要把deletedOnFiles、iSDelField填充进去
			// func (t Tag) Delete(db *gorm.DB) error {
			//return db.Unscoped().Where("id = ? AND is_del = ?", t.ID, 0).Delete(&t).Error //这个是硬删除
			//return db.Where("id = ? AND is_del = ?", t.ID, 0).Delete(&t).Error
			db.Exec(sql)

		} else {
			// 硬删除
			// 硬删除，确保 WHERE 条件
			db.Statement.SQL.Reset()    // 清空 SQL 语句，防止拼接错误，我没有清空就导致把整个表都给删除了
			db.Statement.Build("WHERE") // 重新构建 WHERE 子句

			whereClause := db.Statement.SQL.String()
			if whereClause == "" {
				global.Logger.Errorf(db.Statement.Context, "警告: WHERE 条件为空，避免误删整张表")
				return
			}

			sql := fmt.Sprintf(
				"DELETE FROM %v %v %v",
				db.Statement.Table,
				addExtraSpaceIfExist(whereClause),
				addExtraSpaceIfExist(extraOption),
			)
			db.Exec(sql)

		}
	}
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
