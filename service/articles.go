package service

import (
	"rare_earth_mining_BE/dao"
	"rare_earth_mining_BE/model"
)

func BriefArticles(page model.Page) (briefArticleInformations map[int64]model.BriefArticleInformation, err error) {
	briefArticleInformations, err = dao.BriefArticles(page)
	return
}

func GetClassification() (classification string, err error) {
	classification, err = dao.GetClassification()
	return
}

func GetTags() (tags string, err error) {
	tags, err = dao.GetTags()
	return
}

func DetailArticle(aID int64) (article model.DetailArticle, err error) {
	article, err = dao.DetailArticle(aID)
	return
}
