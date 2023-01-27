package api

import "github.com/gin-gonic/gin"

func InitRouter() {
	r := gin.Default()

	u := r.Group("/user")
	{
		u.POST("/register", Register)
		u.POST("/login", Login)
	}

	a := r.Group("/article")
	{
		a.GET("/brief", BriefArticles)
		a.GET("/detail/:aID", DetailArticle)
	}

	r.GET("/classification", GetClassification)
	r.GET("/tags", GetTags)

	r.Run(":9099")

}
