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

	if global.BLOG_DB == nil {
		panic("数据库初始化失败")
	}
	initialize.RegisterTables(global.BLOG_DB) // 初始化表

	db, err := global.BLOG_DB.DB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	core.RunServer()
}
