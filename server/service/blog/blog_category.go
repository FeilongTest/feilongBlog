package blog

import (
	"errors"
	"feilongBlog/global"
	model "feilongBlog/model/blog"
	"feilongBlog/model/common/request"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CategoryService struct {
}

// GetCategoryList 获取分类列表
func (s *CategoryService) GetCategoryList() (list []model.Category, err error) {
	var child []model.Category
	var root []model.Category
	err = global.BLOG_DB.Where("fid = ?", 0).Order("sort desc").Find(&root).Error
	if err != nil {
		return
	}
	err = global.BLOG_DB.Where("fid != ?", 0).Order("sort desc").Find(&child).Error
	if err != nil {
		return
	}
	list = append(list, root...)
	list = append(list, child...)
	return list, err
}

// CreateCategory 创建分类
func (s *CategoryService) CreateCategory(category model.Category) (err error) {
	return global.BLOG_DB.Transaction(func(tx *gorm.DB) error {
		if err := lockCategoryParent(tx, category); err != nil {
			return err
		}
		return tx.Create(&category).Error
	})
}

// UpdateCategory 更新分类
func (s *CategoryService) UpdateCategory(category model.Category) (err error) {
	return global.BLOG_DB.Transaction(func(tx *gorm.DB) error {
		var existing model.Category
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&existing, category.ID).Error; err != nil {
			return err
		}
		if err := lockCategoryParent(tx, category); err != nil {
			return err
		}
		if category.Fid > 0 {
			var children int64
			if err := tx.Model(&model.Category{}).Where("fid = ?", category.ID).Count(&children).Error; err != nil {
				return err
			}
			if children > 0 {
				return invalidArticle("该分类仍有子分类，请先调整子分类再改变所属层级")
			}
		}
		return tx.Save(&category).Error
	})
}

// Serialize parent changes with deletion so a concurrent write cannot recreate
// a child under a parent that has just been removed.
func lockCategoryParent(tx *gorm.DB, category model.Category) error {
	if category.Fid < 0 || (category.ID > 0 && category.Fid == int(category.ID)) {
		return invalidArticle("上级分类不合法")
	}
	if category.Fid == 0 {
		return nil
	}
	var parent model.Category
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&parent, category.Fid).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return invalidArticle("上级分类不存在，请重新选择")
	}
	if err == nil && parent.Fid != 0 {
		return invalidArticle("目前仅支持两级分类，请选择一级分类作为上级")
	}
	return err
}

// DeleteCategory 删除Category记录
func (s *CategoryService) DeleteCategory(category model.Category) (err error) {
	return s.deleteCategories([]int{int(category.ID)})
}

// DeleteCategoryByIds atomically removes only empty categories. Hidden, scheduled
// and draft articles count too; their editors must keep a valid category.
func (s *CategoryService) DeleteCategoryByIds(ids request.IdsReq) error {
	return s.deleteCategories(ids.Ids)
}

func (s *CategoryService) deleteCategories(ids []int) error {
	if len(ids) == 0 || len(ids) > 100 {
		return invalidArticle("请选择 1 到 100 个分类")
	}
	for _, id := range ids {
		if id <= 0 {
			return invalidArticle("分类编号不合法")
		}
	}
	err := global.BLOG_DB.Transaction(func(tx *gorm.DB) error {
		var categories []model.Category
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id IN ?", ids).Order("id asc").Find(&categories).Error; err != nil {
			return err
		}
		var count int64
		if err := tx.Model(&model.Article{}).Where("fid IN ?", ids).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return invalidArticle("分类仍有文章（包括草稿、隐藏和定时文章），请先迁移或删除文章")
		}
		if err := tx.Model(&model.Category{}).Where("fid IN ? AND id NOT IN ?", ids, ids).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return invalidArticle("分类仍有子分类，请先调整子分类，或将空的父子分类一起删除")
		}
		return tx.Where("id IN ?", ids).Delete(&model.Category{}).Error
	})
	if err == nil {
		seoService.ResetSitemap()
	}
	return err
}
