package mysqldb

import (
	"mall/model"

	"gorm.io/gorm"
)

type CartModel interface {
	Create(cart *model.Cart) error
	ShowCartById(id uint) (cart *model.Cart, err error)
}

type defaultCartModel struct {
	db    *gorm.DB
	table string
	// cache...
}

func NewCartModel(db *gorm.DB, table string) CartModel {
	return &defaultCartModel{
		db:    db,
		table: table,
	}
}

func (m *defaultCartModel) Create(cart *model.Cart) error {
	return m.db.Create(&cart).Error
}

func (m *defaultCartModel) ShowCartById(id uint) (cart *model.Cart, err error) {
	err = m.db.Preload("ItemInfo").Where("id=?", id).First(&cart).Error
	return
}
