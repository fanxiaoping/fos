package file


import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Zip 压缩
func Zip(zipFile string, fileList []string) error {
	// 创建 zip 包文件
	fw, err := os.Create(zipFile)
	if err != nil {
		log.Fatal()
	}
	defer fw.Close()

	// 实例化新的 zip.Writer
	zw := zip.NewWriter(fw)
	defer func() {
		// 检测一下是否成功关闭
		if err := zw.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	for _, fileName := range fileList {
		fr, err := os.Open(fileName)
		if err != nil {
			return err
		}
		fi, err := fr.Stat()
		if err != nil {
			return err
		}
		// 写入文件的头信息
		fh, err := zip.FileInfoHeader(fi)
		w, err := zw.CreateHeader(fh)
		if err != nil {
			return err
		}
		// 写入文件内容
		_, err = io.Copy(w, fr)
		if err != nil {
			return err
		}
	}
	return nil
}

// Unzip 解压zip到文件当前目录（主要用于office转png后解压）
func Unzip(zipFile string) (files []string,err error) {
	zr, err := zip.OpenReader(zipFile)
	defer zr.Close()
	if err != nil {
		return
	}
	outDir := RemoveFileSuffix(zipFile)
	outDir, err = CreateMkdir(outDir)
	if err != nil {
		return
	}
	for _, file := range zr.File {
		fileName := filepath.Join(outDir,file.Name)
		// 如果是目录，则创建目录
		if file.FileInfo().IsDir() {
			if err = os.MkdirAll(fileName, file.Mode());err != nil {
				return
			}
			continue
		}
		// 获取到 Reader
		fr, err := file.Open()
		if err != nil {
			return nil,err
		}
		fw, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
		if err != nil {
			return nil,err
		}
		_, err = io.Copy(fw, fr)
		if err != nil {
			return nil,err
		}
		fw.Close()
		fr.Close()
		files = append(files,fileName)
	}
	return
}