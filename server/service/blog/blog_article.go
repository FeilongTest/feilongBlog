package blog

import (
	"errors"
	"feilongBlog/global"
	model "feilongBlog/model/blog"
	"feilongBlog/model/common/request"
	"time"

	"gorm.io/gorm"
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
		db = db.Where("fid = ?", info.Fid)
	}
	if info.Top {
		db = db.Where("istop = ?", 1)
	}
	if info.Hide {
		db = db.Where("status = ?", 1)
	}
	if info.Keyword != "" {
		db = db.Where("title LIKE ?", "%"+info.Keyword+"%")
	}
	db = db.Where("status = ?", 0)
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("istop desc,ctime desc").Find(&article).Error

	// 统计每篇文章的评论数量
	if err == nil && len(article) > 0 {
		// 收集所有文章ID
		var articleIds []uint
		for _, art := range article {
			articleIds = append(articleIds, art.ID)
		}

		// 批量查询评论数量 (status = 0 表示显示)
		var commentCounts []struct {
			Aid   int   `json:"aid"`
			Count int64 `json:"count"`
		}
		global.BLOG_DB.Model(&model.Comment{}).
			Select("aid, count(*) as count").
			Where("aid IN ? AND status = ?", articleIds, 0).
			Group("aid").
			Find(&commentCounts)

		// 将评论数量映射到文章
		commentMap := make(map[uint]int64)
		for _, cc := range commentCounts {
			commentMap[uint(cc.Aid)] = cc.Count
		}

		var likeCounts []struct {
			Aid   uint
			Count int64
		}
		global.BLOG_DB.Model(&model.ArticleLike{}).
			Select("aid, count(*) as count").
			Where("aid IN ?", articleIds).
			Group("aid").
			Find(&likeCounts)

		likeMap := make(map[uint]int64)
		for _, like := range likeCounts {
			likeMap[like.Aid] = like.Count
		}

		// 设置每篇文章的评论和点赞数量
		for i := range article {
			article[i].CommentCount = int(commentMap[article[i].ID])
			article[i].LikeCount = int(likeMap[article[i].ID])
		}
	}

	return article, total, err
}

// GetArticleListAdmin 管理员获取文章分类
func (s *ArticleService) GetArticleListAdmin(info model.ArticleSearch) (article []model.Article, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.BLOG_DB.Model(&model.Article{})
	//如果拥有搜索条件
	if info.Fid != 0 {
		db = db.Where("fid = ?", info.Fid)
	}
	if info.Top {
		db = db.Where("istop = ?", 1)
	}
	if info.Hide {
		db = db.Where("status = ?", 1)
	}
	if info.Keyword != "" {
		db = db.Where("title LIKE ?", "%"+info.Keyword+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("ctime desc").Find(&article).Error
	return article, total, err
}

// GetArticleSummary 获取文章分类大概情况
func (s *ArticleService) GetArticleSummary(id int) (result any, err error) {
	var categoryGroup []struct {
		Fid   int    `json:"fid" gorm:"fid"`
		Total int    `json:"total" gorm:"total"`
		Name  string `json:"name"`
	}
	err = global.BLOG_DB.Table("blog_article").
		Select("count(*) as total,fid").
		Where("fid in (?) and status = ?",
			global.BLOG_DB.
				Table("blog_category").
				Select("id").
				Where("fid = (?)",
					global.BLOG_DB.
						Table("blog_category").
						Select("fid").
						Where("id = ?", id)), 0).
		Group("fid").Find(&categoryGroup).Error
	//获取分类名称 手动拼接
	var categoryList []model.Category
	if err == nil {
		if err = global.BLOG_DB.Find(&categoryList).Error; err == nil {
			for i := 0; i < len(categoryGroup); i++ {
				for _, category := range categoryList {
					if categoryGroup[i].Fid == int(category.ID) {
						categoryGroup[i].Name = category.Name
					}
				}
			}

		}
	}
	return categoryGroup, err
}

// GetArticle 查看文章
func (s *ArticleService) GetArticle(id int) (article model.Article, err error) {
	db := global.BLOG_DB.Where("id = ?", id).First(&article)
	err = db.Error
	if err == nil {
		var likeCount int64
		if countErr := global.BLOG_DB.Model(&model.ArticleLike{}).Where("aid = ?", id).Count(&likeCount).Error; countErr == nil {
			article.LikeCount = int(likeCount)
		}
	}
	if err := db.Transaction(func(tx *gorm.DB) error {
		return tx.Update("view", gorm.Expr("view + ?", 1)).Error
	}); err == nil {
		article.View = article.View + 1
	}
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

// LikeArticle 点赞或取消点赞
func (s *ArticleService) LikeArticle(req model.ArticleLikeRequest) (liked bool, count int64, err error) {
	err = global.BLOG_DB.Transaction(func(tx *gorm.DB) error {
		var like model.ArticleLike
		result := tx.Where("aid = ? AND visitor_id = ?", req.Aid, req.VisitorID).First(&like)
		if result.Error == nil {
			liked = false
			return tx.Delete(&like).Error
		}
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}
		liked = true
		return tx.Create(&model.ArticleLike{Aid: req.Aid, VisitorID: req.VisitorID, Ctime: int(time.Now().Unix())}).Error
	})
	if err == nil {
		err = global.BLOG_DB.Model(&model.ArticleLike{}).Where("aid = ?", req.Aid).Count(&count).Error
	}
	return
}

// GetArticleLike 获取访客的点赞状态
func (s *ArticleService) GetArticleLike(req model.ArticleLikeRequest) (liked bool, count int64, err error) {
	err = global.BLOG_DB.Model(&model.ArticleLike{}).Where("aid = ?", req.Aid).Count(&count).Error
	if err != nil {
		return
	}
	var likeCount int64
	err = global.BLOG_DB.Model(&model.ArticleLike{}).Where("aid = ? AND visitor_id = ?", req.Aid, req.VisitorID).Count(&likeCount).Error
	liked = likeCount > 0
	return
}
