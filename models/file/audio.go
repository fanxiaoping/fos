package file

import (
	"fmt"
	"github.com/fanxiaoping/fos/util/convert"
	"github.com/fanxiaoping/fos/util/db"
	"github.com/fanxiaoping/fos/util/file"
	"gopkg.in/mgo.v2/bson"
	"io"
	"log"
	"os"
	"path/filepath"
)

type Audio struct {
	Id string `bson:"id"`
	Name string `bson:"name"`
	Type string `bson:"type"`
	Duration string `bson:"duration"`
	BitRate string `bson:"bit_rate"`
}
// run 音频 转码为mp3（64k/128k 比特率）、获取音频基本信息
func (_self Audio) run(info Info) error {
	//原始文件保存到本地进行转码
	fm, err := db.MG().GridFS("fs").OpenId(bson.ObjectIdHex(info.FileId))
	if err != nil {
		return err
	}
	dir,err := file.GetUploadAddress()
	if err != nil {
		return err
	}
	//音频保存到本地地址
	audioPath := filepath.Join(dir,fmt.Sprintf("%s.temp",info.FileId))
	fw, err := os.Create(audioPath)
	if err != nil {
		return err
	}
	_, err = io.Copy(fw, fm)
	if err != nil {
		return err
	}
	fm.Close()
	fw.Close()
	defer func() {
		//删除临时文件
		if err = os.Remove(audioPath);err != nil{
			log.Println(err)
		}
	}()
	var fileList [] interface{}
	var audio Audio
	//获取音频基本信息
	meta,err := convert.MediaInfo(audioPath)
	if err != nil {
		return err
	}
	info.Duration = meta.Format.Duration
	for _,item := range  meta.Streams{
		if item.CodecType == "audio"{
			audio.Duration = item.Duration
		}
	}
	//获取音频转mp3 64k
	err = convert.AudioToMp3(audioPath,"64k", func(stdout io.ReadCloser) error {
		fileName := fmt.Sprintf("%s_64k.mp3",info.FileId)
		fId,err := saveFile(stdout,fileName,"audio/mp3")
		if err != nil {
			return err
		}
		audio.Id = fId
		audio.Name = fileName
		audio.Type = ".mp3"
		audio.BitRate = "64k"

		fileList = append(fileList,audio)
		return nil
	})
	if err != nil {
		return err
	}
	//获取音频转mp3 128k
	err = convert.AudioToMp3(audioPath,"128k", func(stdout io.ReadCloser) error {
		fileName := fmt.Sprintf("%s_128k.mp3",info.FileId)
		fId,err := saveFile(stdout,fileName,"audio/mp3")
		if err != nil {
			return err
		}
		audio.Id = fId
		audio.Name = fileName
		audio.Type = ".mp3"
		audio.BitRate = "128k"

		fileList = append(fileList,audio)
		return nil
	})
	if err != nil {
		return err
	}
	//转码完成 更新数据
	info.TranscodingInfo = fileList
	info.Status = File_TranscodingCompleted
	err = db.MG().C("file_info").UpdateId(info.Id,info)
	if err != nil{
		return  err
	}
	return nil
}
