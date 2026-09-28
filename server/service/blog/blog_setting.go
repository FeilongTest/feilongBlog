package blog

import (
	"feilongBlog/global"
	model "feilongBlog/model/blog"
	"strings"
	"time"

	"gorm.io/gorm"
)

type SettingService struct {
}

// 站点设置固定使用主键为 1 的单行记录
const settingRowID = 1

// defaultSetting 站点设置的默认值，数据库中没有记录时使用
func defaultSetting() model.Setting {
	return model.Setting{
		ID:          settingRowID,
		SiteName:    "飞龙博客",
		SiteSlogan:  "记录技术与生活",
		Description: "一个分享编程技术、折腾经验与日常记录的个人博客。",
		Keywords:    "个人博客,技术博客,编程,Go,Vue",
		Author:      "飞龙",
	}
}

// GetSetting 读取站点设置，数据库无记录时回退到默认值
func (s *SettingService) GetSetting() (setting model.Setting, err error) {
	err = global.BLOG_DB.Where("id = ?", settingRowID).First(&setting).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return defaultSetting(), nil
		}
		return setting, err
	}
	return s.fillDefault(setting), nil
}

// GetPublicSetting 获取前台可以公开访问的站点信息
func (s *SettingService) GetPublicSetting() (setting model.PublicSetting, err error) {
	row, err := s.GetSetting()
	if err != nil {
		return setting, err
	}
	return model.PublicSetting{
		SiteName:    row.SiteName,
		SiteSlogan:  row.SiteSlogan,
		Description: row.Description,
		Keywords:    row.Keywords,
		Author:      row.Author,
		Icp:         row.Icp,
	}, nil
}

// UpdateSetting 更新站点设置，记录不存在时自动创建
func (s *SettingService) UpdateSetting(input model.UpdateSetting) (setting model.Setting, err error) {
	current, err := s.GetSetting()
	if err != nil {
		return setting, err
	}
	current.SiteName = strings.TrimSpace(input.SiteName)
	current.SiteSlogan = strings.TrimSpace(input.SiteSlogan)
	current.Description = strings.TrimSpace(input.Description)
	current.Keywords = strings.TrimSpace(input.Keywords)
	current.Author = strings.TrimSpace(input.Author)
	current.Icp = strings.TrimSpace(input.Icp)
	current.BaiduVerify = strings.TrimSpace(input.BaiduVerify)
	current.GoogleVerify = strings.TrimSpace(input.GoogleVerify)
	current.BingVerify = strings.TrimSpace(input.BingVerify)
	current.EditTime = int(time.Now().Unix())
	if current.Ctime == 0 {
		current.Ctime = current.EditTime
	}
	current = s.fillDefault(current)
	if err = global.BLOG_DB.Save(&current).Error; err != nil {
		return setting, err
	}
	return current, nil
}

// fillDefault 为空字段补齐默认值，避免前台出现空白标题
func (s *SettingService) fillDefault(setting model.Setting) model.Setting {
	fallback := defaultSetting()
	if setting.SiteName == "" {
		setting.SiteName = fallback.SiteName
	}
	if setting.SiteSlogan == "" {
		setting.SiteSlogan = fallback.SiteSlogan
	}
	if setting.Description == "" {
		setting.Description = fallback.Description
	}
	if setting.Keywords == "" {
		setting.Keywords = fallback.Keywords
	}
	if setting.Author == "" {
		setting.Author = fallback.Author
	}
	return setting
}
