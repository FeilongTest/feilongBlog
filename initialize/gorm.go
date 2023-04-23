package initialize

import (
	"feilongBlog/model/blog"
	"os"

	"feilongBlog/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Gorm 初始化数据库并产生数据库全局变量
// Author SliverHorn
func Gorm() *gorm.DB {
	return GormMysql()
}

// RegisterTables 注册数据库表专用
// Author SliverHorn
func RegisterTables(db *gorm.DB) {
	err := db.AutoMigrate(
		// 系统模块表
		blog.User{},
	)
	if err != nil {
		global.BLOG_LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}
	global.BLOG_LOG.Info("register table success")
}
