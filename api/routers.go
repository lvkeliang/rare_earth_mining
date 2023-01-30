package api

import "github.com/gin-gonic/gin"

func InitRouter() {
	r := gin.Default()

	user := r.Group("/user")
	{
		user.POST("/register", Register)
		user.POST("/login", Login)
	}

	article := r.Group("/article")
	{
		article.GET("/brief", BriefArticles)
		article.GET("/detail/:aID", DetailArticle)
		article.POST("/postComment", AuthMiddleware(), PostComment)
	}

	r.GET("/classification", GetClassification)
	r.GET("/tags", GetTags)

	creator := r.Group("/creator")
	{
		creator.POST("/publishArticle", AuthMiddleware(), PublishArticle)
		creator.GET("/information", AuthMiddleware(), CreatorArticleInformation)
	}

	r.Run(":9099")

}
