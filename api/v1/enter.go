package v1

import "feilongBlog/api/v1/blog"

type ApiGroup struct {
	BaseApiGroup     blog.BaseApi
	BlogApiGroup     blog.UserApi
	ArticleApiGroup  blog.ArticleApi
	CategoryApiGroup blog.CategoryApi
}

var ApiGroupApp = new(ApiGroup)
