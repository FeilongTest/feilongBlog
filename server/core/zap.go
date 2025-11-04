package core

import (
	"feilongBlog/core/internal"
	"feilongBlog/global"
	"feilongBlog/utils"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// Zap 获取 zap.Logger
// Author [SliverHorn](https://github.com/SliverHorn)
func Zap() (logger *zap.Logger) {
	if ok, _ := utils.PathExists(global.BLOG_CONFIG.Zap.Director); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", global.BLOG_CONFIG.Zap.Director)
		_ = os.Mkdir(global.BLOG_CONFIG.Zap.Director, os.ModePerm)
	}

	cores := internal.Zap.GetZapCores()
	logger = zap.New(zapcore.NewTee(cores...))

	if global.BLOG_CONFIG.Zap.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger
}
