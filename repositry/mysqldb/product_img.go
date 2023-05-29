package mysqldb

import (
	"mall/model"

	"gorm.io/gorm"
)

type ProductImgModel interface {
	CreateProductImg(pimg *model.ProductImg) error
	ListProductImgByProductId(productId uint) (productImg []*model.ProductImg, err error)
	DeleteProductImgByPid(produtId uint) error
}

type defaultProductImgModel struct {
	db    *gorm.DB
	table string
	// cache...
}

func NewProductImgModel(db *gorm.DB, table string) ProductImgModel {
	return &defaultProductImgModel{
		db:    db,
		table: table,
	}
}

// CreateProductImg 创建商品图片
func (m *defaultProductImgModel) CreateProductImg(pimg *model.ProductImg) error {
	return m.db.Create(pimg).Error
}

func (m *defaultProductImgModel) DeleteProductImgByPid(produtId uint) error {
	return m.db.Where("product_id=?", produtId).Delete(&model.ProductImg{}).Error
}

// ListProductImgByProductId 根据商品id获取商品图片
func (m *defaultProductImgModel) ListProductImgByProductId(productId uint) (productImg []*model.ProductImg, err error) {
	err = m.db.Model(&model.ProductImg{}).Where("product_id=?", productId).Find(&productImg).Error
	return
}
