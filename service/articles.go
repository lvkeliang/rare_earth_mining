package service

import (
	"rare_earth_mining_BE/dao"
	"rare_earth_mining_BE/model"
)

func BriefArticles(page model.Page) (briefArticleInformations map[int64]model.BriefArticleInformation, err error) {
	briefArticleInformations, err = dao.BriefArticles(page)
	return
}
