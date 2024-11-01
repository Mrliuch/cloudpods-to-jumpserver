package jumpserver

import (
	"cloudpods-webhook/pkg/common"
	"cloudpods-webhook/pkg/db"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func CreateVM(details common.ResourceDetails) error {
	method := "/api/v1/assets/hosts/"
	var j JumpServerAPI
	assetID, _ := uuid.NewV4()
	t := genHostInfo(details)
	t.ID = assetID.String()
	buf, err := json.MarshalIndent(t, "", "")
	if err != nil {
		return err
	}
	formattedJsonStr := strings.ReplaceAll(string(buf), ",\n", ", ")
	formattedJsonStr = strings.ReplaceAll(string(formattedJsonStr), "\n", "")
	statusCode, body := j.Post(method, formattedJsonStr)
	if statusCode == 201 {
		var host common.RespJMP
		err = json.Unmarshal(body, &host)
		if err != nil {
			logrus.Errorf("解析jumpserver返回失败: %s", err.Error())
			return err
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		err = db.Rdb.Set(ctx, details.HostID, host.ID, 0).Err()
		if err != nil {
			logrus.Errorf("存储资产信息失败：%s", err.Error())
			return err
		}
		logrus.Infof("新资产创建成功!")
	} else {
		logrus.Info("资产创建失败:" + string(body))
		err = fmt.Errorf("资产创建失败:" + string(body))
	}
	return err
}

func (p JumpServerAPI) Post(method string, msgText string) (int, []byte) {
	var body []byte
	var resp *http.Response
	var req *http.Request
	var err error
	auth := SigAuth{
		KeyID:    viper.GetString("jumpserver.AccessKeyID"),
		SecretID: viper.GetString("jumpserver.AccessKeySecret"),
	}
	url := viper.GetString("jumpserver.address") + method
	logrus.Infof("请求URL：%s", url)

	gmtFmt := "Mon, 02 Jan 2006 15:04:05 GMT"
	client := &http.Client{}
	req, err = http.NewRequest("POST", url, strings.NewReader(msgText))
	req.Header.Add("Date", time.Now().Format(gmtFmt))
	req.Header.Add("Content-Disposition", "inline;filename=file.txt")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-JMS-ORG", "00000000-0000-0000-0000-000000000002")
	err = auth.Sign(req)
	if err != nil {
		logrus.Errorf("堡垒机认证失败：%s", err.Error())
	}
	resp, err = client.Do(req)
	if err != nil {
		logrus.Errorf("堡垒机请求失败：%s", err.Error())
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("获取堡垒机响应失败：%s", err.Error())
	}
	return resp.StatusCode, body
}
