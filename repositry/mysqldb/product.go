package mysqldb

import (
	"mall/model"
	"mall/pkg/utils"

	"gorm.io/gorm"
)

type ProductModel interface {
	GetProductById(pid uint) (p *model.Product, e error)
	GetProductInfoById(pid uint) (p *model.Product, e error)
	ListProductByCondition(condition map[string]interface{}, page *utils.BasePage) (products []*model.Product, err error)
	CreateProduct(p *model.Product) error
	CountProductByCondition(condition map[string]interface{}) (count int64, err error)
	DeleteProduct(pId uint) error
	UpdateProduct(pid uint, p *model.Product) error
	SearchProduct(info string, page *utils.BasePage) (products []*model.Product, count int, err error)
	SubProductNum(productId uint, num int) error
}

type defaultProductModel struct {
	db    *gorm.DB
	table string
	// cache...
}

func NewProductModel(db *gorm.DB, table string) ProductModel {
	return &defaultProductModel{
		db:    db,
		table: table,
	}
}

// GetProductById 通过 id 获取product
func (m *defaultProductModel) GetProductById(pid uint) (p *model.Product, err error) {
	err = m.db.Model(&model.Product{}).Where("id=?", pid).First(&p).Error
	return
}

func (m *defaultProductModel) GetProductInfoById(pid uint) (p *model.Product, err error) {
	err = m.db.Preload("Category").Preload("ProductImgs").Where("id=?", pid).First(&p).Error
	return
}

// ListProductByCondition 根据condition获取商品列表
func (m *defaultProductModel) ListProductByCondition(condition map[string]interface{}, page *utils.BasePage) (products []*model.Product, err error) {
	err = m.db.Preload("Category").Preload("ProductImgs").Where(condition).
		Offset((page.PageNum - 1) * page.PageSize).
		Limit(page.PageSize).Find(&products).Error
	return
}

// CreateProduct 创建商品，创建后会返回主键值
func (m *defaultProductModel) CreateProduct(p *model.Product) error {
	return m.db.Create(p).Error
}

// CountProductByCondition 根据condition获取商品的数量
func (m *defaultProductModel) CountProductByCondition(condition map[string]interface{}) (count int64, err error) {
	err = m.db.Model(&model.Product{}).Where(condition).Count(&count).Error
	return
}

// DeleteProduct 删除商品
func (m *defaultProductModel) DeleteProduct(pId uint) error {
	return m.db.Where("id = ?", pId).Delete(&model.Product{}).Error
}

// UpdateProduct 更新商品
func (m *defaultProductModel) UpdateProduct(pid uint, p *model.Product) error {
	return m.db.Model(&model.Product{}).Where("id=?", pid).Updates(p).Error
}

// SearchProduct 搜索商品
func (m *defaultProductModel) SearchProduct(info string, page *utils.BasePage) (products []*model.Product, count int, err error) {
	searchInfo := "%" + info + "%"
	err = m.db.Preload("Category").Preload("ProductImgs").
		Where("name LIKE ? OR info LIKE ?", searchInfo, searchInfo).
		Offset((page.PageNum - 1) * page.PageSize).Limit(page.PageSize).
		Find(&products).Error
	return products, len(products), err
}

// SubProductNum 扣库存（num为负数）
func (m *defaultProductModel) SubProductNum(productId uint, num int) error {
	return m.db.Model(&model.Product{}).Where("id=?", productId).UpdateColumn("num", gorm.Expr("num-?", num)).Error
}
