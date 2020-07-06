package file

import (
	"github.com/fanxiaoping/fos/util/convert"
	"github.com/fanxiaoping/fos/util/db"
	"log"
)

// Convert 文件转码
type Convert struct {
}

// Execute 执行转码
func (_self Convert) Execute(info Info){
	log.Println("开始处理文件：",info.FileName)
	mode := convert.GetTranscodingMethod(info.FileName)
	var err error
	info.TranscodingType = mode

	switch mode {
	case convert.T_Office:
		err = Office{}.run(info)
	case convert.T_Video:
		err = Video{}.run(info)
	case convert.T_Audio:
		err = Audio{}.run(info)
	}
	//处理转码失败
	if err != nil{
		log.Println("转码失败：",err)
		info.Status = File_TranscodingFailed
		err = db.MG().C("file_info").UpdateId(info.Id,info)
		if err != nil{
			log.Println(err)
		}
	}else{
		log.Println("转码成功：",info.FileName)
	}
}
