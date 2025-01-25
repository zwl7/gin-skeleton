package utils

import "fmt"

// 响应解析，返回统一数据格式
func Parse(data interface{}) map[string]interface{} {
	_code := 200
	_msg := "SUCCESS"
	var _rsData interface{}
	switch _t := data.(type) {
	case UtilsException:
		fmt.Println(_t)
		_code = ToInt(_t.GetCode())
		_msg = _t.GetMessage()
		_rsData = _t.GetData()
	case error:
		_code = 400
		_msg = _t.Error()
		_rsData = ""
	default:
		_rsData = data
	}
	return map[string]interface{}{
		"code":    _code,
		"message": _msg,
		"data":    _rsData,
	}
}
