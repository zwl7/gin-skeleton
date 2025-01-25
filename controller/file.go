package controller

import (
	"gin-skeleton/global/common"
	"gin-skeleton/global/consts"
	"gin-skeleton/logic"
	"github.com/gin-gonic/gin"
	"strings"
)

func FileUpload(c *gin.Context) {

	//图片是否需要从令牌格式中隔离
	//_authStr := ctx.GetHeader("Authorization")
	//_, _jwtParams, _ := checkJwtHander(_authStr)
	//_superDir := ""
	//if _jwtParams.Has("company_id") {
	//	_superDir = utils.ToString(_jwtParams.Get("company_id"), true)
	//}
	//

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		common.ResponseErrorWithMsg(c, consts.InvalidParamsCode, "上传出现错误"+err.Error())
		return
	}

	//允许上传的文件类型
	filename := strings.Split(header.Filename, ".")
	// png jpg txt
	filenameSuffix := filename[len(filename)-1]

	logic.Upload(c, &file, filenameSuffix, header.Size)
	return
}
