package controller

import (
	"github.com/gin-gonic/gin"
)

const CtxUserIDKey = "UserID"

func getUserId(c *gin.Context) (id int64, err error) {
	uid, exist := c.Get(CtxUserIDKey)
	if !exist {
		return 0, err
	}
	var ok bool
	id, ok = uid.(int64)
	if !ok {
		return 0, err
	}
	return
}
