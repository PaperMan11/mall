package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	HttpServerConf  *HttpServerConf  `mapstructure:"httpserver"`
	LogConf         *LogConf         `mapstructure:"log"`
	LocalUploadConf *LocalUploadConf `mapstructure:"local_upload"`
	MysqlConf       *MysqlConf       `mapstructure:"mysql"`
	RedisConf       *RedisConf       `mapstructure:"redis"`
	RabbitMQConf    *RabbitMQConf    `mapstructure:"rabbitmq"`
}

type HttpServerConf struct {
	Port         int    `mapstructure:"port"`
	Model        string `mapstructure:"model"`
	ServiceName  string `mapstructure:"servicename"`
	Key          string `mapstructure:"key"`
	FillInterval int64  `mapstructure:"fill_interval"`
	Cap          int64  `mapstructure:"cap"`
}

type LogConf struct {
	EnableConsole     bool   `mapstructure:"enableConsole"`
	ConsoleJSONFormat bool   `mapstructure:"consoleJSONFormat"`
	ConsoleLevel      string `mapstructure:"consoleLevel"`
	EnableFile        bool   `mapstructure:"enableFile"`
	FileJSONFormat    bool   `mapstructure:"fileJSONFormat"`
	FileLevel         string `mapstructure:"fileLevel"`
	FileLocation      string `mapstructure:"fileLocation"`
	MaxAge            int    `mapstructure:"maxAge"`
	MaxSize           int    `mapstructure:"maxSize"`
	Compress          bool   `mapstructure:"compress"`
}

type LocalUploadConf struct {
	PhotoHost        string `mapstructure:"photoHost"`
	ProductPhotoPath string `mapstructure:"productPhotoPath"`
	AvatarPath       string `mapstructure:"avatarPath"`
}

type MysqlConf struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DbName       string `mapstructure:"db"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConf struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Db       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type RabbitMQConf struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
}

func Init(filePath string) (conf Config, err error) {
	viper.SetConfigFile(filePath) // 指定配置文件以及路径
	err = viper.ReadInConfig()    // 读取配置信息
	if err != nil {
		fmt.Println("viper.ReadInConfig failed err:", err)
		return
	}
	// 把读取到的配置的信息反序列化到 Conf 变量中
	if err = viper.Unmarshal(&conf); err != nil {
		fmt.Println("viper.Unmarshal failed")
	}
	viper.WatchConfig()                            // 热加载
	viper.OnConfigChange(func(in fsnotify.Event) { // hockfunc
		fmt.Println("配置文件修改了")
		if err = viper.Unmarshal(&conf); err != nil {
			fmt.Println("viper.Unmarshal failed")
		}
	})
	return
}
