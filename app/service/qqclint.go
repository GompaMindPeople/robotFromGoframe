package service

import (
	"github.com/gogf/gf/os/glog"
	"net/url"
	"robotFromGoframe/app/model"
)

var HttpClient = model.New(nil)

type QQClient struct {
	Host string
}

//终结点：/send_msg
//
//参数
//
//字段名	数据类型	默认值	说明
//message_type	string	-	消息类型, 支持 private、group , 分别对应私聊、群组, 如不传入, 则根据传入的 *_id 参数判断
//user_id	int64	-	对方 QQ 号 ( 消息类型为 private 时需要 )
//group_id	int64	-	群号 ( 消息类型为 group 时需要 )
//message	message	-	要发送的内容
//auto_escape	boolean	false	消息内容是否作为纯文本发送 ( 即不解析 CQ 码 ) , 只在 message 字段是字符串时有效
//响应数据
//
//字段名	数据类型	说明
//message_id	int32	消息 ID
func (qq *QQClient) SendMessage(messageType, message, userId, groupId string, autoEscape bool) model.ResponseData {

	//encoding := base64.URLEncoding.EncodeToString([]byte(message))
	values := url.Values{}
	values.Add("message", message)
	values.Add("message_type", messageType)
	values.Add("group_id", groupId)
	values.Add("user_id", userId)
	//data := ""
	//data += "message_type="+messageType
	//data += "&user_id="+userId
	//data +=  "&group_id" + groupId
	//data += "&message=" + message
	//data += values.Encode()
	if autoEscape {
		//data += "&auto_escape=true"
		values.Add("auto_escape", "true")
	} else {
		//data += "&auto_escape=false"
		values.Add("auto_escape", "false")
	}

	parse, err := url.Parse(qq.Host + "/send_msg?" + values.Encode())
	if err != nil {
		return model.ResponseData{Code: 500, Data: "变更!" + err.Error()}
	}
	s := parse.String()
	response, err := HttpClient.GET(s, "")
	if err != nil {
		glog.Error("发生错误-->", err)
		return model.ResponseData{Code: 500, Data: "发送消息是发生错误!" + err.Error()}
	}
	return model.ResponseData{Code: 200, Data: response.Data, OriginalData: response}
}
