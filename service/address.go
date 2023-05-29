package service

import (
	"mall/middleware"
	"mall/model"
	"mall/pkg/e"
	"mall/pkg/utils"
	"mall/repositry/mysqldb"
	"time"

	"mall/svc"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AddressLogic struct {
	ctx    *gin.Context
	svcCtx *svc.ServiceContext
}

func NewAddressLogic(ctx *gin.Context, svcCtx *svc.ServiceContext) *AddressLogic {
	return &AddressLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddressLogic) AddressCreateLogic(req *model.AddressCreateReq) (resp *model.Address, code int) {
	userIdV, _ := l.ctx.Get(middleware.CtxUserIdKey)
	userId, ok := userIdV.(uint)
	if !ok {
		l.svcCtx.Log.Errorf("ctx key [%s] get value faild", middleware.CtxUserIdKey)
		return nil, e.ERROR
	}
	address := model.Address{
		Model: model.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserID:  userId,
		Name:    req.Name,
		Phone:   req.Phone,
		Address: req.Address,
	}
	err := l.svcCtx.AddressModel.CreateAddress(&address)
	if err != nil {
		l.svcCtx.Log.Errorf("AddressModel.CreateAddress err: %s", err)
		return nil, e.ErrorDatabase
	}

	return &address, e.SUCCESS
}

func (l *AddressLogic) AddressGetLogic(id uint) (resp *model.Address, code int) {

	// exist?
	exist, err := l.svcCtx.AddressModel.ExistAddressByAid(id)
	if err != nil {
		l.svcCtx.Log.Errorf("AddressModel.ExistAddressByAid err: %s", err)
		return nil, e.ErrorDatabase
	}
	if !exist {
		l.svcCtx.Log.Infof("[%d] address exist", id)
		return nil, e.ErrorNotExistAddress
	}

	// get
	a, err := l.svcCtx.AddressModel.GetAddressByAid(id)
	if err != nil {
		l.svcCtx.Log.Errorf("AddressModel.GetAddressByAid err: %s", err)
		return nil, e.ErrorDatabase
	}
	return a, e.SUCCESS
}

func (l *AddressLogic) AddressListLogic(req *utils.BasePage) (resp *model.AddressList, code int) {
	userIdV, _ := l.ctx.Get(middleware.CtxUserIdKey)
	userId, ok := userIdV.(uint)
	if !ok {
		l.svcCtx.Log.Errorf("ctx key [%s] get value faild", middleware.CtxUserIdKey)
		return nil, e.ERROR
	}
	addresses, err := l.svcCtx.AddressModel.ListAddressByUid(userId)
	if err != nil {
		l.svcCtx.Log.Errorf("AddressModel.ListAddressByUid err: %s", err)
		return nil, e.ErrorDatabase
	}

	return &model.AddressList{
		AddressInfos: addresses,
		Total:        len(addresses),
	}, e.SUCCESS
}

func (l *AddressLogic) AddressDeleteLogic(id uint) (code int) {
	err := l.svcCtx.AddressModel.DeleteAddressById(id)
	if err != nil {
		l.svcCtx.Log.Errorf("AddressModel.DeleteAddressById err: %s", err)
		return e.ErrorDatabase
	}
	return e.SUCCESS
}

func (l *AddressLogic) AddressUpdateLogic(req *model.AddressUpdateReq, id uint) (nAddress *model.Address, code int) {

	// search
	exist, err := l.svcCtx.AddressModel.ExistAddressByAid(id)
	if err != nil {
		l.svcCtx.Log.Errorf("AddressModel.ExistAddressByAid err: %s", err)
		return nil, e.ErrorDatabase
	}
	if !exist {
		l.svcCtx.Log.Infof("[%d] address not exist", id)
		return nil, e.ErrorNotExistAddress
	}

	// update select
	address := model.Address{
		Model: model.Model{
			UpdatedAt: time.Now(),
		},
		Name:    req.Name,
		Phone:   req.Phone,
		Address: req.Address,
	}

	if err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		addressDAO := mysqldb.NewAddressModel(tx, "address")
		err := addressDAO.UpdateAddressById(id, &address)
		if err != nil {
			l.svcCtx.Log.Errorf("Transaction UpdateAddressById err: %s", err)
			return err
		}

		nAddress, err = addressDAO.GetAddressByAid(id)
		if err != nil {
			l.svcCtx.Log.Errorf("Transaction GetAddressByAid err: %s", err)
			return err
		}
		return nil
	}); err != nil {
		l.svcCtx.Log.Errorf("Transaction AddressUpdateSelect err:%s", err)
		return nil, e.ErrorDatabase
	}
	return nAddress, e.SUCCESS
}
