package config

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ymlConfig struct {
	Server    *ServerConfig
	Zookeeper *ZkConfig
}

// ParseConfig 解析命令行参数
func ParseConfig() {

	// 保存配置
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()

	if err != nil {
		logrus.Error("读取配置文件失败: %v", err)
	}
	for _, key := range viper.AllKeys() {
		viper.BindEnv(key, strings.ReplaceAll(key, ".", "_"))
	}

	// 保存配置
	configYml := &ymlConfig{}
	viper.Unmarshal(configYml)
	SetServerConfig(configYml.Server)
	SetZkConfig(configYml.Zookeeper)

	// 验证
	if errs := validator.New().Struct(configYml.Server); errs != nil {
		logrus.Panic(errs)
	}
	if errs := validator.New().Struct(configYml.Zookeeper); errs != nil {
		logrus.Panic(errs)
	}

	// 打印配置
	fmt.Printf("%c[%dm%s%c[m\n", 0x1B, 0x23, fmt.Sprintf("ServerConfig config: %+v", configYml.Server), 0x1B)
	fmt.Printf("%c[%dm%s%c[m\n", 0x1B, 0x23, fmt.Sprintf("ZooKeeper config: %+v", configYml.Zookeeper), 0x1B)

}
