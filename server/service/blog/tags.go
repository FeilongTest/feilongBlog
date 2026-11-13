package blog

import (
	"errors"
	driver "github.com/go-sql-driver/mysql"
	"strings"
	"time"
	"unicode/utf8"

	"feilongBlog/global"
	model "feilongBlog/model/blog"
	"gorm.io/gorm"
)

type TagService struct{}

func (s *TagService) List(public bool) ([]model.Tag, error) {
	tags := []model.Tag{}
	db := global.BLOG_DB.Table("blog_tag AS tag").Select("tag.id, tag.name, tag.description, COUNT(article.id) AS count").
		Joins("LEFT JOIN blog_article_tag AS link ON link.tag_id = tag.id").
		Joins("LEFT JOIN blog_article AS article ON article.id = link.article_id AND article.status = ?", model.ArticlePublished).
		Group("tag.id, tag.name, tag.description").Order("count desc, tag.name asc, tag.id asc")
	if public {
		db = db.Having("COUNT(article.id) > 0")
	}
	err := db.Limit(500).Scan(&tags).Error
	return tags, err
}

func (s *TagService) GetPublic(id uint) (model.Tag, error) {
	var tag model.Tag
	err := global.BLOG_DB.Table("blog_tag AS tag").Select("tag.id, tag.name, tag.description, COUNT(article.id) AS count").
		Joins("JOIN blog_article_tag AS link ON link.tag_id = tag.id").
		Joins("JOIN blog_article AS article ON article.id = link.article_id AND article.status = ?", model.ArticlePublished).
		Where("tag.id = ?", id).Group("tag.id, tag.name, tag.description").Scan(&tag).Error
	if err == nil && tag.ID == 0 {
		err = gorm.ErrRecordNotFound
	}
	return tag, err
}

func (s *TagService) Save(tag *model.Tag) error {
	names, err := normalizeTagNames([]string{tag.Name})
	if err != nil {
		return err
	}
	if len(names) != 1 {
		return invalidArticle("请输入标签名称")
	}
	tag.Name = names[0]
	tag.Description = strings.TrimSpace(tag.Description)
	if utf8.RuneCountInString(tag.Description) > 300 {
		return invalidArticle("标签描述最多 300 字")
	}
	if tag.ID == 0 {
		err = global.BLOG_DB.Create(tag).Error
	} else {
		err = global.BLOG_DB.Transaction(func(tx *gorm.DB) error {
			var existing model.Tag
			if err := tx.First(&existing, tag.ID).Error; err != nil {
				return err
			}
			oldName := existing.Name
			if err := tx.Model(&existing).Updates(map[string]interface{}{"name": tag.Name, "description": tag.Description}).Error; err != nil {
				return err
			}
			if oldName != tag.Name {
				return invalidateTagEditors(tx, tag.ID)
			}
			return nil
		})
	}
	if err == nil {
		seoService.ResetSitemap()
		seoService.ResetFeed()
	}
	var mysqlErr *driver.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return invalidArticle("该标签名称已存在，请使用其他名称")
	}
	return err
}

func (s *TagService) Delete(id uint) error {
	err := global.BLOG_DB.Transaction(func(tx *gorm.DB) error {
		if err := invalidateTagEditors(tx, id); err != nil {
			return err
		}
		if err := tx.Where("tag_id = ?", id).Delete(&model.ArticleTag{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Tag{}, id).Error
	})
	if err == nil {
		seoService.ResetSitemap()
	}
	return err
}

// 全局改名或移除标签后使打开的旧编辑器失效，避免再次保存时悄悄重建旧标签。
func invalidateTagEditors(tx *gorm.DB, id uint) error {
	return tx.Model(&model.Article{}).Where("id IN (?)", tx.Model(&model.ArticleTag{}).Select("article_id").Where("tag_id = ?", id)).Updates(map[string]interface{}{"version": gorm.Expr("version + 1"), "edittime": time.Now().Unix()}).Error
}
