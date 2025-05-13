package main

import (
	"Golang_Programming_Journey/2_blog-serie/global"
	"Golang_Programming_Journey/2_blog-serie/internal/model"
	"Golang_Programming_Journey/2_blog-serie/internal/routers"
	"Golang_Programming_Journey/2_blog-serie/pkg/logger"
	"Golang_Programming_Journey/2_blog-serie/pkg/redisClient"
	"Golang_Programming_Journey/2_blog-serie/pkg/setting"
	"Golang_Programming_Journey/2_blog-serie/pkg/tracer"
	"Golang_Programming_Journey/2_blog-serie/pkg/validator"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/natefinch/lumberjack"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	port      string
	runMode   string
	config    string
	isVersion bool
)

func init() {
	//err := SetupFlag()
	//if err != nil {
	//	log.Fatalf("init.setupFlag err: %v", err)
	//	return
	//}
	err := SetupSetting()
	if err != nil {
		log.Fatal("setup Setting err:", err)
		return
	}
	err = SetupDBEngine()
	if err != nil {
		log.Fatal("setup DBEngine err:", err)
	}

	err = SetupLogger()
	if err != nil {
		log.Fatal("setup Logger err:", err)
	}

	err = setupValidator()
	if err != nil {
		log.Fatalf("init.setupValidator err: %v", err)
	}

	err = SetupTracer()
	if err != nil {
		log.Fatal("setup Tracer err:", err)
	}

	err = SetupRedis()
	if err != nil {
		log.Fatal("setup Redis err:", err)
	}
}

//@title 博客系统
//@version 1.0
//@description Go 语言编程之旅：博客项目

func main() {
	//test01()
	gin.SetMode(global.ServerSetting.RunMode)
	router := routers.NewRouter()
	s := &http.Server{
		Addr:           global.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	//log.Printf("global.ServerSetting %v", global.ServerSetting)
	//log.Printf("global.AppSetting %v", global.AppSetting)
	//log.Printf("global.DatabaseSetting %v", global.DatabaseSetting)

	//测试Logger是否装配好
	//global.Logger.Infof(context.Background(), "%s:go programing-tour-book/%s", "shxu", "blog_service")
	//log.Println(s.ListenAndServe())

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("s.ListenAndServe err: %v\n", err)
		}
	}()

	//等待中断信号
	quit := make(chan os.Signal)
	//接收syscall.SIGINT和 syscall.SIGTERM信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	//最大控制时间，用于通知该服务端只有5秒的时间来处理原有的请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	log.Println("11111111111111")
	defer cancel()

	time.Sleep(1 * time.Second)
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

//func SetupFlag() error {
//	flag.StringVar(&port, "port", "", "启动端口")
//	flag.StringVar(&runMode, "mode", "", "启动模式")
//	flag.StringVar(&config, "config", "configs/", "指定要使用的配置文件路径")
//	flag.BoolVar(&isVersion, "version", false, "编译信息")
//	flag.Parse()
//
//	return nil
//}

func SetupSetting() error {

	//set, err := setting.NewSetting(strings.Split(config, ",")...)
	//if err != nil {
	//	return err
	//}

	set, err := setting.NewSetting()
	if err != nil {
		return err
	}

	err = set.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}

	err = set.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}

	err = set.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}

	err = set.ReadSection("JWT", &global.JWTSetting)
	if err != nil {
		return err
	}

	err = set.ReadSection("Email", &global.EmailSetting)

	if err != nil {
		return err
	}

	err = set.ReadSection("Redis", &global.RedisSetting)

	if err != nil {
		return err
	}

	err = set.ReadSection("Limiter", &global.LimiterSetting)

	if err != nil {
		return err
	}

	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	global.AppSetting.UploadImageMaxSize *= 1024 * 1024
	global.AppSetting.DefaultContextTimeout *= time.Second
	global.JWTSetting.Expire *= time.Second

	global.RedisSetting.ReadTimeout *= time.Second
	global.RedisSetting.WriteTimeout *= time.Second
	global.RedisSetting.DialTimeout *= time.Second

	global.LimiterSetting.FillInterval *= time.Second
	//global.LimiterSetting.Expiration *= time.Second

	if port != "" {
		global.ServerSetting.HttpPort = port
	}
	if runMode != "" {
		global.ServerSetting.RunMode = runMode
	}

	return nil

}

//func setupLogger()error{
//
//}

func SetupDBEngine() error {
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}
	return nil
}

func SetupLogger() error {
	fileName := global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt

	global.Logger = logger.NewLogger(
		&lumberjack.Logger{
			Filename:  fileName,
			MaxSize:   600,
			MaxAge:    10,
			LocalTime: true,
		},
		"",
		log.LstdFlags,
	).WithCaller(2)

	return nil
}

func setupValidator() error {
	global.Validator = validator.NewCustomValidator()
	global.Validator.Engine()
	binding.Validator = global.Validator

	return nil
}

func SetupTracer() error {
	jaegerTracer, _, err := tracer.NewJaegerTracer(
		"blog_service",
		"127.0.0.1:6831",
	)
	if err != nil {
		return err
	}
	global.Tracer = jaegerTracer
	return nil
}

func SetupRedis() error {
	global.RedisClient = redisClient.NewRedisClient(global.RedisSetting)
	return nil
}
