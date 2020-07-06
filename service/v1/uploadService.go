package v1

import (
	"encoding/json"
	"fmt"
	"github.com/fanxiaoping/fos/models/file"
	"github.com/fanxiaoping/fos/util/db"
	"github.com/fanxiaoping/fos/util/fhttp"
	ufile "github.com/fanxiaoping/fos/util/file"
	"gopkg.in/mgo.v2/bson"
	"io"
	"net/http"
	"path"
)

// UploadService 文件上传
type UploadService struct {
}


func (_self *UploadService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	//接收客户端传来的文件 file 与客户端保持一致
	hFile, handler, err := r.FormFile("file")
	if err != nil {
		fhttp.Error(w,err)
		return
	}
	defer hFile.Close()

	//文件写入数据库
	fw, err := db.MG().GridFS("fs").Create(handler.Filename)
	if err != nil {
		fhttp.Error(w, err)
		return
	}
	fId := bson.NewObjectId()
	fw.SetId(fId)
	fw.SetContentType(handler.Header.Get("Content-Type"))
	_, err = io.Copy(fw, hFile)
	if err != nil {
		fhttp.Error(w, err)
		return
	}
	fw.Close()
	//保存文件信息
	fBase := file.Info{
		Id:bson.NewObjectId(),
		FileId:fmt.Sprintf("%x", string(fId)),
		FileName:handler.Filename,
		FileType:path.Ext(handler.Filename),
		Status:file.File_Transcoding,
		Title:ufile.RemoveFileSuffix(handler.Filename),
		FileSize:handler.Size,
	}
	err = db.MG().C("file_info").Insert(&fBase)
	if err != nil {
		fhttp.Error(w, err)
		return
	}
	//加入转码
	go file.Convert{}.Execute(fBase)

	b, _ := json.Marshal(map[string]interface{}{"id":fBase.Id})
	fmt.Fprintln(w,string(b))
}

