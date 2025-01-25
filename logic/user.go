package logic

import (
	"context"
	"fmt"
	"gin-skeleton/database/redis"
	"gin-skeleton/global/common"
	"gin-skeleton/global/consts"
	"gin-skeleton/global/variable"
	"gin-skeleton/models"
	"gin-skeleton/utils/jwt"
	"gin-skeleton/utils/md5_encrypt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type UserList map[string]interface{}

func CreateUser(c *gin.Context, name, password string, age int64) {

	//逻辑处理
	userModel := new(models.User)

	userModel.Name = name
	userModel.Age = age
	userModel.IsDel = 2
	userModel.Password = md5_encrypt.MD5(password)

	// 3.注册用户
	if err := userModel.CreateUser(userModel); err != nil {
		zap.L().Fatal("错误原因: " + err.Error())
		common.ResponseErrorWithMsg(c, consts.ErrorCode, consts.ErrorMsg)
		return
	}

	common.ResponseSuccess(c, true, "注册成功")
	return

}

func Login(c *gin.Context, name, password string) {
	//逻辑处理
	_userModel := new(models.User)

	_isExist, _err := _userModel.GetOneByName(name, _userModel)
	if _err != nil {
		zap.L().Fatal("错误原因: " + _err.Error())
		common.ResponseErrorWithMsg(c, consts.ErrorCode, consts.ErrorMsg)
		return
	}

	if !_isExist {
		//错误的用户名 用户名不存在
		common.ResponseErrorWithMsg(c, consts.ErrorCode, "账号或密码错误。")
		return
	}

	//密码错误
	if _userModel.Password != md5_encrypt.MD5(password) {
		common.ResponseErrorWithMsg(c, consts.InvalidLogicCode, "账号或密码错误.")
		return
	}

	aToken, _err := jwt.GetToken(_userModel.ID)
	if _err != nil {
		zap.L().Fatal("错误原因: " + _err.Error())
		common.ResponseErrorWithMsg(c, consts.ErrorCode, consts.ErrorMsg)
		return

	}

	//test
	userlists := []UserList{
		UserList{"name": "zwl", "age": 10},
		UserList{"name": "yxj", "age": 11},
		UserList{"name": "kll", "age": 12},
	}
	common.ResponseSuccess(c, variable.Mp{"token": aToken, "userName": _userModel.Name, "userId": _userModel.ID, "array": []int{1, 2, 3}, "arrayTwo": userlists}, "登录成功")
	return
}

func GetUserInfo(c *gin.Context, userId int64) {
	//

	_userModel := new(models.User)
	isExist, _err := _userModel.GetOneById(userId, _userModel)
	if _err != nil {
		zap.L().Fatal("错误原因: " + _err.Error())
		common.ResponseErrorWithMsg(c, consts.ErrorCode, consts.ErrorMsg)
		return
	}

	if !isExist {
		common.ResponseErrorWithMsg(c, consts.InvalidLogicCode, "无效的用户id")
		return
	}

	ctx := context.Background()
	fmt.Println(redis.Client.Set(ctx, redis.KeyDemo, "zwl", -1).Result())
	fmt.Println(redis.Client.Set(ctx, redis.KeyTest, "zwltest", 60*time.Second).Result())

	exist, val := redis.GetKey(ctx, redis.KeyDemo)
	if !exist {
		common.ResponseErrorWithMsg(c, consts.InvalidLogicCode, "无效的redis key")
		return
	}
	fmt.Println("KeyDemo : ", val)

	exist, val = redis.GetKey(ctx, redis.KeyTest)
	if !exist {
		common.ResponseErrorWithMsg(c, consts.InvalidLogicCode, "无效的redis KeyTest")
		return
	}
	fmt.Println("KeyTest : ", val)

	common.ResponseSuccess(c, variable.Mp{"name": _userModel.Name, "userId": _userModel.ID, "age": _userModel.Age}, "ok")
	return
}

func Test(c *gin.Context) {
	//查询一条记录
	_userModel := new(models.User)
	_userModel.GetOne(_userModel)

	ctx := context.Background()
	exist, val := redis.GetKey(ctx, "qwe")
	if !exist {
		common.ResponseErrorWithMsg(c, consts.InvalidLogicCode, "无效的redis key")
		return
	}
	fmt.Println("zwl : ", val)
	common.ResponseSuccess(c, "success", "ok")
	return
}

func TestMemory(c *gin.Context, times int64) {
	//查询一条记录

	//var CompanyAreas []models.CompanyAreas
	//
	//CompanyAreasModel := new(models.CompanyAreas)
	//
	//var i int64 = 0
	//
	//for i <= times {
	//
	//	_ = CompanyAreasModel.GetAll(&CompanyAreas)
	//
	//	i++
	//}
	//
	//fmt.Println(len(CompanyAreas))
	//
	//time.Sleep(3 * time.Second)

	var i int64 = 0

	var str string

	for i < times {

		str += "kkkkkkkkkkkkk"

		i++
	}

	//fmt.Println(len(str))
	common.ResponseSuccess(c, str, "ok")
	return
}
