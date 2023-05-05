package blog

import (
	"feilongBlog/global"
	model "feilongBlog/model/blog"
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
