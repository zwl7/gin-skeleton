package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LogWriter struct {
	logDir           string //日志根目录地址。
	module           string //模块 名
	curFileName      string //当前被指定的filename
	curBaseFileName  string //在使用中的file
	turnCateDuration time.Duration
	mutex            sync.RWMutex
	outFh            *os.File
}

func (w *LogWriter) Write(p []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if out, err := w.getWriter(); err != nil {
		return 0, errors.New("failed to fetch target io.Writer")
	} else {
		return out.Write(p)
	}
}

func (w *LogWriter) getFileName() string {
	base := time.Now().Truncate(w.turnCateDuration)
	return fmt.Sprintf("%s/%s/%s_%s", w.logDir, base.Format("2006-01-02"), w.module, base.Format("15"))
}

func (w *LogWriter) getWriter() (io.Writer, error) {
	fileName := w.curBaseFileName
	//判断是否有新的文件名
	//会出现新的文件名
	baseFileName := w.getFileName()
	if baseFileName != fileName {
		fileName = baseFileName
	}

	dirname := filepath.Dir(fileName)
	if err := os.MkdirAll(dirname, 0755); err != nil {
		return nil, fmt.Errorf("error:%s , failed to create directory %s", err.Error(), dirname)
	}

	fileHandler, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s", err.Error())
	}
	w.outFh.Close()
	w.outFh = fileHandler
	w.curBaseFileName = fileName
	w.curFileName = fileName

	return fileHandler, nil
}

func NewLoggerWriter(logPath, module string, duration time.Duration) *LogWriter {
	return &LogWriter{
		logDir:           logPath,
		module:           module,
		turnCateDuration: duration,
		curFileName:      "",
		curBaseFileName:  "",
	}
}
