package blog

// Setting 站点设置，全站仅保存一行记录
type Setting struct {
	ID           uint   `gorm:"primarykey"` // 主键ID
	SiteName     string `json:"siteName" gorm:"type:varchar(60);column:site_name;comment:站点名称"`
	SiteSlogan   string `json:"siteSlogan" gorm:"type:varchar(120);column:site_slogan;comment:站点副标题"`
	Description  string `json:"description" gorm:"type:varchar(200);column:description;comment:站点描述"`
	Keywords     string `json:"keywords" gorm:"type:varchar(200);column:keywords;comment:站点关键词"`
	Author       string `json:"author" gorm:"type:varchar(60);column:author;comment:作者"`
	Icp          string `json:"icp" gorm:"type:varchar(80);column:icp;comment:备案信息"`
	BaiduVerify  string `json:"baiduVerify" gorm:"type:varchar(120);column:baidu_verify;comment:百度站长平台验证码"`
	GoogleVerify string `json:"googleVerify" gorm:"type:varchar(120);column:google_verify;comment:Google搜索控制台验证码"`
	BingVerify   string `json:"bingVerify" gorm:"type:varchar(120);column:bing_verify;comment:Bing网站管理员验证码"`
	Ctime        int    `json:"ctime" gorm:"ctime;comment:创建时间"`
	EditTime     int    `json:"editTime" gorm:"column:edittime;comment:编辑时间"`
}

func (Setting) TableName() string {
	return "blog_setting"
}

// UpdateSetting 更新站点设置的入参
type UpdateSetting struct {
	SiteName     string `json:"siteName"`
	SiteSlogan   string `json:"siteSlogan"`
	Description  string `json:"description"`
	Keywords     string `json:"keywords"`
	Author       string `json:"author"`
	Icp          string `json:"icp"`
	BaiduVerify  string `json:"baiduVerify"`
	GoogleVerify string `json:"googleVerify"`
	BingVerify   string `json:"bingVerify"`
}

// PublicSetting 前台公开的站点信息
type PublicSetting struct {
	SiteName    string `json:"siteName"`
	SiteSlogan  string `json:"siteSlogan"`
	Description string `json:"description"`
	Keywords    string `json:"keywords"`
	Author      string `json:"author"`
	Icp         string `json:"icp"`
}
