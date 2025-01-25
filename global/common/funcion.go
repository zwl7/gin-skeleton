package common

import (
	"gin-skeleton/global/consts"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseData struct {
	Code    consts.MyCode `json:"code"`
	Message string        `json:"message"`
	Data    interface{}   `json:"data"`
}

func ResponseErrorWithMsg(ctx *gin.Context, code consts.MyCode, errMsg string) {
	rd := &ResponseData{
		Code:    code,
		Message: errMsg,
		Data:    nil,
	}
	ctx.JSON(http.StatusOK, rd)
}

func ResponseSuccess(ctx *gin.Context, data interface{}, msg ...string) {

	if len(msg) <= 0 {
		msg = append(msg, consts.SuccessMsg)
	}

	rd := &ResponseData{
		Code:    consts.SuccessCode,
		Message: msg[0],
		Data:    data,
	}
	ctx.JSON(http.StatusOK, rd)
}
