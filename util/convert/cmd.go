package convert

import (
	"bytes"
	"os/exec"
	"runtime"
	"strings"
)

func FFmpegExec()(string,error){
	outFFmpeg, err := TestCmd(pathExec(), "ffmpeg")
	if err != nil {
		return "",err
	}
	return strings.Replace(outFFmpeg.String(), lineSeparator(), "", -1),nil
}

func FFprobeExec()(string,error){
	outProbe, err := TestCmd(pathExec(), "ffprobe")
	if err != nil {
		return "",err
	}
	return strings.Replace(outProbe.String(), lineSeparator(), "", -1),nil
}

func pathExec()string{
	var platform = runtime.GOOS
	var command string
	switch platform {
	case "windows":
		command = "where"
		break
	default:
		command = "which"
		break
	}
	return command
}

func lineSeparator() string {
	switch runtime.GOOS {
	case "windows":
		return "\r\n"
	default:
		return "\n"
	}
}

// TestCmd 输出打印值
func TestCmd(name string, arg... string)(bytes.Buffer, error){
	var out bytes.Buffer

	cmd := exec.Command(name, arg...)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return out, err
	}
	return out, nil
}
