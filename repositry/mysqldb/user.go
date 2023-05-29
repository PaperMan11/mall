package mysqldb

import (
	"mall/model"

	"gorm.io/gorm"
)

type UserModel interface {
	GetUserById(uid uint) (u *model.User, e error)
	UpdateUserById(userId uint, u *model.User) error
	ExistByUserName(userName string) (user *model.User, exist bool, err error)
	CreateUser(user *model.User) error
}

type defaultUserModel struct {
	db    *gorm.DB
	table string
	// cache...
}

func NewUserModel(db *gorm.DB, table string) UserModel {
	return &defaultUserModel{
		db:    db,
		table: table,
	}
}

// GetUserById 根据 id 获取用户
func (m *defaultUserModel) GetUserById(uid uint) (u *model.User, e error) {
	e = m.db.Model(&model.User{}).Where("id=?", uid).First(&u).Error
	return
}

// UpdateUserById 根据 id 更新用户信息
func (m *defaultUserModel) UpdateUserById(userId uint, u *model.User) error {
	return m.db.Model(&model.User{}).Where("id=?", userId).Updates(&u).Error
}

// ExistByUserName 根据username判断是否存在该名字
func (m *defaultUserModel) ExistByUserName(userName string) (user *model.User, exist bool, err error) {
	var count int64
	err = m.db.Model(&model.User{}).Where("user_name=?", userName).Count(&count).Error
	if count == 0 {
		return nil, false, err
	}
	if err = m.db.Model(&model.User{}).Where("user_name=?", userName).First(&user).Error; err != nil {
		return nil, false, err
	}
	return user, true, nil
}

// CreateUser 创建用户
func (m *defaultUserModel) CreateUser(user *model.User) error {
	return m.db.Create(user).Error
}
