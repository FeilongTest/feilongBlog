package blog

import (
	"feilongBlog/global"
	"feilongBlog/model/blog"
)

type ArticleService struct {
}

func (s *ArticleService) GetArticleList() (article []blog.Article, err error) {
	err = global.BLOG_DB.Limit(10).Find(&article).Error
	return article, err
}
