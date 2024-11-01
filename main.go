package main

import (
	"cloudpods-webhook/cmd/api"
	_ "cloudpods-webhook/pkg/config"
	_ "cloudpods-webhook/pkg/db"
	_ "cloudpods-webhook/pkg/jumpserver"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	port := viper.GetString("api_port")
	logrus.Infof("启动API端口：%s", port)
	if err := api.RunApi().Run(fmt.Sprintf(":%s", port)); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
