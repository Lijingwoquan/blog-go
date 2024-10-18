package controller

import (
	"blog/logic"
	"blog/models"
	"github.com/gin-gonic/gin"
)

func ResponseIndexDataHandler(c *gin.Context) {
	var data = new(models.DataAboutIndex)
	var err error
	data, err = logic.GetIndexData()
	if err != nil {
		ResponseError(c, CodeServeBusy)
		return
	}
	ResponseSuccess(c, *data)
}
