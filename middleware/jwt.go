package middleware

import (
	"mall/pkg/e"
	"mall/pkg/utils/jwt"
	"mall/svc"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	AccessTokenHeader    = "access_token"
	RefreshTokenHeader   = "refresh_token"
	HeaderForwardedProto = "X-Forwarded-Proto"
	MaxAge               = 3600 * 24

	CtxUserIdKey = "userId"
)

// AuthMiddleware token验证中间件
func AuthMiddleware(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		code = e.SUCCESS
		accessToken := c.GetHeader("access_token")
		refreshToken := c.GetHeader("refresh_token")
		if accessToken == "" {
			code = e.InvalidParams
			c.JSON(http.StatusBadRequest, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
				"data":   "Token不能为空",
			})
			svc.Log.Errorf("Token empty")
			c.Abort()
			return
		}
		nAtoken, nRtoken, err := jwt.ParseRefreshToken(accessToken, refreshToken)
		if err != nil {
			svc.Log.Errorf("jwt.ParseRefreshToken err: %s", err)
			code = e.ErrorAuthCheckTokenFail
		}
		if code != e.SUCCESS {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
				"data":   "鉴权失败",
			})
			c.Abort()
			return
		}
		claims, err := jwt.ParseToken(accessToken)
		if err != nil {
			code = e.ErrorAuthCheckTokenFail
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
				"data":   err.Error(),
			})
			svc.Log.Errorf("jwt.ParseToken err: %s", err)
			c.Abort()
			return
		}
		SetToken(c, nAtoken, nRtoken)
		c.Set(CtxUserIdKey, claims.UserId)
		// c.Request = c.Request.WithContext(ctl.NewContext(c.Request.Context(), &ctl.UserInfo{Id: claims.ID}))
		// ctl.InitUserInfo(c.Request.Context())
		c.Next()
	}
}

func SetToken(c *gin.Context, accessToken, refreshToken string) {
	secure := IsHttps(c)
	c.Header(AccessTokenHeader, accessToken)
	c.Header(RefreshTokenHeader, refreshToken)
	c.SetCookie(AccessTokenHeader, accessToken, MaxAge, "/", "", secure, true)
	c.SetCookie(RefreshTokenHeader, refreshToken, MaxAge, "/", "", secure, true)
}

// 判断是否https
func IsHttps(c *gin.Context) bool {
	if c.GetHeader(HeaderForwardedProto) == "https" || c.Request.TLS != nil {
		return true
	}
	return false
}
