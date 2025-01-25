package utils

import (
	"net/http"
	"os"
)

// GetMimeTypeByPath  获取文件的MIME类型
func GetMimeTypeByPath(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	//只读模块打开，不需要close
	//defer file.Close()

	// 读取文件的前 512 个字节用于检测 MIME 类型
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}

	// 获取文件的 MIME 类型
	mimeType := http.DetectContentType(buffer)

	return mimeType, nil
}

// GetMimeTypeBySlice 获取文件的MIME类型
func GetMimeTypeBySlice(slice []byte) (string, error) {

	// 获取文件的 MIME 类型
	mimeType := http.DetectContentType(slice)
	return mimeType, nil
}

// 验证文件MIME类型的函数
func IsValidFileType(fileType string) bool {
	// 在这里添加允许的MIME类型
	allowedFileTypes := map[string]bool{
		//img
		"image/jpeg":    true,
		"image/png":     true,
		"image/gif":     true,
		"image/x-icon":  true,
		"image/svg+xml": true,

		//video
		"video/mp4":       true,
		"video/x-msvideo": true,
		"audio/x-wav":     true,
		"video/x-ms-wmv":  true,

		//音频
		"audio/mpeg":  true, //mp3
		"audio/x-m4a": true, //mp4

		//文本
		"application/pdf":    true, //pdf
		"text/plain":         true, //txt
		"application/msword": true, //doc
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true, //docx
		"application/x-gzip":       true, //GZ 压缩文件格式
		"application/zip":          true, //ZIP 压缩文件格式
		"application/rar":          true, //RAR 压缩文件格式
		"application/vnd.ms-excel": true, //微软 Office Excel 格式（Microsoft Excel 97 - 2004 Workbook
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         true, //微软 Office Excel 文档格式
		"application/vnd.ms-powerpoint":                                             true, //微软 Office PowerPoint 格式（Microsoft PowerPoint 97 - 2003 演示文稿）
		"application/vnd.openxmlformats-officedocument.presentationml.presentation": true, //微软 Office PowerPoint 文稿格式

		"application/kswps": true, //wps 金山 Office 文字排版文件格式
		"application/ksdps": true, //dps	金山 Office 演示文稿格式
		"application/kset":  true, //金山 Office 表格文件格式

		//其他
		"application/octet-stream": true,
	}
	return allowedFileTypes[fileType]
}
