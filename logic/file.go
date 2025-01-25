package logic

import (
	"fmt"
	"gin-skeleton/global/common"
	"gin-skeleton/global/consts"
	"gin-skeleton/global/variable"
	"gin-skeleton/utils/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func Upload(ctx *gin.Context, file *multipart.File, fileSuffixType string, fileSize int64) {

	//fmt.Println(variable.ServerDomain)
	//验证---start

	//大小判断
	if fileSize > consts.FileMaxSize {
		common.ResponseErrorWithMsg(ctx, consts.InvalidLogicCode, "上传的图像大小不能超过"+utils.SizeFormat(consts.FileMaxSize))
		return
	}

	//判断文件后缀是否允许上传
	if strings.Contains(consts.AllowFileSuffix, fileSuffixType) == false {
		//
		common.ResponseErrorWithMsg(ctx, consts.InvalidLogicCode, "上传格式不允许，只允许上传上传："+consts.AllowFileSuffix)
		return
	}
	//GetMimeType
	headSlice := make([]byte, 512)
	length, _err := (*file).Read(headSlice)
	fmt.Println(length)
	if _err != nil {
		//fmt.Println(_err)
		common.ResponseErrorWithMsg(ctx, consts.InvalidLogicCode, "无效的文件格式:"+_err.Error())
		return
	}

	//fmt.Println(headSlice)
	mimeType, err := utils.GetMimeTypeBySlice(headSlice)
	fmt.Println(mimeType)
	if err != nil {
		common.ResponseErrorWithMsg(ctx, consts.InvalidLogicCode, "无效的mime文件格式:"+_err.Error())
		return
	}

	if !utils.IsValidFileType(mimeType) {
		common.ResponseErrorWithMsg(ctx, consts.InvalidLogicCode, "不支持的mime格式:"+_err.Error())
		return
	}

	//重置文件的偏移量
	_, err = (*file).Seek(0, 0)
	if err != nil {
		common.ResponseErrorWithMsg(ctx, consts.InvalidLogicCode, "不支持的mime格式 "+_err.Error())
		return
	}

	//验证---end

	//文件的名称
	newFilename := utils.UniqueId() + "." + fileSuffixType

	//创建文件夹 txt/20240104
	_filePath := fileSuffixType + "/" + time.Now().Format("20060102")

	///1232131.txt
	_fileName := "/" + newFilename

	///storage/upload/txt/20240104
	_fileUrl := "/storage/upload/" + _filePath

	///Users/zwl/go/src/gin-skeleton
	dir, _ := os.Getwd()

	///Users/zwl/go/src/gin-skeleton/storage/upload/jpeg/20240104/eb0aa33d4a5e4533b3d7bd175e5e3dd6.jpeg
	dirFilePath := dir + _fileUrl
	isEx, err := utils.PathExists(dirFilePath)
	if err != nil {
		//
		zap.L().Error("系统繁忙:" + err.Error())
		common.ResponseErrorWithMsg(ctx, consts.ErrorCode, consts.ErrorMsg)
		return
	}
	if !isEx {
		err = os.MkdirAll(dirFilePath, 0777)
		if err != nil {
			zap.L().Error("系统繁忙:" + err.Error())
			common.ResponseErrorWithMsg(ctx, consts.ErrorCode, consts.ErrorMsg)
			return
		}
	}
	//创建文件
	out, err := os.Create(dirFilePath + _fileName)
	if err != nil {
		//获取错误文件和错误行
		_, errorFileInfo, line, _ := runtime.Caller(0)
		common.ResponseErrorWithMsg(ctx, consts.ErrorCode, errorFileInfo+":"+strconv.Itoa(line)+",上传错误：%s"+err.Error())
		return
	}

	defer out.Close()
	_, err = io.Copy(out, *file)
	if err != nil {
		_, errorFileInfo, line, _ := runtime.Caller(0) //获取错误文件和错误行
		common.ResponseErrorWithMsg(ctx, consts.ErrorCode, errorFileInfo+":"+strconv.Itoa(line)+",上传错误：%s"+err.Error())
		return
	}

	common.ResponseSuccess(ctx, variable.Mp{"size": utils.SizeFormat(fileSize), "file_path": _fileUrl + _fileName, "file_url": variable.ServerDomain + _fileUrl + _fileName}, "success")
	return
}
