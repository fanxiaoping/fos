package convert

import (
	"github.com/fanxiaoping/fos/util/validator"
	"path/filepath"
)

type Transcoding int32

const (
	T_Office  Transcoding = 1
	T_Pdf     Transcoding = 2
	T_Image   Transcoding = 3
	T_Video   Transcoding = 4
	T_Audio   Transcoding = 5
	T_Unknown Transcoding = 6
)

// 文件转码方式
var mode map[string]Transcoding

func init() {
	mode = map[string]Transcoding{
		".mp3":	T_Audio,
		".mp4":	T_Video,
		".doc":  T_Office,
		".docm": T_Office,
		".docx": T_Office,
		".dot":  T_Office,
		".dotm": T_Office,
		".dotx": T_Office,
		".epub": T_Office,
		".fodt": T_Office,
		".html": T_Office,
		".mht":  T_Office,
		".odt":  T_Office,
		".ott":  T_Office,
		".pdf":  T_Office,
		".rtf":  T_Office,
		".txt":  T_Office,
		".xps":  T_Office,
		".csv":  T_Office,
		".fods": T_Office,
		".ods":  T_Office,
		".ots":  T_Office,
		".xls":  T_Office,
		".xlsm": T_Office,
		".xlsx": T_Office,
		".xlt":  T_Office,
		".xltm": T_Office,
		".xltx": T_Office,
		".fodp": T_Office,
		".odp":  T_Office,
		".otp":  T_Office,
		".pot":  T_Office,
		".potm": T_Office,
		".potx": T_Office,
		".pps":  T_Office,
		".ppsm": T_Office,
		".ppsx": T_Office,
		".ppt":  T_Office,
		".pptm": T_Office,
		".pptx": T_Office,
	}
}

// GetTranscodingMethod 后去转码方式
func GetTranscodingMethod(fileName string) Transcoding {
	if validator.IsEmpty(fileName){
		return T_Unknown
	}
	suffix := filepath.Ext(fileName)
	if v,ok := mode[suffix];ok{
		return v
	}
	return T_Unknown
}
