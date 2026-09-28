package initialize

import (
	v1 "feilongBlog/api/v1"
	"feilongBlog/global"
	"feilongBlog/middleware"
	"feilongBlog/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 初始化总路由

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.MaxMultipartMemory = 10 << 20
	Router.Use(middleware.CorsByRules())
	blogRouter := router.RouterGroupApp.Blog

	// 跨域，如需跨域可以打开下面的注释
	//Router.Use(middleware.Cors()) // 直接放行全部跨域请求
	//Router.Use(middleware.CorsByRules()) // 按照配置的规则放行跨域请求
	//global.GVA_LOG.Info("use middleware cors")
	Router.StaticFS(global.BLOG_CONFIG.Local.Path, http.Dir(global.BLOG_CONFIG.Local.StorePath)) // 为用户头像和文件提供静态地址

	PublicGroup := Router.Group("")
	{
		// 健康监测
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(200, "ok")
		})
	}
	{
		blogRouter.InitBaseRouter(PublicGroup)    //注册公用路由
		blogRouter.InitArticleRouter(PublicGroup) //注册文章路由
		blogRouter.InitSeoRouter(PublicGroup)     //注册站点地图、爬虫规则与站点信息
	}

	PrivateGroup := Router.Group("/admin")
	PrivateGroup.Use(middleware.JWTAuth()).Use()
	{
		blogRouter.InitUserRouter(PrivateGroup)
		blogRouter.InitFileRouter(PrivateGroup)
		blogRouter.InitCategoryRouter(PrivateGroup)
		blogRouter.InitArticleRouter(PrivateGroup)
		blogRouter.InitCommentRouter(PrivateGroup)
		blogRouter.InitFriendlinkRouter(PrivateGroup)
		blogRouter.InitStatisticRouter(PrivateGroup)
		blogRouter.InitSettingRouter(PrivateGroup)
	}

	// 未命中接口的请求交由服务端渲染，为前端路由注入 SEO 信息
	Router.NoRoute(v1.ApiGroupApp.SeoApiGroup.RenderPage)
	global.BLOG_LOG.Info("router register success")
	return Router
}
