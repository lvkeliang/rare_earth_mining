package model

type DetailArticle struct {
	User     User              `json:"publisher"`
	Article  Article           `json:"article"`
	Comments map[int64]Comment `json:"comments"`
}
