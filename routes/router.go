package routes

import (
	"mall/middleware"
	"mall/svc"
	"net/http"
	"time"

	api "mall/api/v1"
	_ "mall/docs"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title mall API
// @version 1.0
// @description mall server.
// @license.name Apache 2.0
// @host localhost:8000
// @BasePath /api/v1
func NewRouter(svc *svc.ServiceContext) *gin.Engine {
	e := gin.Default()
	e.Use(middleware.RecordMetricsMiddleware(), middleware.CorsMiddleware())
	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	e.StaticFS("/static", http.Dir("./static"))
	e.GET("/metrics", gin.WrapH(promhttp.Handler()))

	v1 := e.Group("api/v1")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(200, "success")
		})
		// 用户操作
		v1.POST("user/register", api.UserRegisterHandler(svc))
		v1.POST("user/login", api.UserLoginHandler(svc))
		// 商品操作（显示等）
		// http://localhost:8000/api/v1/1/products?offset=1&limit=10
		v1.GET("/products-categories/:category_id", api.ProductListHandler(svc))
		v1.GET("/products/:id", api.ProductShowHandler(svc))
		// http://localhost:8000/api/v1/products?offset=1&limit=10&seach_info="in"
		v1.GET("/products", api.ProductSearchHandler(svc))
		v1.GET("/products/imgs/:id", api.ProductImgListHandler(svc)) // 商品图片
		// 商品分类
		v1.GET("/categories", api.CategoryListHandler(svc))
		// v1.GET("carousels", api.ListCarouselsHandler()) // 轮播图

		authed := v1.Group("/")
		authed.Use(middleware.AuthMiddleware(svc))
		{
			{
				authed.GET("/user", api.UserShowHandler(svc))
			}
			// 商品操作
			{
				authed.POST("/products", api.ProductCreateHandler(svc))
				authed.PUT("/products/:id", api.ProductUpdateHandler(svc))
				authed.DELETE("/products/:id", api.ProductDeleteHandler(svc))
			}
			// 种类
			{
				authed.POST("/categories", api.CategoryCreateHandler(svc))
				authed.DELETE("/categories/:id", api.CategoryDeleteHandler(svc))
			}
			// 地址
			{
				authed.POST("/addresses", api.AddressCreateHandler(svc))
				authed.GET("/addresses/:id", api.AddressGetHandler(svc))
				// http://localhost:8000/api/v1/address?offset=1&limit=10
				authed.GET("/addresses", api.AddressListHandler(svc))
				authed.DELETE("/addresses/:id", api.AddressDeleteHandler(svc))
				authed.PUT("/addresses/:id", api.AddressUpdateHandler(svc)) // 非全量更新（全量更新 PUT）
			}
			// 购物车
			{
				authed.POST("/carts/:cart_id/item-infos/products/:product_id", api.ItemInfoAddHandler(svc))
				authed.PUT("/carts/item-infos/:id", api.ItemInfoUpdateHandler(svc))
				authed.DELETE("/carts/item-infos/:id", api.ItemInfoDeleteHandler(svc))
				authed.GET("/carts/item-infos/:id", api.ItemInfoShowHandler(svc))
				authed.GET("/carts/:cart_id", api.CartInfoHandler(svc))
			}
			{
				authed.POST("/orders", api.OrderCreateHandler(svc))
				authed.DELETE("/orders/:order_id", api.OrderDeleteHandler(svc))
				authed.GET("/orders/:order_id", api.OrderGetHandler(svc))
				authed.GET("/orders", api.OrderListHandler(svc))
				authed.POST("/orders/:order_id/payment", api.OrderPaymentHandler(svc))
			}
			{
				authed.POST("/skill-goods", api.SkillProductImportHandler(svc))
				authed.GET("/skill-goods", api.SkillStartHandler(svc))
				authed.POST("/skill-goods/:product_id",
					middleware.RatelimitMiddleware(time.Second*time.Duration(svc.Config.HttpServerConf.FillInterval), svc.Config.HttpServerConf.Cap),
					api.SkillHandler(svc),
				) // 限流
			}
		}
	}
	return e
}
