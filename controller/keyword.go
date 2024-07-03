package controller

import (
	"blog/logic"
	"blog/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func IncreaseSearchKeywordHandler(c *gin.Context) {
	//1.参数检验
	keyword := new(models.KeywordParam)
	if err := c.ShouldBindJSON(keyword); err != nil {
		zap.L().Error("c.ShouldBindJSON(keyword) failed", zap.Error(err))
		ResponseError(c, CodeParamInvalid)
		return
	}

	//2.逻辑处理
	if err := logic.IncreaseSearchKeyword(keyword); err != nil {
		zap.L().Error("logic.IncreaseSearchKeyword(keyword) failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}

	//3.返回响应
	ResponseSuccess(c, nil)
}

func GetSearchKeywordRankHandel(c *gin.Context) {
	rankKind := new(models.KeywordRankKind)
	// 逻辑处理
	if err := logic.GetSearchKeywordRank(rankKind); err != nil {
		zap.L().Error("logic.GetSearchKeywordRank(rankKind) failed,err:", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	//	返回响应
	ResponseSuccess(c, rankKind)
}
