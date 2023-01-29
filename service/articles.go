package service

import (
	"bytes"
	"rare_earth_mining_BE/dao"
	"rare_earth_mining_BE/model"
	"rare_earth_mining_BE/util"
	"strconv"
	"strings"
)

func BriefArticles(page model.Page) (briefArticleInformations map[int64]model.BriefArticleInformation, err error) {
	briefArticleInformations, err = dao.BriefArticles(page)
	return
}

func GetClassification() (classification string, err error) {
	classification, err = dao.GetClassification()
	return
}

func GetTags() (tags string, err error) {
	tags, err = dao.GetTags()
	return
}

func DetailArticle(aID int64) (article model.DetailArticle, err error) {
	article, err = dao.DetailArticle(aID)
	return
}

func SaveArticle(information model.Article) (err error) {
	err = dao.SaveArticle(information)
	return
}

func PostComment(comment model.Comment) (err error) {

	//fmt.Println("执行到1了")
	//根据oID前缀的不同来查询文章或评论是否存在
	aIDPrefix := "aID"
	cIDPrefix := "cID"

	//fmt.Printf("comment.OID: %v\n", comment.OID)
	if strings.HasPrefix(comment.OID, aIDPrefix) {
		//fmt.Println("执行到2了")
		//处理字符串
		aID, err := strconv.Atoi(string(bytes.TrimPrefix([]byte(comment.OID), []byte("aID"))))
		if err != nil || aID <= 0 {
			return err
		}

		//查询文章,如果有error则说明文章不存在
		_, err = dao.QueryArticleByaID(int64(aID))
		if err != nil {
			return err
		}

		comment.Layer = 1
		//fmt.Println("执行到3了")
	} else if strings.HasPrefix(comment.OID, cIDPrefix) {
		//fmt.Println("执行到4了")
		//处理字符串
		cID, err := strconv.Atoi(string(bytes.TrimPrefix([]byte(comment.OID), []byte("cID"))))
		if err != nil || cID <= 0 {
			return err
		}

		//查询文章,如果有error则说明文章不存在
		tempComment, err := dao.QueryCommentBycID(int64(cID))
		if err != nil {
			return err
		}

		comment.Layer = tempComment.Layer + 1
		//fmt.Println("执行到5了")
	} else {
		//fmt.Println("执行到6了")
		return util.FormError
	}
	//fmt.Println("执行到7了")
	err = dao.AddComment(comment)
	//fmt.Println("执行到8了")
	return
}
