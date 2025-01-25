package controller

import (
	"gin-skeleton/global/variable"
	"gin-skeleton/middleware/jwt"
	"github.com/gin-gonic/gin"
)

func getCurrentUserID(c *gin.Context) (userID uint64, err error) {
	_userID, ok := c.Get(jwt.ContextUserIDKey)
	if !ok {
		err = variable.ErrorUserNotLogin
		return
	}
	userID, ok = _userID.(uint64)
	if !ok {
		err = variable.ErrorUserNotLogin
		return
	}
	return
}
