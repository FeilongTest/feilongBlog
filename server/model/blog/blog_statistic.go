package blog

type BlogStatistic struct {
	Date string `json:"date" gorm:"primaryKey;size:10"`
	Pv   int64  `json:"pv"`
	Uv   int64  `json:"uv"`
}

type BlogVisitor struct {
	ID        uint   `gorm:"primarykey"`
	Date      string `json:"date" gorm:"size:10;uniqueIndex:idx_date_visitor"`
	VisitorID string `json:"visitorId" gorm:"size:64;uniqueIndex:idx_date_visitor"`
}

type VisitRequest struct {
	VisitorID string `json:"visitorId"`
}

func (BlogStatistic) TableName() string {
	return "blog_statistic"
}

func (BlogVisitor) TableName() string {
	return "blog_visitor"
}
