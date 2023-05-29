package v1

import (
	"mall/model"
	"mall/pkg/e"
	"mall/service"
	"mall/svc"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary 	购物车信息
// @Description	购物车信息
// @Tags		cart
// @Produce		json
// @Param 		cart_id path uint true "cart id"
// @Param		access_token header string true "access_token"
// @Param		refresh_token header string true "refresh_token"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/carts/{cart_id} [get]
func CartInfoHandler(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		idStr := ctx.Param("cart_id")
		itemInfoID, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			svc.Log.Errorf("[%s] InvalidParams", idStr)
			model.RespError(ctx, e.InvalidParams)
			return
		}

		l := service.NewCartLogic(ctx, svc)
		resp, code := l.CartInfoLogic(uint(itemInfoID))
		model.RespWithCode(ctx, resp, code)
	}
}
