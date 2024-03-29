package api

import (
	"fmt"
	"github.com/go-ll/ginx/examples/mmo_game/core"
	"github.com/go-ll/ginx/examples/mmo_game/pb"
	"github.com/go-ll/ginx/giface"
	"github.com/go-ll/ginx/gnet"
	"github.com/golang/protobuf/proto"
)

type WorldChatApi struct {
	gnet.BaseRouter
}

func (*WorldChatApi) Handle(request giface.Request) {
	//1. 将客户端传来的proto协议解码
	msg := &pb.Talk{}
	err := proto.Unmarshal(request.GetData(), msg)
	if err != nil {
		fmt.Println("Talk Unmarshal error ", err)
		return
	}

	//2. 得知当前的消息是从哪个玩家传递来的,从连接属性pID中获取
	pID, err := request.GetConn().GetProperty("pID")
	if err != nil {
		fmt.Println("GetProperty pID error", err)
		request.GetConn().Stop()
		return
	}
	//3. 根据pID得到player对象
	player := core.WorldMgrObj.GetPlayerByPID(pID.(int32))

	//4. 让player对象发起聊天广播请求
	player.Talk(msg.Content)
}

func (*WorldChatApi) PreHandle(request giface.Request) {}

func (*WorldChatApi) PostHandle(request giface.Request) {}
