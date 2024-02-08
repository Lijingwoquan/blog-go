package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const CtxUserIDKey = "UserID"

func getUserId(c *gin.Context) (id int64) {
	uid, exist := c.Get(CtxUserIDKey)
	if !exist {
		zap.L().Error("c.Get(CtxUserIDKey) failed", zap.Error(errors.New(userIDInvalid)))
		return
	}
	id, ok := uid.(int64)
	if !ok {
		fmt.Println(id)
		return
	}
	return
}
