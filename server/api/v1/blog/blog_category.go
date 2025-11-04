package blog

import (
	"feilongBlog/global"
	model "feilongBlog/model/blog"
	"feilongBlog/model/common/request"
	"feilongBlog/model/common/response"
	"feilongBlog/service/blog"
	"feilongBlog/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CategoryApi struct {
}

var categoryApi = blog.CategoryService{}

// GetCategoryList 获取分类列表
func (a *CategoryApi) GetCategoryList(c *gin.Context) {
	categoryList, err := categoryApi.GetCategoryList()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(categoryList, "获取成功", c)
	}
	return
}

// CreateCategory 创建分类
func (a *CategoryApi) CreateCategory(c *gin.Context) {
	var category model.Category
	err := c.ShouldBindJSON(&category)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(category, utils.CategoryVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = categoryApi.CreateCategory(category)
	if err != nil {
		global.BLOG_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// UpdateCategoryById 更新分类
func (a *CategoryApi) UpdateCategoryById(c *gin.Context) {
	var category model.Category
	err := c.ShouldBindJSON(&category)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(category, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err = categoryApi.UpdateCategory(category); err != nil {
		global.BLOG_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// DelCategoryById 删除Category记录
func (a *CategoryApi) DelCategoryById(c *gin.Context) {
	var category model.Category
	err := c.ShouldBindJSON(&category)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(category, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = categoryApi.DeleteCategory(category)
	if err != nil {
		global.BLOG_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DelCategoryByIds 批量删除Category记录
func (a *CategoryApi) DelCategoryByIds(c *gin.Context) {
	var ids request.IdsReq
	err := c.ShouldBindJSON(&ids)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = categoryApi.DeleteCategoryByIds(ids)
	if err != nil {
		global.BLOG_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
