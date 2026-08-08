package core

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"feilongBlog/global"
	"feilongBlog/initialize"
	"go.uber.org/zap"
)

type server interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

func RunServer() {
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
	errCh := make(chan error, 1)
	go func() { errCh <- s.ListenAndServe() }()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	select {
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			global.BLOG_LOG.Error("服务启动失败", zap.Error(err))
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.Shutdown(shutdownCtx); err != nil {
			global.BLOG_LOG.Error("服务退出失败", zap.Error(err))
		}
	}
}
