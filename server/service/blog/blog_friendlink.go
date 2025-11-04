package blog

import (
	"feilongBlog/global"
	model "feilongBlog/model/blog"
	"time"
)

type FriendlinkService struct {
}

// GetFriendlinkList 获取友链列表
func (s *FriendlinkService) GetFriendlinkList(info model.FriendlinkSearch) (friendlink []model.Friendlink, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.BLOG_DB.Model(&model.Friendlink{})

	// 如果有状态过滤
	if info.Status != nil {
		db = db.Where("status = ?", *info.Status)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Order("id desc").Find(&friendlink).Error
	return friendlink, total, err
}

// GetAllFriendlinks 获取所有显示的友链（前台使用）
func (s *FriendlinkService) GetAllFriendlinks() (friendlink []model.Friendlink, err error) {
	err = global.BLOG_DB.Where("status = ?", 0).Order("id desc").Find(&friendlink).Error
	return friendlink, err
}

// CreateFriendlink 创建友链
func (s *FriendlinkService) CreateFriendlink(friendlink model.Friendlink) (err error) {
	friendlink.Ctime = int(time.Now().Unix())
	return global.BLOG_DB.Create(&friendlink).Error
}

// UpdateFriendlink 更新友链
func (s *FriendlinkService) UpdateFriendlink(friendlink model.Friendlink) (err error) {
	return global.BLOG_DB.Save(&friendlink).Error
}

// DeleteFriendlink 删除友链
func (s *FriendlinkService) DeleteFriendlink(friendlink model.Friendlink) (err error) {
	return global.BLOG_DB.Delete(&friendlink).Error
}

