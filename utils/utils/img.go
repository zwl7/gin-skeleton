package utils

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var GIF = []byte("GIF")
var BMP = []byte("BM")
var JPG = []byte{0xff, 0xd8, 0xff}
var PNG = []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}

const (
	GIF_TYPE = "image/gif"
	BMP_TYPE = "image/x-ms-bmp"
	JPG_TYPE = "image/jpeg"
	PNG_TYPE = "image/png"
)

func NewImgData() *imgData {
	return &imgData{}
}

type imgData struct {
	Reader io.Reader
}

func (that *imgData) WgetToBase64(imgUrl string) (string, error) {
	_byte, _err := that.Wget(imgUrl)
	if _err != nil {
		return "", errors.New("wget error:" + _err.Error())
	}
	return base64.StdEncoding.EncodeToString(_byte), nil
}

func (that *imgData) WgetToBase64AndHeader(imgUrl string) (string, error) {
	_byte, _err := that.Wget(imgUrl)
	if _err != nil {
		return "", errors.New("wget error:" + _err.Error())
	}
	_prefix, _err := that.getBase64Prefix(_byte)
	if _err != nil {
		return "", _err
	}
	return _prefix + "," + base64.StdEncoding.EncodeToString(_byte), nil
}

func (that *imgData) WgetToDecode(imgUrl string) (image.Image, error) {
	_byte, _err := that.Wget(imgUrl)
	if _err != nil {
		return nil, _err
	}
	return that.getImage(_byte)
}

func (that *imgData) WgetToFile(imgUrl, dirName string) (string, error) {
	_byte, _err := that.Wget(imgUrl)
	if _err != nil {
		return "", _err
	}
	//
	_suffix, _err := that.getImageSuffix(_byte)
	if _err != nil {
		return "", _err
	}
	_activeName := UniqueId()
	_fileName := that.path(dirName, _activeName) + _suffix
	_out, _err := os.Create(_fileName)
	if _err != nil {
		return "", _err
	}
	if _, _err := io.Copy(_out, bytes.NewReader(_byte)); _err != nil {
		return "", _err
	}
	return _fileName, nil
}

func (that *imgData) Wget(imgUrl string) ([]byte, error) {
	if imgUrl == "" {
		return nil, errors.New("no image source detected")
	}
	//从网络下载
	_resp, _err := http.Get(imgUrl)
	if _err != nil {
		return nil, errors.New("no image source detected")
	}
	defer _resp.Body.Close()
	return ioutil.ReadAll(_resp.Body)
}

func (that *imgData) Decode(imgName string) (image.Image, error) {
	_file, _err := os.Open(imgName)
	defer _file.Close()
	if _err != nil {
		return nil, _err
	}
	_byte, _err := ioutil.ReadFile(imgName)
	if _err != nil {
		return nil, _err
	}
	return that.getImage(_byte)
}

func (that *imgData) getImage(bt []byte) (image.Image, error) {
	_type, _err := that.getImageType(bt)
	if _err != nil {
		return nil, _err
	}
	_reader := bytes.NewReader(bt)
	switch _type {
	case GIF_TYPE:
		return gif.Decode(_reader)
	case JPG_TYPE:
		return jpeg.Decode(_reader)
	case PNG_TYPE:
		return png.Decode(_reader)
	default:
		return nil, errors.New("undefined type")
	}
}

func (that *imgData) getBase64Prefix(bt []byte) (string, error) {
	_type, _err := that.getImageType(bt)
	if _err != nil {
		return "", _err
	}
	_prefix := ""
	switch _type {
	case PNG_TYPE:
		_prefix = "data:image/png;base64"
	case JPG_TYPE:
		_prefix = "data:image/jpeg;base64"
	case GIF_TYPE:
		_prefix = "data:image/gif;base64"
	case BMP_TYPE:
		_prefix = "data:image/bmp;base64"
	}
	return _prefix, nil
}

func (that *imgData) getImageSuffix(bt []byte) (string, error) {
	_type, _err := that.getImageType(bt)
	if _err != nil {
		return "", _err
	}
	_suffix := ""
	switch _type {
	case PNG_TYPE:
		_suffix = ".png"
	case JPG_TYPE:
		_suffix = ".jpg"
	case GIF_TYPE:
		_suffix = ".gif"
	case BMP_TYPE:
		_suffix = ".bmp"
	}
	return _suffix, nil
}

func (that *imgData) getImageType(bt []byte) (string, error) {
	var _type string
	if len(bt) < 8 {
		return "", errors.New("undefined type")
	}
	if bytes.Equal(PNG, bt[0:8]) {
		_type = PNG_TYPE
	}
	if bytes.Equal(GIF, bt[0:3]) {
		_type = GIF_TYPE
	}
	if bytes.Equal(BMP, bt[0:2]) {
		_type = BMP_TYPE
	}
	if bytes.Equal(JPG, bt[0:3]) {
		_type = JPG_TYPE
	}
	if _type == "" {
		return _type, errors.New("undefined type")
	} else {
		return _type, nil
	}
}

func (that imgData) path(dirName, activeName string) string {
	_imgNameSlice := []string{
		dirName,
		TimeStr("Ymd"),
		activeName[0:2],
		activeName[2:4],
		activeName + ".png",
	} //图像名称的前4位做目录
	return strings.Join(_imgNameSlice, "/") //保存图像目录
}
