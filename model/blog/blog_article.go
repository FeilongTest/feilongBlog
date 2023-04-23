package blog

type Article struct {
	//global.GVA_MODEL
	ID      uint   `gorm:"primarykey"` // 主键ID
	Title   string `json:"title" gorm:"title;comment:标题"`
	Content string `json:"content" gorm:"type:text;column:content;comment:文章内容"`
	Pic     string `json:"pic" gorm:"pic;comment:图片"`
	//Uid     string`json:"uid" gorm:"uid;comment:标题"`
	//Fid     string`json:"fid" gorm:"fid;comment:标题"`
	//View    int`json:"content" gorm:"content;comment:标题"`
}

func (Article) TableName() string {
	return "blog_article"
}
