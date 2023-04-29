package blog

import (
	model "feilongBlog/model/blog"
	"feilongBlog/model/common/response"
	"feilongBlog/service/blog"
	"github.com/gin-gonic/gin"
)

type ArticleApi struct {
}

var articleApi = blog.ArticleService{}

func (a *ArticleApi) GetArticleList(c *gin.Context) {
	var pageInfo model.ArticleSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := articleApi.GetArticleList(pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
	return
}
