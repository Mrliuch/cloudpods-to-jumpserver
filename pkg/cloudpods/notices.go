package cloudpods

import (
	"cloudpods-webhook/pkg/common"
	"cloudpods-webhook/pkg/db"
	"cloudpods-webhook/pkg/jumpserver"
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

type Notices struct {
	Notices *common.Notice
}

// 处理WebHook事件

func (n *Notices) DealEvent() error {
	if n.Notices.ResourceType != "server" {
		logrus.Infof("资源类型：%s,不予处理！", n.Notices.ResourceType)
		return nil
	}
	if n.Notices.Result != "succeed" {
		logrus.Infof("资源状态：%s,不予处理！", n.Notices.Result)
		return nil
	}
	switch n.Notices.Action {
	case "create":
		return n.CreateVM()
	case "delete":
		return n.DeleteVM()
	}
	return nil
}

// 创建虚拟机

func (n *Notices) CreateVM() error {
	logrus.Infof("获取到虚拟机ID：%s，IP：%s，Name：%s", n.Notices.ResourceDetails.HostID, n.Notices.ResourceDetails.IPS, n.Notices.ResourceDetails.Name)
	return jumpserver.CreateVM(n.Notices.ResourceDetails)
}

// 删除虚拟机

func (n *Notices) DeleteVM() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	jmpID, err := db.Rdb.Get(ctx, n.Notices.ResourceDetails.HostID).Result()
	if err != nil {
		logrus.Errorf("获取虚拟机IP失败，HostID：%s，错误原因：%s", n.Notices.ResourceDetails.HostID, err.Error())
		return err
	}
	n.Notices.ResourceDetails.JMPID = jmpID
	logrus.Infof("删除虚拟机资源ID：%s", n.Notices.ResourceDetails.HostID)
	err = jumpserver.DeleteVM(n.Notices.ResourceDetails)
	if err != nil {
		logrus.Errorf("删除虚拟机失败：%s", err.Error())
	}
	return db.Rdb.Del(ctx, n.Notices.ResourceDetails.HostID).Err()
}
