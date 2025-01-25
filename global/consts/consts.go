package consts

// 这里定义的常量，一般是具有错误代码+错误说明组成，一般用于接口返回
type MyCode int64

const (
	SuccessCode       MyCode = 200
	InvalidParamsCode MyCode = 410

	InvalidTokenCode MyCode = 411
	InvalidLogicCode MyCode = 412

	ErrorCode MyCode = 500

	SuccessMsg       string = "成功"
	ErrorMsg         string = "系统繁忙请稍后再试"
	InvalidTokenMsg  string = "token无效"
	InvalidParamsMsg string = "参数错误"

	AllowFileSuffix string = "doc,docx,xls,xlsx,pdf,mp4,txt,jpg,jpeg,png,gif"
	FileMaxSize     int64  = 1024 * 1024 * 5 //5m
)
