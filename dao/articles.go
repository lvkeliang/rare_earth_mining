package dao

import (
	"fmt"
	"io"
	"os"
	"rare_earth_mining_BE/model"
	"rare_earth_mining_BE/util"
	"strconv"
	"strings"
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

	//处理tags字符串
	page.Tags = "%" + strings.Replace(page.Tags, ",", "%%", -1) + "%"
	//处理classification字符串
	page.Classification = "%" + page.Classification + "%"

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
		rows, err := DB.Query("select users.ID, users.nickname, articles.ID, articles.uID, articles.title, articles.publishTime, articles.viewerNum, articles.likeNum, articles.commentNum, articles.classification, articles.tags from users, articles where articles.classification LIKE ? and articles.tags LIKE ? and users.ID = articles.uID and state = 2 and articles.ID <= ? order by articles.ID desc LIMIT ?", page.Classification, page.Tags, page.FirstaID-page.Count*(page.PageNumber-1), page.Count)
		//rows, err := DB.Query("select users.ID, users.nickname, articles.ID, articles.uID, articles.title, articles.publishTime, articles.viewerNum, articles.likeNum, articles.commentNum, articles.classification, articles.tags from users, articles where articles.classification LIKE ? ", "%"+page.Classification+"%")
		//fmt.Println("'%" + page.Classification + "%'")
		//fmt.Printf("page.FirstaID: %v\n", page.FirstaID)
		//fmt.Printf("page.FirstaID-page.Count*(page.PageNumber-1): %v\n", page.FirstaID-page.Count*(page.PageNumber-1))

		if err != nil {
			//fmt.Println("执行到2了")
			return briefArticleInformations, err
		}

		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&publisher.ID, &publisher.Nickname, &article.ID, &article.UID, &article.Title, &article.PublishTime, &article.ViewerNum, &article.LikeNum, &article.CommentNum, &article.Classification, &article.Tags)
			if err != nil {
				return briefArticleInformations, err
			}
			briefArticleInformations[num] = model.BriefArticleInformation{User: publisher, Article: article}
			num++
		}

	case "popularity":
		rows, err := DB.Query("select users.ID, users.nickname, articles.ID, articles.uID, articles.title, articles.publishTime, articles.viewerNum, articles.likeNum, articles.commentNum, articles.classification, articles.tags from users, articles where articles.classification LIKE ? and articles.tags LIKE ? and users.ID = articles.uID and state = 2 order by articles.popularityValue desc LIMIT ?, ?", page.Classification, page.Tags, (page.PageNumber-1)*page.Count, page.PageNumber*page.Count)
		if err != nil {
			//fmt.Println("执行到2了")
			return briefArticleInformations, err
		}

		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&publisher.ID, &publisher.Nickname, &article.ID, &article.UID, &article.Title, &article.PublishTime, &article.ViewerNum, &article.LikeNum, &article.CommentNum, &article.Classification, &article.Tags)
			if err != nil {
				return briefArticleInformations, err
			}
			briefArticleInformations[num] = model.BriefArticleInformation{User: publisher, Article: article}
			num++
		}

	case "publisher":
		rows, err := DB.Query("select users.ID, users.nickname, articles.ID, articles.uID, articles.title, articles.publishTime, articles.viewerNum, articles.likeNum, articles.commentNum, articles.classification, articles.tags from users, articles where articles.classification LIKE ? and articles.tags LIKE ? and articles. uID = ? and users.ID = articles.uID and state = 2 order by articles.ID desc LIMIT ?, ?", page.Classification, page.Tags, page.PublisheruID, (page.PageNumber-1)*page.Count, page.PageNumber*page.Count)
		if err != nil {
			return briefArticleInformations, err
		}

		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&publisher.ID, &publisher.Nickname, &article.ID, &article.UID, &article.Title, &article.PublishTime, &article.ViewerNum, &article.LikeNum, &article.CommentNum, &article.Classification, &article.Tags)
			if err != nil {
				return briefArticleInformations, err
			}
			briefArticleInformations[num] = model.BriefArticleInformation{User: publisher, Article: article}
			num++
		}

	default:
		err = util.FormError
		return nil, err

	}

	return
}

func GetClassification() (classification string, err error) {
	//临时储存分类名
	var tempClassName string

	rows, err := DB.Query("select className from articleClassification")
	defer rows.Close()
	if err != nil {
		return classification, err
	}
	for rows.Next() {
		err = rows.Scan(&tempClassName)
		if err != nil {
			return classification, err
		}
		classification += tempClassName + ","
	}

	classification = strings.TrimSuffix(classification, ",")

	return
}

func GetTags() (tags string, err error) {
	//临时储存分类名
	var tempTag string

	rows, err := DB.Query("select tag from articleTags")
	defer rows.Close()
	if err != nil {
		return tags, err
	}
	for rows.Next() {
		err = rows.Scan(&tempTag)
		if err != nil {
			return tags, err
		}
		tags += tempTag + ","
	}

	tags = strings.TrimSuffix(tags, ",")

	return
}

func DetailArticle(aID int64) (article model.DetailArticle, err error) {

	//fmt.Println("执行到1了")
	//临时储存内容
	var tempContent []byte

	row := DB.QueryRow("select users.ID, users.nickname, articles.ID, articles.uID, articles.title, articles.publishTime, articles.viewerNum, articles.likeNum, articles.commentNum, articles.classification, articles.tags from users, articles where articles.ID = ?", aID)

	if err = row.Err(); row.Err() != nil {
		fmt.Printf("err: --- : %v\n", err)
		return
	}

	err = row.Scan(&article.User.ID, &article.User.Nickname, &article.Article.ID, &article.Article.UID, &article.Article.Title, &article.Article.PublishTime, &article.Article.ViewerNum, &article.Article.LikeNum, &article.Article.CommentNum, &article.Article.Classification, &article.Article.Tags)

	if err != nil {
		return
	}

	fmt.Printf("article: --- : %v\n", article)
	//fmt.Println("执行到2了")
	//读取文件内容
	file, err := os.Open("data/articles/" + strconv.Itoa(int(aID)) + ".html")
	if err != nil {
		return
		//log.Fatal(err)
	}
	defer file.Close()

	tempContent, err = io.ReadAll(file)
	if err != nil {
		return
		//log.Fatal(err)
	}

	/*tempContent, err = ioutil.ReadFile("data/articles/" + strconv.Itoa(int(aID)))
	if err != nil {
		return
	}
	*/

	article.Article.Content = string(tempContent)

	//fmt.Println("执行到3了")
	//fmt.Println(article)
	//获取评论
	if article.Article.CommentNum > 0 {
		article.Comments, err = GetCommentsByoID("aID" + strconv.Itoa(int(aID)))
	}

	//fmt.Printf("aID" + strconv.Itoa(int(aID)))
	//fmt.Println(article)
	//fmt.Println("执行到4了")
	return
}

func GetCommentsByoID(oID string) (comments map[int64]model.Comment, err error) {
	//初始化
	comments = make(map[int64]model.Comment)

	//fmt.Println("执行到1了")
	//计数
	var num int64 = 0

	//临时储存单个评论
	var comment = model.Comment{}

	//fmt.Println("执行到2了")
	//查询
	rows, err := DB.Query("select ID, uID, oID, publishTime, likeNum, commentNum, layer, content from comments where oID = ? order by likeNum desc", oID)
	defer rows.Close()
	if err != nil {
		return map[int64]model.Comment{}, err
	}

	//fmt.Println("执行到3了")
	for rows.Next() {
		err = rows.Scan(&comment.ID, &comment.UID, &comment.OID, &comment.PublishTime, &comment.LikeNum, &comment.CommentNum, &comment.Layer, &comment.Content)
		if err != nil {
			return map[int64]model.Comment{}, err
		}
		//存入comments
		comments[num] = comment
		num++
	}

	//fmt.Println("执行到4了")
	//fmt.Printf("comments: %v\n", comments)
	//fmt.Printf("comment: %v\n", comment)

	//fmt.Println("执行到5了")
	//递归获取每层评论
	for n, comment := range comments {
		if comment.CommentNum > 0 {
			comment.NextLayerComments, err = GetCommentsByoID("cID" + strconv.Itoa(int(comment.ID)))
			comments[n] = comment
		}
	}

	//fmt.Println("执行到6了")
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
