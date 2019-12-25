package jmessage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/franela/goreq"
)

// JMUser 用户信息 返回信息
type JMUser struct {
	Username  string `json:"username"`            //（必填 Byte(4~128) 开头：字母或者数字 字母、数字、下划线英文点、减号、@
	Password  string `json:"password"`            //（必填） Byte(4~128) 用户密码。极光IM服务器会MD5加密保存
	Nickname  string `json:"nickname,omitempty"`  //Byte(0~64) 不支持的字符：英文字符： \n \r\n
	Avatar    string `json:"avatar,omitempty"`    //（选填）头像 需要填上从文件上传接口获得的media_id
	Birthday  string `json:"star,omitempty"`      //(选填）生日 example: 1990-01-24
	Gender    int    `json:"gender,omitempty"`    //性别 0 - 未知， 1 - 男 ，2 - 女
	Signature string `json:"signature,omitempty"` //用户签名 Byte(0~250)
	Region    string `json:"region,omitempty"`    //（选填）地区 用户所属地区 Byte(0~250)
	Address   string `json:"address,omitempty"`   //（选填）地址 用户详细地址	Byte(0~250)
	Mtime     int    `json:"mtime,omitempty"`     //用户最后修改时间
	Ctime     int    `json:"ctime,omitempty"`     //用户创建时间

	Error *JMError `json:"error,omitempty"`
}
type JMResponse struct {
	Error *JMError `json:"error,omitempty"`
}
type JMError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *JMError) Error() string {
	return fmt.Sprintf("%d - %s", e.Code, e.Message)
}

//RegisterUsers 批量注册用户
func (jclient *JMessageClient) RegisterUsers(users []*JMUser) ([]*JMUser, error) {

	rusers := []*JMUser{}

	if nil == users || len(users) == 0 {
		return rusers, errors.New("no user to register")
	}
	if len(jclient.appKey) != KEY_LENGTH || len(jclient.masterSecret) != KEY_LENGTH {
		return rusers, fmt.Errorf("invalidate appkey/masterSecret")
	}

	req := goreq.Request{
		Method:            "POST",
		Uri:               JMESSAGE_IM_URL + REGIST_USER_URL,
		Accept:            "application/json",
		ContentType:       "application/json",
		UserAgent:         "JMessage-API-GO-Client",
		BasicAuthUsername: jclient.appKey,
		BasicAuthPassword: jclient.masterSecret,
		Timeout:           30 * time.Second, //30s
	}
	req.Body = users
	req.ShowDebug = jclient.showDebug

	res, err := req.Do()

	if err != nil {
		return rusers, err
	}

	defer res.Body.Close()

	ibytes, err := ioutil.ReadAll(res.Body)
	if nil != err {
		return rusers, err
	}

	if jclient.showDebug {
		fmt.Println("respone:", string(ibytes))
	}

	err = json.Unmarshal(ibytes, &rusers)
	if nil != err {
		return rusers, err
	}
	return rusers, nil
}

//RegisterUser 注册用户
func (jclient *JMessageClient) RegisterUser(username string,
	nickName string,
	password string,
	avatar string) (*JMUser, *JMError, error) {

	user := JMUser{
		Username: username,
		Password: password,
		Nickname: nickName,
		Avatar:   avatar,
	}
	users, err := jclient.RegisterUsers([]*JMUser{&user})
	if nil != err {
		return nil, nil, err
	}

	if len(users) == 1 {
		jmuser := users[0]

		if nil == jmuser.Error || jmuser.Error.Code == 0 {
			return jmuser, nil, nil
		}

		return nil, jmuser.Error, jmuser.Error
	}

	return nil, nil, errors.New("response failed")
}

func (jclient *JMessageClient) RegisterAdmin(username string,
	nickName string,
	password string,
	avatar string) (*JMUser, *JMError, error) {

	user := JMUser{
		Username: username,
		Password: password,
		Nickname: nickName,
		Avatar:   avatar,
	}
	if len(jclient.appKey) != KEY_LENGTH || len(jclient.masterSecret) != KEY_LENGTH {
		return nil, nil, fmt.Errorf("invalidate appkey/masterSecret")
	}

	req := goreq.Request{
		Method:            "POST",
		Uri:               JMESSAGE_IM_URL + REGIST_ADMIN_URL,
		Accept:            "application/json",
		ContentType:       "application/json",
		UserAgent:         "JMessage-API-GO-Client",
		BasicAuthUsername: jclient.appKey,
		BasicAuthPassword: jclient.masterSecret,
		Timeout:           30 * time.Second, //30s
	}
	req.Body = user
	req.ShowDebug = jclient.showDebug

	res, err := req.Do()

	if err != nil {
		return nil, nil, err
	}

	defer res.Body.Close()

	ibytes, err := ioutil.ReadAll(res.Body)
	if nil != err {
		return nil, nil, err
	}
	if jclient.showDebug {
		fmt.Println("respone:", string(ibytes))
	}

	if string(ibytes) == "" {
		return &user, nil, nil
	}
	jmResult := JMResponse{}

	err = json.Unmarshal(ibytes, &jmResult)
	if nil != err {
		return nil, nil, err
	}

	if nil != jmResult.Error {
		return nil, jmResult.Error, jmResult.Error
	}

	return &user, nil, nil
}

func (jclient *JMessageClient) UpdatePasswd(username string, passwd string) error {

	req := goreq.Request{
		Method:            "PUT",
		Uri:               JMESSAGE_IM_URL + REGIST_USER_URL + username + "/password",
		Accept:            "application/json",
		ContentType:       "application/json",
		UserAgent:         "JMessage-API-GO-Client",
		BasicAuthUsername: jclient.appKey,
		BasicAuthPassword: jclient.masterSecret,
		Timeout:           30 * time.Second, //30s
	}
	req.ShowDebug = jclient.showDebug
	req.Body = map[string]string{"new_password": passwd}
	res, err := req.Do()

	if err != nil {
		return err
	}

	defer res.Body.Close()

	ibytes, err := ioutil.ReadAll(res.Body)
	if nil != err {
		return err
	}
	if jclient.showDebug {
		fmt.Println("respone:", string(ibytes))
	}

	if string(ibytes) == "" {
		return nil
	}
	jmResult := JMResponse{}

	err = json.Unmarshal(ibytes, &jmResult)
	if nil != err {
		return err
	}

	if nil != jmResult.Error {
		return jmResult.Error
	}

	return nil
}
func (jclient *JMessageClient) UpdateProfile(username string, nickname, avatar, birthday string,
	signature, gender, region, address string,
	extras string) error {

	req := goreq.Request{
		Method:            "PUT",
		Uri:               JMESSAGE_IM_URL + REGIST_USER_URL + username,
		Accept:            "application/json",
		ContentType:       "application/json",
		UserAgent:         "JMessage-API-GO-Client",
		BasicAuthUsername: jclient.appKey,
		BasicAuthPassword: jclient.masterSecret,
		Timeout:           30 * time.Second, //30s
	}
	req.ShowDebug = jclient.showDebug

	params := map[string]string{}
	if nickname != "" {
		params["nickname"] = nickname
	}
	if avatar != "" {
		params["avatar"] = avatar
	}
	if birthday != "" {
		params["birthday"] = birthday
	}

	if signature != "" {
		params["signature"] = signature
	}
	if gender != "" {
		params["gender"] = gender
	}
	if region != "" {
		params["region"] = region
	}
	if address != "" {
		params["address"] = address
	}

	if extras != "" {
		params["extras"] = extras
	}

	req.Body = params
	res, err := req.Do()

	if err != nil {
		return err
	}

	defer res.Body.Close()

	ibytes, err := ioutil.ReadAll(res.Body)
	if nil != err {
		return err
	}
	if jclient.showDebug {
		fmt.Println("respone:", string(ibytes))
	}

	if string(ibytes) == "" {
		return nil
	}
	jmResult := JMResponse{}

	err = json.Unmarshal(ibytes, &jmResult)
	if nil != err {
		return err
	}

	if nil != jmResult.Error {
		return jmResult.Error
	}

	return nil
}

func (jclient *JMessageClient) DeleteUser(username string) error {

	req := goreq.Request{
		Method:            "DELETE",
		Uri:               JMESSAGE_IM_URL + REGIST_USER_URL + username,
		Accept:            "application/json",
		ContentType:       "application/json",
		UserAgent:         "JMessage-API-GO-Client",
		BasicAuthUsername: jclient.appKey,
		BasicAuthPassword: jclient.masterSecret,
		Timeout:           30 * time.Second, //30s
	}
	req.ShowDebug = jclient.showDebug
	res, err := req.Do()

	if err != nil {
		return err
	}

	defer res.Body.Close()

	ibytes, err := ioutil.ReadAll(res.Body)
	if nil != err {
		return err
	}
	if jclient.showDebug {
		fmt.Println("respone:", string(ibytes))
	}

	if string(ibytes) == "" {
		return nil
	}
	jmResult := JMResponse{}

	err = json.Unmarshal(ibytes, &jmResult)
	if nil != err {
		return err
	}

	if nil != jmResult.Error {
		return jmResult.Error
	}

	return nil
}

//BlackUsers 添加黑名单
func (jclient *JMessageClient) BlackUsers(fromUsername string, blackUserNames []string) error {

	if len(blackUserNames) == 0 {
		return errors.New("empty blackUserNames")
	}
	res, err := jclient.request(JMESSAGE_IM_URL+REGIST_USER_URL+fromUsername+"/blacklist", "PUT", blackUserNames)
	if nil != err {
		return err
	}
	defer res.Body.Close()

	ibytes, err := ioutil.ReadAll(res.Body)
	if nil != err {
		return err
	}

	if jclient.showDebug {
		fmt.Println("respone:", string(ibytes))
	}

	return nil
}

//DeleteBlackUsers 删除黑名单
func (jclient *JMessageClient) DeleteBlackUsers(fromUsername string, blackUserNames []string) error {

	if len(blackUserNames) == 0 {
		return errors.New("empty blackUserNames")
	}
	res, err := jclient.request(JMESSAGE_IM_URL+REGIST_USER_URL+fromUsername+"/blacklist", "DELETE", blackUserNames)
	if nil != err {
		return err
	}
	defer res.Body.Close()

	ibytes, err := ioutil.ReadAll(res.Body)
	if nil != err {
		return err
	}

	if jclient.showDebug {
		fmt.Println("respone:", string(ibytes))
	}

	return nil
}

//GetBlackUsers 获取黑名单
func (jclient *JMessageClient) GetBlackUsers(fromUsername string) ([]string, error) {

	users := []string{}
	res, err := jclient.request(JMESSAGE_IM_URL+REGIST_USER_URL+fromUsername+"/blacklist", "Get", nil)
	if nil != err {
		return users, err
	}
	defer res.Body.Close()

	ibytes, err := ioutil.ReadAll(res.Body)
	if nil != err {
		return users, err
	}

	if jclient.showDebug {
		fmt.Println("respone:", string(ibytes))
	}

	err = json.Unmarshal(ibytes, &users)
	return users, err
}

//ForbiddenUser 禁用/解除禁用
func (jclient *JMessageClient) ForbiddenUser(username string, forbidden bool) error {

	url := JMESSAGE_IM_URL + REGIST_USER_URL + username + "/forbidden?disable="
	if forbidden {
		url += "true"
	} else {
		url += "false"
	}
	res, err := jclient.request(url, "PUT", nil)
	if nil != err {
		return err
	}
	defer res.Body.Close()

	ibytes, err := ioutil.ReadAll(res.Body)
	if nil != err {
		return err
	}

	if jclient.showDebug {
		fmt.Println("respone:", string(ibytes))
	}

	return nil
}
