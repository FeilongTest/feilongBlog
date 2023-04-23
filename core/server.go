package core

import (
	"fmt"
	"time"

	"feilongBlog/global"
	"feilongBlog/initialize"
	"go.uber.org/zap"
)

type server interface {
	ListenAndServe() error
}

func RunWindowsServer() {
	Router := initialize.Routers()
	address := fmt.Sprintf(":%d", global.BLOG_CONFIG.System.Addr)
	s := initServer(address, Router)
	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)
	global.BLOG_LOG.Info("server run success on ", zap.String("address", address))

	fmt.Printf(`
	欢迎使用 feilongBlog
	当前版本:v1.0.0
	默认前端文件运行地址:http://127.0.0.1:8080
`)
	global.BLOG_LOG.Error(s.ListenAndServe().Error())
}
