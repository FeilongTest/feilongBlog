package initialize

import (
	"feilongBlog/global"
	"feilongBlog/router"
	"github.com/gin-gonic/gin"
)

// 初始化总路由

func Routers() *gin.Engine {
	Router := gin.Default()
	blogRouter := router.RouterGroupApp.Blog

	// 跨域，如需跨域可以打开下面的注释
	//Router.Use(middleware.Cors()) // 直接放行全部跨域请求
	//Router.Use(middleware.CorsByRules()) // 按照配置的规则放行跨域请求
	//global.GVA_LOG.Info("use middleware cors")
	PublicGroup := Router.Group("")
	{
		// 健康监测
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(200, "ok")
		})
	}
	{
		blogRouter.InitBaseRouter(PublicGroup)    //注册公用路由
		blogRouter.InitUserRouter(PublicGroup)    //注册用户路由
		blogRouter.InitArticleRouter(PublicGroup) //注册文章路由
	}
	global.BLOG_LOG.Info("router register success")
	return Router
}
