package blog

import (
	"feilongBlog/global"
	model "feilongBlog/model/blog"
	"feilongBlog/model/common/response"
	"feilongBlog/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BaseApi struct {
}

// Login 后台登录
func (a *BaseApi) Login(c *gin.Context) {
	var l model.Login
	err := c.ShouldBindJSON(&l)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(l, utils.LoginVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	u := &model.User{Username: l.Username, Password: l.Password}
	user, err := userService.Login(u)
	if err != nil {
		global.BLOG_LOG.Error("登陆失败! 用户名不存在或者密码错误!", zap.Error(err))
		response.FailWithMessage("用户名不存在或者密码错误", c)
		return
	}
	if user.Status != 0 {
		global.BLOG_LOG.Error("登陆失败! 用户被禁止登录!")
		response.FailWithMessage("用户被禁止登录", c)
		return
	}
	if user.Admin != 1 {
		global.BLOG_LOG.Error("登陆失败! 非管理员用户!")
		response.FailWithMessage("非管理员用户被禁止登录", c)
		return
	}
	a.TokenNext(c, *user)
	return
}

// TokenNext 登录以后签发jwt
func (a *BaseApi) TokenNext(c *gin.Context, user model.User) {
	j := &utils.JWT{SigningKey: []byte(global.BLOG_CONFIG.JWT.SigningKey)} // 唯一签名
	claims := j.CreateClaims(model.BaseClaims{
		ID:       user.ID,
		Username: user.Username,
		TrueName: user.TrueName,
		Admin:    user.Admin,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		global.BLOG_LOG.Error("获取token失败!", zap.Error(err))
		response.FailWithMessage("获取token失败", c)
		return
	}
	if !global.BLOG_CONFIG.System.UseMultipoint {
		response.OkWithDetailed(model.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		}, "登录成功", c)
		return
	}
}

// VeryToken 验证Token
func (a *BaseApi) VeryToken(c *gin.Context) {
	//TODO 是否需要验证Token
}

// GetArticleList 获取文章列表
func (a *BaseApi) GetArticleList(c *gin.Context) {
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
