package v1

import (
	"mall/middleware"
	"mall/model"
	"mall/pkg/e"
	"mall/service"
	"mall/svc"

	"github.com/gin-gonic/gin"
)

/*
注解	描述
@Summary	摘要
@Produce	API 可以产生的 MIME 类型的列表，MIME 类型你可以简单的理解为响应类型，例如：json、xml、html 等等
@Param	参数格式，从左到右分别为：参数名、入参类型、数据类型、是否必填、注释
@Success	响应成功，从左到右分别为：状态码、参数类型、数据类型、注释
@Failure	响应失败，从左到右分别为：状态码、参数类型、数据类型、注释
@Router	路由，从左到右分别为：路由地址，HTTP 方法
*/

// @Summary 	用户注册
// @Description	用户注册
// @Tags		user
// @Accept		json
// @Produce		json
// @Param 		req body model.UserRegisterReq true "用户注册信息"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/user/register [post]
func UserRegisterHandler(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req model.UserRegisterReq
		if err := ctx.ShouldBind(&req); err != nil {
			svc.Log.Errorf("ShouldBind err: %s", err)
			model.RespError(ctx, nil, e.InvalidParams)
			return
		}
		l := service.NewUserLogic(ctx, svc)
		code := l.UserRegisterLogic(&req)
		model.RespWithCode(ctx, nil, code)
	}
}

// @Summary 	用户登录
// @Description	用户登录
// @Tags		user
// @Accept		json
// @Produce		json
// @Param 		req body model.UserLoginReq true "用户登录信息"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/user/login [post]
func UserLoginHandler(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req model.UserLoginReq
		if err := ctx.ShouldBind(&req); err != nil {
			svc.Log.Errorf("ShouldBind err: %s", err)
			model.RespError(ctx, nil, e.InvalidParams)
			return
		}
		l := service.NewUserLogic(ctx, svc)
		resp, code := l.UserLoginLogic(&req)
		model.RespWithCode(ctx, resp, code)
	}
}

// @Summary 	用户登录
// @Description	用户登录
// @Tags		user
// @Produce		json
// @Param		access_token header string true "access_token"
// @Param		refresh_token header string true "refresh_token"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/user [get]
func UserShowHandler(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIdV, _ := ctx.Get(middleware.CtxUserIdKey)
		userId, ok := userIdV.(uint)
		if !ok {
			svc.Log.Errorf("key [%s] get value faild", middleware.CtxUserIdKey)
			model.RespError(ctx, nil, e.ErrorAdminFindUser)
			return
		}
		l := service.NewUserLogic(ctx, svc)
		resp, code := l.UserShowLogic(userId)
		model.RespWithCode(ctx, resp, code)
	}
}
