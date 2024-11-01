package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	logrus.Infof("初始化日志服务！")
	log := logrus.New()
	log.Infof("加载配置！")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Errorf("读取配置文件失败: %v", err)
	}
	viper.AutomaticEnv()
	logrus.Infof("配置加载完成！")
}
