package blog

import (
	"feilongBlog/global"
	model "feilongBlog/model/blog"
	"feilongBlog/model/common/request"
	"feilongBlog/model/common/response"
	"feilongBlog/service/blog"
	"feilongBlog/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
)

type ArticleApi struct {
}

var articleApi = blog.ArticleService{}

// LikeArticle 点赞或取消点赞
func (a *ArticleApi) LikeArticle(c *gin.Context) {
	var req model.ArticleLikeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.VisitorID = strings.TrimSpace(req.VisitorID)
	if req.Aid == 0 || req.VisitorID == "" || len(req.VisitorID) > 64 {
		response.FailWithMessage("点赞参数不合法", c)
		return
	}
	liked, count, err := articleApi.LikeArticle(req)
	if err != nil {
		response.FailWithMessage("点赞失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"liked": liked, "count": count}, "操作成功", c)
}

// GetArticleLike 获取点赞状态
func (a *ArticleApi) GetArticleLike(c *gin.Context) {
	var req model.ArticleLikeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.VisitorID = strings.TrimSpace(req.VisitorID)
	if req.Aid == 0 || req.VisitorID == "" || len(req.VisitorID) > 64 {
		response.FailWithMessage("点赞参数不合法", c)
		return
	}
	liked, count, err := articleApi.GetArticleLike(req)
	if err != nil {
		response.FailWithMessage("获取点赞状态失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"liked": liked, "count": count}, "获取成功", c)
}

// GetArticleList 获取文章列表
func (a *ArticleApi) GetArticleList(c *gin.Context) {
	var pageInfo model.ArticleSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := articleApi.GetArticleList(pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
	return
}

// GetArticleListAdmin 管理员获取文章列表
func (a *ArticleApi) GetArticleListAdmin(c *gin.Context) {
	var pageInfo model.ArticleSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := articleApi.GetArticleListAdmin(pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
	return
}

// GetArticle 查看文章
func (a *ArticleApi) GetArticle(c *gin.Context) {
	var req request.GetById
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if article, err := articleApi.GetArticle(req.ID); err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(article, "获取成功", c)
	}
}

// CreateArticle 创建文章
func (a *ArticleApi) CreateArticle(c *gin.Context) {
	var article model.Article
	err := c.ShouldBindJSON(&article)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(article, utils.ArticleVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = articleApi.CreateArticle(article)
	if err != nil {
		global.BLOG_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// UpdateArticleById 更新文章
func (a *ArticleApi) UpdateArticleById(c *gin.Context) {
	var article model.Article
	err := c.ShouldBindJSON(&article)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(article, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err = articleApi.UpdateArticle(article); err != nil {
		global.BLOG_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// DelArticleById 删除Article记录
func (a *ArticleApi) DelArticleById(c *gin.Context) {
	var article model.Article
	err := c.ShouldBindJSON(&article)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(article, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = articleApi.DeleteArticle(article)
	if err != nil {
		global.BLOG_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DelArticleByIds 批量删除Article记录
func (a *ArticleApi) DelArticleByIds(c *gin.Context) {
	var ids request.IdsReq
	err := c.ShouldBindJSON(&ids)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = articleApi.DeleteArticleByIds(ids)
	if err != nil {
		global.BLOG_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// GetArticleSummary 获取文章分类概况
func (a *ArticleApi) GetArticleSummary(c *gin.Context) {
	var req request.GetById
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var article any
	if article, err = articleApi.GetArticleSummary(req.ID); err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(article, "获取成功", c)
	}
}
