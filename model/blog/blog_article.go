package blog

import (
	"feilongBlog/model/common/request"
)

type Article struct {
	//global.GVA_MODEL
	ID      uint   `gorm:"primarykey"` // 主键ID
	Title   string `json:"title" gorm:"title;comment:标题"`
	Content string `json:"content" gorm:"type:text;column:content;comment:文章内容"`
	Pic     string `json:"pic" gorm:"pic;comment:图片"`
	Uid     int    `json:"uid" gorm:"uid;comment:用户id"`
	Fid     int    `json:"fid" gorm:"fid;comment:分类id"`
	View    int    `json:"view" gorm:"view;comment:浏览次数"`
}

type ArticleSearch struct {
	Fid int `json:"fid" form:"fid"`
	request.PageInfo
}

func (Article) TableName() string {
	return "blog_article"
}
