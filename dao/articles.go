package dao

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"rare_earth_mining_BE/model"
	"rare_earth_mining_BE/util"
	"strconv"
	"strings"
	"time"
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

	/*
		row := DB.QueryRow("select users.ID, users.nickname, articles.ID, articles.uID, articles.title, articles.publishTime, articles.viewerNum, articles.likeNum, articles.commentNum, articles.classification, articles.tags from users, articles where articles.ID = ?", aID)

		if err = row.Err(); row.Err() != nil {
			fmt.Printf("err: --- : %v\n", err)
			return
		}

		err = row.Scan(&article.User.ID, &article.User.Nickname, &article.Article.ID, &article.Article.UID, &article.Article.Title, &article.Article.PublishTime, &article.Article.ViewerNum, &article.Article.LikeNum, &article.Article.CommentNum, &article.Article.Classification, &article.Article.Tags)
	*/

	//查询文章信息
	article.Article, err = QueryArticleByaID(aID)
	if err != nil {
		return
	}

	//查询用户信息
	row := DB.QueryRow("select ID, nickname from users where ID = ?", article.Article.UID)

	if err = row.Err(); row.Err() != nil {
		fmt.Printf("err: --- : %v\n", err)
		return
	}

	err = row.Scan(&article.User.ID, &article.User.Nickname)

	if err != nil {
		return
	}

	//fmt.Printf("article: --- : %v\n", article)
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

	//更新文章查看数
	_, err = DB.Exec("UPDATE articles SET viewerNum = viewerNum + 1 WHERE ID = ?", aID)
	if err != nil {
		fmt.Println("保存文章时更新用户数据出错：", err.Error())
		return
	}

	//fmt.Printf("aID" + strconv.Itoa(int(aID)))
	//fmt.Println(article)
	//fmt.Println("执行到4了")
	return
}

// 保存文章
func SaveArticle(information model.Article) (err error) {

	//保存文章信息
	result, err := DB.Exec("insert into articles (uID ,title ,classification, tags) value (?,?,?,?)", information.UID, information.Title, information.Classification, information.Tags)
	if err != nil {
		fmt.Println("保存文章信息出错：", err.Error())
		return
	}

	// 返回新插入数据的id
	aID, err := result.LastInsertId()
	if err != nil {
		return
	}

	_, err = DB.Exec("UPDATE users SET articleNum = articleNum + 1 WHERE ID = ?", information.UID)
	if err != nil {
		fmt.Println("保存文章时更新用户数据出错：", err.Error())
		return
	}

	//保存文章内容为文件
	file, err := os.OpenFile("./data/articles/"+strconv.Itoa(int(aID))+".html", os.O_RDWR|os.O_CREATE, 0644)
	defer file.Close()
	//err := ioutil.WriteFile("./data/articles/"+strconv.Itoa(int(aID))+".html", []byte(information.Content), 0644)
	if err != nil {
		return
	}

	if _, err := file.Write([]byte(information.Content)); err != nil {
		// handle error
		return err
	}

	return
}

func QueryArticleByaID(aID int64) (article model.Article, err error) {
	row := DB.QueryRow("select ID, uID, title, publishTime, viewerNum, likeNum, commentNum, classification, tags from articles where ID = ?", aID)

	if err = row.Err(); row.Err() != nil {
		fmt.Printf("err: --- : %v\n", err)
		return
	}

	err = row.Scan(&article.ID, &article.UID, &article.Title, &article.PublishTime, &article.ViewerNum, &article.LikeNum, &article.CommentNum, &article.Classification, &article.Tags)

	if err != nil {
		return
	}
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

func QueryCommentBycID(cID int64) (comment model.Comment, err error) {
	row := DB.QueryRow("select ID, uID, oID, publishTime, likeNum, commentNum, layer, content from comments where ID = ?", cID)

	if err = row.Err(); row.Err() != nil {
		fmt.Printf("err: --- : %v\n", err)
		return
	}

	err = row.Scan(&comment.ID, &comment.UID, &comment.OID, &comment.PublishTime, &comment.LikeNum, &comment.CommentNum, &comment.Layer, &comment.Content)

	if err != nil {
		return
	}

	return
}

func AddComment(comment model.Comment) (err error) {
	//存入评论
	_, err = DB.Exec("insert into comments (uID, oID, layer, content) values (?,?,?,?)", comment.UID, comment.OID, comment.Layer, comment.Content)
	if err != nil {
		return err
	}
	if comment.Layer == 1 {
		//更新文章的评论数
		aID, err := strconv.Atoi(string(bytes.TrimPrefix([]byte(comment.OID), []byte("aID"))))
		if err != nil || aID <= 0 {
			return err
		}

		_, err = DB.Exec("UPDATE articles SET commentNum = commentNum + 1 WHERE ID = ?", aID)
		if err != nil {
			return err
		}
	} else {
		//更新评论的评论数
		cID, err := strconv.Atoi(string(bytes.TrimPrefix([]byte(comment.OID), []byte("cID"))))
		if err != nil || cID <= 0 {
			return err
		}

		_, err = DB.Exec("UPDATE comments SET commentNum = commentNum + 1 WHERE ID = ?", cID)
		if err != nil {
			return err
		}
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

func QueryArticlesByuID(uID int64) (articles map[int64]model.Article, err error) {

	//初始化
	articles = make(map[int64]model.Article)

	rows, err := DB.Query("select ID, uID, title, publishTime, viewerNum, likeNum, commentNum, classification, tags, state from articles where uID = ? order by ID desc", uID)

	if err != nil {
		//fmt.Println("执行到2了")
		return articles, err
	}

	defer rows.Close()

	num := 0
	var tempArticle model.Article
	for rows.Next() {

		err = rows.Scan(&tempArticle.ID, &tempArticle.UID, &tempArticle.Title, &tempArticle.PublishTime, &tempArticle.ViewerNum, &tempArticle.LikeNum, &tempArticle.CommentNum, &tempArticle.Classification, &tempArticle.Tags, &tempArticle.State)
		if err != nil {
			return articles, err
		}

		articles[int64(num)] = tempArticle
		num++
	}

	return
}

// 获取文章列表中，所有文章从指定日期往前指定天数的每日的总评论数
func GetArticleDailyCommentsNum(oIDList []string, latestDate time.Time, day int64) (dailyCommentsNum map[int64]int64, err error) {

	//初始化
	dailyCommentsNum = make(map[int64]int64)

	oIDListStr := "\"" + strings.Trim(strings.Replace(fmt.Sprint(oIDList), " ", "\",\"", -1), "[]") + "\""
	//fmt.Println("oIDListStr: ---- : ", oIDListStr)

	//select date(publishTime) AS date, COUNT(*) from comments where oID IN ("aID21","aID20","aID19","aID17","aID16","aID14") and date(publishTime) <= "2023-01-30" GROUP BY date order by date desc LIMIT 30
	//rows, err := DB.Query("select date(publishTime) AS date, COUNT(*) from comments where oID IN (?) and date(publishTime) <= ? GROUP BY date order by date desc LIMIT ?", "\""+oIDListStr+"\"", latestDate.Format("2006-01-02"), day)

	//查询文章列表中，所有文章从指定日期往前指定天数的每日的总评论数
	rows, err := DB.Query("select date(publishTime) AS date, COUNT(*) from comments where oID IN ("+oIDListStr+") and date(publishTime) <= ? GROUP BY date order by date desc LIMIT ?", latestDate.Format("2006-01-02"), day)

	if err != nil {
		return dailyCommentsNum, err
	}

	defer rows.Close()

	var tempNum int64
	var tempTime string
	result := make(map[string]int64)

	for rows.Next() {
		err = rows.Scan(&tempTime, &tempNum)
		if err != nil {
			return dailyCommentsNum, err
		}
		//fmt.Println("tempTime: ", tempTime)
		result[tempTime] = tempNum
		//fmt.Println("result[tempTime]: ", result[tempTime])
	}

	//比对日期,result中如果存在相同的日期,则数量就取reusult中的值,否则就是0
	for ; day >= 0; day-- {
		//fmt.Println("day-: ", day)
		d, _ := time.ParseDuration(strconv.Itoa(int(day*-24)) + "h")
		tempDate := latestDate.Add(d).Format("2006-01-02") + "T00:00:00Z"
		//fmt.Println("tempDate: ", tempDate)
		//fmt.Println("result[tempDate]: ", result[tempDate])
		if result[tempDate] > 0 {
			dailyCommentsNum[day] = result[tempDate]
		} else {
			dailyCommentsNum[day] = 0
		}
	}
	//fmt.Println("dailyCommentsNum-----: ", dailyCommentsNum)

	return
}

// 获取文章列表中，所有文章从指定日期往前指定天数的每日的总点赞数
func GetArticleDailyLikesNum(oIDList []string, latestDate time.Time, day int64) (dailyLikesNum map[int64]int64, err error) {

	//初始化
	dailyLikesNum = make(map[int64]int64)

	oIDListStr := "\"" + strings.Trim(strings.Replace(fmt.Sprint(oIDList), " ", "\",\"", -1), "[]") + "\""
	//fmt.Println("oIDListStr: ---- : ", oIDListStr)

	//select date(publishTime) AS date, COUNT(*) from comments where oID IN ("aID21","aID20","aID19","aID17","aID16","aID14") and date(publishTime) <= "2023-01-30" GROUP BY date order by date desc LIMIT 30
	//rows, err := DB.Query("select date(publishTime) AS date, COUNT(*) from comments where oID IN (?) and date(publishTime) <= ? GROUP BY date order by date desc LIMIT ?", "\""+oIDListStr+"\"", latestDate.Format("2006-01-02"), day)

	//查询文章列表中，所有文章从指定日期往前指定天数的每日的总评论数
	rows, err := DB.Query("select date(time) AS date, COUNT(*) from likes where oID IN ("+oIDListStr+") and date(time) <= ? GROUP BY date order by date desc LIMIT ?", latestDate.Format("2006-01-02"), day)

	if err != nil {
		return dailyLikesNum, err
	}

	defer rows.Close()

	var tempNum int64
	var tempTime string
	result := make(map[string]int64)

	for rows.Next() {
		err = rows.Scan(&tempTime, &tempNum)
		if err != nil {
			return dailyLikesNum, err
		}
		//fmt.Println("tempTime: ", tempTime)
		result[tempTime] = tempNum
		//fmt.Println("result[tempTime]: ", result[tempTime])
	}

	//比对日期,result中如果存在相同的日期,则数量就取reusult中的值,否则就是0
	for ; day >= 0; day-- {
		//fmt.Println("day-: ", day)
		d, _ := time.ParseDuration(strconv.Itoa(int(day*-24)) + "h")
		tempDate := latestDate.Add(d).Format("2006-01-02") + "T00:00:00Z"
		//fmt.Println("tempDate: ", tempDate)
		//fmt.Println("result[tempDate]: ", result[tempDate])
		if result[tempDate] > 0 {
			dailyLikesNum[day] = result[tempDate]
		} else {
			dailyLikesNum[day] = 0
		}
	}
	//fmt.Println("dailyLikesNum-----: ", dailyLikesNum)

	return
}

// 获取文章列表中，所有文章从指定日期往前指定天数的每日的总收藏数
func GetArticleDailyCollectNum(oIDList []string, latestDate time.Time, day int64) (dailyCollectNum map[int64]int64, err error) {

	//初始化
	dailyCollectNum = make(map[int64]int64)

	oIDListStr := "\"" + strings.Trim(strings.Replace(fmt.Sprint(oIDList), " ", "\",\"", -1), "[]") + "\""
	//fmt.Println("oIDListStr: ---- : ", oIDListStr)

	//select date(publishTime) AS date, COUNT(*) from comments where oID IN ("aID21","aID20","aID19","aID17","aID16","aID14") and date(publishTime) <= "2023-01-30" GROUP BY date order by date desc LIMIT 30
	//rows, err := DB.Query("select date(publishTime) AS date, COUNT(*) from comments where oID IN (?) and date(publishTime) <= ? GROUP BY date order by date desc LIMIT ?", "\""+oIDListStr+"\"", latestDate.Format("2006-01-02"), day)

	//查询文章列表中，所有文章从指定日期往前指定天数的每日的总评论数
	rows, err := DB.Query("select date(time) AS date, COUNT(*) from collections where oID IN ("+oIDListStr+") and date(time) <= ? GROUP BY date order by date desc LIMIT ?", latestDate.Format("2006-01-02"), day)

	if err != nil {
		return dailyCollectNum, err
	}

	defer rows.Close()

	var tempNum int64
	var tempTime string
	result := make(map[string]int64)

	for rows.Next() {
		err = rows.Scan(&tempTime, &tempNum)
		if err != nil {
			return dailyCollectNum, err
		}
		//fmt.Println("tempTime: ", tempTime)
		result[tempTime] = tempNum
		//fmt.Println("result[tempTime]: ", result[tempTime])
	}

	//比对日期,result中如果存在相同的日期,则数量就取reusult中的值,否则就是0
	for ; day >= 0; day-- {
		//fmt.Println("day-: ", day)
		d, _ := time.ParseDuration(strconv.Itoa(int(day*-24)) + "h")
		tempDate := latestDate.Add(d).Format("2006-01-02") + "T00:00:00Z"
		//fmt.Println("tempDate: ", tempDate)
		//fmt.Println("result[tempDate]: ", result[tempDate])
		if result[tempDate] > 0 {
			dailyCollectNum[day] = result[tempDate]
		} else {
			dailyCollectNum[day] = 0
		}
	}
	//fmt.Println("dailyLikesNum-----: ", dailyLikesNum)

	return
}

/*
func GetArticleDailyCommentsNum(oID string, latestDate time.Time, day int64) (dailyCommentsNum map[int64]int64, err error) {

	//初始化
	dailyCommentsNum = make(map[int64]int64)

	//fmt.Println("oID: ", oID)

	var tempNum int64
	var i int64 = 0

	for ; i < day; i++ {

		d, _ := time.ParseDuration(strconv.Itoa(int(i*-24)) + "h")
		date := latestDate.Add(d).Format("2006-01-02")

		//fmt.Println("date: ", date)

		row := DB.QueryRow("select COUNT(*) from comments where oID = ? and date(publishTime) = ?", oID, date)
		//row := DB.QueryRow("select COUNT(*) from comments where oID = ? order by ID desc LIMIT ?", oID, date)

		if err = row.Err(); row.Err() != nil {
			fmt.Printf("err: --- : %v\n", err)
			return
		}

		err = row.Scan(&tempNum)

		//fmt.Println("tempNum: ", tempNum)

		if err != nil {
			return
		}

		dailyCommentsNum[i] = tempNum
		//fmt.Printf("dailyCommentsNum[%v]: %v\n", i, dailyCommentsNum)

	}

	return
}*/
/*
func GetArticleDailyLikeNum(oID string, latestDate time.Time, day int64) (dailyLikesNum map[int64]int64, err error) {

	//初始化
	dailyLikesNum = make(map[int64]int64)

	var tempNum int64
	var i int64 = 0

	for ; i < day; i++ {

		d, _ := time.ParseDuration(strconv.Itoa(int(i*-24)) + "h")
		date := latestDate.Add(d).Format("2006-01-02")

		row := DB.QueryRow("select COUNT(*) from likes where oID = ? and date(time) = ?", oID, date)

		if err = row.Err(); row.Err() != nil {
			fmt.Printf("err: --- : %v\n", err)
			return
		}

		err = row.Scan(&tempNum)

		if err != nil {
			return
		}

		dailyLikesNum[i] = tempNum

	}

	return
}*/
/*
func GetArticleDailyCollectNum(oID string, latestDate time.Time, day int64) (dailyCollectNum map[int64]int64, err error) {

	//初始化
	dailyCollectNum = make(map[int64]int64)

	var tempNum int64
	var i int64 = 0

	for ; i < day; i++ {

		d, _ := time.ParseDuration(strconv.Itoa(int(i*-24)) + "h")
		date := latestDate.Add(d).Format("2006-01-02")

		row := DB.QueryRow("select COUNT(*) from collections where oID = ? and date(time) = ?", oID, date)

		if err = row.Err(); row.Err() != nil {
			fmt.Printf("err: --- : %v\n", err)
			return
		}

		err = row.Scan(&tempNum)

		if err != nil {
			return
		}

		dailyCollectNum[i] = tempNum

	}

	return
}
*/
