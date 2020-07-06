package file

import (
	"github.com/fanxiaoping/fos/util/convert"
	"gopkg.in/mgo.v2/bson"
)
type fileUrlSort []string

func (s fileUrlSort) Len() int { return len(s) }

func (s fileUrlSort) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

//排序规则
//先比较字符串长度，如果长度相同就比较值
func (s fileUrlSort) Less(i, j int) bool {
	//比较长度
	if len(s[i]) < len(s[j]) {
		return true
	}
	//比较值
	if len(s[i]) == len(s[j]) {
		return s[i] < s[j]
	}
	return false
}


type FileStatus int32

const (
	File_Transcoding FileStatus = 1
	File_TranscodingFailed FileStatus = 2
	File_TranscodingCompleted FileStatus = 3
)

type Info struct {
	//mg_id
	Id bson.ObjectId `bson:"_id"`
	//标题
	Title string `bson:"title"`
	//原始文件编号
	FileId string `bson:"file_id"`
	//原始文件名(带后缀)
	FileName string `bson:"file_name"`
	//原始文件类型 例：.doc .docx .png .jpg ...
	FileType string `bson:"file_type"`
	//原始文件大小
	FileSize int64 `bson:"file_size"`
	//状态 1=转码中 2=转码失败 3=转码完成
	Status FileStatus `bson:"status"`
	//预览图
	Preview string `bson:"preview"`
	//多媒体时长
	Duration string `bson:"duration"`
	//转码类型
	TranscodingType convert.Transcoding `bson:"transcoding_type"`
	//转码信息
	TranscodingInfo []interface{ } `bson:"transcoding_info"`
}

