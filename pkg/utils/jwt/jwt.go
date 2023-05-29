package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	// AccessTokenExpireDuration  = 24 * time.Hour
	// RefreshTokenExpireDuration = 10 * 24 * time.Hour
	// 调试
	AccessTokenExpireDuration  = 360 * time.Hour
	RefreshTokenExpireDuration = 360 * 24 * time.Hour
)

var jwtSecret = []byte("ZiQing")

type Claims struct {
	UserId   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateToken 签发用户Token
func GenerateToken(id uint, username string) (accessToken, refreshToken string, err error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(AccessTokenExpireDuration)
	rtExpireTime := nowTime.Add(RefreshTokenExpireDuration)
	claims := Claims{
		UserId:   id,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "mall",
		},
	}
	// 加密并获得完整的编码后的字符串token
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: rtExpireTime.Unix(),
		Issuer:    "mall",
	}).SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, err
}

// ParseToken 验证用户token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

// ParseRefreshToken 验证用户token
/*
在access_token里加入refresh_token标识，给access_token设置短时间的期限（例如一天），给refresh_token设置长时间的期限（例如七天）。当活动用户（拥有access_token）发起request时，在权限验证里，对于requeset的header包含的access_token、refresh_token分别进行验证：

1、access_token没过期，即通过权限验证；

2、access_token过期,refresh_token没过期，则返回权限验证失败，并在返回的response的header中加入标识状态的key，在request方法的catch中通过webException来获取标识的key，获取新的token（包含新的access_token和refresh_token），再次发起请求，并返回给客户端请求结果以及新的token，再在客户端更新公共静态token模型；

3、access_token过期,refresh_token过期即权限验证失败。
*/
func ParseRefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	if _, err = jwt.Parse(rToken, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, err
	}); err != nil {
		return
	}

	var claims Claims
	_, err = jwt.ParseWithClaims(aToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, err
	})

	v, ok := err.(*jwt.ValidationError)

	if ok && v.Errors == jwt.ValidationErrorExpired {
		return GenerateToken(claims.UserId, claims.Username)
	}

	return
}

// EmailClaims
type EmailClaims struct {
	UserID        uint   `json:"user_id"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	OperationType uint   `json:"operation_type"`
	jwt.StandardClaims
}

// GenerateEmailToken 签发邮箱验证Token
func GenerateEmailToken(userID, Operation uint, email, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(15 * time.Minute)
	claims := EmailClaims{
		UserID:        userID,
		Email:         email,
		Password:      password,
		OperationType: Operation,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "cmall",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// ParseEmailToken 验证邮箱验证token
func ParseEmailToken(token string) (*EmailClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &EmailClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*EmailClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
