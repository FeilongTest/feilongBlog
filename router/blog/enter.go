package blog

import (
	v1 "feilongBlog/api/v1"
	"github.com/gin-gonic/gin"
)

type ApiRouter struct{}

func (s *ApiRouter) InitBaseRouter(Router *gin.RouterGroup) {
	apiRouterWithoutRecord := Router.Group("base")
	baseRouterApi := v1.ApiGroupApp.BaseApiGroup
	{
		apiRouterWithoutRecord.GET("getCategoryList", baseRouterApi.GetCategoryList) // 创建Api
		apiRouterWithoutRecord.GET("getArticleList", baseRouterApi.GetArticleList)
		//TODO Login、GetArticle、
	}
}

func (s *ApiRouter) InitUserRouter(Router *gin.RouterGroup) {
	//apiRouterWithoutRecord := Router.Group("user")
	//apiRouterApi := v1.ApiGroupApp.BlogApiGroup
	//{
	//	apiRouterWithoutRecord.POST("getUserList", apiRouterApi.GetUserList) // 创建Api
	//}
}

func (s *ApiRouter) InitArticleRouter(Router *gin.RouterGroup) {
	//apiRouterWithoutRecord := Router.Group("article")
	//apiRouterApi := v1.ApiGroupApp.BlogApiGroup
	//{
	//	apiRouterWithoutRecord.POST("getArticleList", apiRouterApi.) // 创建Api
	//}
}
