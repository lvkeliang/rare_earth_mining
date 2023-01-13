package model

type BriefArticleInformation struct {
	User    User    `json:"publisher"`
	Article Article `json:"article"`
}
