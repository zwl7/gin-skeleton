package router

import (
	"gin-skeleton/controller"
	"gin-skeleton/logger"
	"gin-skeleton/middleware/jwt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func SetupRouter() *gin.Engine {
	//gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	//r.Run()
	r.Use(logger.GinLogger())
	//r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/", func(ctx *gin.Context) {
		time.Sleep(5 * time.Second)
		ctx.String(http.StatusOK, "ok")
	})

	r.GET("/version", func(ctx *gin.Context) {
		ctx.Request.Context()
		ctx.String(http.StatusOK, viper.GetString("version"))
	})

	//生产环境，最好用nginx处理静态资源
	r.Static("/storage/upload", "./storage/upload") //  定义静态资源路由与实际目录映射关系

	//r.POST("/signup", controller.SignUpHandler)
	//r.GET("/getUserInfo", middlewares.JWTAuthMiddleware(), controller.GetUserInfo)

	//  创建一个api路由组
	vApi := r.Group("/api/")
	{
		// 模拟一个用户路由
		user := vApi.Group("user/")
		{
			user.GET("test", controller.Test)

			user.GET("testMemory", controller.TestMemory)

			user.POST("login", controller.Login)

			user.POST("signup", controller.SignUp)

			user.GET("refresh_token", controller.RefreshToken)
		}

		vApi.POST("file/upload", controller.FileUpload)

		//需要jwt验证的url
		vApi.Use(jwt.JWTAuthMiddleware())
		{
			vApi.GET("user/getList", controller.GetList)
			vApi.GET("user/get", controller.Get)

			vApi.POST("user/update", controller.Update)
			vApi.POST("user/del", controller.Del)
		}

	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
