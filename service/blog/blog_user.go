package blog

import (
	"errors"
	"feilongBlog/global"
	model "feilongBlog/model/blog"
	"feilongBlog/utils"
	"fmt"
)

type UserService struct {
}

func (s *UserService) GetUserList() (user model.User, err error) {
	err = global.BLOG_DB.First(&user).Error
	return user, err
}

func (s *UserService) Login(u *model.User) (userInter *model.User, err error) {
	if nil == global.BLOG_DB {
		return nil, fmt.Errorf("db not init")
	}
	var user model.User
	err = global.BLOG_DB.Where("username = ?", u.Username).First(&user).Error
	if err == nil {
		if utils.MD5V([]byte(u.Password)) != user.Password {
			return nil, errors.New("密码错误")
		}
	}
	return &user, err
}
