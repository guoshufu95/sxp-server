package config

import (
	"github.com/spf13/viper"
	"log"
)

// Config
// @Description: 配置文件解析
type Config struct {
	Server struct {
		Host string
		Mode string
		Port string
	}
	Logger struct {
		Path  string
		Level string
	}
	Mysql struct {
		Host     string
		Port     string
		Mode     string
		Db       string
		UserName string
		Password string
	}
	Redis struct {
		Addr     string
		Password string
	}
	Jwt struct {
		Secret  string
		Timeout int64
	}
	Grpc struct {
		Addr    string
		Retry   int
		TimeOut int
	}
	Jaeger struct {
		Addr string
	}
	Kafka struct {
		Brokers         []string
		ProducerTimeOut int
		ConsumerTimeOut int
		Ack             int
	}
}

var Conf *Config

// ReadConfig
//
//	@Description: 加载配置文件
//	@param path
//	@return *Config
func ReadConfig(path string) *Config {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		panic("配置文件读取失败: " + err.Error())
	}
	if err := v.Unmarshal(&Conf); err != nil {
		panic("配置文件反序列化失败: " + err.Error())
	}
	log.Println("配置文件内容加载成功: ", path)
	return Conf
}
