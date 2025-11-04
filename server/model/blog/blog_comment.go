package blog

type Comment struct {
	ID      uint   `gorm:"primarykey"` // 主键ID
	Name    string `json:"name" form:"name"`
	Email   string `json:"email" form:"email"`
	Content string `json:"content" form:"content"`
	Uid     int    `json:"uid" form:"uid"`
	Ctime   int    `json:"ctime" form:"ctime" gorm:"ctime;comment:创建时间"`
	Aid     int    `json:"aid" form:"aid" gorm:"aid;comment:文章id"`
	Status  int    `json:"status" form:"status" gorm:"status;column:status;comment:状态"`
}

func (Comment) TableName() string {
	return "blog_comment"
}
