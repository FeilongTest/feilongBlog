package blog

import (
	model "feilongBlog/model/blog"
	"feilongBlog/model/common/response"
	"feilongBlog/service/blog"
	"strings"

	"github.com/gin-gonic/gin"
)

type StatisticApi struct{}

var statisticApi = blog.StatisticService{}

// RecordVisit 记录网站访问量
func (a *StatisticApi) RecordVisit(c *gin.Context) {
	var req model.VisitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.VisitorID = strings.TrimSpace(req.VisitorID)
	if req.VisitorID == "" || len(req.VisitorID) > 64 {
		response.FailWithMessage("访客标识不合法", c)
		return
	}
	if err := statisticApi.RecordVisit(req.VisitorID); err != nil {
		response.FailWithMessage("记录访问失败", c)
		return
	}
	response.OkWithMessage("记录成功", c)
}

// GetDashboard 获取后台统计数据
func (a *StatisticApi) GetDashboard(c *gin.Context) {
	result, err := statisticApi.GetDashboard()
	if err != nil {
		response.FailWithMessage("获取统计数据失败", c)
		return
	}
	response.OkWithDetailed(result, "获取成功", c)
}
