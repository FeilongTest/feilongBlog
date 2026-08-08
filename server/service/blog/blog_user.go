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

func (s *UserService) ChangePassword(userID uint, currentPassword, newPassword string) error {
	if global.BLOG_DB == nil {
		return errors.New("数据库未初始化")
	}
	if len(newPassword) < 8 || len(newPassword) > 72 {
		return errors.New("新密码长度需要在 8 到 72 位之间")
	}

	var user model.User
	if err := global.BLOG_DB.First(&user, userID).Error; err != nil {
		return err
	}
	if !passwordMatches(user.Password, currentPassword) {
		return errors.New("当前密码不正确")
	}
	if currentPassword == newPassword {
		return errors.New("新密码不能与当前密码相同")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return global.BLOG_DB.Model(&user).Update("password", string(hash)).Error
}

func (s *UserService) GetProfile(userID uint) (user model.User, err error) {
	err = global.BLOG_DB.First(&user, userID).Error
	return user, err
}

func (s *UserService) UpdateProfile(userID uint, input model.UpdateProfile) (user model.User, err error) {
	updates := map[string]interface{}{
		"truename": input.TrueName,
		"email":    input.Email,
		"pic":      input.Pic,
		"bio":      input.Bio,
	}
	if err = global.BLOG_DB.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		return user, err
	}
	err = global.BLOG_DB.First(&user, userID).Error
	return user, err
}

func (s *UserService) GetPublicContact() (contact model.PublicContact, err error) {
	var user model.User
	err = global.BLOG_DB.Select("truename", "email", "pic", "bio").
		Where("admin = ? AND status = ?", 1, 0).
		Order("id ASC").
		First(&user).Error
	if err == nil {
		contact = model.PublicContact{
			TrueName: user.TrueName,
			Email:    user.Email,
			Pic:      user.Pic,
			Bio:      user.Bio,
		}
	}
	return contact, err
}
