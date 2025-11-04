package blog

import (
	"feilongBlog/global"
	"feilongBlog/model/blog"
	"feilongBlog/model/common/response"
	"time"

	"github.com/gin-gonic/gin"
)

type CommentApi struct{}

// GetCommentList 获取评论列表
func (c *CommentApi) GetCommentList(ctx *gin.Context) {
	var comment blog.Comment
	err := ctx.ShouldBindQuery(&comment)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	// 查询指定文章的评论，按时间倒序 (status = 0 表示显示)
	var comments []blog.Comment
	err = global.BLOG_DB.Where("aid = ? AND status = ?", comment.Aid, 0).
		Order("ctime desc").
		Find(&comments).Error

	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.OkWithData(comments, ctx)
}

// CreateComment 创建评论
func (c *CommentApi) CreateComment(ctx *gin.Context) {
	var comment blog.Comment
	err := ctx.ShouldBindJSON(&comment)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	// 设置创建时间和状态
	comment.Ctime = int(time.Now().Unix())
	comment.Status = 0 // 0表示显示/已审核通过

	// 创建评论
	err = global.BLOG_DB.Create(&comment).Error
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.OkWithMessage("评论成功", ctx)
}

// DeleteComment 删除评论（管理员）
func (c *CommentApi) DeleteComment(ctx *gin.Context) {
	var comment blog.Comment
	err := ctx.ShouldBindJSON(&comment)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	err = global.BLOG_DB.Delete(&blog.Comment{}, comment.ID).Error
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.OkWithMessage("删除成功", ctx)
}

// GetCommentListAdmin 获取所有评论（管理员）
func (c *CommentApi) GetCommentListAdmin(ctx *gin.Context) {
	var pageInfo struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}
	err := ctx.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	// 设置默认值
	if pageInfo.Page <= 0 {
		pageInfo.Page = 1
	}
	if pageInfo.PageSize <= 0 {
		pageInfo.PageSize = 10
	}

	var comments []blog.Comment
	var total int64

	db := global.BLOG_DB.Model(&blog.Comment{})

	// 获取总数
	err = db.Count(&total).Error
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	// 分页查询
	limit := pageInfo.PageSize
	offset := pageInfo.PageSize * (pageInfo.Page - 1)
	err = db.Limit(limit).Offset(offset).Order("ctime desc").Find(&comments).Error
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     comments,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", ctx)
}

// UpdateCommentStatus 更新评论状态（管理员）
func (c *CommentApi) UpdateCommentStatus(ctx *gin.Context) {
	var comment blog.Comment
	err := ctx.ShouldBindJSON(&comment)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	err = global.BLOG_DB.Model(&blog.Comment{}).
		Where("id = ?", comment.ID).
		Update("status", comment.Status).Error

	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.OkWithMessage("更新成功", ctx)
}
