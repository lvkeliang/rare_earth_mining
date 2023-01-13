package dao

import (
	"rare_earth_mining_BE/model"
	"rare_earth_mining_BE/util"
	"strconv"
)

func BriefArticles(page model.Page) (briefArticleInformations map[int64]model.BriefArticleInformation, err error) {
	//初始化
	briefArticleInformations = make(map[int64]model.BriefArticleInformation)

	//存储查询出来的每一个文章
	var article = model.Article{}
	//储存相应的作者的信息
	var publisher = model.User{}
	//存储序号
	var num int64 = 1

	//临时变量
	var tempstr string
	var tempint int

	switch page.Mode {
	case "recommend":

	case "newest":
		/*row := DB.QueryRow("select publishTime from articles where ID = ?", lastaID)
		var lastParam string
		err = row.Scan(&lastParam)
		if err != nil {
			return
		}*/

		//如果是第一次请求(没带firstaID),则将其分配为最晚的文章的aID
		if page.FirstaID == 0 {

			//找到最晚的文章的ID
			row := DB.QueryRow("SELECT max(ID) ID from articles")
			err = row.Scan(&tempstr)
			if err != nil {
				return
			}

			/*tempstr, err = QueryArticleMaximun("ID", "ID")
			if err != nil {
				return
			}*/

			tempint, err = strconv.Atoi(tempstr)
			if err != nil {
				return
			}

			page.FirstaID = int64(tempint)
			//fmt.Printf("page.FirstaID: %v\n", page.FirstaID)
		}
		//fmt.Println("执行到1了")
		rows, err := DB.Query("select users.ID, users.nickname, articles.ID, articles.uID, articles.title, articles.publishTime, articles.viewerNum, articles.likeNum, articles.commentNum, articles.classification, articles.tags from users, articles where users.ID = articles.uID and state = 2 and articles.ID <= ? order by articles.publishTime desc LIMIT ?", page.FirstaID-page.Count*(page.PageNumber-1), page.Count)
		if err != nil {
			//fmt.Println("执行到2了")
			return briefArticleInformations, err
		}

		//fmt.Println("执行到1了")
		defer rows.Close()
		//fmt.Println("执行到2了")
		for rows.Next() {
			//fmt.Println("执行到3了")
			err = rows.Scan(&publisher.ID, &publisher.Nickname, &article.ID, &article.UID, &article.Title, &article.PublishTime, &article.ViewerNum, &article.LikeNum, &article.CommentNum, &article.Classification, &article.Tags)
			if err != nil {
				return briefArticleInformations, err
			}

			//查询发布者的信息

			//将结果储存到输出的map中
			//fmt.Printf("article: %v\n", article)
			//fmt.Printf("num: %v\n", num)
			briefArticleInformations[num] = model.BriefArticleInformation{User: publisher, Article: article}
			//fmt.Printf("articles: %v\n", articles)
			num++
		}

	case "popularity":
		rows, err := DB.Query("select users.ID, users.nickname, articles.ID, articles.uID, articles.title, articles.publishTime, articles.viewerNum, articles.likeNum, articles.commentNum, articles.classification, articles.tags from users, articles where users.ID = articles.uID and state = 2 order by articles.popularityValue desc LIMIT ?, ?", (page.PageNumber-1)*page.Count, page.PageNumber*page.Count)
		if err != nil {
			//fmt.Println("执行到2了")
			return briefArticleInformations, err
		}

		//fmt.Println("执行到1了")
		defer rows.Close()
		//fmt.Println("执行到2了")
		for rows.Next() {
			//fmt.Println("执行到3了")
			err = rows.Scan(&publisher.ID, &publisher.Nickname, &article.ID, &article.UID, &article.Title, &article.PublishTime, &article.ViewerNum, &article.LikeNum, &article.CommentNum, &article.Classification, &article.Tags)
			if err != nil {
				return briefArticleInformations, err
			}
			briefArticleInformations[num] = model.BriefArticleInformation{User: publisher, Article: article}
			num++
		}
	case "user":

	default:
		err = util.FormError

	}
	return
}

/*
func SearchApprovedArticles() (articles map[int64]model.Article, err error) {
	//初始化
	articles = make(map[int64]model.Article)

	//存储查询出来的每一个文章
	var article = model.Article{}
	//存储序号
	var num int64 = 1

	rows, err := DB.Query("select ID, uID, title, publishTime, viewerNum, likeNum, commentNum, classification, tags from articles where state = 2 and ? < ? order by publishTime desc LIMIT ?", Param, lastParam, page.Count)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&article.ID, &article.UID, &article.Title, &article.PublishTime, &article.ViewerNum, &article.LikeNum, &article.CommentNum, &article.Classification, &article.Tags)
		if err != nil {
			return
		}

		//将结果储存到输出的map中
		articles[num] = article
		num++
	}
}*/