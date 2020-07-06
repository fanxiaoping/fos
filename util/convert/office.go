package convert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fanxiaoping/fos/util/config"
	"github.com/fanxiaoping/fos/util/fhttp"
	"github.com/fanxiaoping/fos/util/file"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

type onlyOffice struct {
	Percent int32  `json:"percent"`
	FileUrl string `json:"fileUrl"`
}

// office 转为图片 采用ONLYOFFICE文件服务器转码
func OfficeToPNG(key,fileName string) (files []string,err error) {
	convertUrl := fmt.Sprintf("%s%s",config.String("file_down_url"),key)
	//office转png
	onlyO, err := onlyOfficeToPNG(fileName, convertUrl, filepath.Ext(fileName), key)
	if err != nil {
		return
	}
	path,err  := file.GetUploadAddressDate()
	dirPath, err := file.CreateMkdir(path)
	if err != nil {
		return
	}
	zipPath,err := fhttp.DownloadFile(onlyO.FileUrl,dirPath)
	if err != nil {
		return
	}
	return file.Unzip(zipPath)
}

// onlyOfficeToPNG only office转png
func onlyOfficeToPNG(title, url, fileType, key string) (res onlyOffice, err error) {
	params := make(map[string]interface{})
	params["key"] = key
	params["async"] = false
	params["filetype"] = fileType
	params["outputtype"] = "png"
	params["thumbnail"] = map[string]interface{}{
		"first": false,
	}
	params["title"] = title
	params["url"] = url

	return onlyConvert(params)
}

// onlyConvert only转码
func onlyConvert(params map[string]interface{}) (res onlyOffice, err error) {
	client := &http.Client{}

	jsonStr, _ := json.Marshal(params)
	fmt.Println(string(jsonStr))
	req, err := http.NewRequest("POST",config.String("only_converter_url"),bytes.NewBuffer(jsonStr))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	res = onlyOffice{}
	fmt.Println(string(body))
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}
	return res, nil
}

