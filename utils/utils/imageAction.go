package utils

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// NewDownImage ...
func NewDownImage(appID int64, tagName, pathName string) DownImage {
	return DownImage{
		tagName:  tagName,
		appID:    appID,
		pathName: pathName,
	}
}

// DownImage 下载图像初始化
type DownImage struct {
	appID    int64
	tagName  string
	pathName string
}

func (that DownImage) path(activeName string) string {
	_imgNameSlice := []string{
		"upload",
		strings.Trim(that.tagName, "/"),
		ToString(that.appID),
		strings.Trim(that.pathName, "/"),
		TimeStr("Ymd"),
		activeName[0:2],
		activeName[2:4],
		//activeName + ".png",
	} //图像名称的前4位做目录
	return strings.Join(_imgNameSlice, "/") //保存图像目录
}

// GetBASE64ToImg ...
func (that DownImage) GetBASE64ToImg(imagBASE64 string) (string, error) {
	activeName := UniqueId()
	dirPath := that.path(activeName)
	faceImageName := dirPath + "/" + activeName + ".png" //图像名称路径
	//
	if _err := ImageFileExits(dirPath); _err != nil {
		return faceImageName, _err
	}
	dist, _ := base64.StdEncoding.DecodeString(imagBASE64)
	f, err := os.OpenFile(faceImageName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer f.Close()
	f.Write(dist)
	return faceImageName, err
}

// GetUrlToImg  远程下载图片   Author:tang
func (that DownImage) GetUrlToImg(imagPath string) (string, error) {
	activeName := UniqueId()
	dirPath := that.path(activeName)
	faceImageName := dirPath + "/" + activeName + ".png" //图像名称路径
	//
	if _err := ImageFileExits(dirPath); _err != nil {
		return faceImageName, _err
	}
	resp, _err := http.Get(imagPath)
	if _err != nil {
		return faceImageName, _err
	}
	body, _err := ioutil.ReadAll(resp.Body)
	if _err != nil {
		return faceImageName, _err
	}
	out, err := os.Create(faceImageName)
	if err != nil {
		return faceImageName, err
	}
	io.Copy(out, bytes.NewReader(body))
	return faceImageName, err

}

// 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// region   检测目录   Author:tang
func ImageFileExits(_dir string) error {
	exist, err := PathExists(_dir)
	if err != nil {
		return errors.New("get dir error! " + err.Error())
	}
	if !exist {
		// 创建文件夹
		err := os.MkdirAll(_dir, os.ModePerm)
		if err != nil {
			return errors.New("mkdir failed! " + err.Error())
		}
	}
	return nil
}

//endregion

// 远端图片转 Base64 格式
func ImageToBase64(imgUrl string) (string, error) {
	//获取远端图片
	res, err := http.Get(imgUrl)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	// 读取获取的[]byte数据
	data, _ := ioutil.ReadAll(res.Body)
	return base64.StdEncoding.EncodeToString(data), nil
}

// 远端图片转 Byte 格式
func ImageToByte(imgUrl string) ([]byte, error) {
	//获取远端图片
	res, err := http.Get(imgUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	// 读取获取的[]byte数据
	data, _err := ioutil.ReadAll(res.Body)
	if _err != nil {
		return nil, _err
	}
	return data, nil
}
