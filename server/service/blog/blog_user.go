package blog

import (
	"errors"
	"feilongBlog/global"
	model "feilongBlog/model/blog"
	"feilongBlog/utils"
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
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
		if !passwordMatches(user.Password, u.Password) {
			return nil, errors.New("密码错误")
		}
		if regexp.MustCompile(`^[a-fA-F0-9]{32}$`).MatchString(user.Password) {
			hash, hashErr := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
			if hashErr != nil {
				return nil, hashErr
			}
			if err = global.BLOG_DB.Model(&user).Update("password", string(hash)).Error; err != nil {
				return nil, err
			}
		}
	}
	return &user, err
}

func passwordMatches(stored, password string) bool {
	if regexp.MustCompile(`^[a-fA-F0-9]{32}$`).MatchString(stored) {
		return utils.MD5V([]byte(password)) == stored
	}
	return bcrypt.CompareHashAndPassword([]byte(stored), []byte(password)) == nil
}
