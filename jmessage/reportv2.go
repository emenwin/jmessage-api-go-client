package jmessage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//GetMessages 获取消息
//count （必填）每次查询的总条数 一次最多1000
//begin_time (必填) 记录开始时间 格式 yyyy-MM-dd HH:mm:ss 设置筛选条件大于等于begin time
//end_time (必填) 记录结束时间 格式 yyyy-MM-dd HH:mm:ss 设置筛选条件小于等于end time
//begin_time end_time 之间最大范围不得超过7天 cursor  当第一次请求后如果后面有数据，会返回一个cursor回来用这个获取接下来的消息 (cursor 有效时间是120s，过期后需要重新通过第一个请求获得cursor，重新遍历)
//查询的消息按发送时间升序排序
func (jclient *JMessageClient) GetMessages(count int, beginTime, endTime string, cursor string) (*JPMessageList, error) {

	body := map[string]interface{}{}
	if count > 0 {
		body["count"] = count
	}
	if beginTime != "" {
		body["begin_time"] = beginTime
	}
	if endTime != "" {
		body["end_time"] = endTime
	}

	if cursor != "" {
		body["cursor"] = cursor
	}
	res, err := jclient.request(JMESSAGE_REPORT_V2_URL+REPORT_MESSAGE, "GET", body)
	defer res.Body.Close()

	ibytes, err := ioutil.ReadAll(res.Body)
	if nil != err {
		return nil, err
	}

	if ShowDebug {
		fmt.Println("respone:", string(ibytes))
	}

	var jmessageList JPMessageList
	err = json.Unmarshal(ibytes, &jmessageList)
	if nil != err {
		return nil, err
	}
	return &jmessageList, nil
}

func (jclient *JMessageClient) GetUserMessages(userName string, count int, beginTime, endTime string, cursor string) (*JPMessageList, error) {

	body := map[string]interface{}{}
	if count > 0 {
		body["count"] = count
	}
	if beginTime != "" {
		body["begin_time"] = beginTime
	}
	if endTime != "" {
		body["end_time"] = endTime
	}

	if cursor != "" {
		body["cursor"] = cursor
	}
	res, err := jclient.request(JMESSAGE_REPORT_V2_URL+"/"+userName+REPORT_MESSAGE, "GET", body)
	defer res.Body.Close()

	ibytes, err := ioutil.ReadAll(res.Body)
	if nil != err {
		return nil, err
	}

	if ShowDebug {
		fmt.Println("respone:", string(ibytes))
	}

	var jmessageList JPMessageList
	err = json.Unmarshal(ibytes, &jmessageList)
	if nil != err {
		return nil, err
	}
	return &jmessageList, nil
}
