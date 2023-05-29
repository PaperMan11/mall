package mysqldb

import (
	"mall/model"

	"gorm.io/gorm"
)

type SkillModel interface {
	CreateByList(in []*model.SkillProduct) error
	ListSkillGoods() (skillGoods []*model.SkillProduct, err error)
}

type defaultSkillModel struct {
	db    *gorm.DB
	table string
	// cache...
}

func NewSkillModel(db *gorm.DB, table string) SkillModel {
	return &defaultSkillModel{
		db:    db,
		table: table,
	}
}

func (m *defaultSkillModel) CreateByList(in []*model.SkillProduct) error {
	return m.db.Model(&model.SkillProduct{}).Create(&in).Error
}

func (m *defaultSkillModel) ListSkillGoods() (skillGoods []*model.SkillProduct, err error) {
	err = m.db.Where("num>0").Find(&skillGoods).Error
	return
}
