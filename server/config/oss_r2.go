package config

type R2 struct {
	AccountID       string `mapstructure:"account-id" json:"account-id" yaml:"account-id"`
	AccessKeyID     string `mapstructure:"access-key-id" json:"-" yaml:"access-key-id"`
	SecretAccessKey string `mapstructure:"secret-access-key" json:"-" yaml:"secret-access-key"`
	Bucket          string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	PublicURL       string `mapstructure:"public-url" json:"public-url" yaml:"public-url"`
	MaxSizeMB       int64  `mapstructure:"max-size-mb" json:"max-size-mb" yaml:"max-size-mb"`
}
