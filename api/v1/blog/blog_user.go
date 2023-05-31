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

type UserApi struct {
}

var userService = blog.UserService{}

// Login 后台登录
func (a *UserApi) Login(c *gin.Context) {
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
func (a *UserApi) TokenNext(c *gin.Context, user model.User) {
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
func (a *UserApi) VeryToken(c *gin.Context) {
	//TODO 是否需要验证Token
}

// GetUserList 获取用户列表
func (a *UserApi) GetUserList(c *gin.Context) {
	user, err := userService.GetUserList()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(user, "获取成功", c)
	}
	return
}
