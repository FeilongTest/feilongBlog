package blog

import (
	"errors"
	"feilongBlog/global"
	model "feilongBlog/model/blog"
	"feilongBlog/model/common/request"
	"gorm.io/gorm"
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
	return global.BLOG_DB.Create(&category).Error
}

// UpdateCategory 更新分类
func (s *CategoryService) UpdateCategory(category model.Category) (err error) {
	if errors.Is(global.BLOG_DB.Where("id = ?", category.ID).First(&model.Category{}).Error, gorm.ErrRecordNotFound) {
		return gorm.ErrRecordNotFound
	}
	err = global.BLOG_DB.Save(&category).Error
	return err
}

// DeleteCategory 删除Category记录
func (s *CategoryService) DeleteCategory(category model.Category) (err error) {
	err = global.BLOG_DB.Delete(&category).Error
	return err
}

// DeleteCategoryByIds 批量删除Category记录
func (s *CategoryService) DeleteCategoryByIds(ids request.IdsReq) (err error) {
	err = global.BLOG_DB.Delete(&[]model.Category{}, "id in ?", ids.Ids).Error
	return err
}
