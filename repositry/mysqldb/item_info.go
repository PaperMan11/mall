package mysqldb

import (
	"mall/model"

	"gorm.io/gorm"
)

type ItemInfoModel interface {
	ExistByCartID(cartId, productId uint) (item *model.ItemInfo, exist bool, err error)
	ExistByID(id uint) (item *model.ItemInfo, exist bool, err error)
	UpdateByID(id uint, itemInfo *model.ItemInfo) error
	Create(itemInfo *model.ItemInfo) error
	Delete(id uint) error
	ShowById(id uint) (item *model.ItemInfo, err error)
}

type defaultItemInfoModel struct {
	db    *gorm.DB
	table string
	// cache...
}

func NewItemInfoModel(db *gorm.DB, table string) ItemInfoModel {
	return &defaultItemInfoModel{
		db:    db,
		table: table,
	}
}

func (m *defaultItemInfoModel) ExistByCartID(cartId, productId uint) (item *model.ItemInfo, exist bool, err error) {
	var count int64
	err = m.db.Model(&model.ItemInfo{}).Where("cart_id=? AND product_id=?", cartId, productId).Count(&count).Error
	if err != nil || count == 0 {
		return nil, false, err
	}

	err = m.db.Model(&model.ItemInfo{}).Where("cart_id=? AND product_id=?", cartId, productId).First(&item).Error
	if err != nil {
		return nil, false, err
	}
	return item, true, nil
}

func (m *defaultItemInfoModel) ExistByID(id uint) (item *model.ItemInfo, exist bool, err error) {
	var count int64
	err = m.db.Model(&model.ItemInfo{}).Where("id=?", id).Count(&count).Error
	if err != nil || count == 0 {
		return nil, false, err
	}

	err = m.db.Model(&model.ItemInfo{}).Where("id=?", id).First(&item).Error
	if err != nil {
		return nil, false, err
	}

	return item, true, nil
}

func (m *defaultItemInfoModel) UpdateByID(id uint, itemInfo *model.ItemInfo) error {
	return m.db.Model(model.ItemInfo{}).Where("id=?", id).Updates(&itemInfo).Error
}

func (m *defaultItemInfoModel) Create(itemInfo *model.ItemInfo) error {
	return m.db.Create(itemInfo).Error
}

func (m *defaultItemInfoModel) Delete(id uint) error {
	return m.db.Delete(&model.ItemInfo{}, id).Error
}

func (m *defaultItemInfoModel) ShowById(id uint) (item *model.ItemInfo, err error) {
	err = m.db.Where("id=?", id).First(&item).Error
	return
}
