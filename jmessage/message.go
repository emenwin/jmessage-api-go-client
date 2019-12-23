package jmessage

import (
	"encoding/json"
)

// JPMessage 消息
type JPMessage struct {
	Version    int    `json:"version"`     //版本号 目前是1 （必填）
	TargetType string `json:"target_type"` //发送目标类型 single - 个人，group - 群组 chatroom - 聊天室（必填）
	TargetId   string `json:"target_id"`   //目标id single填username group 填Group id chatroom 填chatroomid（必填）
	FromType   string `json:"from_type"`   //发送消息者身份 当前只限admin用户，必须先注册admin用户 （必填）
	FromId     string `json:"from_id"`     //发送者的username （必填）

	MsgType string                 `json:"msg_type"` //发消息类型 text - 文本，image - 图片, custom - 自定义消息（msg_body为json对象即可，服务端不做校验）voice - 语音 （必填）
	MsgBody map[string]interface{} `json:"msg_body"`

	FromName       string          `json:"from_name,omitempty"`       //发送者展示名（选填）
	TargetName     string          `json:"target_name,omitempty"`     //接受者展示名（选填）
	NoOffline      string          `json:"no_offline,omitempty"`      //消息是否离线存储 true或者false，默认为false，表示需要离线存储（选填）
	NoNotification bool            `json:"no_notification,omitempty"` //消息是否在通知栏展示 true或者false，默认为false，表示在通知栏展示（选填）
	Notification   *JPNotification `json:"notification,omitempty"`    //自定义通知栏展示（选填）

	MsgId    int64 `json:"msgid,omitempty"`     //消息id
	MsgLevel uint8 `json:"msg_level,omitempty"` // 0代表应用内消息 1代表跨应用消息
	MsgCtime int64 `json:"msg_ctime,omitempty"` // 服务器接收到消息的时间，单位毫秒

}

// JPNotification 消息
type JPNotification struct {
	Title string `json:"title,omitempty"` //通知的标题（选填）
	Alert string `json:"alert,omitempty"` //通知的内容（选填）
}

// JMsg RCMsg接口
type JPMessageBody interface {
	toString() (string, error)
	toMap() (map[string]interface{}, error)
}

// JPTxtMsg 消息
type JPTxtMsg struct {
	Text  string                 `json:"text"`            //消息内容 （必填）
	Extra map[string]interface{} `json:"extra,omitempty"` //选填的json对象 开发者可以自定义extras里面的key value（选填）
}

// JPIMGMsg 消息
type JPIMGMsg struct {
	MediaId    string `json:"media_id"`       //String 文件上传之后服务器端所返回的key，用于之后生成下载的url（必填）
	MediaCrc32 int64  `json:"media_crc32"`    //long 文件的crc32校验码，用于下载大图的校验 （必填）
	width      int    `json:"width"`          //int 图片原始宽度（必填）
	height     int    `json:"height"`         //int 图片原始高度（必填）
	Format     string `json:"format"`         //String 图片格式（必填）
	Hash       string `json:"hash,omitempty"` //String 图片hash值（可选）
	Fsize      int    `json:"fsize"`          //（必填）
}

// JPVoiceMsg 消息
type JPVoiceMsg struct {
	MediaId    string `json:"media_id"`       //String 文件上传之后服务器端所返回的key，用于之后生成下载的url（必填）
	MediaCrc32 int64  `json:"media_crc32"`    //long 文件的crc32校验码，用于下载大图的校验 （必填）
	Duration   int    `json:"duration"`       //int 音频时长（必填）
	Hash       string `json:"hash,omitempty"` //String 音频hash值（可选）
	Fsize      int    `json:"fsize"`          //（必填）
}

func objectToMap(msg interface{}) (map[string]interface{}, error) {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{}
	err = json.Unmarshal(bytes, &result)
	return result, err
}

// toString TXTMsg
func (msg JPTxtMsg) toString() (string, error) {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
func (msg JPTxtMsg) toMap() (map[string]interface{}, error) {
	return objectToMap(&msg)
}

// toString ImgMsg
func (msg JPIMGMsg) toString() (string, error) {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
func (msg JPIMGMsg) toMap() (map[string]interface{}, error) {
	return objectToMap(&msg)
}

// toString InfoNtf
func (msg JPVoiceMsg) toString() (string, error) {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
func (msg JPVoiceMsg) toMap() (map[string]interface{}, error) {
	return objectToMap(&msg)
}

//JPMessageList 消息列表
type JPMessageList struct {
	Total    int64       `json:"total"`
	Cursor   string      `json:"cursor"`
	Count    int64       `json:"count"`
	Messages []JPMessage `json:"messages,omitempty"`
}
