package upload

import (
	"Golang_Programming_Journey/2_blog-serie/global"
	"Golang_Programming_Journey/2_blog-serie/pkg/util"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

type FileType int

const TypeImage FileType = iota + 1

// 将文件名统一转换为 什么形式
func GetFileName(name string) string {
	ext := GetFileExt(name) //获取文件扩展名

	fileName := strings.TrimSuffix(name, ext) //删除指定的字符尾串，其中第一个参数s是原字符串，第二个参数suffix是需要删除的后缀子串
	fileName = util.EncodeMD5(fileName)
	return fileName + ext
}

// 返回文件的扩展名
func GetFileExt(fileName string) string {
	return path.Ext(fileName)
}
func GetSavePath() string {
	return global.AppSetting.UploadSavePath
}

func CheckSavePath(savePath string) bool {
	_, err := os.Stat(savePath) //主要用于获取指定路径文件或目录的相关信息
	return os.IsNotExist(err)   //如果有err判断文件是不是不存在导致的错误
}

func CheckContainExt(t FileType, name string) bool {
	ext := GetFileExt(name)
	ext = strings.ToUpper(ext)
	switch t {
	case TypeImage:
		for _, allowExt := range global.AppSetting.UploadImageAllowExts {
			if strings.ToUpper(allowExt) == ext {
				return true
			}
		}
	}
	return false
}

func CheckMaxSize(t FileType, f multipart.File) bool {
	content, _ := io.ReadAll(f)
	size := len(content)
	switch t {
	case TypeImage:
		if size >= global.AppSetting.UploadImageMaxSize {
			return true
		}
	}
	return false
}

func CheckPermission(dst string) bool {
	_, err := os.Stat(dst)
	return os.IsPermission(err)
}

func CreateSavePath(dst string, perm os.FileMode) error {
	err := os.MkdirAll(dst, perm)
	if err != nil {
		return err
	}
	return nil
}

// 将一个通过 HTTP 表单上传的文件保存到指定的目标路径,所以要先读提交过来的信息，然后打开本地文件，然后将服务器端的写入到本地文件
func SaveFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
