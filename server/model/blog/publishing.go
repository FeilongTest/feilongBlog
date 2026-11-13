package blog

const (
	ArticlePublished = 0
	ArticleHidden    = 1
	ArticleDraft     = 2
	ArticleScheduled = 3
)

// Tag 为公开内容提供跨分类的主题组织；文章关联单独维护，避免隐式联级写入。
type Tag struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"type:varchar(64);not null;uniqueIndex"`
	Description string `json:"description" gorm:"type:varchar(300);not null;default:''"`
	Count       int64  `json:"count" gorm:"column:count;->;-:migration"`
}

func (Tag) TableName() string { return "blog_tag" }

type ArticleTag struct {
	ArticleID uint `json:"articleId" gorm:"primaryKey;autoIncrement:false"`
	TagID     uint `json:"tagId" gorm:"primaryKey;autoIncrement:false;index"`
}

func (ArticleTag) TableName() string { return "blog_article_tag" }

// ArticleRevision 保存可恢复的内容快照，访问仅限管理员。
type ArticleRevision struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	ArticleID uint   `json:"articleId" gorm:"not null;uniqueIndex:idx_article_revision"`
	Version   uint64 `json:"version" gorm:"not null;uniqueIndex:idx_article_revision"`
	Title     string `json:"title" gorm:"type:varchar(200);not null"`
	SavedAt   int64  `json:"savedAt" gorm:"not null"`
	Reason    string `json:"reason" gorm:"type:varchar(32);not null"`
	Snapshot  string `json:"-" gorm:"type:longtext;not null"`
}

func (ArticleRevision) TableName() string { return "blog_article_revision" }

// ArticleSnapshot 不包含阅读数、评论、作者凭据，载入历史版本不会修改发布状态。
type ArticleSnapshot struct {
	Title          string   `json:"title"`
	Content        string   `json:"content"`
	ContentFormat  string   `json:"contentFormat"`
	Pic            string   `json:"pic"`
	Fid            int      `json:"fid"`
	Type           int      `json:"type"`
	Excerpt        string   `json:"excerpt"`
	SeoTitle       string   `json:"seoTitle"`
	SeoDescription string   `json:"seoDescription"`
	TagNames       []string `json:"tagNames"`
}

type ArticleStatePatch struct {
	ID      uint   `json:"id" binding:"required"`
	Version uint64 `json:"version" binding:"required"`
	Status  *int   `json:"status"`
	IsTop   *int   `json:"isTop"`
}
