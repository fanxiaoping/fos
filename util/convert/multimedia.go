package convert

import (
	"bytes"
	"encoding/json"
	"io"
	"os/exec"
)

// VideoToH264 视频转换为h264
func VideoToH264(fp string,callBack func(stdout io.ReadCloser)error)error{
	arg := []string{"-i",fp,"-vcodec","h264","-f","h264","pipe:1"}

	return ffmpeg(callBack,arg...)
}

// VideoToMp4 视频转换为不同码率的mp4 (360k/720k/1550k)
func VideoToMp4(fp,bitRate string,callBack func(stdout io.ReadCloser)error)error{
	arg := []string{"-i",fp,"-b:v",bitRate,"-c:a","copy","-c:v","libx264","-movflags","frag_keyframe+empty_moov","-f", "mp4","pipe:1"}

	return ffmpeg(callBack,arg...)
}

// AudioToMp3 音频转换为不同码率的mp3 (64k | 128k)
func AudioToMp3(fp,bitRate string,callBack func(stdout io.ReadCloser)error)error{
	arg := []string{"-i",fp,"-ab",bitRate,"-f", "mp3","pipe:1"}

	return ffmpeg(callBack,arg...)
}

// VideoPreview 获取视频预览图
func VideoPreview(fp string,callBack func(stdout io.ReadCloser) error)error{
	//取第一秒为预览图
	arg := []string{"-i", fp, "-ss", "1",  "-vframes", "1", "-f","mjpeg", "pipe:1"}

	return ffmpeg(callBack,arg...)
}

// 获取多媒体信息
func MediaInfo(fp string)(Metadata,error){
	var meata  Metadata
	p,err := FFprobeExec()
	if err != nil {
		return  meata,err
	}
	arg := []string{"-i", fp, "-print_format", "json", "-show_format", "-show_streams", "-show_error"}
	out,err := TestCmd(p,arg...)
	if err != nil {
		return  meata,err
	}
	if err = json.Unmarshal([]byte(out.String()), &meata); err != nil {
		return meata,err
	}
	return meata,nil
}

// ffmpeg 转换
func ffmpeg(callBack func(stdout io.ReadCloser) error,arg... string)error{
	f,err := FFmpegExec()
	if err != nil {
		return  err
	}
	cmd := exec.Command(f,arg...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return  err
	}
	defer stdout.Close()

	err = cmd.Start()
	if err != nil {
		return  err
	}

	//对转换后数据进行处理
	err = callBack(stdout)
	if err != nil {
		return  err
	}
	err = cmd.Wait()
	if err != nil {
		return  err
	}
	return nil
}

// populateStdin 数据写入stdin
func populateStdin(stdin io.WriteCloser,file []byte) error {
	_,err := io.Copy(stdin, bytes.NewReader(file))
	if err != nil{
		return err
	}
	return  nil
}