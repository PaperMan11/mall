package skillconsumer

import (
	"context"
	"encoding/json"
	"mall/model"
	"mall/service"
	"mall/svc"
)

func SkillReqHandleTask(ctx context.Context, svc *svc.ServiceContext) error {
	svc.Log.Infof("SkillReqHandleTask Start...")
	ch, reqCh, err := svc.SkAmqp.ReceiveFromQueue("sk-goods")
	if err != nil {
		svc.Log.Errorf("SkAmqp.ReceiveFromQueue err: %s", err)
		return err
	}
	defer ch.Close()

	var req model.SkillReq2MQ
	for {
		select {
		case data := <-reqCh:
			err := json.Unmarshal(data.Body, &req)
			if err != nil {
				svc.Log.Errorf("SkillReqHandleTask: json.Unmarshal err: %s", err)
				return err
			}

			// get address by id
			a, err := svc.AddressModel.GetAddressByAid(req.AddressId)
			if err != nil {
				svc.Log.Errorf("SkillReqHandleTask: AddressModel.GetAddressByAid err: %s", err)
				return err
			}

			// 生成订单
			order := service.NewOrderLogic(nil, svc)
			_, _ = order.OrderCreateLogic(req.CustomerId, &model.OrderCreateReq{
				ProductID:  req.ProductId,
				ProductNum: 1,
				Address:    a.Address,
				Money:      req.Money,
			})
		case <-ctx.Done():
			svc.Log.Infof("SkillReqHandleTask End...")
			return nil
		}
	}
}
