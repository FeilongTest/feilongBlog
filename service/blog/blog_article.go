package blog

import (
	"errors"
	"feilongBlog/global"
	model "feilongBlog/model/blog"
	"feilongBlog/model/common/request"
	"gorm.io/gorm"
	"time"
)

type ArticleService struct {
}

// GetArticleList 获取文章列表
func (s *ArticleService) GetArticleList(info model.ArticleSearch) (article []model.Article, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.BLOG_DB.Model(&model.Article{})
	//如果拥有搜索条件
	if info.Fid != 0 {
		db.Where("fid = ?", info.Fid)
	}
	if info.Top {
		db.Where("istop = ?", 1)
	}
	if info.Hide {
		db.Where("status = ?", 1)
	}
	if info.Keyword != "" {
		db.Where("title LIKE ?", "%"+info.Keyword+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("ctime desc").Find(&article).Error
	return article, total, err
}

// GetArticle 查看文章
func (s *ArticleService) GetArticle(id int) (article model.Article, err error) {
	err = global.BLOG_DB.Where("id = ?", id).First(&article).Error
	return article, err
}

// CreateArticle 创建文章
func (s *ArticleService) CreateArticle(article model.Article) (err error) {
	article.Ctime = int(time.Now().Unix())
	article.EditTime = int(time.Now().Unix())
	article.Uid = 1
	article.IsTop = 0
	article.Status = 0
	article.View = 0
	return global.BLOG_DB.Create(&article).Error
}

// UpdateArticle 更新文章
func (s *ArticleService) UpdateArticle(article model.Article) (err error) {
	if errors.Is(global.BLOG_DB.Where("id = ?", article.ID).First(&model.Article{}).Error, gorm.ErrRecordNotFound) {
		return gorm.ErrRecordNotFound
	}
	article.EditTime = int(time.Now().Unix())
	err = global.BLOG_DB.Save(&article).Error
	return err
}

// DeleteArticle 删除文章
func (s *ArticleService) DeleteArticle(article model.Article) (err error) {
	err = global.BLOG_DB.Delete(&article).Error
	return err
}

// DeleteArticleByIds 批量删除文章
func (s *ArticleService) DeleteArticleByIds(ids request.IdsReq) (err error) {
	err = global.BLOG_DB.Delete(&[]model.Article{}, "id in ?", ids.Ids).Error
	return err
}
