package service

import (
	"bytes"
	"database/sql"
	"fmt"
	"rare_earth_mining_BE/dao"
	"rare_earth_mining_BE/model"
	"rare_earth_mining_BE/util"
	"strconv"
	"strings"
	"time"
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

// 点赞
func Like(like model.Comment) (err error) {

	//fmt.Println("执行到1了")
	//根据oID前缀的不同来查询文章或评论是否存在
	aIDPrefix := "aID"
	cIDPrefix := "cID"

	//fmt.Printf("comment.OID: %v\n", comment.OID)
	if strings.HasPrefix(like.OID, aIDPrefix) {
		//fmt.Println("执行到2了")
		//处理字符串
		aID, err := strconv.Atoi(string(bytes.TrimPrefix([]byte(like.OID), []byte("aID"))))
		if err != nil || aID <= 0 {
			return err
		}

		//查询文章,如果有error则说明文章不存在
		_, err = dao.QueryArticleByaID(int64(aID))
		if err != nil {
			//fmt.Printf("err: %T\n--- %v\n", err, err.Error())
			return util.NoArticleExistsError
		}

		//查询是否点赞过了
		isLiked, err := dao.IsLiked(like)
		if err != sql.ErrNoRows && err != nil {
			return err
		}

		if !isLiked {
			return util.AreadyLikedError
		}
		//like.Layer = 1
		//fmt.Println("执行到3了")
	} else if strings.HasPrefix(like.OID, cIDPrefix) {
		fmt.Println("执行到4了")
		//处理字符串
		cID, err := strconv.Atoi(string(bytes.TrimPrefix([]byte(like.OID), []byte("cID"))))
		if err != nil || cID <= 0 {
			return err
		}

		//查询文章,如果有error则说明文章不存在
		//tempComment, err := dao.QueryCommentBycID(int64(cID))
		_, err = dao.QueryCommentBycID(int64(cID))
		if err != nil {
			return util.NoCommectExistsError
		}

		//查询是否点赞过了
		isLiked, err := dao.IsLiked(like)
		if err != sql.ErrNoRows && err != nil {
			return err
		}

		if !isLiked {
			return util.AreadyLikedError
		}

		//like.Layer = tempComment.Layer + 1
		//fmt.Println("执行到5了")
	} else {
		//fmt.Println("执行到6了")
		return util.FormError
	}
	//fmt.Println("执行到7了")
	err = dao.Like(like)
	//fmt.Println("执行到8了")
	return
}

// 收藏
func Collect(collect model.Comment) (err error) {

	//fmt.Println("执行到1了")
	//根据oID前缀的不同来查询文章或评论是否存在
	aIDPrefix := "aID"
	cIDPrefix := "cID"

	//fmt.Printf("comment.OID: %v\n", comment.OID)
	if strings.HasPrefix(collect.OID, aIDPrefix) {
		//fmt.Println("执行到2了")
		//处理字符串
		aID, err := strconv.Atoi(string(bytes.TrimPrefix([]byte(collect.OID), []byte("aID"))))
		if err != nil || aID <= 0 {
			return err
		}

		//查询文章,如果有error则说明文章不存在
		_, err = dao.QueryArticleByaID(int64(aID))
		if err != nil {
			return util.NoArticleExistsError
		}

		//查询是否收藏过了
		isCollected, err := dao.IsCollected(collect)
		if err != sql.ErrNoRows && err != nil {
			return err
		}

		if !isCollected {
			return util.AreadyCollectedError
		}

		//collect.Layer = 1
		//fmt.Println("执行到3了")
	} else if strings.HasPrefix(collect.OID, cIDPrefix) {
		//fmt.Println("执行到4了")
		//处理字符串
		cID, err := strconv.Atoi(string(bytes.TrimPrefix([]byte(collect.OID), []byte("cID"))))
		if err != nil || cID <= 0 {
			return err
		}

		//查询文章,如果有error则说明文章不存在
		_, err = dao.QueryCommentBycID(int64(cID))
		if err != nil {
			return util.NoCommectExistsError
		}

		//查询是否收藏过了
		isCollected, err := dao.IsLiked(collect)
		if err != sql.ErrNoRows && err != nil {
			return err
		}

		if !isCollected {
			return util.AreadyCollectedError
		}

		//collect.Layer = tempComment.Layer + 1
		//fmt.Println("执行到5了")
	} else {
		//fmt.Println("执行到6了")
		return util.FormError
	}
	//fmt.Println("执行到7了")
	err = dao.Collect(collect)
	//fmt.Println("执行到8了")
	return
}

func CreatorArticleInformation(uID int64, latestDate time.Time, day int64) (information model.CreatorArticleInformation, err error) {

	//初始化
	information = model.CreatorArticleInformation{
		ViewerNum:  make(map[int64]int64),
		LikeNum:    make(map[int64]int64),
		CommentNum: make(map[int64]int64),
		CollectNum: make(map[int64]int64),
	}

	articles, err := dao.QueryArticlesByuID(uID)
	if err != nil {
		return
	}

	//fmt.Println("articles: ", articles)
	var tempInformation model.CreatorArticleInformation

	tempInformation.CommentNum = make(map[int64]int64)
	tempInformation.LikeNum = make(map[int64]int64)
	tempInformation.CollectNum = make(map[int64]int64)
	tempInformation.ViewerNum = make(map[int64]int64)

	var oIDList []string
	for _, article := range articles {
		oIDList = append(oIDList, "aID"+strconv.Itoa(int(article.ID)))
	}

	information.CommentNum, err = dao.GetArticleDailyCommentsNum(oIDList, latestDate, day)
	if err != nil {
		return
	}

	information.LikeNum, err = dao.GetArticleDailyLikesNum(oIDList, latestDate, day)
	if err != nil {
		return
	}

	information.CollectNum, err = dao.GetArticleDailyCollectNum(oIDList, latestDate, day)
	if err != nil {
		return
	}

	/*
		for _, article := range articles {

			//for i := 0; articles[int64(i)] != nil; i++ {
			//article := articles[int64(i)]

			if article.State == 2 {

				tempStrcID := "aID" + strconv.Itoa(int(article.ID))
				fmt.Println("cID: ", tempStrcID)

				/*
					//获取每天的评论数
					tempInformation.CommentNum, err = dao.GetArticleDailyCommentsNum(tempStrcID, latestDate, day)
					if err != nil {
						return
					}

				//fmt.Println("CommentNum: ", tempInformation.CommentNum)
				//获取每天的点赞数


				tempInformation.LikeNum, err = dao.GetArticleDailyLikeNum(tempStrcID, latestDate, day)
				if err != nil {
					return
				}

				//fmt.Println("LikeNum: ", tempInformation.LikeNum)
				//获取每天的收藏数
				tempInformation.CollectNum, err = dao.GetArticleDailyCollectNum(tempStrcID, latestDate, day)
				if err != nil {
					return
				}

				//fmt.Println("CollectNum: ", tempInformation.CollectNum)
				//将每篇文章的数据累计,得到所有文章的总数据

					for key, _ := range tempInformation.CommentNum {
						fmt.Println("文章的数据累计key: ", key)
						information.CommentNum[key] += tempInformation.CommentNum[key]
					}

				for key, _ := range tempInformation.LikeNum {
					information.LikeNum[key] += tempInformation.LikeNum[key]
				}

				for key, _ := range tempInformation.CollectNum {
					information.CollectNum[key] += tempInformation.CollectNum[key]
				}

				for key,_ := range tempInformation.CommentNum{
					information.CommentNum[key] += tempInformation.CommentNum[key]
				}
			}

			fmt.Println("lastinformation.Comment: ", information.CommentNum)
		}
	*/

	return
}

// 获取一个用户的文章,state为任意
func QueryArticleByuID(uID int64) (article model.Article, err error) {
	article, err = dao.QueryArticleByuID(uID)
	return
}
