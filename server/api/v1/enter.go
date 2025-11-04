package v1

import "feilongBlog/api/v1/blog"

type ApiGroup struct {
	UserApiGroup       blog.UserApi
	ArticleApiGroup    blog.ArticleApi
	CategoryApiGroup   blog.CategoryApi
	FileApiGroup       blog.FileApi
	CommentApiGroup    blog.CommentApi
	FriendlinkApiGroup blog.FriendlinkApi
}

var ApiGroupApp = new(ApiGroup)
