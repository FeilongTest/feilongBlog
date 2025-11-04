package blog

import "feilongBlog/model/common/request"

type Friendlink struct {
	ID     uint   `gorm:"primarykey"` // 主键ID
	Name   string `json:"name" form:"name" gorm:"name;comment:友链名称"`
	Url    string `json:"url" form:"url" gorm:"url;comment:友链地址"`
	Logo   string `json:"logo" form:"logo" gorm:"logo;comment:友链图标"`
	Desc   string `json:"desc" form:"desc" gorm:"desc;comment:友链描述"`
	Status int    `json:"status" form:"status" gorm:"status;comment:状态 0显示 1隐藏"`
	Ctime  int    `json:"ctime" form:"ctime" gorm:"ctime;comment:创建时间"`
}

type FriendlinkSearch struct {
	Status *int `json:"status" form:"status"`
	request.PageInfo
}

func (Friendlink) TableName() string {
	return "blog_friendlink"
}
