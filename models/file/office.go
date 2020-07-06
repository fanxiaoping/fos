package file

import (
	"fmt"
	"github.com/fanxiaoping/fos/util/convert"
	"github.com/fanxiaoping/fos/util/db"
	"gopkg.in/mgo.v2/bson"
	"io"
	"os"
	"path/filepath"
	"sort"
)

type Office struct {
	Id string `bson:"id"`
	Name string `bson:"name"`
	Type string `bson:"type"`
	Size int64 `bson:"size"`
}
// run office文档转为图片、获取预览图
func (_self Office) run(info Info) error{
	res,err := convert.OfficeToPNG(info.FileId,info.FileName)
	if err != nil{
		return  err
	}
	//对文件排序
	sort.Sort(fileUrlSort(res))

	var fileList [] interface{}
	for n,f := range res{
		id := bson.NewObjectId()
		item := Office{
			Id:fmt.Sprintf("%x", string(id)),
			Name:filepath.Base(f),
			Type:filepath.Ext(f),
		}
		//获取文件信息
		fileInfo, err := os.Stat(f)
		if err != nil {
			return err
		}
		item.Size = fileInfo.Size()

		//文件写入mongodb
		fw, err := db.MG().GridFS("fs").Create(item.Name)
		if err != nil {
			return err
		}
		fw.SetId(id)
		fw.SetContentType("image/png")

		out, _ := os.OpenFile(f, os.O_RDWR, 0666)
		_, err = io.Copy(fw, out)
		if err != nil {
			return err
		}
		err = out.Close()
		if err != nil {
			return err
		}
		fw.Close()
		out.Close()
		//取第一张作为预览图
		if n == 0{
			info.Preview = item.Id
		}
		fileList = append(fileList,item)
	}
	info.TranscodingInfo = fileList
	info.Status = File_TranscodingCompleted
	err = db.MG().C("file_info").UpdateId(info.Id,info)
	if err != nil{
		return  err
	}
	return nil
}
