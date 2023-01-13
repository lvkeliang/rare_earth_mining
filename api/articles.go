package api

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"rare_earth_mining_BE/model"
	"rare_earth_mining_BE/service"
	"rare_earth_mining_BE/util"
	"strconv"
)

func BriefArticles(c *gin.Context) {

	var err error
	var temp int

	page := model.Page{}
	//确认需要回复哪些文章，以及文章顺序
	//获取模式
	page.Mode = c.PostForm("mode")

	//获取页数
	temp, err = strconv.Atoi(c.PostForm("pageNumber"))
	if err != nil {
		util.RespFormatError(c)
		return
	}
	page.PageNumber = int64(temp)

	//获取每页显示文章个数
	temp, err = strconv.Atoi(c.PostForm("count"))
	if err != nil {
		util.RespFormatError(c)
		return
	}
	page.Count = int64(temp)

	//获取第一页第一个文章aID
	firstaID := c.PostForm("firstaID")

	//未传入firstaID的话，即第一次请求
	temp = 0

	if len(firstaID) >= 1 {
		temp, err = strconv.Atoi(string(bytes.TrimPrefix([]byte(firstaID), []byte("aID"))))
		if err != nil {
			util.RespAIDError(c)
			return
		}
	}

	page.FirstaID = int64(temp)

	//获取文章发布者的uID
	publisheruID := c.PostForm("publisheruID")

	temp = 0

	if len(publisheruID) >= 1 {
		temp, err = strconv.Atoi(string(bytes.TrimPrefix([]byte(publisheruID), []byte("uID"))))
		if err != nil {
			util.RespUserNotExist(c)
			return
		}
	}

	page.PublisheruID = int64(temp)

	//检验提交的表单数据合法性
	if ( /*page.Mode != "recommend" && */ page.Mode != "newest" && page.Mode != "popularity" && page.Mode != "publisher") || page.PageNumber <= 0 || page.Count < 0 || page.FirstaID < 0 || (page.Mode == "publisher" && page.PublisheruID <= 0) {
		util.RespFormatError(c)
		return
	}

	pagesMap, err := service.BriefArticles(page)
	if err != nil {
		fmt.Println(err)
		fmt.Printf("err: %T\n", err)
		if err == sql.ErrNoRows {
			util.RespNoArticleExit(c)
		} else {
			util.RespUnexceptedError(c)
		}
		return
	}

	/*jsonPagesMap, err := json.Marshal(pagesMap)
	if err != nil {
		util.RespUnexceptedError(c)
		return
	}

	fmt.Println(string(jsonPagesMap))*/

	//回复内容
	util.RespQuerySuccess(c, pagesMap)
	return
}

func DetailArticle(c *gin.Context) {

}

func MyArticle(c *gin.Context) {

}
