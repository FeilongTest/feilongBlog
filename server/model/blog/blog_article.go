package blog

import (
	"feilongBlog/model/common/request"
)

type Article struct {
	ID           uint      `gorm:"primarykey"` // 主键ID
	Title        string    `json:"title" gorm:"title;comment:标题"`
	Content      string    `json:"content" gorm:"type:text;column:content;comment:文章内容"`
	Pic          string    `json:"pic" gorm:"pic;comment:图片"`
	Uid          int       `json:"uid" gorm:"uid;comment:用户id"`
	Fid          int       `json:"fid" gorm:"fid;comment:分类id"`
	View         int       `json:"view" gorm:"view;comment:浏览次数"`
	Ctime        int       `json:"ctime" gorm:"ctime;comment:创建时间"`
	EditTime     int       `json:"editTime" gorm:"editTime;column:edittime;comment:编辑时间"`
	File         string    `json:"file" gorm:"pic;file:附件"`
	Type         int       `json:"type" gorm:"type;comment:文章类型"`
	IsTop        int       `json:"isTop" gorm:"isTop;column:istop;comment:置顶"`
	Status       int       `json:"status" gorm:"status;column:status;comment:状态"`
	Comment      []Comment `json:"comment" gorm:"foreignKey:Aid;references:ID;comment:评论"`
	CommentCount int       `json:"commentCount" gorm:"-"` // 评论数量，不存储在数据库
	LikeCount    int       `json:"likeCount" gorm:"-"`    // 点赞数量，不存储在数据库
}

type ArticleLike struct {
	ID        uint   `gorm:"primarykey"`
	Aid       uint   `json:"aid" gorm:"uniqueIndex:idx_article_visitor"`
	VisitorID string `json:"visitorId" gorm:"size:64;uniqueIndex:idx_article_visitor"`
	Ctime     int    `json:"ctime"`
}

type ArticleLikeRequest struct {
	Aid       uint   `json:"aid"`
	VisitorID string `json:"visitorId"`
}

type ArticleSearch struct {
	Fid     int    `json:"fid" form:"fid"`
	Hide    bool   `json:"hide" form:"hide"`
	Top     bool   `json:"top" form:"top"`
	Keyword string `json:"keyword" form:"keyword"`
	request.PageInfo
}

func (Article) TableName() string {
	return "blog_article"
}

func (ArticleLike) TableName() string {
	return "blog_article_like"
}
