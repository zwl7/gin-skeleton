/**
 * @desc		消息中间件：消费推送
 * ---------------------------------------------------------------------
 * @author		unphp <unphp@qq.com>
 * @date		2018-11-27
 * @copyright	rabbitmq-consumer 0.1
 * ---------------------------------------------------------------------
 */

package utils

import (
	"fmt"
	//base "unphp/go-utils/base"
	//"os"
	//"runtime"
	//"strconv"
	//"strings"
	//"time"
)

// NewLog ...
func NewLog(mode int) *Log {
	SysLog := &Log{
		LogChan:        NewStackChanPool(),
		Mode:           mode,
		LogHanderSlice: make([]HanderLogInterface, 0),
	}
	if mode == 1 {
		SysLog.AddLogHander(SysLog)
	}
	return SysLog
}

// HanderLogInterface sss
type HanderLogInterface interface {
	LogHander(log interface{})
}

// Log 系统服务日志
type Log struct {
	LogChan        StackChanInterface
	Mode           int
	LogHanderSlice []HanderLogInterface
}

// AddLogHander sss
func (that *Log) AddLogHander(h HanderLogInterface) *Log {
	that.LogHanderSlice = append(that.LogHanderSlice, h)
	return that
}

// LogHander 日志
func (that *Log) LogHander(log interface{}) {

	switch data := log.(type) {
	case []byte:
		fmt.Println(string(data))
	case string:
		fmt.Println(data)
	case []interface{}:
		switch _data := data[0].(type) {
		case interface{}:
			fmt.Println(_data)
		}
	case interface{}:
		fmt.Println(data)
	}
	//fmt.Println(log)
}

// Start 开始
func (that *Log) Start(master bool) {
	_func := func() {
		for {
			logSlice := that.LogChan.Get()
			for _, log := range logSlice {
				for _, hander := range that.LogHanderSlice {
					hander.LogHander(log)
				}
			}

		}

	}
	if master {
		_func()
	} else {
		go _func()
	}
}

// Println 日志记录
func (that *Log) Println(log ...interface{}) *Log {
	that.LogChan.Stack(log)
	return that
}
