package blog

import (
	v1 "feilongBlog/api/v1"
	"github.com/gin-gonic/gin"
)

type ApiRouter struct{}

func (s *ApiRouter) InitBaseRouter(Router *gin.RouterGroup) {
	apiRouterWithoutRecord := Router.Group("base")
	baseRouterApi := v1.ApiGroupApp.BaseApiGroup
	articleApi := v1.ApiGroupApp.ArticleApiGroup
	categoryApi := v1.ApiGroupApp.CategoryApiGroup
	{
		apiRouterWithoutRecord.POST("login", baseRouterApi.Login)                  // 登录
		apiRouterWithoutRecord.GET("getCategoryList", categoryApi.GetCategoryList) // 获取分类
		apiRouterWithoutRecord.GET("getArticleList", articleApi.GetArticleList)    //获取文章列表
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

func (s *ApiRouter) InitCategoryRouter(Router *gin.RouterGroup) {
	apiRouterWithoutRecord := Router.Group("category")
	apiRouterApi := v1.ApiGroupApp.CategoryApiGroup
	{
		apiRouterWithoutRecord.GET("getCategoryList", apiRouterApi.GetCategoryList) //查
		//TODO 增删改
	}
}
