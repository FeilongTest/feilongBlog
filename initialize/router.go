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

	PublicGroup := Router.Group("")
	{
		// 健康监测
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(200, "ok")
		})
	}
	{
		blogRouter.InitUserRouter(PublicGroup)    //注册用户路由
		blogRouter.InitArticleRouter(PublicGroup) //注册文章路由
	}
	global.BLOG_LOG.Info("router register success")
	return Router
}
