package svc

import (
	"mall/config"
	"mall/pkg/utils/encryption"
	"mall/pkg/utils/logger"
	"mall/repositry/cache"
	"mall/repositry/mysqldb"
	"mall/repositry/rabbitmq"

	"mall/pkg/utils/snowflake"

	"github.com/redis/go-redis/v9"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config    *config.Config
	Log       logger.Logger
	SnowFlake *snowflake.IDWorker
	Encrypt   *encryption.Encryption // 金钱加解密

	Redis             *redis.Client
	ProductRedisCache cache.ProductRedisCache
	SkRedisCache      cache.SkRedisCache

	DB              *gorm.DB
	UserModel       mysqldb.UserModel
	ProductModel    mysqldb.ProductModel
	ProductImgModel mysqldb.ProductImgModel
	CategoryModel   mysqldb.CategoryModel
	AddressModel    mysqldb.AddressModel
	CartModel       mysqldb.CartModel
	ItemInfoModel   mysqldb.ItemInfoModel
	OrderModel      mysqldb.OrderModel
	SkillModel      mysqldb.SkillModel

	Amqp   *amqp.Connection
	SkAmqp rabbitmq.SkAmqp
}

func NewServiceContext(c config.Config) *ServiceContext {
	// log
	log := logger.InitLogger(logger.Configuration{
		EnableConsole:     c.LogConf.EnableConsole,
		ConsoleJSONFormat: c.LogConf.ConsoleJSONFormat,
		ConsoleLevel:      c.LogConf.ConsoleLevel,
		EnableFile:        c.LogConf.EnableFile,
		FileJSONFormat:    c.LogConf.FileJSONFormat,
		FileLevel:         c.LogConf.FileLevel,
		FileLocation:      c.LogConf.FileLocation,
		MaxAge:            c.LogConf.MaxAge,
		MaxSize:           c.LogConf.MaxSize,
		Compress:          c.LogConf.Compress,
	}, logger.InstanceZapLogger)

	// mysql
	dbConn, err := mysqldb.NewMysqlConn(c.MysqlConf)
	if err != nil {
		panic(err)
	}

	// redis
	redisConn, err := cache.NewRedisConn(c.RedisConf)
	if err != nil {
		panic(err)
	}

	// rabbitmq
	amqpConn, err := rabbitmq.NewRabbitMQConn(c.RabbitMQConf)
	if err != nil {
		panic(err)
	}

	// snowflake
	snowflaker, err := snowflake.NewIDWorker(1)
	if err != nil {
		panic(err)
	}

	// encrypt
	encrypt := encryption.NewEncryption()
	encrypt.SetKey(c.HttpServerConf.Key)

	return &ServiceContext{
		Config:    &c,
		Log:       log,
		SnowFlake: snowflaker,
		Encrypt:   encrypt,

		Redis:             redisConn,
		ProductRedisCache: cache.NewProductRedisCache(redisConn),
		SkRedisCache:      cache.NewSkRedisCache(redisConn),

		DB:              dbConn,
		UserModel:       mysqldb.NewUserModel(dbConn, "user"),
		ProductModel:    mysqldb.NewProductModel(dbConn, "product"),
		ProductImgModel: mysqldb.NewProductImgModel(dbConn, "product_img"),
		CategoryModel:   mysqldb.NewCategoryModel(dbConn, "category"),
		AddressModel:    mysqldb.NewAddressModel(dbConn, "address"),
		CartModel:       mysqldb.NewCartModel(dbConn, "cart"),
		ItemInfoModel:   mysqldb.NewItemInfoModel(dbConn, "item_info"),
		OrderModel:      mysqldb.NewOrderModel(dbConn, "order"),
		SkillModel:      mysqldb.NewSkillModel(dbConn, "skill_porduct"),

		Amqp:   amqpConn,
		SkAmqp: rabbitmq.NewSkAmqp(amqpConn),
	}
}

func (svc *ServiceContext) Close() {
	var err error
	if err = svc.Redis.Close(); err != nil {
		svc.Log.Errorf("Redis.Close() err: %s", err)
	}

	sqlDB, err := svc.DB.DB()
	if err != nil {
		svc.Log.Errorf("Get sqlDB err: %s", err)
	}
	if sqlDB != nil {
		if err = sqlDB.Close(); err != nil {
			svc.Log.Errorf("sqlDB.Close() err: %s", err)
		}
	}

	if err = svc.Amqp.Close(); err != nil {
		svc.Log.Errorf("Amqp.Close() err: %s", err)
	}
}
