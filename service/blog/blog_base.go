package blog

import (
	"feilongBlog/global"
	"feilongBlog/model/blog"
)

type BaseService struct {
}

// GetCategoryList 获取分类列表
func (s *BaseService) GetCategoryList() (list []blog.Category, err error) {
	var child []blog.Category
	var root []blog.Category
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
