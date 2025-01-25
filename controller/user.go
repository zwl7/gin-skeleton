package controller

import (
	"fmt"
	"gin-skeleton/database/redis"
	"gin-skeleton/global/common"
	"gin-skeleton/global/consts"
	"gin-skeleton/logic"
	"gin-skeleton/middleware/jwt"
	"gin-skeleton/translator"
	"github.com/gin-gonic/gin"
)

type SignUpForm struct {
	Name     string `form:"name" validate:"required,min=4,max=12" label:"用户名"`
	Age      int64  `form:"age" validate:"required" label:"年龄"`
	Password string `form:"password" validate:"required,min=6,max=20" label:"密码"`
}

func SignUp(c *gin.Context) {
	// 1.获取请求参数 2.校验数据有效性

	var params SignUpForm

	//if err := c.ShouldBind(&params); err != nil {
	//	//提示语优化
	//	common.ResponseErrorWithMsg(c, consts.InvalidParamsCode, "参数解析失败："+err.Error())
	//	return
	//}

	_ = c.ShouldBind(&params)

	msg, code := translator.Validate(&params)
	//验证报错，返回错误信息给前端
	if code != 200 {
		common.ResponseErrorWithMsg(c, consts.InvalidParamsCode, msg)
		return
	}

	//开始逻辑处理
	logic.CreateUser(c, params.Name, params.Password, params.Age)
	return
}

// `json:"name"  只可以接受json body ，`form:"name" 可以接受form和json格式
type LoginForm struct {
	Name     string `form:"name" validate:"required,min=4,max=12" label:"用户名"`
	Password string `form:"password" validate:"required,min=6,max=20" label:"密码"`
}

func Login(c *gin.Context) {

	var params LoginForm
	_ = c.ShouldBind(&params)

	msg, code := translator.Validate(&params)
	//验证报错，返回错误信息给前端
	if code != 200 {
		common.ResponseErrorWithMsg(c, consts.InvalidParamsCode, msg)
		return
	}

	//逻辑层处理
	logic.Login(c, params.Name, params.Password)

	return

}

func RefreshToken(c *gin.Context) {
	//rt := c.Query("refresh_token")
	//// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
	//// 这里假设Token放在Header的Authorization中，并使用Bearer开头
	//// 这里的具体实现方式要依据你的实际业务情况决定
	//authHeader := c.Request.Header.Get("Authorization")
	//if authHeader == "" {
	//	common.ResponseErrorWithMsg(c, consts.InvalidTokenCode, "请求头缺少Auth Token")
	//	c.Abort()
	//	return
	//}
	//// 按空格分割
	//parts := strings.SplitN(authHeader, " ", 2)
	//if !(len(parts) == 2 && parts[0] == "Bearer") {
	//	common.ResponseErrorWithMsg(c, consts.InvalidTokenCode, "Token格式不对")
	//	c.Abort()
	//	return
	//}
	//aToken, rToken, err := jwt.RefreshToken(parts[1], rt)
	//fmt.Println(err)
	//c.JSON(http.StatusOK, gin.H{
	//	"access_token":  aToken,
	//	"refresh_token": rToken,
	//})
}

func GetList(c *gin.Context) {

}

func Get(c *gin.Context) {

	fmt.Println()
	userID, isExist := c.Get(jwt.ContextUserIDKey)
	if !isExist {
		common.ResponseErrorWithMsg(c, consts.ErrorCode, consts.ErrorMsg)
	}

	//getUserInfo
	logic.GetUserInfo(c, userID.(int64))
	fmt.Println(redis.Client.PoolStats())
	return
}

func Update(c *gin.Context) {

}

func Del(c *gin.Context) {

}

func Test(c *gin.Context) {
	//
	fmt.Println("test-------------")
	logic.Test(c)
	return
}

type TestForm struct {
	Time int64 `form:"time" validate:"required" label:"次数"`
}

func TestMemory(c *gin.Context) {
	//

	var params TestForm
	_ = c.ShouldBind(&params)

	msg, code := translator.Validate(&params)
	//验证报错，返回错误信息给前端
	if code != 200 {
		common.ResponseErrorWithMsg(c, consts.InvalidParamsCode, msg)
		return
	}
	fmt.Println("TestMemory-------------")
	logic.TestMemory(c, params.Time)
	return
}
