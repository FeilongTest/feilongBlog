package blog

type Category struct {
	ID   uint   `gorm:"primarykey"` // 主键ID
	Name string `gorm:"name" json:"name"`
	Fid  int    `gorm:"fid" json:"fid"`
	Type int    `gorm:"type" json:"type"`
	Sort int    `gorm:"sort" json:"sort"`
}

func (Category) TableName() string {
	return "blog_category"
}
