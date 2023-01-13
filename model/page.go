package model

type Page struct {
	Mode         string `json:"mode"`      //模式
	PageNumber   int64  `json:"start"`     //页数
	Count        int64  `json:"count"`     //本页显示几个数据
	FirstaID     int64  `json:"firstaID"`  //第一页的第一个文章aID(仅mode = newest时生效)
	PublisheruID int64  `json:"publisher"` //发布文章的用户(仅mode = user时生效)
}
