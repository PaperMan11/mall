package model

type Address struct {
	Model
	UserID  uint   `json:"user_id"`
	Name    string `gorm:"type:varchar(20) not null" json:"name"`
	Phone   string `gorm:"type:varchar(11) not null" json:"phone"`
	Address string `gorm:"type:varchar(50) not null" json:"address"`
}

// -------------------------------------------------------------------

type AddressCreateReq struct {
	Name    string `form:"name" json:"name"`
	Phone   string `form:"phone" json:"phone"`
	Address string `form:"address" json:"address"`
}

type AddressList struct {
	AddressInfos []*Address `json:"address_infos"`
	Total        int        `json:"total"`
}

type AddressUpdateReq struct {
	Name    string `form:"name" json:"name"`
	Phone   string `form:"phone" json:"phone"`
	Address string `form:"address" json:"address"`
}
