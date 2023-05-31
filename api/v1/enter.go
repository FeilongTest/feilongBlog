package v1

import "feilongBlog/api/v1/blog"

type ApiGroup struct {
	UserApiGroup     blog.UserApi
	ArticleApiGroup  blog.ArticleApi
	CategoryApiGroup blog.CategoryApi
	FileApiGroup     blog.FileApi
}

var ApiGroupApp = new(ApiGroup)
