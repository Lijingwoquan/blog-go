package controller

import (
	"blog/logic"
	"blog/models"
	"github.com/gin-gonic/gin"
)

func ResponseIndexDataHandler(c *gin.Context) {
	var data = new(models.IndexData)
	logic.GetIndexData(&data)
	ResponseSuccess(c, *data)
}
