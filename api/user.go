package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"rare_earth_mining_BE/model"
	"rare_earth_mining_BE/service"
	"rare_earth_mining_BE/util"
	"regexp"
)

// 用于匹配mail的正则表达式
var mailreg = ".+@.+[.].+"
var mailregexp *regexp.Regexp = regexp.MustCompile(mailreg)

func Register(c *gin.Context) {

	mail := c.PostForm("mail")
	password := c.PostForm("password")
	nickname := c.PostForm("nickname")

	//mailregexp := regexp.MustCompile(mailreg)

	//判断提交的表单中的内容的格式是否符合要求
	//1.正则表达式判断是否是邮箱
	//2.昵称不应大于60字,要么就不填
	//3.密码不应小于6个字符
	if mailregexp.FindAllStringSubmatch(mail, 1) == nil || len(nickname) > 60 || len(password) < 6 {
		util.RespFormatError(c)
		return
	}

	//确保邮箱不重复
	_, err := service.SearchUser("mail", mail)

	//处理除了没查询到以外的错误
	if err != sql.ErrNoRows {
		if err == util.FieldsError {
			//处理数据库查询字段不符的错误
			util.RespFieldsMatchError(c)
		} else if err == nil {
			//处理查询出数据的结果(即邮箱重复)
			util.RespMailRepeated(c)
		} else if err != sql.ErrNoRows {
			//处理意料之外的错误
			util.RespUnexceptedError(c)
			fmt.Println(err)
		}
		return
	}

	//如果输入了昵称，确保昵称不重复
	if len(nickname) > 0 {
		_, err := service.SearchUser("nickname", nickname)

		//处理除了没查询到以外的错误
		if err != sql.ErrNoRows {
			if err == util.FieldsError {
				//处理数据库查询字段不符的错误
				util.RespFieldsMatchError(c)
			} else if err == nil {
				//处理查询出数据的结果(即昵称重复)
				util.RespNicknameRepeated(c)
			} else if err != sql.ErrNoRows {
				//处理意料之外的错误
				util.RespUnexceptedError(c)
			}
			return
		}
	} else {
		//若没输入昵称，则将"用户" + mail作为昵称
		nickname = "用户" + mail
	}

	err = service.CreateUser(model.User{
		Mail:     mail,
		Nickname: nickname,
		Password: password,
	})

	if err != nil {
		fmt.Println(err)
		util.RespCreateUserError(c)
		return
	}
	util.RespOK(c)
}

func Login(c *gin.Context) {
	mail := c.PostForm("mail")
	password := c.PostForm("password")

	//判断提交的表单中的内容的格式是否符合要求
	//1.正则表达式判断是否是邮箱
	//3.密码不应为空
	if mailregexp.FindAllStringSubmatch(mail, 1) == nil || len(password) < 1 {
		util.RespFormatError(c)
		return
	}

	_, userpassword, err := service.SearchUserPassword("mail", mail)

	if err != nil {
		if err == sql.ErrNoRows {
			//处理该用户不存在(回复邮箱或密码错误)
			util.RespIncorrectPassword(c)
		} else if err == util.FieldsError {
			//处理数据库查询字段不符的错误
			util.RespFieldsMatchError(c)
		} else if err != sql.ErrNoRows {
			//处理意料之外的错误
			util.RespUnexceptedError(c)
		}
		return
	}

	//处理密码不正确(回复邮箱或密码错误)
	if userpassword != password {
		//fmt.Println(userpassword)
		//fmt.Println(password)
		util.RespIncorrectPassword(c)
		return
	}

	//设置Cookie保持登录
	u, err := service.SearchUser("mail", mail)
	//fmt.Println(u.ID)
	//fmt.Println(strconv.Itoa(int(u.ID)))
	//c.SetCookie("uID", "uID"+strconv.Itoa(int(u.ID)), 604800, "", "/", false, false)

	//设置token
	SetToken(u.Mail, userpassword, c)

	//util.RespOK(c)
}

// 中间件
// 必须在登录状态才能访问
// 未登录则重定向到登录页面
func LoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {

		cookie, err := c.Cookie("uID")

		if err != nil {
			if err == http.ErrNoCookie {
				//处理未登录
				util.RespDidNotLogin(c)
			} else {
				//处理意外情况
				util.RespUnexceptedError(c)
			}
			//终止
			c.Abort()
			return
		}

		_, err = service.SearchUser("uID", cookie)

		if err != nil {
			if err == sql.ErrNoRows {
				//处理cookie有问题
				util.RespDidNotLogin(c)
			} else {
				//处理意外情况
				util.RespUnexceptedError(c)
			}
			//终止
			c.Abort()
			return
		}

		//若已登录则执行下个中间件
		c.Next()
		return
	}
}

// 中间件
// 必须在不登录状态才能访问
// 已登录则重定向到首页页面
func NotLoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {

		cookie, err := c.Cookie("uID")
		if err == nil {
			_, err = service.SearchUser("uID", cookie)
			if err == sql.ErrNoRows {
				//处理cookie有问题
				util.RespDidNotLogin(c)
			} else {
				//处理已登录
				util.RespLoggedin(c)
			}
			c.Abort()
			return
		}

		if err != http.ErrNoCookie {
			//处理未预料到的情况
			util.RespUnexceptedError(c)
			c.Abort()
			return
		}

		//若未登录则执行下个中间件
		c.Next()
		return
	}
}
