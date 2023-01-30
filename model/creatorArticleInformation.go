package model

type CreatorArticleInformation struct {
	ViewerNum  map[int64]int64 `json:"viewerNum"`
	LikeNum    map[int64]int64 `json:"likeNum"`
	CommentNum map[int64]int64 `json:"commentNum"`
	CollectNum map[int64]int64 `json:"collectNum"`
}
