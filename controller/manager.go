package controller

import (
	"blog/logic"
	"blog/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

const (
	addClassifySuccess    = "添加分类成功"
	addEssaySuccess       = "添加文章成功"
	updateEssaySuccess    = "修改文章成功"
	updateClassifySuccess = "修改分类成功"
	updateKindSuccess     = "修改大纲成功"
	deleteEssaySuccess    = "删除文章成功"
)

func AddClassifyHandler(c *gin.Context) {
	//1.参数处理
	var classify = new(models.ClassifyParams)
	fmt.Println(classify)
	if err := c.ShouldBindJSON(classify); err != nil {
		zap.L().Error("c.ShouldBindJSON(classify) failed", zap.Error(err))
		ResponseError(c, CodeParamInvalid)
		return
	}
	//2.业务处理
	err := logic.AddClassify(classify)
	if err != nil {
		zap.L().Error("mysql.AddClassify(classify) failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, addClassifySuccess)
}

func AddEssayHandler(c *gin.Context) {
	//1.参数处理
	var essay = new(models.EssayParams)
	if err := c.ShouldBindJSON(essay); err != nil {
		zap.L().Error("c.ShouldBindJSON(essay) failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	//2.业务处理
	if err := logic.CreateEssay(essay); err != nil {
		zap.L().Error("mysql.CreateEssay(essay) failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, addEssaySuccess)
}

func UpdateKindHandler(c *gin.Context) {
	//1.参数检验
	var k = new(models.UpdateKindParams)
	err := c.ShouldBindJSON(k)
	if err != nil {
		zap.L().Error("c.ShouldBindJSON(k) failed,err:", zap.Error(err))
		ResponseError(c, CodeParamInvalid)
		return
	}
	//2.逻辑处理
	err = logic.UpdateKind(k)
	if err != nil {
		zap.L().Error("logic.UpdateKind(k) failed,err:", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, updateKindSuccess)
}

func UpdateClassifyHandler(c *gin.Context) {
	//1.参数处理
	var classify = new(models.UpdateClassifyParams)
	err := c.ShouldBindJSON(classify)
	if err != nil {
		zap.L().Error("c.ShouldBindJSON(classify) failed,err:", zap.Error(err))
		ResponseError(c, CodeParamInvalid)
		return
	}
	//2.业务处理
	err = logic.UpdateClassify(classify)
	if err != nil {
		zap.L().Error("logic.UpdateClassify(classify) failed err:", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, updateClassifySuccess)
}

func UpdateEssayMSgHandler(c *gin.Context) {
	//1.获取参数
	var data = new(models.UpdateEssayMSg)
	if err := c.ShouldBindJSON(data); err != nil {
		zap.L().Error("c.ShouldBindJSON(data) failed", zap.Error(err))
		ResponseError(c, CodeParamInvalid)
		return
	}
	//2.业务处理
	if err := logic.UpdateEssayMsg(data); err != nil {
		zap.L().Error("mysql.UpdateEssay(data) failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, updateEssaySuccess)
}

func DeleteEssayHandler(c *gin.Context) {
	//1.获取参数
	idS := c.Query("id")
	id, err := strconv.Atoi(idS)
	if err != nil {
		zap.L().Error(" strconv.Atoi(idS) failed,err:", zap.Error(err))
		ResponseError(c, CodeParamInvalid)
		return
	}
	//2.逻辑处理
	err = logic.DeleteEssay(id)
	if err != nil {
		zap.L().Error("logic.DeleteEssay(id) failed,err:", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, deleteEssaySuccess)
}
