package config

import "errors"

type Server struct {
	JWT    JWT    `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Zap    Zap    `mapstructure:"zap" json:"zap" yaml:"zap"`
	System System `mapstructure:"system" json:"system" yaml:"system"`
	Mysql  Mysql  `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Local  Local  `mapstructure:"local" json:"local" yaml:"local"`
	R2     R2     `mapstructure:"r2" json:"r2" yaml:"r2"`
	// 跨域配置
	Cors CORS `mapstructure:"cors" json:"cors" yaml:"cors"`
}

// Validate 检查服务启动所需配置
func (s Server) Validate() error {
	if len(s.JWT.SigningKey) < 32 {
		return errors.New("BLOG_JWT_SIGNING_KEY长度不能小于32位")
	}
	if s.System.OssType == "r2" {
		if s.R2.AccountID == "" || s.R2.AccessKeyID == "" || s.R2.SecretAccessKey == "" || s.R2.Bucket == "" || s.R2.PublicURL == "" {
			return errors.New("已启用R2，但BLOG_R2_*配置不完整")
		}
	}
	return nil
}
