package jmessage

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/franela/goreq"
)

const (
	KEY_LENGTH   = 24
	CONN_TIMEOUT = 5
	RW_TIMEOUT   = 30
)

type JMessageClient struct {
	appKey       string
	masterSecret string
}

//NewJMessageClient 创建client
func NewJMessageClient(appkey, master_secret string) *JMessageClient {
	return &JMessageClient{appkey, master_secret}
}

func (jclient *JMessageClient) buildReq(uri string, method string, body interface{}) (*goreq.Request, error) {
	if len(jclient.appKey) != KEY_LENGTH || len(jclient.masterSecret) != KEY_LENGTH {
		return nil, fmt.Errorf("invalidate appkey/masterSecret")
	}

	req := goreq.Request{
		Method:            method, //"POST",
		Uri:               uri,
		Accept:            "application/json",
		ContentType:       "application/json",
		UserAgent:         "JMessage-API-GO-Client",
		BasicAuthUsername: jclient.appKey,
		BasicAuthPassword: jclient.masterSecret,
		Timeout:           30 * time.Second, //30s
	}
	req.Body = body
	req.ShowDebug = ShowDebug

	return &req, nil
}

func (jclient *JMessageClient) request(uri string, method string, body interface{}) (*goreq.Response, error) {
	req, err := jclient.buildReq(uri, method, body)
	if nil != err {
		return nil, err
	}

	if req.ShowDebug {
		fmt.Printf("request:%s, %s\n", uri, method)
	}
	res, err := req.Do()

	if err != nil {
		return nil, err
	}

	return res, nil
}

//SentSystemTxtMsg 发送系统文本消息
func (jclient *JMessageClient) SentSystemTxtMsg(fromId string,
	targetType string,
	targetId string,
	message string, ext map[string]interface{}) error {

	jpMessage := JPMessage{Version: 1, FromType: "admin", MsgType: "text"}

	jpMessage.FromId = fromId
	jpMessage.TargetType = targetType
	jpMessage.TargetId = targetId

	txtMsg := JPTxtMsg{}
	txtMsg.Text = message
	txtMsg.Extra = ext

	jpMessage.MsgBody = txtMsg

	res, err := jclient.request(JMESSAGE_IM_URL+MESSAGES_URL, "POST", jpMessage)
	if nil != err {
		return err
	}
	defer res.Body.Close()

	ibytes, err := ioutil.ReadAll(res.Body)
	if nil != err {
		return err
	}

	if ShowDebug {
		fmt.Println("respone:", string(ibytes))
	}

	return nil
}
