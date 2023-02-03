package api

import "github.com/gin-gonic/gin"

func InitRouter() {
	r := gin.Default()

	user := r.Group("/user")
	{
		user.POST("/register", Register)
		user.POST("/login", Login)
		user.GET("/information", UserInformation)
		user.PUT("/profile", AuthMiddleware(), UserProfile)
	}

	article := r.Group("/article")
	{
		article.GET("/brief", BriefArticles)
		article.GET("/detail/:aID", DetailArticle)
		article.POST("/postComment", AuthMiddleware(), PostComment)
		article.POST("/like", AuthMiddleware(), Like)
		article.POST("/collect", AuthMiddleware(), Collect)
	}

	r.GET("/classification", GetClassification)
	r.GET("/tags", GetTags)

	creator := r.Group("/creator")
	{
		creator.POST("/publishArticle", AuthMiddleware(), PublishArticle)
		creator.GET("/information", AuthMiddleware(), CreatorArticleInformation)
		creator.GET("/myArticles", AuthMiddleware(), MyArticles)
	}

	r.Run(":9099")

}
