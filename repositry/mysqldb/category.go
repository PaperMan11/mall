package mysqldb

import (
	"mall/model"
	"mall/pkg/utils"

	"gorm.io/gorm"
)

type CategoryModel interface {
	ListCategory(page *utils.BasePage) (categories []*model.Category, err error)
	AddCategory(category *model.Category) error
	RemoveCategoryById(id uint) error
	ExistCategoryByName(cname string) (exist bool, err error)
	GetCatoryById(id uint) (category *model.Category, err error)
}

type defaultCategoryModel struct {
	db    *gorm.DB
	table string
	// cache...
}

func NewCategoryModel(db *gorm.DB, table string) CategoryModel {
	return &defaultCategoryModel{
		db:    db,
		table: table,
	}
}

func (m *defaultCategoryModel) ListCategory(page *utils.BasePage) (categories []*model.Category, err error) {
	err = m.db.Model(&model.Category{}).Find(&categories).Offset((page.PageNum - 1) * page.PageSize).Limit(page.PageSize).Error
	return
}

func (m *defaultCategoryModel) AddCategory(category *model.Category) error {
	return m.db.Create(category).Error
}

func (m *defaultCategoryModel) RemoveCategoryById(id uint) error {
	return m.db.Where("id=?", id).Delete(&model.Category{}).Error
}

func (m *defaultCategoryModel) ExistCategoryByName(cname string) (exist bool, err error) {
	var count int64
	err = m.db.Model(&model.Category{}).Where("category_name=?", cname).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

func (m *defaultCategoryModel) GetCatoryById(id uint) (category *model.Category, err error) {
	err = m.db.Model(&model.Category{}).Where("id=?", id).First(&category).Error
	return
}
