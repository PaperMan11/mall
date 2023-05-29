package model

// 单个商品订单
type Order struct {
	Model
	OrderId   int64 `gorm:"unique"` // 订单号
	UserId    uint
	BossId    uint
	ProductId uint
	// Product    *Product `gorm:"foreignKey:ProductId" json:"product,omitempty"` // 不设置外键
	ProductNum int
	Money      float64
	Address    string
	State      bool // 支付状态 true 完成
}

// -------------------------------------------------------------------

type OrderCreateReq struct {
	ProductID  uint    `json:"product_id" binding:"required"`
	ProductNum int     `json:"product_num" binding:"gte=1"`
	Address    string  `json:"address" binding:"required"`
	Money      float64 `json:"-"` // 非必须
}
