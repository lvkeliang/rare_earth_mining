package model

type comment struct {
	ID          int64  `json:"ID"`
	UID         int64  `json:"uID"`
	OID         int64  `json:"oID"` //"aID"\"cID"+对应ID
	PublishTime string `json:"publishTime"`
	LikeNum     int64  `json:"likeNum"`
	CommentNum  int64  `json:"commentNum"`
	Layer       int64  `json:"layer"`
	Content     int64  `json:"content"`
}
