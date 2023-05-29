package mysqldb

import (
	"mall/model"

	"gorm.io/gorm"
)

type OrderModel interface {
	Create(order *model.Order) error
	DeleteByOrderId(oId int64) error
	ExistByCondition(condition map[string]interface{}) (exist bool, err error)
	GetByCondition(condition map[string]interface{}) (order *model.Order, err error)
	ListByUserId(userId uint) (orders []*model.Order, err error)
	UpdateById(oId int64, order *model.Order) error
}

type defaultOrderModel struct {
	db    *gorm.DB
	table string
	// cache...
}

func NewOrderModel(db *gorm.DB, table string) OrderModel {
	return &defaultOrderModel{
		db:    db,
		table: table,
	}
}

func (m *defaultOrderModel) Create(order *model.Order) error {
	return m.db.Create(&order).Error
}

func (m *defaultOrderModel) DeleteByOrderId(oId int64) error {
	return m.db.Where("order_id=?", oId).Delete(&model.Order{}).Error
}

func (m *defaultOrderModel) ExistByCondition(condition map[string]interface{}) (exist bool, err error) {
	var count int64
	err = m.db.Model(&model.Order{}).Where(condition).Count(&count).Error
	if err != nil || count == 0 {
		return false, err
	}
	return true, nil
}

func (m *defaultOrderModel) GetByCondition(condition map[string]interface{}) (order *model.Order, err error) {
	err = m.db.Where(condition).First(&order).Error
	return
}

func (m *defaultOrderModel) ListByUserId(userId uint) (orders []*model.Order, err error) {
	err = m.db.Where("user_id=?", userId).Find(&orders).Error
	return
}

func (m *defaultOrderModel) UpdateById(oId int64, order *model.Order) error {
	return m.db.Model(&model.Order{}).Where("order_id=?", oId).Updates(&order).Error
}
