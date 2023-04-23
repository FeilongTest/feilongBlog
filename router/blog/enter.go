package blog

import (
	v1 "feilongBlog/api/v1"
	"github.com/gin-gonic/gin"
)

type ApiRouter struct{}

func (s *ApiRouter) InitUserRouter(Router *gin.RouterGroup) {
	apiRouterWithoutRecord := Router.Group("user")
	apiRouterApi := v1.ApiGroupApp.BlogApiGroup
	{
		apiRouterWithoutRecord.POST("getUserList", apiRouterApi.GetUserList) // 创建Api
	}
}

func (s *ApiRouter) InitArticleRouter(Router *gin.RouterGroup) {
	apiRouterWithoutRecord := Router.Group("article")
	apiRouterApi := v1.ApiGroupApp.BlogApiGroup
	{
		apiRouterWithoutRecord.POST("getArticleList", apiRouterApi.GetArticleList) // 创建Api
	}
}
