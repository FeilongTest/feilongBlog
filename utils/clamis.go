package utils

import (
	"feilongBlog/global"
	model "feilongBlog/model/blog"
	"github.com/gin-gonic/gin"
)

func GetClaims(c *gin.Context) (*model.CustomClaims, error) {
	token := c.Request.Header.Get("x-token")
	j := NewJWT()
	claims, err := j.ParseToken(token)
	if err != nil {
		global.BLOG_LOG.Error("从Gin的Context中获取从jwt解析信息失败, 请检查请求头是否存在x-token且claims是否为规定结构")
	}
	return claims, err
}

// GetUserID 从Gin的Context中获取从jwt解析出来的用户ID
func GetUserID(c *gin.Context) uint {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return 0
		} else {
			return cl.ID
		}
	} else {
		waitUse := claims.(*model.CustomClaims)
		return waitUse.ID
	}
}

// GetUserInfo 从Gin的Context中获取从jwt解析出来的用户角色id
func GetUserInfo(c *gin.Context) *model.CustomClaims {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return nil
		} else {
			return cl
		}
	} else {
		waitUse := claims.(*model.CustomClaims)
		return waitUse
	}
}

// IsAdmin 是否为管理员
func IsAdmin(c *gin.Context) bool {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return false
		} else {
			return cl.Admin == 1
		}
	} else {
		waitUse := claims.(*model.CustomClaims)
		return waitUse.Admin == 1
	}
}
