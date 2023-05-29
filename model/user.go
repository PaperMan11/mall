package model

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	PassWordCost        = 12       // 密码加密难度
	Active       string = "active" // 激活用户
)

type User struct {
	Model
	UserName       string    `gorm:"unique" json:"user_name"`
	Email          string    `json:"email"`
	PasswordDigest string    `json:"-"`
	NickName       string    `json:"nick_name"`
	Status         string    `json:"status"`
	Avatar         string    `gorm:"size:1000" json:"avatar"`
	Money          string    `json:"money"`
	Address        *Address  `gorm:"foreignKey:UserID" json:"address,omitempty"`
	Product        []Product `gorm:"foreignKey:BossId" json:"product,omitempty"`
	Cart           *Cart     `gorm:"foreignKey:UserID" json:"cart,omitempty"`
	Orders         []Order   `gorm:"foreignKey:UserId" json:"orders,omitempty"`
}

func (u *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PassWordCost)
	if err != nil {
		return err
	}
	u.PasswordDigest = string(bytes)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordDigest), []byte(password))
	return err == nil
}

// -------------------------------------------------------------------
type UserRegisterReq struct {
	NickName string `form:"nick_name" json:"nick_name"`
	UserName string `form:"user_name" json:"user_name"`
	Password string `form:"password" json:"password"`
	Email    string `form:"email" json:"email"`
}

type UserLoginReq struct {
	UserName string `form:"user_name" json:"user_name"`
	Password string `form:"password" json:"password"`
}

type UserLoginResp struct {
	UserInfo  *UserInfo  `json:"user_info"`
	TokenInfo *TokenInfo `json:"token_info"`
}

type UserInfo struct {
	Id       uint   `json:"id"`
	UserName string `json:"user_name"`
	NickName string `json:"nick_name"`
	Type     int    `json:"type"`
	Email    string `json:"email"`
	Status   string `json:"status"`
	Avatar   string `json:"avatar"`
	CreateAt int64  `json:"create_at"`
}

type TokenInfo struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
