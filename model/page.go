package model

type Page struct {
	Mode       string `json:"mode"`     //模式
	PageNumber int64  `json:"start"`    //页数
	Count      int64  `json:"count"`    //本页显示几个数据
	FirstaID   int64  `json:"firstaID"` //第一页的第一个文章aID
}
