package v1

import "feilongBlog/api/v1/blog"

type ApiGroup struct {
	UserApiGroup       blog.UserApi
	ArticleApiGroup    blog.ArticleApi
	CategoryApiGroup   blog.CategoryApi
	FileApiGroup       blog.FileApi
	CommentApiGroup    blog.CommentApi
	FriendlinkApiGroup blog.FriendlinkApi
	StatisticApiGroup  blog.StatisticApi
	SeoApiGroup        blog.SeoApi
	SettingApiGroup    blog.SettingApi
}

var ApiGroupApp = new(ApiGroup)
