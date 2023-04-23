package blog

import (
	"feilongBlog/model/common/response"
	"feilongBlog/service/blog"
	"github.com/gin-gonic/gin"
)

type UserApi struct {
}

var userService = blog.UserService{}

func (u *UserApi) GetUserList(c *gin.Context) {
	user, err := userService.GetUserList()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(user, "获取成功", c)
	}
	return
}
