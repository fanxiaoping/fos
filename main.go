package main

import (
	"github.com/fanxiaoping/fos/router"
	"github.com/fanxiaoping/fos/util/config"
	"github.com/fanxiaoping/gfly"
	"io/ioutil"
	"log"
	"path/filepath"
)

func main(){
	//获取证书
	key, err := ioutil.ReadFile(filepath.Join(config.GetCurrentDirectory(), "conf", "server.key"))
	if err != nil {
		log.Fatalln("read server.key:%s", err)
	}
	crt, err := ioutil.ReadFile(filepath.Join(config.GetCurrentDirectory(), "conf", "server.crt"))
	if err != nil {
		log.Fatalln("read server.crt:%s", err)
	}
	//初始化配置
	gfly.SetConfigCertOverride([]byte(crt), []byte(key), config.String("addr"), "yxt.com")

	//静态文件访问
	//gfly.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(config.GetCurrentDirectory(), "static")))))
	//文件处理
	gfly.Handle("/file/", router.RegisterMux())

	gfly.Run()

}
