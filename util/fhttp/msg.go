package fhttp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

/**
*	接口返回报文
 */
type HttpErrMsg struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func ErrorMsg(w http.ResponseWriter, msg string, data interface{}) {
	w.Header().Set("Status","500")
	e := HttpErrMsg{
		Code: 21,
		Msg:  msg,
		Data: data,
	}
	b, _ := json.Marshal(e)

	log.Println(msg)

	fmt.Fprintln(w, string(b))
}

func Error(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	e := HttpErrMsg{Code: 20, Data: ""}
	//判断是否是用户错误
	//if _, ok := err.(*userError.SyntaxError); ok {
	//	e.Code = 21
	//}
	e.Msg = err.Error()
	b, _ := json.Marshal(e)
	//写入日志
	log.Println(err.Error())

	fmt.Fprintln(w, string(b))
}
