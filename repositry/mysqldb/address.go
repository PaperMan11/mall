package mysqldb

import (
	"mall/model"

	"gorm.io/gorm"
)

type AddressModel interface {
	GetAddressByAid(addreseId uint) (address *model.Address, err error)
	ListAddressByUid(userId uint) (addresses []*model.Address, err error)
	CreateAddress(address *model.Address) error
	DeleteAddressById(addressId uint) error
	UpdateAddressById(aId uint, address *model.Address) (err error)
	ExistAddressByAid(addreseId uint) (exist bool, err error)
}

type defaultAddressModel struct {
	db    *gorm.DB
	table string
	// cache...
}

func NewAddressModel(db *gorm.DB, table string) AddressModel {
	return &defaultAddressModel{
		db:    db,
		table: table,
	}
}

func (m *defaultAddressModel) ExistAddressByAid(addreseId uint) (exist bool, err error) {
	var count int64
	err = m.db.Model(&model.Address{}).Where("id=?", addreseId).Count(&count).Error
	return count != 0, err
}

// GetAddressByAid 根据 AddressId 获取 model.Address
func (m *defaultAddressModel) GetAddressByAid(addreseId uint) (address *model.Address, err error) {
	err = m.db.Where("id=?", addreseId).First(&address).Error
	return
}

// ListAddressByUid 根据 User Id 获取 model.Address
func (m *defaultAddressModel) ListAddressByUid(userId uint) (addresses []*model.Address, err error) {
	err = m.db.Model(&model.Address{}).Where("user_id=?", userId).
		Order("created_at desc").
		Find(&addresses).Error
	return
}

// CreateAddress 创建地址
func (m *defaultAddressModel) CreateAddress(address *model.Address) error {
	return m.db.Create(address).Error
}

// DeleteAddressById 根据 id 删除地址
func (m *defaultAddressModel) DeleteAddressById(addressId uint) error {
	return m.db.Where("id=?", addressId).Delete(&model.Address{}).Error
}

// UpdateAddressById 通过 id 修改地址信息
func (m *defaultAddressModel) UpdateAddressById(aId uint, address *model.Address) (err error) {
	// 只更新非零字段
	err = m.db.Model(&model.Address{}).Where("id=?", aId).Updates(address).Error
	return
}
