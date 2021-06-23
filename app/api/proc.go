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
	"time"
)

var Proc = proc{}

type proc struct{}

var DataQueue = gqueue.New(65535)

var HeartbeatQueue = gqueue.New(1000)

var path = ""

var mapping = map[string]func(content *qc.Content){
	"help":      help,
	"quake":     quake,
	"dirsearch": dirSearch,
}

var qqClient = service.QQClient{Host: "http://127.0.0.1:5700"}

func init() {
	path, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	go func() {
		// 每过3分钟 就清空心跳包的数据...感觉心跳包暂时没用..
		time.Sleep(3 * time.Minute)
		for i := 0; i < HeartbeatQueue.Size(); i++ {
			HeartbeatQueue.Pop()
		}
	}()
	//处理消息
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
		err := json.Unmarshal([]byte(reqBody), result)
		if err != nil {
			glog.Error("数据转换成好友消息时报错")
			r.Response.Writeln("数据转换成好友消息时报错")
		} else {
			DataQueue.Push(result)
		}

	}
	r.Response.Writeln("Hello proc!")

}

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
			if s == "#" {
				for k, v := range mapping {
					if strings.Index(message, k) != -1 {
						v(&qc.Content{Param: pop})
					}
				}
				//SendPrivateMessage("----开始处理请求---",strconv.FormatInt(pop.UserId,10))
				//strings.Index(s,)
			}

		}
	}
}

func SendPrivateMessage(message, qqNumber string) model.ResponseData {
	return qqClient.SendMessage("private", message, qqNumber, "", false)
}

func dirSearch(content *qc.Content) {

	fm := content.Param.(*model.FriendPrivateMessage)

	//message := parseMessage(fm.Message[len("#dirsearch"):])

	cmd := []string{"cd tools/dirsearch-master", fm.Message[len("#dirsearch")+1:]}

	execute(fm, cmd...)

}

// quake 的主要函数 入口
func quake(content *qc.Content) {
	fm := content.Param.(*model.FriendPrivateMessage)

	replace := strings.Replace(fm.Message, "#quake ", "", -1)

	parseMessage(replace)
	execute(fm, "/tools/quake/InformationGatheringTool")

	return
}

// 帮助的函数入口
func help(content *qc.Content) {

	help1 := "quake---360网络空间资产测绘工具\n"
	id := content.Param.(*model.FriendPrivateMessage).UserId
	SendPrivateMessage(help1, strconv.FormatInt(id, 10))
}

// 执行本地指令的函数
func execute(fm *model.FriendPrivateMessage, cmd ...string) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	QQNumber := strconv.FormatInt(fm.UserId, 10)
	msgChan := make(chan int)
	// 需要对字符串通过空格进行分割一下.
	command := exec.Command("dir")
	command.Stdout = stdout
	command.Stderr = stderr
	pipe, err2 := command.StdinPipe()

	if err2 != nil {
		log.Println("打开输入流管道失败." + err2.Error())
		return
	}
	if err := command.Start(); err != nil {
		log.Println(err)
		return
	}

	for _, v := range cmd {
		if _, err22 := pipe.Write([]byte(v)); err22 != nil {
			log.Println("写入命错误->", err22.Error())
			continue
		}

		go func(mChan chan int) {
		A:
			for true {
				time.Sleep(time.Second * 2)
				SendPrivateMessage(stdout.String(), QQNumber)
				errString := stderr.String()
				if errString != "" {
					time.Sleep(time.Second * 1)
					SendPrivateMessage(errString, QQNumber)
				}
				select {
				case <-mChan:
					break A
				default:
					stderr.Reset()
					stdout.Reset()
					continue
				}
			}
		}(msgChan)
		err := command.Wait()
		msgChan <- 1
		if err != nil {
			log.Println(err)
			return
		}
	}

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
