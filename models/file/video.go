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

type Video struct {
	Id string `bson:"id"`
	Name string `bson:"name"`
	Type string `bson:"type"`
	Width  int `bson:"width"`
	Height int `bson:"height"`
	Duration string `bson:"duration"`
	BitRate string `bson:"bit_rate"`
}
//240P 320×240 //Mobile iPhone MP4
//360P 640×360 //SD FLV
//480P 864×480 //HD MP4
//720P 960×720 //HD MP4
// run 视频 获取预览图、转码为mp4（360k/720k/1550k比特率）、获取视频基本信息
func (_self Video) run(info Info) error {
	//原始文件保存到本地进行转码
	//尝试过直接传入byte数组进行转码，有时候能转码成功，有时候就不行，希望有对cmd比较熟悉的同学进行重写
	fm, err := db.MG().GridFS("fs").OpenId(bson.ObjectIdHex(info.FileId))
	if err != nil {
		return err
	}
	dir,err := file.GetUploadAddress()
	if err != nil {
		return err
	}
	//视频保存到本地地址
	videoPath := filepath.Join(dir,fmt.Sprintf("%s.temp",info.FileId))
	fw, err := os.Create(videoPath)
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
		if err = os.Remove(videoPath);err != nil{
			log.Println(err)
		}
	}()
	var fileList [] interface{}
	var video Video
	//获取视频基本信息
	meta,err := convert.MediaInfo(videoPath)
	if err != nil {
		return err
	}
	info.Duration = meta.Format.Duration
	for _,item := range  meta.Streams{
		if item.CodecType == "video"{
			video.Width = item.Width
			video.Height = item.Height
			video.Duration = item.Duration
		}
	}
	//获取视频预览图
	err = convert.VideoPreview(videoPath, func(stdout io.ReadCloser) error {
		fId,err := saveFile(stdout,fmt.Sprintf("%s.jpg",info.FileId),"image/jpeg")
		if err != nil {
			return err
		}
		info.Preview = fId
		return nil
	})
	if err != nil {
		return err
	}
	//获取视频转mp4 360k
	err = convert.VideoToMp4(videoPath,"360k", func(stdout io.ReadCloser) error {
		fileName := fmt.Sprintf("%s_360k.mp4",info.FileId)
		fId,err := saveFile(stdout,fileName,"video/mp4")
		if err != nil {
			return err
		}
		video.Id = fId
		video.Name = fileName
		video.Type = ".mp4"
		video.BitRate = "360k"

		fileList = append(fileList,video)
		return nil
	})
	if err != nil {
		return err
	}
	//获取视频转mp4 720k
	err = convert.VideoToMp4(videoPath,"720k", func(stdout io.ReadCloser) error {
		fileName := fmt.Sprintf("%s_720k.mp4",info.FileId)
		fId,err := saveFile(stdout,fileName,"video/mp4")
		if err != nil {
			return err
		}
		video.Id = fId
		video.Name = fileName
		video.Type = ".mp4"
		video.BitRate = "720k"

		fileList = append(fileList,video)
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

// saveFile 保存文件
func saveFile(stdout io.ReadCloser,fileName,contentType string) (string,error){
	defer stdout.Close()
	fId := bson.NewObjectId()
	//文件写入数据库
	fw, err := db.MG().GridFS("fs").Create(fileName)
	if err != nil {
		return "",err
	}
	defer fw.Close()
	fw.SetId(fId)
	fw.SetContentType(contentType)
	_, err = io.Copy(fw, stdout)
	if err != nil {
		return "",err
	}
	return fmt.Sprintf("%x", string(fId)),nil
}
