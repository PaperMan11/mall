package model

type ItemInfo struct {
	Model
	CartId    uint `json:"cart_id"`
	ProductId uint `gorm:"not null" json:"product_id"`
	// Product   *Product `gorm:"foreignKey:ProductId" json:"product,omitempty"` // 不设置外键
	Num int `json:"num"`
}

// -------------------------------------------------------------------

type ItemInfoAddReq struct {
	Num int `json:"num" binding:"required,gte=1"`
}

type ItemInfoUpdateReq struct {
	Num int `json:"num" binding:"required"`
}
