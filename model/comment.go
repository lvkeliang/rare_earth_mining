package model

type Comment struct {
	ID                int64             `json:"ID"`
	UID               int64             `json:"uID"`
	OID               string            `json:"oID"` //"aID"\"cID"+对应ID
	PublishTime       string            `json:"publishTime"`
	LikeNum           int64             `json:"likeNum"`
	CommentNum        int64             `json:"commentNum"`
	Layer             int64             `json:"layer"`
	Content           string            `json:"content"`
	NextLayerComments map[int64]Comment `json:"nextLayerComments"`
}
