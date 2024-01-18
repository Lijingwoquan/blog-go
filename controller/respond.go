package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type codes int64

const (
	CodeSuccess codes = 1000 + iota
	CodeParamInvalid
	CodeUserExist
	CodeLoginFailed
	CodeInvalidToken
	CodeServeBusy
)

var CodeMsgMap = map[codes]string{
	CodeSuccess:      "请求成功",
	CodeParamInvalid: "参数错误",
	CodeUserExist:    "用户已存在",
	CodeLoginFailed:  "登录失败",
	CodeInvalidToken: "token失效",
	CodeServeBusy:    "服务繁忙",
}

type Response struct {
	Code codes
	Msg  interface{}
	Data interface{}
}

// Msg --> code 的msg方法返回字符串
func (code codes) Msg() (msg string) {
	var ok bool
	msg, ok = CodeMsgMap[code]
	if !ok {
		msg = CodeMsgMap[CodeServeBusy]
		return
	}
	return
}

func ResponseError(c *gin.Context, code codes) {
	rd := &Response{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	}
	c.JSON(http.StatusOK, rd)
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	rd := &Response{
		Code: CodeSuccess,
		Msg:  CodeMsgMap[CodeSuccess],
		Data: data,
	}
	c.JSON(http.StatusOK, rd)
}
