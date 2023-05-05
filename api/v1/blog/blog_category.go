package blog

import (
	"feilongBlog/model/common/response"
	"feilongBlog/service/blog"
	"github.com/gin-gonic/gin"
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
