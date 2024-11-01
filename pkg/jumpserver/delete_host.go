package jumpserver

import (
	"cloudpods-webhook/pkg/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"time"
)

func DeleteVM(details common.ResourceDetails) error {
	method := "/api/v1/assets/assets/"
	var j JumpServerAPI
	_ = j.Delete(method, details)
	return nil
}

func (p JumpServerAPI) Delete(method string, n common.ResourceDetails) []byte {
	var body []byte
	auth := SigAuth{
		KeyID:    viper.GetString("jumpserver.AccessKeyID"),
		SecretID: viper.GetString("jumpserver.AccessKeySecret"),
	}
	url := viper.GetString("jumpserver.address") + method + n.JMPID + "/"
	gmtFmt := "Mon, 02 Jan 2006 15:04:05 GMT"
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Date", time.Now().Format(gmtFmt))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-JMS-ORG", "00000000-0000-0000-0000-000000000002")
	if err != nil {
		log.Fatal(err)
	}
	err = auth.Sign(req)
	if err != nil {
		logrus.Info(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Info(err)
	}
	//	defer resp.Body.Close()
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			logrus.Info(err)
		}
	}()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		logrus.Info(err)
	}
	return body
}
