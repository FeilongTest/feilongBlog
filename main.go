package main

import (
	"feilongBlog/core"
	"feilongBlog/global"
	"feilongBlog/initialize"
	"go.uber.org/zap"
)

func main() {
	//初始化
	global.BLOG_VP = core.Viper() // 初始化Viper
	global.BLOG_LOG = core.Zap()  // 初始化zap日志库
	zap.ReplaceGlobals(global.BLOG_LOG)
	global.BLOG_DB = initialize.Gorm() // gorm连接数据库

	if global.BLOG_DB != nil {
		initialize.RegisterTables(global.BLOG_DB) // 初始化表
		// 程序结束前关闭数据库链接
		db, _ := global.BLOG_DB.DB()
		defer db.Close()
	}
	core.RunWindowsServer()
}
