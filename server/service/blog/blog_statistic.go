package blog

import (
	"feilongBlog/global"
	model "feilongBlog/model/blog"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StatisticService struct{}

// RecordVisit 记录网站访问量
func (s *StatisticService) RecordVisit(visitorID string) (err error) {
	date := time.Now().Format("2006-01-02")
	return global.BLOG_DB.Transaction(func(tx *gorm.DB) error {
		visitor := model.BlogVisitor{Date: date, VisitorID: visitorID}
		result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&visitor)
		if result.Error != nil {
			return result.Error
		}
		uv := int64(0)
		if result.RowsAffected > 0 {
			uv = 1
		}
		statistic := model.BlogStatistic{Date: date, Pv: 1, Uv: uv}
		return tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "date"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"pv": gorm.Expr("pv + ?", 1),
				"uv": gorm.Expr("uv + ?", uv),
			}),
		}).Create(&statistic).Error
	})
}

// GetDashboard 获取后台统计数据
func (s *StatisticService) GetDashboard() (result map[string]interface{}, err error) {
	var articleCount, categoryCount, commentCount, friendlinkCount, likeCount, totalPv, totalUv int64
	queries := []struct {
		model interface{}
		count *int64
	}{
		{&model.Article{}, &articleCount}, {&model.Category{}, &categoryCount},
		{&model.Comment{}, &commentCount}, {&model.Friendlink{}, &friendlinkCount},
		{&model.ArticleLike{}, &likeCount},
	}
	for _, query := range queries {
		if err = global.BLOG_DB.Model(query.model).Count(query.count).Error; err != nil {
			return
		}
	}
	if err = global.BLOG_DB.Model(&model.BlogStatistic{}).Select("COALESCE(SUM(pv), 0)").Scan(&totalPv).Error; err != nil {
		return
	}
	if err = global.BLOG_DB.Model(&model.BlogVisitor{}).Distinct("visitor_id").Count(&totalUv).Error; err != nil {
		return
	}

	start := time.Now().AddDate(0, 0, -6).Format("2006-01-02")
	var rows []model.BlogStatistic
	if err = global.BLOG_DB.Where("date >= ?", start).Order("date asc").Find(&rows).Error; err != nil {
		return
	}
	rowMap := make(map[string]model.BlogStatistic)
	for _, row := range rows {
		rowMap[row.Date] = row
	}
	trend := make([]model.BlogStatistic, 0, 7)
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		row := rowMap[date]
		row.Date = date
		trend = append(trend, row)
	}
	result = map[string]interface{}{
		"articleCount": articleCount, "categoryCount": categoryCount,
		"commentCount": commentCount, "friendlinkCount": friendlinkCount,
		"likeCount": likeCount, "totalPv": totalPv, "totalUv": totalUv, "trend": trend,
	}
	return
}
