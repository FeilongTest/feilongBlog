package blog

import (
	"errors"
	"feilongBlog/global"
	"go.uber.org/zap"
	"net/http"
	"strconv"

	model "feilongBlog/model/blog"
	"feilongBlog/model/common/response"
	service "feilongBlog/service/blog"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func articleFailure(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	message := "服务暂时不可用，请稍后重试"
	var validation *service.ValidationError
	switch {
	case errors.Is(err, service.ErrArticleConflict):
		status = http.StatusConflict
		message = err.Error()
	case errors.Is(err, gorm.ErrRecordNotFound):
		status = http.StatusNotFound
		message = "文章或标签不存在，或暂不公开"
	case errors.As(err, &validation):
		status = http.StatusBadRequest
		message = validation.Error()
	default:
		if global.BLOG_LOG != nil {
			global.BLOG_LOG.Error("内容服务请求失败", zap.Error(err))
		}
	}
	c.JSON(status, response.Response{Code: response.ERROR, Data: gin.H{}, Msg: message})
}

func requestFailure(c *gin.Context) {
	c.JSON(http.StatusBadRequest, response.Response{Code: response.ERROR, Data: gin.H{}, Msg: "请求参数不合法，请检查输入"})
}

func (a *ArticleApi) PatchState(c *gin.Context) {
	var input model.ArticleStatePatch
	if err := c.ShouldBindJSON(&input); err != nil {
		requestFailure(c)
		return
	}
	article, err := articleApi.PatchArticleState(input)
	if err != nil {
		articleFailure(c, err)
		return
	}
	response.OkWithData(article, c)
}

func (a *ArticleApi) Revisions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil || id == 0 {
		response.FailWithMessage("文章编号不合法", c)
		return
	}
	revisions, err := articleApi.ListRevisions(uint(id))
	if err != nil {
		articleFailure(c, err)
		return
	}
	response.OkWithData(revisions, c)
}

func (a *ArticleApi) Revision(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	revisionID, revisionErr := strconv.ParseUint(c.Query("revisionId"), 10, 64)
	if err != nil || revisionErr != nil || id == 0 || revisionID == 0 {
		response.FailWithMessage("版本编号不合法", c)
		return
	}
	snapshot, err := articleApi.GetRevision(uint(id), uint(revisionID))
	if err != nil {
		articleFailure(c, err)
		return
	}
	c.Header("Cache-Control", "no-store")
	response.OkWithData(snapshot, c)
}

// Tag endpoints share admin middleware for mutations and return only public counts anonymously.
func (a *ArticleApi) PublicTags(c *gin.Context) {
	tags, err := (&service.TagService{}).List(true)
	if err != nil {
		articleFailure(c, err)
		return
	}
	response.OkWithData(tags, c)
}
func (a *ArticleApi) PublicTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil || id == 0 {
		response.FailWithMessage("标签编号不合法", c)
		return
	}
	tag, err := (&service.TagService{}).GetPublic(uint(id))
	if err != nil {
		articleFailure(c, err)
		return
	}
	response.OkWithData(tag, c)
}
func (a *ArticleApi) AdminTags(c *gin.Context) {
	tags, err := (&service.TagService{}).List(false)
	if err != nil {
		articleFailure(c, err)
		return
	}
	response.OkWithData(tags, c)
}
func (a *ArticleApi) SaveTag(c *gin.Context) {
	var tag model.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		requestFailure(c)
		return
	}
	if err := (&service.TagService{}).Save(&tag); err != nil {
		articleFailure(c, err)
		return
	}
	response.OkWithData(tag, c)
}
func (a *ArticleApi) DeleteTag(c *gin.Context) {
	var input struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		requestFailure(c)
		return
	}
	if err := (&service.TagService{}).Delete(input.ID); err != nil {
		articleFailure(c, err)
		return
	}
	response.OkWithMessage("标签已移除，文章仍然保留", c)
}

func (a *ArticleApi) ReadingArticles(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 8192)
	var input struct {
		IDs []uint `json:"ids" binding:"required,min=1,max=100,dive,gt=0"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		requestFailure(c)
		return
	}
	articles, err := articleApi.ReadingArticles(input.IDs)
	if err != nil {
		articleFailure(c, err)
		return
	}
	response.OkWithData(articles, c)
}
