package model

type Article struct {
	ID             int64  `json:"ID"`
	UID            int64  `json:"uID"`
	Title          string `json:"title"`
	PublishTime    string `json:"publishTime"`
	UpdateTime     string `json:"updateTime"`
	ViewerNum      int64  `json:"viewerNum"`
	LikeNum        int64  `json:"likeNum"`
	CommentNum     int64  `json:"commentNum"`
	Classification string `json:"classification"`
	Tags           string `json:"tags"`
	State          string `json:"state"`
	Content        string `json:"content"`
}
