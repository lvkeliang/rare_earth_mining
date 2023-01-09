package api

import "github.com/gin-gonic/gin"

func InitRouter() {
	r := gin.Default()

	u := r.Group("/user")
	{
		u.POST("/register", Register)
		u.POST("/login", Login)
	}

	r.Run(":9099")
}
