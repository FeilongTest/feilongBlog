package global

import (
	"feilongBlog/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

var (
	BLOG_DB                  *gorm.DB
	BLOG_CONFIG              config.Server
	BLOG_VP                  *viper.Viper
	BLOG_LOG                 *zap.Logger
	BLOG_Concurrency_Control = &singleflight.Group{}
)
