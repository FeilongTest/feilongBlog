package blog

import (
	v1 "feilongBlog/api/v1"

	"github.com/gin-gonic/gin"
)

type FriendlinkRouter struct{}

func (a *FriendlinkRouter) InitFriendlinkRouter(Router *gin.RouterGroup) {
	friendlinkRouter := Router.Group("friendlink")
	friendlinkApi := v1.ApiGroupApp.FriendlinkApiGroup
	{
		friendlinkRouter.GET("getFriendlinkList", friendlinkApi.GetFriendlinkList)  // 获取友链列表
		friendlinkRouter.POST("createFriendlink", friendlinkApi.CreateFriendlink)   // 创建友链
		friendlinkRouter.PUT("updateFriendlink", friendlinkApi.UpdateFriendlink)    // 更新友链
		friendlinkRouter.DELETE("deleteFriendlink", friendlinkApi.DeleteFriendlink) // 删除友链
	}
}

func (a *FriendlinkRouter) InitFriendlinkPublicRouter(Router *gin.RouterGroup) {
	friendlinkRouter := Router.Group("base")
	friendlinkApi := v1.ApiGroupApp.FriendlinkApiGroup
	{
		friendlinkRouter.GET("getAllFriendlinks", friendlinkApi.GetAllFriendlinks) // 获取所有友链（前台）
	}
}
