package v1

import (
	"mall/model"
	"mall/pkg/e"
	"mall/service"
	"mall/svc"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary 	秒杀商品导入
// @Description	秒杀商品导入
// @Tags		skill
// @Accept		mpfd
// @Produce		json
// @Param		file formData file true "file"
// @Param		access_token header string true "access_token"
// @Param		refresh_token header string true "refresh_token"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/skill-goods [post]
func SkillProductImportHandler(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		f, _, err := ctx.Request.FormFile("file")
		if err != nil {
			svc.Log.Errorf("ctx.Request.FormFile err: %s", err)
			model.RespError(ctx, nil, e.InvalidParams)
			return
		}
		l := service.NewSkillProductLogic(ctx, svc)
		code := l.SkillProductImportLogic(f)
		model.RespWithCode(ctx, nil, code)
	}
}

// @Summary 	开启秒杀
// @Description	开启秒杀
// @Tags		skill
// @Produce		json
// @Param		access_token header string true "access_token"
// @Param		refresh_token header string true "refresh_token"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/skill-goods [get]
func SkillStartHandler(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		l := service.NewSkillProductLogic(ctx, svc)
		code := l.SkillStartLogic()
		model.RespWithCode(ctx, nil, code)
	}
}

// @Summary 	秒杀
// @Description	秒杀
// @Tags		skill
// @Accept		json
// @Produce		json
// @Param		product_id path uint true "product id"
// @Param 		req body model.SkillReq true "秒杀请求信息"
// @Param		access_token header string true "access_token"
// @Param		refresh_token header string true "refresh_token"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/skill-goods/{product_id} [post]
func SkillHandler(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s := ctx.Param("product_id")
		pId, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			svc.Log.Errorf("[%s] InvalidParams", pId)
			model.RespError(ctx, e.InvalidParams)
			return
		}

		var req model.SkillReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			svc.Log.Errorf("ShouldBind err: %s", err)
			model.RespError(ctx, nil, e.InvalidParams)
			return
		}
		l := service.NewSkillProductLogic(ctx, svc)
		code := l.SkillLogic(uint(pId), &req)
		model.RespWithCode(ctx, nil, code)
	}
}
