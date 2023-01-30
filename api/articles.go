package api

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"rare_earth_mining_BE/model"
	"rare_earth_mining_BE/service"
	"rare_earth_mining_BE/util"
	"regexp"
	"strconv"
	"time"
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

	//获取文章分类
	page.Classification = c.PostForm("classification")

	//获取文章标签
	page.Tags = c.PostForm("tags")

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

// 获取全部分类
func GetClassification(c *gin.Context) {
	classification, err := service.GetClassification()
	if err != nil {
		util.RespUnexceptedError(c)
		return
	}

	util.RespQuerySuccess(c, classification)
	return
}

// 获取全部标签
func GetTags(c *gin.Context) {
	tags, err := service.GetTags()
	if err != nil {
		util.RespUnexceptedError(c)
		return
	}

	util.RespQuerySuccess(c, tags)
	return
}

func DetailArticle(c *gin.Context) {
	aIDin := c.Param("aID")

	//处理输入参数并校验合法性
	if len(aIDin) < 1 {
		util.RespAIDError(c)
		return
	}

	aID, err := strconv.Atoi(string(bytes.TrimPrefix([]byte(aIDin), []byte("aID"))))
	if err != nil || aID <= 0 {
		util.RespAIDError(c)
		return
	}

	//获取文章
	article, err := service.DetailArticle(int64(aID))

	if err != nil {
		fmt.Println(err)
		fmt.Printf("err: %T\n", err)
		if err == sql.ErrNoRows {
			util.RespAIDError(c)
		} else {
			util.RespUnexceptedError(c)
		}
		return

	}

	//响应文章
	util.RespQuerySuccess(c, article)
}

// 发布文章
func PublishArticle(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	classification := c.PostForm("classification")
	tags := c.PostForm("tags")

	//从token获取uID,token已经保证了用户必须登录
	publisheruID, exists := c.Get("uID")

	if !exists {
		util.RespDidNotLogin(c)
		return
	}

	//验证标题使其不能为空且不能大于100字符，验证正文使其不能为空
	if len(title) < 1 || len(title) > 100 || len(content) < 1 {
		util.RespFormatError(c)
		return
	}

	//验证分类是否存在
	if len(classification) < 1 {
		util.RespFormatError(c)
		return
	}

	existsClassification, err := service.GetClassification()
	if err != nil {
		util.RespUnexceptedError(c)
		return
	}

	classificationArr := regexp.MustCompile(",").Split(classification, -1)
	for _, str := range classificationArr {
		reg := regexp.MustCompile(str)
		if !reg.MatchString(existsClassification) {
			util.RespFormatError(c)
			return
		}
	}

	//验证标签是否都存在
	existsTags, err := service.GetTags()
	if err != nil {
		util.RespUnexceptedError(c)
		return
	}

	tagsArr := regexp.MustCompile(",").Split(tags, -1)
	for _, str := range tagsArr {
		reg := regexp.MustCompile(str)
		if !reg.MatchString(existsTags) {
			util.RespFormatError(c)
			return
		}
	}

	//验证通过后,保存文章(此处gin.any的转换，必须先用断言转换成string,再转换成int)
	tempStruID, ok := publisheruID.(string)
	if !ok {
		fmt.Println(ok)
		fmt.Println(publisheruID)
		fmt.Println(tempStruID)
		util.RespUnexceptedError(c)
		return
	}

	tempIntuID, err := strconv.Atoi(tempStruID)
	if err != nil {
		fmt.Println("uID转换出错：", err.Error())
		util.RespUnexceptedError(c)
	}

	err = service.SaveArticle(model.Article{
		UID:            int64(tempIntuID),
		Title:          title,
		Classification: classification,
		Tags:           tags,
		Content:        content,
	})

	if err != nil {
		fmt.Println("保存文章出错：", err.Error())
		util.RespUnexceptedError(c)
		return
	}

	util.RespOK(c)
}

// 发表评论
func PostComment(c *gin.Context) {
	content := c.PostForm("content")
	oID := c.PostForm("oID")

	//从token获取uID,token已经保证了用户必须登录
	publisheruID, exists := c.Get("uID")

	//fmt.Println(publisheruID)

	if !exists {
		util.RespDidNotLogin(c)
		return
	}

	//验证内容和oID使其不能为空,内容不能大于300个字符
	if len(content) < 1 || len(content) > 300 || len(oID) < 1 {
		util.RespFormatError(c)
		return
	}

	//验证通过后,保存文章(此处gin.any的转换，必须先用断言转换成string,再转换成int)
	tempStruID, ok := publisheruID.(string)
	if !ok {
		fmt.Println(ok)
		fmt.Println(publisheruID)
		fmt.Println(tempStruID)
		util.RespUnexceptedError(c)
		return
	}

	tempIntuID, err := strconv.Atoi(tempStruID)
	if err != nil {
		fmt.Println("uID转换出错：", err.Error())
		util.RespUnexceptedError(c)
	}

	err = service.PostComment(model.Comment{
		UID:     int64(tempIntuID),
		OID:     oID,
		Content: content,
	})

	if err != nil {
		fmt.Println("发表评论出错：", err.Error())
		util.RespUnexceptedError(c)
		return
	}

	util.RespOK(c)
}

func CreatorArticleInformation(c *gin.Context) {

	//从token获取uID,token已经保证了用户必须登录
	publisheruID, exists := c.Get("uID")

	if !exists {
		util.RespDidNotLogin(c)
		return
	}

	tempStruID, ok := publisheruID.(string)
	if !ok {
		fmt.Println(ok)
		fmt.Println(publisheruID)
		fmt.Println(tempStruID)
		util.RespUnexceptedError(c)
		return
	}

	tempIntuID, err := strconv.Atoi(tempStruID)
	if err != nil {
		fmt.Println("uID转换出错：", err.Error())
		util.RespUnexceptedError(c)
		return
	}

	information, err := service.CreatorArticleInformation(int64(tempIntuID), time.Now(), 30)

	if err != nil {
		fmt.Println("查询信息错误", err.Error())
		util.RespUnexceptedError(c)
		return
	}

	util.RespQuerySuccess(c, information)

	return
}
