package blog

import (
	"feilongBlog/global"
	model "feilongBlog/model/blog"
	"feilongBlog/model/common/response"
	"feilongBlog/service/blog"
	"feilongBlog/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FriendlinkApi struct {
}

var friendlinkService = blog.FriendlinkService{}

// GetFriendlinkList 获取友链列表（管理员）
func (a *FriendlinkApi) GetFriendlinkList(c *gin.Context) {
	var pageInfo model.FriendlinkSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, total, err := friendlinkService.GetFriendlinkList(pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// GetAllFriendlinks 获取所有友链（前台）
func (a *FriendlinkApi) GetAllFriendlinks(c *gin.Context) {
	if list, err := friendlinkService.GetAllFriendlinks(); err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(list, c)
	}
}

// CreateFriendlink 创建友链
func (a *FriendlinkApi) CreateFriendlink(c *gin.Context) {
	var friendlink model.Friendlink
	err := c.ShouldBindJSON(&friendlink)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = friendlinkService.CreateFriendlink(friendlink)
	if err != nil {
		global.BLOG_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// UpdateFriendlink 更新友链
func (a *FriendlinkApi) UpdateFriendlink(c *gin.Context) {
	var friendlink model.Friendlink
	err := c.ShouldBindJSON(&friendlink)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = utils.Verify(friendlink, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err = friendlinkService.UpdateFriendlink(friendlink); err != nil {
		global.BLOG_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// DeleteFriendlink 删除友链
func (a *FriendlinkApi) DeleteFriendlink(c *gin.Context) {
	var friendlink model.Friendlink
	err := c.ShouldBindJSON(&friendlink)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = utils.Verify(friendlink, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = friendlinkService.DeleteFriendlink(friendlink)
	if err != nil {
		global.BLOG_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

