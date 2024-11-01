package jumpserver

import (
	"cloudpods-webhook/pkg/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/twindagger/httpsig.v1"
	"net/http"
)

func init() {
	logrus.Infof("获取到堡垒机地址：%s", viper.GetString("jumpserver.address"))
}

type SigAuth struct {
	KeyID    string
	SecretID string
}

type JumpServerAPI struct{}

func (auth *SigAuth) Sign(r *http.Request) error {
	headers := []string{"(request-target)", "date"}
	signer, err := httpsig.NewRequestSigner(auth.KeyID, auth.SecretID, "hmac-sha256")
	if err != nil {
		return err
	}
	return signer.SignRequest(r, headers, nil)
}

func genHostInfo(details common.ResourceDetails) *common.Asset {
	//passwd, err := utils.EncryptPassword(details.Password)
	//if err != nil {
	//	logrus.Error("堡垒机加密密码失败：" + err.Error())
	//}
	host := &common.Asset{
		ID:       details.HostID,
		Name:     details.Name,
		Address:  details.IPS,
		Platform: "2",
		Protocols: []common.Protocol{
			{
				Name: "ssh",
				Port: 22,
			}, {
				Name: "sftp",
				Port: 22,
			},
		},
		Nodes:    []string{viper.GetString("jumpserver.node_id")},
		IsActive: true,
		Accounts: []common.Accounts{
			{
				Name:       details.LoginAccount,
				Username:   details.LoginAccount,
				Secret:     details.Password,
				Privileged: true,
				IsActive:   true,
				SecretType: "password",
			},
		},
	}
	switch details.OSType {
	case "l24", "l26", "solaris":
		host.Platform = viper.GetString("jumpserver.linuxPlatform")
		host.Protocols = []common.Protocol{
			{
				Name: "ssh",
				Port: 22,
			}, {
				Name: "sftp",
				Port: 22,
			},
		}
	case "wxp", "w2k", "w2k3", "w2k8", "wvista", "win7", "win8", "win10", "win11":
		host.Platform = viper.GetString("jumpserver.winPlatform")
		host.Protocols = []common.Protocol{
			{
				Name: "rdp",
				Port: 3389,
			},
		}
	default:
		host.Platform = viper.GetString("jumpserver.linuxPlatform")
		host.Protocols = []common.Protocol{
			{
				Name: "ssh",
				Port: 22,
			}, {
				Name: "sftp",
				Port: 22,
			},
		}
	}
	return host
}
