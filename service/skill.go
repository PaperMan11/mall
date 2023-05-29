package service

import (
	"mall/middleware"
	"mall/model"
	"mall/pkg/e"
	"mall/svc"
	"mime/multipart"
	"strconv"

	"mall/pkg/utils/uuid"

	xlsx "github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
)

type SkillProductLogic struct {
	ctx    *gin.Context
	svcCtx *svc.ServiceContext
}

func NewSkillProductLogic(ctx *gin.Context, svcCtx *svc.ServiceContext) *SkillProductLogic {
	return &SkillProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SkillProductLogic) SkillProductImportLogic(file multipart.File) (code int) {
	xlFile, err := xlsx.OpenReader(file)
	if err != nil {
		l.svcCtx.Log.Errorf("xlsx.OpenReader err: %s", err)
		return e.ERROR
	}

	// 获取一个 sheet1 上的所有单元格
	rows := xlFile.GetRows("Sheet1")
	length := len(rows[1:])
	skillGoods := make([]*model.SkillProduct, length)
	for index, colCell := range rows {
		if index == 0 {
			continue
		}
		pId, _ := strconv.Atoi(colCell[0])
		bId, _ := strconv.Atoi(colCell[1])
		num, _ := strconv.Atoi(colCell[3])
		money, _ := strconv.ParseFloat(colCell[4], 64)
		skillGood := &model.SkillProduct{
			ProductId: uint(pId),
			BossId:    uint(bId),
			Title:     colCell[2],
			Money:     money,
			Num:       num,
		}
		skillGoods[index-1] = skillGood
	}
	err = l.svcCtx.SkillModel.CreateByList(skillGoods)
	if err != nil {
		l.svcCtx.Log.Errorf("SkillModel.CreateByList err: %s", err)
		return e.ErrorDatabase
	}
	return e.SUCCESS
}

func (l *SkillProductLogic) SkillStartLogic() (code int) {
	// 加载到内存
	skillGoods, err := l.svcCtx.SkillModel.ListSkillGoods()
	if err != nil {
		l.svcCtx.Log.Errorf("SkillModel.ListSkillGoods err: %s", err)
		return e.ErrorDatabase
	}
	for _, p := range skillGoods {
		_, err = l.svcCtx.SkRedisCache.HMSetProductNumMoney(p.ProductId, p.Num, p.Money)
		if err != nil {
			l.svcCtx.Log.Infof("good [%d] load err: %s", p.ProductId, err)
		}
	}
	return e.SUCCESS
}

func (l *SkillProductLogic) SkillLogic(pId uint, req *model.SkillReq) (code int) {
	userIdV, _ := l.ctx.Get(middleware.CtxUserIdKey)
	userId, ok := userIdV.(uint)
	if !ok {
		l.svcCtx.Log.Errorf("key [%s] get value faild", middleware.CtxUserIdKey)
		return e.ERROR
	}

	// 1 get num cache
	num, money, err := l.svcCtx.SkRedisCache.HMGetProductNumMoney(pId)
	if err != nil {
		l.svcCtx.Log.Errorf("SkRedisCache.HMGetProductNumMoney err: %s", err)
		return e.ErrorDatabase
	}
	if num == 0 {
		l.svcCtx.Log.Infof("skill product [%d] sold out", pId)
		return e.ErrorSoldOut
	}

	skillReq2MQ := model.SkillReq2MQ{
		ProductId:  pId,
		CustomerId: userId,
		AddressId:  req.AddressId,
		Money:      money,
	}
	// 2 获取分布式锁，失败直接返回
	pUUID := uuid.UUID()
	ok, err = l.svcCtx.SkRedisCache.Lock(pId, pUUID)
	if err != nil {
		l.svcCtx.Log.Errorf("SkRedisCache.Lock err: %s", err)
		return e.ErrorDatabase
	}
	if !ok {
		l.svcCtx.Log.Infof("user [%d] get lock failed", userId)
		return e.ERROR
	}
	// 3 send to rabbit mq
	err = l.svcCtx.SkAmqp.SendToQueue("sk-goods", &skillReq2MQ)
	if err != nil {
		l.svcCtx.Log.Errorf("SkAmqp.SendToQueue err: %s", err)
		return e.ErrorDatabase
	}
	// 4 扣库存
	_, err = l.svcCtx.SkRedisCache.HDecrByNum(pId)
	if err != nil {
		l.svcCtx.Log.Errorf("SkRedisCache.HDecrByNum err: %s", err)
		return e.ErrorDatabase
	}
	// 5 unlock
	_, err = l.svcCtx.SkRedisCache.UnLock(pId)
	if err != nil {
		l.svcCtx.Log.Errorf("SkRedisCache.UnLock err: %s", err)
		return e.ErrorDatabase
	}

	return e.SUCCESS
}
