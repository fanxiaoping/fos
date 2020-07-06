package fhttp

import (
	"errors"
	"fmt"
	"github.com/fanxiaoping/fos/util/file"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

// DownloadFile 下载文件
func DownloadFile(durl string, outDir string) (outPath string,err error) {
	var (
		buf     = make([]byte, 32*1024)
		written int64
	)
	uri, err := url.ParseRequestURI(durl)
	if err != nil {
		return
	}
	dirPath, err := file.CreateMkdir(outDir)
	if err != nil {
		return
	}
	filePath :=  filepath.Join(dirPath,path.Base(uri.Path))
	fmt.Println(filePath)
	//创建一个http client
	client := new(http.Client)
	//client.Timeout = time.Second * 60 //设置超时时间
	//get方法获取资源
	resp, err := client.Get(durl)
	if err != nil {
		return
	}

	//创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer file.Close()
	if resp.Body == nil {
		return "",errors.New("body is null")
	}
	defer resp.Body.Close()
	//下面是 io.copyBuffer() 的简化版本
	for {
		//读取bytes
		nr, er := resp.Body.Read(buf)
		if nr > 0 {
			//写入bytes
			nw, ew := file.Write(buf[0:nr])
			//数据长度大于0
			if nw > 0 {
				written += int64(nw)
			}
			//写入出错
			if ew != nil {
				err = ew
				break
			}
			//读取是数据长度不等于写入的数据长度
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	if err != nil{
		return
	}
	return filePath,nil
}

