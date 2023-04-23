package blog

import (
	"feilongBlog/global"
	"feilongBlog/model/blog"
)

type UserService struct {
}

func (s *UserService) GetUserList() (user blog.User, err error) {
	err = global.BLOG_DB.First(&user).Error
	return user, err
}
