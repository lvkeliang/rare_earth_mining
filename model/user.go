package model

type User struct {
	ID           int64  `json:"ID"`
	Mail         string `json:"mail"`
	Nickname     string `json:"nickname"`
	Password     string `json:"password"`
	Position     string `json:"position"`
	Company      string `json:"company"`
	Introduction string `json:"introduction"`
	ArticleNum   string `json:"articleNum"`
	CollectNum   string `json:"collectNum"`
	LikeNum      string `json:"likeNum"`
}
