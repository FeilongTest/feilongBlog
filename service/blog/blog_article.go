package blog

import (
	"feilongBlog/global"
	model "feilongBlog/model/blog"
)

type ArticleService struct {
}

func (s *ArticleService) GetArticleList(info model.ArticleSearch) (article []model.Article, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.BLOG_DB.Model(&model.Article{})
	//如果拥有搜索条件
	if info.Fid != 0 {
		db.Where("fid = ?", info.Fid)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("ctime desc").Find(&article).Error
	return article, total, err
}
