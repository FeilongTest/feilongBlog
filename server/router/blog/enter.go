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
	var commentApi = v1.ApiGroupApp.CommentApiGroup
	var friendlinkApi = v1.ApiGroupApp.FriendlinkApiGroup
	var statisticApi = v1.ApiGroupApp.StatisticApiGroup
	{
		apiRouterWithoutRecord.POST("login", userRouterApi.Login)                        // 登录
		apiRouterWithoutRecord.GET("getCategoryList", categoryApi.GetCategoryList)       // 获取分类
		apiRouterWithoutRecord.GET("getArticleList", articleApi.GetArticleList)          //获取文章列表
		apiRouterWithoutRecord.POST("getArticle", articleApi.GetArticle)                 //获取文章信息
		apiRouterWithoutRecord.POST("getSummary", articleApi.GetArticleSummary)          //获取分类概述
		apiRouterWithoutRecord.POST("likeArticle", articleApi.LikeArticle)               //点赞或取消点赞
		apiRouterWithoutRecord.GET("getArticleLike", articleApi.GetArticleLike)          //获取点赞状态
		apiRouterWithoutRecord.GET("getCommentList", commentApi.GetCommentList)          //获取评论列表
		apiRouterWithoutRecord.POST("createComment", commentApi.CreateComment)           //创建评论
		apiRouterWithoutRecord.GET("getAllFriendlinks", friendlinkApi.GetAllFriendlinks) //获取所有友链（前台）
		apiRouterWithoutRecord.POST("visit", statisticApi.RecordVisit)                   //记录网站访问
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

func (s *ApiRouter) InitCommentRouter(Router *gin.RouterGroup) {
	apiRouterWithoutRecord := Router.Group("comment")
	commentApi := v1.ApiGroupApp.CommentApiGroup
	{
		apiRouterWithoutRecord.GET("getCommentList", commentApi.GetCommentListAdmin)      // 获取所有评论（管理员）
		apiRouterWithoutRecord.DELETE("deleteComment", commentApi.DeleteComment)          // 删除评论
		apiRouterWithoutRecord.PUT("updateCommentStatus", commentApi.UpdateCommentStatus) // 更新评论状态
	}
}

func (s *ApiRouter) InitFriendlinkRouter(Router *gin.RouterGroup) {
	apiRouterWithoutRecord := Router.Group("friendlink")
	friendlinkApi := v1.ApiGroupApp.FriendlinkApiGroup
	{
		apiRouterWithoutRecord.GET("getFriendlinkList", friendlinkApi.GetFriendlinkList)  // 获取友链列表
		apiRouterWithoutRecord.POST("createFriendlink", friendlinkApi.CreateFriendlink)   // 创建友链
		apiRouterWithoutRecord.PUT("updateFriendlink", friendlinkApi.UpdateFriendlink)    // 更新友链
		apiRouterWithoutRecord.DELETE("deleteFriendlink", friendlinkApi.DeleteFriendlink) // 删除友链
	}
}

func (s *ApiRouter) InitStatisticRouter(Router *gin.RouterGroup) {
	apiRouterWithoutRecord := Router.Group("statistics")
	statisticApi := v1.ApiGroupApp.StatisticApiGroup
	{
		apiRouterWithoutRecord.GET("dashboard", statisticApi.GetDashboard)
	}
}
