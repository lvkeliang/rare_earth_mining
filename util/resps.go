package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type respTemplate struct {
	Status int    `json:"status"`
	Info   string `json:"info"`
	Data   []byte `json:"data"` //json序列化后的字符串
}

// 未预料到的错误(90000)
var UnexceptedError = respTemplate{
	Status: 90000,
	Info:   "There is an unexcepted error",
}

func RespUnexceptedError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, UnexceptedError)
}

// 正常(1xxxx)
var OK = respTemplate{
	Status: 10000,
	Info:   "success",
}

// 回复正常
func RespOK(c *gin.Context) {
	c.JSON(http.StatusOK, OK)
}

// 客户端操作有误(2xxxx)
// 输入格式有误
var FormatError = respTemplate{
	Status: 20001,
	Info:   "The input format is incorrect",
}

// 回复输入格式错误
func RespFormatError(c *gin.Context) {
	c.JSON(http.StatusForbidden, FormatError)
}

// 个人信息错误(3xxxx)
// 创建的昵称重复
var NicknameRepeated = respTemplate{
	Status: 30001,
	Info:   "The nickname is repeated",
}

// 回复创建的昵称重复
func RespNicknameRepeated(c *gin.Context) {
	c.JSON(http.StatusOK, NicknameRepeated)
}

// 邮箱重复
var MailRepeated = respTemplate{
	Status: 30002,
	Info:   "The mail is already in use",
}

// 回复邮箱重复
func RespMailRepeated(c *gin.Context) {
	c.JSON(http.StatusOK, MailRepeated)
}

// 用户不存在
var UserNotExist = respTemplate{
	Status: 30003,
	Info:   "The user does not exist",
}

// 回复用户不存在
func RespUserNotExist(c *gin.Context) {
	c.JSON(http.StatusOK, UserNotExist)
}

// 密码不正确
var IncorrectPassword = respTemplate{
	Status: 30004,
	Info:   "The mail or password is incorrect",
}

func RespIncorrectPassword(c *gin.Context) {
	c.JSON(http.StatusOK, IncorrectPassword)
}

// 服务器错误(5xxxx)
// 数据库查询错误(查询的字段不符）
var FieldsMatchError = respTemplate{
	Status: 50001,
	Info:   "An error occurred with the server",
}

// 回复数据库查询错误
func RespFieldsMatchError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, FieldsMatchError)
}

// 创建用户出错
var CreateUserError = respTemplate{
	Status: 50002,
	Info:   "An error occurred with the server when creating user",
}

// 回复创建用户出错
func RespCreateUserError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, CreateUserError)
}

var ParamError = respTemplate{
	Status: 300,
	Info:   "params error",
}

var InternalError = respTemplate{
	Status: 500,
	Info:   "internal error",
}

/*func NormError(c *gin.Context, status int, info string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"status": status,
		"info":   info,
	})
}*/
