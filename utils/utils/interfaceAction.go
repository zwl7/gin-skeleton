/**
 * 工具
 * @desc 模块对外接口、异常类
 * ---------------------------------------------------------------------
 * @author		Super <super@papa.com.cn>
 * @date		2019-11-07
 * @copyright	cooper
 * ---------------------------------------------------------------------
 */

package utils

import (
	"strings"
)

func NewUtilsDeferFunc(f ...func()) UtilsDeferFunc {
	_slice := []func(){}
	if len(f) > 0 {
		_slice = append(_slice, f...)
	}
	return _slice
}

type UtilsDeferFunc []func()

func (that UtilsDeferFunc) Add(f func()) UtilsDeferFunc {
	that = append(that, f)
	return that
}

func (that UtilsDeferFunc) Run(asyncRun ...bool) {
	if len(that) == 0 {
		return
	}
	//
	//fmt.Println("========= UtilsDeferFunc 001 =======", len(that))
	_run := func() {
		_len := len(that)
		for _i := 1; _i <= _len; _i++ {
			_func := that[_len-_i]
			//fmt.Println("========= UtilsDeferFunc 002 =======", _i)
			func() {
				defer func() {
					if _err := recover(); _err != nil {
					}
				}()
				_func()
			}()
		}
	}
	if len(asyncRun) > 0 && asyncRun[0] {
		_run()
		return
	}
	go _run()
}

// UtilsAction 模块统一遵循的 对外接口
type UtilsAction interface {
	Action(appID int64, params ...interface{}) (interface{}, UtilsException)
	LogicAction() interface{}
}

// UtilsMiddleware ...
type UtilsMiddleware interface {
	//Action(appID int64) (Mp, UtilsException)
	SetRequestMiddleware(func(appID int64, params ...interface{}) interface{}) UtilsMiddleware
	SetResponMiddleware(func(appID int64, mp interface{}, excp UtilsException, params ...interface{}) (interface{}, UtilsException)) UtilsMiddleware
}

// ActionStruct 上游协议层调用的接口
type ActionStruct struct {
	//AppID             int64
	logicMiddleware   func() interface{}
	requestMiddleware func(appID int64, params ...interface{}) interface{}                                                        //模块包自定义方法（处理协议层请求，并返回响应结果）
	responMiddleware  func(appID int64, mp interface{}, excp UtilsException, params ...interface{}) (interface{}, UtilsException) //模块包自定义中间件方法(对返回响应进行处理)
}

func (that *ActionStruct) LogicAction() interface{} {
	return that.logicMiddleware()
}

// Action 上游协议层最终调用的接口方法
func (that *ActionStruct) Action(appID int64, params ...interface{}) (interface{}, UtilsException) {
	//that.AppID = appID
	//模块包 处理请求
	_respon := that.requestMiddleware(appID, params...)
	_mp, _excp := that._parseRespon(_respon)
	//模块包 响应处理中间件
	if that.responMiddleware != nil {
		return that.responMiddleware(appID, _mp, _excp, params...)
	}
	return _mp, _excp
}

//func (that *ActionStruct) Do() interface{} {
//	//that.AppID = appID
//	//模块包 处理请求
//	_respon := that.requestMiddleware(appID, params...)
//	_mp, _excp := that._parseRespon(_respon)
//	//模块包 响应处理中间件
//	if that.responMiddleware != nil {
//		return that.responMiddleware(appID, _mp, _excp, params...)
//	}
//	return _mp, _excp
//}

func (that *ActionStruct) SetLogicMiddleware(logicMiddleware func() interface{}) UtilsMiddleware {
	that.logicMiddleware = logicMiddleware
	return that
}

// SetRequestMiddleware 设置对协议层请求处理的句柄方法（模块包）
func (that *ActionStruct) SetRequestMiddleware(requestMiddleware func(appID int64, params ...interface{}) interface{}) UtilsMiddleware {
	that.requestMiddleware = requestMiddleware
	return that
}

// SetResponMiddleware 设置对协议层响应的中间件方法
func (that *ActionStruct) SetResponMiddleware(responMiddleware func(appID int64, mp interface{}, excp UtilsException, params ...interface{}) (interface{}, UtilsException)) UtilsMiddleware {
	that.responMiddleware = responMiddleware
	return that
}

// _parseRespon 解析响应，生成UtilsAction.Action规范的返回格式
func (that *ActionStruct) _parseRespon(rs interface{}) (mp interface{}, excp UtilsException) {
	switch t := rs.(type) {
	case error:
		excp = ThrowException(40003001, t)
		return
	case Exception:
		excp = ThrowException(40003001, t)
		return
	case Mp:
		mp = t
		return
	case map[string]interface{}:
		mp = t
		return
	case []interface{}:
		mp = t
		return
	case string:
		mp = NewMP().Set("result", t)
		return
	case bool:
		mp = NewMP().Set("result", t)
		return
	default:
		//_jsonByte, _err := json.Marshal(t)
		//fmt.Println(reflect.TypeOf(t), "##################")
		//if _err != nil {
		//	_msg := "action marshal error: " + _err.Error()
		//	excp = ThrowException(40003002, errors.New(_msg))
		//}
		//_jsonMap := make(map[string]interface{}, 0)
		//if _err := json.Unmarshal([]byte(_jsonByte), &_jsonMap); _err != nil {
		//	_msg := "action unmarshal error: " + _err.Error()
		//	excp = ThrowException(40003002, errors.New(_msg))
		//	return
		//}
		mp = t
		return
	}
}

// UtilsException 模块异常接口
type UtilsException interface {
	GetCode() int64
	GetMessage() string
	GetData() interface{}
	GetTrace() []UtilsException
	GetTraceAsString() string
}

type ExceptionData struct {
	Data interface{}
}

func NewExceptionData(data interface{}) ExceptionData {
	return ExceptionData{
		data,
	}
}

// ThrowException 抛出异常
func ThrowException(code int64, errSlice ...interface{}) UtilsException {
	var _excp Exception
	_excp.code = code
	for _i, _param := range errSlice {
		switch _t := _param.(type) {
		case string:
			_excp.msg = _t
			_excp.tract = []UtilsException{
				_excp,
			}
		case error:
			_excp.msg = _t.Error()
			_excp.tract = []UtilsException{
				_excp,
			}
		case Exception:
			if _i == 0 {
				_excp = _t
			} else {
				_excp.code = code
				_excp.msg = "内部异常"
				_excp.tract = append(_t.tract, _excp)
			}
		case ExceptionData:
			_excp.data = _t.Data
		case nil:
			_excp.msg = ""
		}
	}
	return _excp
}

// Exception 异常
type Exception struct {
	code  int64
	msg   string
	data  interface{}
	tract []UtilsException
}

// GetCode 获取异常错误码
func (that Exception) GetCode() int64 {
	return that.code
}

// GetMessage 获取异常消息
func (that Exception) GetMessage() string {
	return that.msg
}

func (that Exception) GetData() interface{} {
	return that.data
}

// GetTrace ...
func (that Exception) GetTrace() []UtilsException {
	return that.tract
}

// GetTraceAsString ...
func (that Exception) GetTraceAsString() string {
	_traceSlice := make([]string, 0)
	for _, _e := range that.tract {
		_traceSlice = append(_traceSlice, ToString(_e.GetCode())+":"+_e.GetMessage())
	}
	return strings.Join(_traceSlice, " | ")
}
