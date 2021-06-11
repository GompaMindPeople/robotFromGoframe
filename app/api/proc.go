package api

import (
	"github.com/gogf/gf/net/ghttp"
	"strings"
)


var Proc = proc{}
type proc struct {}


var Data = make([]string,10)

//var heartbeat = make([]string,10)

//  主要用来接收 机器人上传上来的消息
func (*proc) Index(r *ghttp.Request) {
	reqBody := r.GetFormMap()
	// 数据是心跳包,不做处理
	if strings.Index(reqBody["meta_event_type"].(string),"heartbeat") != -1 {
		//heartbeat = append(heartbeat,reqBody)
		r.Response.Writeln("Hello proc!")
	}


}



func  (*proc) Message(r *ghttp.Request){
	r.Response.Writeln(strings.Join(Data,""))
	Data = make([]string,10)
}
