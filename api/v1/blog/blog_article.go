package blog

import (
	"feilongBlog/model/common/response"
	"feilongBlog/service/blog"
	"github.com/gin-gonic/gin"
)

type ArticleApi struct {
}

var articleApi = blog.ArticleService{}

func (u *UserApi) GetArticleList(c *gin.Context) {
	article, err := articleApi.GetArticleList()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(article, "获取成功", c)
	}
	return
}
