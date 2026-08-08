package blog

import (
	"feilongBlog/global"
	model "feilongBlog/model/blog"
	"feilongBlog/model/common/response"
	"feilongBlog/service/blog"
	"feilongBlog/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/mail"
	"strings"
)

type UserApi struct {
}

// GetProfile 获取当前管理员资料
func (a *UserApi) GetProfile(c *gin.Context) {
	if !utils.IsAdmin(c) {
		response.FailWithMessage("仅管理员可以查看账户资料", c)
		return
	}
	user, err := userService.GetProfile(utils.GetUserID(c))
	if err != nil {
		response.FailWithMessage("获取账户资料失败", c)
		return
	}
	response.OkWithData(user, c)
}

// UpdateProfile 更新当前管理员资料
func (a *UserApi) UpdateProfile(c *gin.Context) {
	if !utils.IsAdmin(c) {
		response.FailWithMessage("仅管理员可以修改账户资料", c)
		return
	}
	var input model.UpdateProfile
	if err := c.ShouldBindJSON(&input); err != nil {
		response.FailWithMessage("请求参数不正确", c)
		return
	}
	input.TrueName = strings.TrimSpace(input.TrueName)
	input.Email = strings.TrimSpace(input.Email)
	input.Pic = strings.TrimSpace(input.Pic)
	input.Bio = strings.TrimSpace(input.Bio)
	if len([]rune(input.TrueName)) > 50 {
		response.FailWithMessage("显示名称不能超过 50 个字符", c)
		return
	}
	if len(input.Email) > 100 {
		response.FailWithMessage("邮箱不能超过 100 个字符", c)
		return
	}
	if input.Email != "" {
		address, err := mail.ParseAddress(input.Email)
		if err != nil || address.Address != input.Email {
			response.FailWithMessage("请输入有效的邮箱地址", c)
			return
		}
	}
	if len(input.Pic) > 500 {
		response.FailWithMessage("头像地址过长", c)
		return
	}
	if len([]rune(input.Bio)) > 120 {
		response.FailWithMessage("个人简介不能超过 120 个字符", c)
		return
	}

	user, err := userService.UpdateProfile(utils.GetUserID(c), input)
	if err != nil {
		response.FailWithMessage("更新账户资料失败", c)
		return
	}
	response.OkWithDetailed(user, "账户资料已更新", c)
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

// GetPublicContact 获取前台公开联系信息
func (a *UserApi) GetPublicContact(c *gin.Context) {
	contact, err := userService.GetPublicContact()
	if err != nil {
		response.FailWithMessage("获取联系信息失败", c)
		return
	}
	response.OkWithData(contact, c)
}

// ChangePassword 修改当前管理员密码
func (a *UserApi) ChangePassword(c *gin.Context) {
	if !utils.IsAdmin(c) {
		response.FailWithMessage("仅管理员可以修改登录密码", c)
		return
	}

	var input model.ChangePassword
	if err := c.ShouldBindJSON(&input); err != nil {
		response.FailWithMessage("请求参数不正确", c)
		return
	}
	if len(input.CurrentPassword) == 0 {
		response.FailWithMessage("请输入当前密码", c)
		return
	}
	if len(input.NewPassword) < 8 || len(input.NewPassword) > 72 {
		response.FailWithMessage("新密码长度需要在 8 到 72 位之间", c)
		return
	}

	userID := utils.GetUserID(c)
	if userID == 0 {
		response.FailWithMessage("无法获取当前用户信息", c)
		return
	}
	if err := userService.ChangePassword(userID, input.CurrentPassword, input.NewPassword); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("密码修改成功，请重新登录", c)
}
