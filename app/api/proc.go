package api

import (
	"bytes"
	"encoding/json"
	"github.com/gogf/gf/container/gqueue"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"robotFromGoframe/app/model"
	"robotFromGoframe/app/model/qc"
	"robotFromGoframe/app/service"
	"strconv"
	"strings"
	"sync"
	"time"
)

var Proc = proc{}

type proc struct{}

var DataQueue = gqueue.New(65535)

var HeartbeatQueue = gqueue.New(1000)

var waitGroup = sync.WaitGroup{}

var path = ""

// 指令映射表
var mapping = map[string]func(content *qc.Content){
	"help":      help,
	"quake":     quake,
	"dirsearch": dirSearch,
	"oneForall": oneForall,
}

// 机器人的通信实例.
var qqClient = service.QQClient{Host: "http://127.0.0.1:5700"}

func init() {
	//  获取执行文件的绝对路径
	path, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	go func() {
		// 每过3分钟 就清空心跳包的数据...感觉心跳包暂时没用..
		time.Sleep(3 * time.Minute)
		for i := 0; i < HeartbeatQueue.Size(); i++ {
			HeartbeatQueue.Pop()
		}
	}()
	//处理消息
	//协程处理 消息
	go func() {
		ProcessMessage()
	}()

}

//  主要用来接收 机器人上传上来的消息
func (*proc) Index(r *ghttp.Request) {
	reqBody := r.GetBodyString()

	// 数据是心跳包,不做处理
	if strings.Index(reqBody, "meta_event_type") != -1 {
		if strings.Index(reqBody, "heartbeat") != -1 {
			HeartbeatQueue.Push(reqBody)
		}
	} else {

		result := &model.FriendPrivateMessage{}
		//非心跳包数据尝试转换成好友消息.
		err := json.Unmarshal([]byte(reqBody), result)
		if err != nil {
			glog.Error("数据转换成好友消息时报错")
			r.Response.Writeln("数据转换成好友消息时报错")
		} else {
			// 压队列
			DataQueue.Push(result)
		}

	}
	r.Response.Writeln("ok")

}

// 测试路由
func (*proc) Message(r *ghttp.Request) {
	message := SendPrivateMessage("xiaoxi", "894799178")
	r.Response.Writeln(message.Data)
}

// 处理队列的消息....最主要的函数
func ProcessMessage() {
	for {
		if DataQueue.Size() > 0 {
			pop := DataQueue.Pop().(*model.FriendPrivateMessage)
			//if pop.UserId != 894799178{
			//	continue
			//}
			message := pop.Message
			if message == "" {
				continue
			}
			s := message[0:1]
			//第一个字符是# 表示 该消息是指令
			if s == "#" {
				for k, v := range mapping {
					if strings.Index(message, k) != -1 {
						v(&qc.Content{Param: pop})
					}
				}
			}

		}
	}
}

//发生私聊消息..
func SendPrivateMessage(message, qqNumber string) model.ResponseData {
	return qqClient.SendMessage("private", message, qqNumber, "", false)
}

// dirsearch工具的指令调用实现
func dirSearch(content *qc.Content) {
	fm := content.Param.(*model.FriendPrivateMessage)
	message := parseMessage(fm.Message[len("#dirsearch "):])
	message[1] = path + "/tools/dirsearch-master/" + message[1]
	execute(fm, message...)

}

// oneForall指令
func oneForall(content *qc.Content) {
	fm := content.Param.(*model.FriendPrivateMessage)

	message := parseMessage(fm.Message[len("#oneForall "):])
	message[1] = path + "/tools/OneForAll-master/" + message[1]
	execute(fm, message...)
}

// quake 指令
func quake(content *qc.Content) {
	fm := content.Param.(*model.FriendPrivateMessage)
	replace := strings.Replace(fm.Message, "#quake ", "", -1)
	message := parseMessage(replace)
	newMessage := make([]string, 0)
	newMessage = append(newMessage, path+"/tools/quake/InformationGatheringTool")
	for _, v := range message {
		newMessage = append(newMessage, v)
	}
	execute(fm, newMessage...)
	return
}

// 帮助的函数入口
func help(content *qc.Content) {

	help1 := "quake---360网络空间资产测绘工具,指令演示:\n #quake -t quake -c \"search domain=baidu.com\"\n\n" +
		"OnForAll子域名收集工具,指令演示:\n #oneForall python Oneforall.py --target baidu.com run\n\n" +
		"dirsearch目录遍历工具,指令演示:\n #dirsearch python dirsearch.py -u baidu.com\n\n"
	id := content.Param.(*model.FriendPrivateMessage).UserId
	SendPrivateMessage(help1, strconv.FormatInt(id, 10))
}

// 执行本地指令的函数
func execute(fm *model.FriendPrivateMessage, cmd ...string) {
	SendPrivateMessage("开始执行指令!-->>", cmd[0])
	stderr := &bytes.Buffer{}
	QQNumber := strconv.FormatInt(fm.UserId, 10)
	msgChan := make(chan int)
	// 需要对字符串通过空格进行分割一下.
	var command *exec.Cmd
	command = exec.Command(cmd[0], cmd[1:]...)
	command.Stdout = os.Stdout
	command.Stderr = stderr
	go func(mChan chan int, wait *sync.WaitGroup) {
	A:
		for true {
			time.Sleep(time.Second * 2)
			errString := stderr.String()
			if errString != "" {
				SendPrivateMessage(errString, QQNumber)
			}
			select {
			case <-mChan:
				break A
			default:
				stderr.Reset()
				continue
			}
		}
	}(msgChan, &waitGroup)
	if errC := command.Run(); errC != nil {
		log.Println("发生错误->", errC.Error())
		return
	}
	msgChan <- 1
	SendPrivateMessage(cmd[0]+"--->执行完成!", QQNumber)
}

func parseMessage(msg string) []string {
	index := strings.Index(msg, "\"")
	if index == -1 {
		return strings.Split(msg, " ")
	}
	left := msg[:index-1]
	right := msg[index+1 : len(msg)-1]
	split := strings.Split(left, " ")

	return append(split, right)
}
