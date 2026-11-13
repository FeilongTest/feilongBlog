package blog

import (
	"feilongBlog/model/common/request"
)

type Article struct {
	ID             uint      `gorm:"primarykey"` // 主键ID
	Title          string    `json:"title" gorm:"type:varchar(200);not null;comment:标题"`
	Content        string    `json:"content" gorm:"type:longtext;not null;column:content;comment:文章内容"`
	Excerpt        string    `json:"excerpt" gorm:"type:varchar(300);not null;default:''"`
	SeoTitle       string    `json:"seoTitle" gorm:"type:varchar(200);not null;default:''"`
	SeoDescription string    `json:"seoDescription" gorm:"type:varchar(300);not null;default:''"`
	PublishAt      int64     `json:"publishAt" gorm:"not null;default:0;index:idx_article_schedule,priority:2"`
	Version        uint64    `json:"version" gorm:"not null;default:1"`
	Tags           []Tag     `json:"tags" gorm:"-"`
	TagNames       []string  `json:"tagNames" gorm:"-"`
	Summary        string    `json:"summary" gorm:"-"` // 列表摘要，不写数据库
	ContentFormat  string    `json:"contentFormat" gorm:"type:varchar(16);column:content_format;default:html;comment:内容格式"`
	Pic            string    `json:"pic" gorm:"pic;comment:图片"`
	Uid            int       `json:"uid" gorm:"uid;comment:用户id"`
	Fid            int       `json:"fid" gorm:"fid;comment:分类id"`
	View           int       `json:"view" gorm:"view;comment:浏览次数"`
	Ctime          int       `json:"ctime" gorm:"type:bigint;not null;column:ctime;comment:创建时间"`
	EditTime       int       `json:"editTime" gorm:"type:bigint;not null;column:edittime;comment:编辑时间"`
	File           string    `json:"file" gorm:"pic;file:附件"`
	Type           int       `json:"type" gorm:"column:type;comment:文章类型"`
	IsTop          int       `json:"isTop" gorm:"isTop;column:istop;comment:置顶"`
	Status         int       `json:"status" gorm:"status;column:status;index:idx_article_schedule,priority:1;comment:状态"`
	Comment        []Comment `json:"comment" gorm:"foreignKey:Aid;references:ID;comment:评论"`
	CommentCount   int       `json:"commentCount" gorm:"-"` // 评论数量，不存储在数据库
	LikeCount      int       `json:"likeCount" gorm:"-"`    // 点赞数量，不存储在数据库
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
	TagID   uint   `json:"tagId" form:"tagId"`
	State   string `json:"state" form:"state"`
	Scope   string `json:"scope" form:"scope"`
	Fid     int    `json:"fid" form:"fid"`
	Hide    bool   `json:"hide" form:"hide"`
	Top     bool   `json:"top" form:"top"`
	Keyword string `json:"keyword" form:"keyword"`
	request.PageInfo
}

// ArchiveEntry 归档页使用的文章条目，不包含正文
type ArchiveEntry struct {
	ID    uint   `json:"id"`
	Fid   int    `json:"fid"`
	Title string `json:"title"`
	Ctime int    `json:"ctime"`
}

// ArticleBrief 文章简要信息，用于详情页侧栏展示
type ArticleBrief struct {
	ID    uint   `json:"ID"`
	Fid   int    `json:"fid"`
	Title string `json:"title"`
	Pic   string `json:"pic"`
	View  int    `json:"view"`
	Ctime int    `json:"ctime"`
}

// PostAside 文章详情页侧栏数据
type PostAside struct {
	Related []ArticleBrief `json:"related"`
	Hot     []ArticleBrief `json:"hot"`
}

// PostAsideRequest 侧栏数据查询参数
type PostAsideRequest struct {
	Fid     int  `json:"fid" form:"fid"`
	Exclude uint `json:"exclude" form:"exclude"`
	Limit   int  `json:"limit" form:"limit"`
}

func (Article) TableName() string {
	return "blog_article"
}

func (ArticleLike) TableName() string {
	return "blog_article_like"
}
