package blog

type User struct {
	ID       uint   `gorm:"primarykey"` // 主键ID
	Username string `json:"userName" gorm:"index;comment:用户登录名"`
	Password string `json:"-"  gorm:"comment:用户登录密码"`
	Pic      string `json:"pic" gorm:"pic;comment:用户头像"`
	Email    string `json:"email" gorm:"email;comment:邮箱"`
	Ctime    int    `json:"ctime" gorm:"ctime;comment:创建时间"`
	LastTime int    `json:"lasttime" gorm:"column:lasttime;comment:最后一次登录时间"`
	Ip       string `json:"ip" gorm:"ip;comment:登录ip"`
	Status   int    `json:"status" gorm:"status;comment:状态"`
	TrueName string `json:"truename" gorm:"column:truename;comment:真实用户名"`
	Admin    int    `json:"admin" gorm:"admin;comment:是否管理员"`
}

func (User) TableName() string {
	return "blog_user"
}

type Login struct {
	Username  string `json:"username"`  // 用户名
	Password  string `json:"password"`  // 密码
	Captcha   string `json:"captcha"`   // 验证码
	CaptchaId string `json:"captchaId"` // 验证码ID
}

type LoginResponse struct {
	User      User   `json:"user"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
}
