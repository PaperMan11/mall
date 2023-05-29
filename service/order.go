package service

import (
	"fmt"
	"mall/middleware"
	"mall/model"
	"mall/pkg/e"
	"mall/repositry/mysqldb"
	"mall/svc"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderLogic struct {
	ctx    *gin.Context
	svcCtx *svc.ServiceContext
}

func NewOrderLogic(ctx *gin.Context, svcCtx *svc.ServiceContext) *OrderLogic {
	return &OrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderLogic) OrderCreateLogic(userId uint, req *model.OrderCreateReq) (order *model.Order, code int) {

	// 1 商品是否存在
	count, err := l.svcCtx.ProductModel.CountProductByCondition(map[string]interface{}{
		"id": req.ProductID,
	})
	if err != nil {
		l.svcCtx.Log.Errorf("ProductModel.CountProductByCondition err: %s", err)
		return nil, e.ErrorDatabase
	}
	if count == 0 {
		l.svcCtx.Log.Infof("product [%d] not exist", req.ProductID)
		return nil, e.ErrorNotExistProduct
	}
	// 2 库存
	p, err := l.svcCtx.ProductModel.GetProductById(req.ProductID)
	if err != nil {
		l.svcCtx.Log.Errorf("ProductModel.GetProductById err: %s", err)
		return nil, e.ErrorDatabase
	}
	if p.Num < req.ProductNum {
		l.svcCtx.Log.Infof("Product [%d] quantity is not enough", err)
		return nil, e.ErrorNotEnough
	}
	// 3 生成订单
	m, err := strconv.ParseFloat(p.Price, 64)
	if err != nil {
		l.svcCtx.Log.Errorf("%s convert to float err: %s", m, err)
		return nil, e.ERROR
	}
	money := m * float64(req.ProductNum)
	req.Money = money

	oId, err := l.svcCtx.SnowFlake.NextID()
	if err != nil {
		l.svcCtx.Log.Errorf("gen OrderId err: %s", err)
		return nil, e.ERROR
	}
	o := model.Order{
		Model: model.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		OrderId:    oId,
		UserId:     userId,
		BossId:     p.BossId,
		ProductId:  req.ProductID,
		ProductNum: req.ProductNum,
		Money:      req.Money,
		Address:    req.Address,
	}
	if err = l.svcCtx.OrderModel.Create(&o); err != nil {
		l.svcCtx.Log.Errorf("OrderModel.Create err: %s", err)
		return nil, e.ErrorDatabase
	}
	return &o, e.SUCCESS
}

func (l *OrderLogic) OrderDeleteLogic(oId int64) (code int) {
	exist, err := l.svcCtx.OrderModel.ExistByCondition(map[string]interface{}{
		"order_id": oId,
	})
	if err != nil {
		l.svcCtx.Log.Errorf("ProductModel.ExistByOrderId err: %s", err)
		return e.ErrorDatabase
	}
	if !exist {
		return e.ErrorNotExistOrder
	}

	if err = l.svcCtx.OrderModel.DeleteByOrderId(oId); err != nil {
		l.svcCtx.Log.Errorf("ProductModel.DeleteByOrderId err: %s", err)
		return e.ErrorDatabase
	}

	return e.SUCCESS
}

func (l *OrderLogic) OrderGetLogic(oId int64) (order *model.Order, code int) {
	exist, err := l.svcCtx.OrderModel.ExistByCondition(map[string]interface{}{
		"order_id": oId,
	})
	if err != nil {
		l.svcCtx.Log.Errorf("ProductModel.ExistByOrderId err: %s", err)
		return nil, e.ErrorDatabase
	}
	if !exist {
		return nil, e.ErrorNotExistOrder
	}

	o, err := l.svcCtx.OrderModel.GetByCondition(map[string]interface{}{
		"order_id": oId,
	})
	if err != nil {
		l.svcCtx.Log.Errorf("OrderModel.GetByOrderId err: %s", err)
		return nil, e.ErrorDatabase
	}

	return o, e.SUCCESS
}

func (l *OrderLogic) OrderListLogic() (order []*model.Order, code int) {
	userIdV, _ := l.ctx.Get(middleware.CtxUserIdKey)
	userId, ok := userIdV.(uint)
	if !ok {
		l.svcCtx.Log.Errorf("key [%s] get value faild", middleware.CtxUserIdKey)
		return nil, e.ERROR
	}

	exist, err := l.svcCtx.OrderModel.ExistByCondition(map[string]interface{}{
		"user_id": userId,
	})
	if err != nil {
		l.svcCtx.Log.Errorf("OrderModel.ExistByUserId err: %s", err)
		return nil, e.ErrorDatabase
	}
	if !exist {
		return nil, e.ErrorNotExistOrder
	}

	orders, err := l.svcCtx.OrderModel.ListByUserId(userId)
	if err != nil {
		l.svcCtx.Log.Errorf("OrderModel.ListByUserId err: %s", err)
		return nil, e.ErrorDatabase
	}

	return orders, e.SUCCESS
}

func (l *OrderLogic) OrderPaymentLogic(oId int64) (code int) {
	// 1 userID
	userIdV, _ := l.ctx.Get(middleware.CtxUserIdKey)
	userId, ok := userIdV.(uint)
	if !ok {
		l.svcCtx.Log.Errorf("key [%s] get value faild", middleware.CtxUserIdKey)
		return e.ERROR
	}

	// 2 order exist?
	exist, err := l.svcCtx.OrderModel.ExistByCondition(map[string]interface{}{
		"order_id": oId,
		"user_id":  userId,
	})
	if err != nil {
		l.svcCtx.Log.Errorf("OrderModel.ExistByUserId err: %s", err)
		return e.ErrorDatabase
	}
	if !exist {
		l.svcCtx.Log.Infof("order [%d] not exist", oId)
		return e.ErrorNotExistOrder
	}

	order, err := l.svcCtx.OrderModel.GetByCondition(map[string]interface{}{
		"order_id": oId,
		"user_id":  userId,
	})
	if err != nil {
		l.svcCtx.Log.Errorf("OrderModel.GetByOrderId err: %s", err)
		return e.ErrorDatabase
	}
	// 3 payment
	if err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		orderDAO := mysqldb.NewOrderModel(tx, "order")
		userDAO := mysqldb.NewUserModel(tx, "user")
		productDAO := mysqldb.NewProductModel(tx, "product")

		user, err := userDAO.GetUserById(userId)
		if err != nil {
			l.svcCtx.Log.Errorf("Transaction GetUserById err: %s", err)
			return err
		}
		boss, err := userDAO.GetUserById(order.BossId)
		if err != nil {
			l.svcCtx.Log.Errorf("Transaction GetBossById err: %s", err)
			return err
		}
		// 扣库存
		err = productDAO.SubProductNum(order.ProductId, order.ProductNum)
		if err != nil {
			l.svcCtx.Log.Errorf("Transaction SubProductNum err: %s", err)
			return err
		}
		// - user money
		moneyStr1 := l.svcCtx.Encrypt.AesDecoding(user.Money)
		userMoney, err := strconv.ParseFloat(moneyStr1, 64)
		if err != nil {
			l.svcCtx.Log.Errorf("Transaction Decode User [%d] Money [%s] err: %s", userId, moneyStr1, err)
			return err
		}
		userMoney -= order.Money
		user.Money = l.svcCtx.Encrypt.AesEncoding(fmt.Sprintf("%f", userMoney))
		if err = l.svcCtx.UserModel.UpdateUserById(userId, user); err != nil {
			l.svcCtx.Log.Errorf("Transaction UpdateUserById err: %s", err)
			return err
		}
		// + boss money
		moneyStr2 := l.svcCtx.Encrypt.AesDecoding(boss.Money)
		bossMoney, err := strconv.ParseFloat(moneyStr2, 64)
		if err != nil {
			l.svcCtx.Log.Errorf("Transaction Decode Boss [%d] Money [%s] err: %s", order.BossId, moneyStr2, err)
			return err
		}
		bossMoney += order.Money
		boss.Money = l.svcCtx.Encrypt.AesEncoding(fmt.Sprintf("%f", bossMoney))
		if err = l.svcCtx.UserModel.UpdateUserById(order.BossId, boss); err != nil {
			l.svcCtx.Log.Errorf("Transaction UpdateBossById err: %s", err)
			return err
		}
		// 更新订单状态
		order.State = true
		if err = orderDAO.UpdateById(oId, order); err != nil {
			l.svcCtx.Log.Errorf("Transaction UpdateById err: %s", err)
			return err
		}
		return nil
	}); err != nil {
		l.svcCtx.Log.Errorf("Transaction Payment err:%s", err)
		return e.ErrorPayment
	}

	return e.SUCCESS
}
