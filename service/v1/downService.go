package v1

import (
	"errors"
	"fmt"
	"github.com/fanxiaoping/fos/util/db"
	"github.com/fanxiaoping/fos/util/fhttp"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

// 文件下载
type DownService struct {
}

func (_self *DownService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if !bson.IsObjectIdHex(id) {
		fhttp.Error(w, errors.New("mongodb key格式不正确"))
		return
	}
	file, err := db.MG().GridFS("fs").OpenId(bson.ObjectIdHex(id))
	if err != nil {
		fhttp.Error(w, err)
		return
	}
	defer file.Close()

	b := make([]byte, file.Size())
	_, err = file.Read(b)
	if err != nil {
		fhttp.Error(w, err)
		return
	}
	w.Header().Set("Content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.Name()))
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", file.Size()))
	w.Header().Set("Content-Type", file.ContentType())
	w.Write(b)
}
