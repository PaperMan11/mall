package model

type SkillProduct struct {
	Model
	ProductId uint `gorm:"unique"`
	BossId    uint `gorm:"not null"`
	Title     string
	Money     float64
	Num       int
}

// -------------------------------------------------------------------

type SkillReq struct {
	AddressId uint `json:"address_id"`
}

type SkillReq2MQ struct {
	ProductId  uint
	CustomerId uint
	AddressId  uint
	Money      float64
}
