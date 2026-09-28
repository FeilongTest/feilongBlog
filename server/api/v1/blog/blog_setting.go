package blog

import (
	model "feilongBlog/model/blog"
	"feilongBlog/model/common/response"
	"feilongBlog/service/blog"
	"feilongBlog/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

type SettingApi struct {
}

var siteSettingService = blog.SettingService{}

// GetPublicSetting 获取前台站点信息
func (a *SettingApi) GetPublicSetting(c *gin.Context) {
	setting, err := siteSettingService.GetPublicSetting()
	if err != nil {
		response.FailWithMessage("获取站点信息失败", c)
		return
	}
	response.OkWithData(setting, c)
}

// GetSetting 获取站点设置（管理员）
func (a *SettingApi) GetSetting(c *gin.Context) {
	if !utils.IsAdmin(c) {
		response.FailWithMessage("仅管理员可以查看站点设置", c)
		return
	}
	setting, err := siteSettingService.GetSetting()
	if err != nil {
		response.FailWithMessage("获取站点设置失败", c)
		return
	}
	response.OkWithData(setting, c)
}

// UpdateSetting 更新站点设置（管理员）
func (a *SettingApi) UpdateSetting(c *gin.Context) {
	if !utils.IsAdmin(c) {
		response.FailWithMessage("仅管理员可以修改站点设置", c)
		return
	}
	var input model.UpdateSetting
	if err := c.ShouldBindJSON(&input); err != nil {
		response.FailWithMessage("请求参数不正确", c)
		return
	}
	if len([]rune(strings.TrimSpace(input.SiteName))) > 60 {
		response.FailWithMessage("站点名称不能超过 60 个字符", c)
		return
	}
	if len([]rune(strings.TrimSpace(input.Description))) > 200 {
		response.FailWithMessage("站点描述不能超过 200 个字符", c)
		return
	}
	if len([]rune(strings.TrimSpace(input.Keywords))) > 200 {
		response.FailWithMessage("站点关键词不能超过 200 个字符", c)
		return
	}
	setting, err := siteSettingService.UpdateSetting(input)
	if err != nil {
		response.FailWithMessage("更新站点设置失败", c)
		return
	}
	response.OkWithDetailed(setting, "站点设置已更新", c)
}
