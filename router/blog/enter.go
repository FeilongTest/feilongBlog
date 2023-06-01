package blog

import (
	v1 "feilongBlog/api/v1"
	"github.com/gin-gonic/gin"
)

type ApiRouter struct{}

func (s *ApiRouter) InitBaseRouter(Router *gin.RouterGroup) {
	apiRouterWithoutRecord := Router.Group("base")
	var userRouterApi = v1.ApiGroupApp.UserApiGroup
	var articleApi = v1.ApiGroupApp.ArticleApiGroup
	var categoryApi = v1.ApiGroupApp.CategoryApiGroup
	{
		apiRouterWithoutRecord.POST("login", userRouterApi.Login)                  // 登录
		apiRouterWithoutRecord.GET("getCategoryList", categoryApi.GetCategoryList) // 获取分类
		apiRouterWithoutRecord.GET("getArticleList", articleApi.GetArticleList)    //获取文章列表
		apiRouterWithoutRecord.POST("getArticle", articleApi.GetArticle)           //获取文章信息
		apiRouterWithoutRecord.POST("getSummary", articleApi.GetArticleSummary)    //获取分类概述
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
	apiRouterWithoutRecord := Router.Group("article")
	apiRouterApi := v1.ApiGroupApp.ArticleApiGroup
	{
		apiRouterWithoutRecord.POST("getArticle", apiRouterApi.GetArticle)
		apiRouterWithoutRecord.GET("getArticleList", apiRouterApi.GetArticleListAdmin)
		apiRouterWithoutRecord.POST("createArticle", apiRouterApi.CreateArticle)
		apiRouterWithoutRecord.DELETE("delArticle", apiRouterApi.DelArticleById)
		apiRouterWithoutRecord.DELETE("delArticleByIds", apiRouterApi.DelArticleByIds)
		apiRouterWithoutRecord.PUT("updateArticle", apiRouterApi.UpdateArticleById)
	}
}

func (s *ApiRouter) InitCategoryRouter(Router *gin.RouterGroup) {
	apiRouterWithoutRecord := Router.Group("category")
	apiRouterApi := v1.ApiGroupApp.CategoryApiGroup
	{
		apiRouterWithoutRecord.GET("getCategoryList", apiRouterApi.GetCategoryList)
		apiRouterWithoutRecord.POST("createCategory", apiRouterApi.CreateCategory)
		apiRouterWithoutRecord.DELETE("delCategory", apiRouterApi.DelCategoryById)
		apiRouterWithoutRecord.DELETE("delCategoryByIds", apiRouterApi.DelCategoryByIds)
		apiRouterWithoutRecord.PUT("updateCategory", apiRouterApi.UpdateCategoryById)
	}
}

func (s *ApiRouter) InitFileRouter(Router *gin.RouterGroup) {
	apiRouterWithoutRecord := Router.Group("file")
	apiRouterApi := v1.ApiGroupApp.FileApiGroup
	{
		apiRouterWithoutRecord.POST("upload", apiRouterApi.UploadFile) // 上传文件
	}
}
