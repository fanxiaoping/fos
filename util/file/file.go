package file

import (
	"crypto/md5"
	"fmt"
	"github.com/fanxiaoping/fos/util/config"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// CreateMkdir 创建目录
func  CreateMkdir(path string) (string, error) {
	isExist, err := PathExists(path)
	if err != nil {
		return "", err
	}
	if !isExist {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	return path, nil
}

// PathExists 判断文件夹是否存在
func  PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// RemoveFileSuffix 去掉文件地址后缀
func RemoveFileSuffix(fielPath string) string{
	fileSuffix := path.Ext(fielPath)
	return strings.TrimSuffix(fielPath, fileSuffix)
}

// GetAllFile 获取文件夹中所有文件
func GetAllFile(dirPath string)(files []string,err error){
	rd, err := ioutil.ReadDir(dirPath)
	if err != nil{
		return
	}
	for _,file := range rd{
		if !file.IsDir(){
			files = append(files,filepath.Join(dirPath,file.Name()))
		}
	}
	return files,nil
}
// GetUploadAddress 获取文件上传地址
func GetUploadAddress() (path string, err error){
	path =  filepath.Join(config.GetCurrentDirectory(), "static", "upload")
	_,err = CreateMkdir(path)
	if err != nil{
		return
	}
	return
}

// GetUploadAddress 获取文件上传地址并加上当天时间
func GetUploadAddressDate() (path string, err error){
	path =  filepath.Join(config.GetCurrentDirectory(), "static", "upload", time.Now().Format("2006-01-02"))
	_,err = CreateMkdir(path)
	if err != nil{
		return
	}
	return
}

// GetFileMd5 获取文件md5值
func GetFileMd5(filePath string) (md5Str string,err error){
	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()

	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		return "",err
	}
	return fmt.Sprintf("%x", md5hash.Sum([]byte(""))),nil
}

// GetFileSize 获取文件大小
func GetFileSize(filePath string)(int64,error){
	fileInfo, err := os.Stat(filePath)
	if err != nil{
		return 0,err
	}
	return fileInfo.Size(),nil
}