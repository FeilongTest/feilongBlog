package middleware

import (
	"feilongBlog/global"
	model "feilongBlog/model/blog"
	"feilongBlog/model/common/response"
	"feilongBlog/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localStorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := c.Request.Header.Get("x-token")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{Code: response.ERROR, Data: gin.H{"reload": true}, Msg: "未登录或非法访问"})
			c.Abort()
			return
		}
		j := utils.NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == utils.TokenExpired {
				c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{Code: response.ERROR, Data: gin.H{"reload": true}, Msg: "授权已过期"})
				c.Abort()
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{Code: response.ERROR, Data: gin.H{"reload": true}, Msg: err.Error()})
			c.Abort()
			return
		}

		// 已登录用户被管理员禁用 需要使该用户的jwt失效 此处比较消耗性能 如果需要 请自行打开
		// 用户被删除的逻辑 需要优化 此处比较消耗性能 如果需要 请自行打开

		//if user, err := userService.FindUserByUuid(claims.UUID.String()); err != nil || user.Enable == 2 {
		//	_ = jwtService.JsonInBlacklist(system.JwtBlacklist{Jwt: token})
		//	c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{Code: response.ERROR, Data: gin.H{"reload": true}, Msg: err.Error()})
		//	c.Abort()
		//}
		if claims.ExpiresAt-time.Now().Unix() < claims.BufferTime {
			dr, _ := utils.ParseDuration(global.BLOG_CONFIG.JWT.ExpiresTime)
			claims.ExpiresAt = time.Now().Add(dr).Unix()
			newToken, err := j.CreateTokenByOldToken(token, *claims)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{Code: response.ERROR, Data: gin.H{"reload": true}, Msg: "刷新授权失败"})
				c.Abort()
				return
			}
			newClaims, err := j.ParseToken(newToken)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{Code: response.ERROR, Data: gin.H{"reload": true}, Msg: "刷新授权失败"})
				c.Abort()
				return
			}
			c.Header("new-token", newToken)
			c.Header("new-expires-at", strconv.FormatInt(newClaims.ExpiresAt, 10))
		}
		c.Set("claims", claims)
		c.Next()
	}
}

// AdminOnly 在所有后台接口执行前校验已验证的 JWT 管理员身份。
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		user, ok := claims.(*model.CustomClaims)
		if !exists || !ok || user == nil || user.Admin != 1 {
			c.AbortWithStatusJSON(http.StatusForbidden, response.Response{Code: response.ERROR, Msg: "仅管理员可以访问"})
			return
		}
		c.Next()
	}
}
