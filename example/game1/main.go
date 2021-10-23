package main

import (
	"errors"
	"fmt"
	"game1/handlermessage"
	_ "net/http/pprof"
	"strconv"
	"time"

	"github.com/jackma8ge8/pine/application"
	"github.com/jackma8ge8/pine/rpc"
	"github.com/jackma8ge8/pine/rpc/client"
	"github.com/jackma8ge8/pine/rpc/context"
	"github.com/jackma8ge8/pine/rpc/message"
	"github.com/jackma8ge8/pine/service/channelservice"
	"github.com/jackma8ge8/pine/service/compressservice"
	"github.com/jackma8ge8/pine/service/sessionservice"
	"github.com/sirupsen/logrus"
)

func main() {

	app := application.CreateApp()

	compressservice.Event.AddRecords("onMsg", "onMsgJSON") // 需要压缩的Event

	app.AsConnector(func(uid string, token string, sessionData map[string]string) error {

		if uid == "" || token == "" {
			return errors.New("Invalid token")
		}
		sessionData[token] = token

		return nil
	})

	app.RegisteHandler("handler", func(rpcCtx *context.RPCCtx, data *handlermessage.Handler) {

		channelservice.PushMessageBySession(rpcCtx.Session, "onMsg", &handlermessage.OnMsg{
			Name: "From onMsg",
			Data: "哈哈哈哈哈",
		})

		// logrus.Warn(data)

		handlerResp := &handlermessage.HandlerResp{
			Name: "HandlerResp",
		}

		rpcCtx.Response(handlerResp)

	})

	app.RegisteHandler("handlerJSON", func(rpcCtx *context.RPCCtx, data map[string]interface{}) {

		// 直接通过session推送消息
		channelservice.PushMessageBySession(rpcCtx.Session, "onMsg1", "hahahah")

		// 广播给所有人
		channelservice.BroadCast("onMsg2", "==========广播广播广播广播广播==========")

		// 创建channel。通过channel推送消息
		channel := channelservice.CreateChannel("101")
		channel.Add(rpcCtx.Session.UID, rpcCtx.Session)

		// 推送给所有在当前channel中的玩家
		channel.PushMessage("onMsg1", map[string]string{
			"Name": "onMsg",
			"Data": "啊哈哈傻法师打上发发",
		})
		// 推送给除了切片内的channel中的玩家
		channel.PushMessageToOthers([]string{rpcCtx.Session.UID}, "onMsg2", map[string]string{
			"Name": "onMsg",
			"Data": "啊哈哈傻法师打上发发",
		})
		// 只推送给当前玩家
		channel.PushMessageToUser(rpcCtx.Session.UID, "onMsg3", map[string]string{
			"Name": "onMsg",
			"Data": "啊哈哈傻法师打上发发",
		})
		// 只推送给切片的指定的玩家
		channel.PushMessageToUsers([]string{rpcCtx.Session.UID}, "onMsg4", map[string]string{
			"Name": "onMsg",
			"Data": "啊哈哈傻法师打上发发",
		})

		rpcMsg := &message.RPCMsg{
			Handler: "getOneRobot",
			RawData: []byte{},
		}

		rpc.Request.ByKind("connector", rpcMsg, func(data map[string]interface{}) {
			logrus.Info("收到Rpc的回复：", fmt.Sprint(data))
		})

		rpcCtx.Response(map[string]interface{}{
			"Route":     "onMsgJSON",
			"heiheihie": "heiheihei",
		})
	})

	app.RegisteRemoter("getOneRobot", func(rpcCtx *context.RPCCtx, data interface{}) {

		rpcCtx.Response(map[string]interface{}{
			"name": "盖伦",
			"sex":  1,
			"age":  18,
		})
	})

	app.RegisteHandlerBeforeFilter(func(rpcCtx *context.RPCCtx) (next bool) {

		if rpcCtx.GetHandler() == "enterRoom" {
			lastEnterRoomTimeInterface := rpcCtx.Session.Data["lastEnterRoomTime"]
			if lastEnterRoomTimeInterface != "" {
				timestamp, e := strconv.ParseInt(lastEnterRoomTimeInterface, 10, 64)
				if e != nil {
					logrus.Error("不能将", lastEnterRoomTimeInterface, "转换成时间戳")
				} else if time.Now().Sub(time.Unix(timestamp, 0)) < time.Second {
					logrus.Error("操作太频繁") // 返回结果
					return false          // 停止执行下个before filter以及hanler
				}
			}

			rpcCtx.Session.Set("lastEnterRoomTime", fmt.Sprint(time.Now().Unix()))
			sessionservice.UpdateSession(rpcCtx.Session, "lastEnterRoomTime")
		}
		return true // 继续执行下个before filter直到执行handler
	})

	app.RegisteHandlerAfterFilter(func(rpcResp *message.PineMsg) (next bool) {
		return true // return true继续执行下个after filter, return false停止执行下个after filter
	})

	app.RegisteRouter("ddz", func(rpcMsg *message.RPCMsg, clients []*client.RPCClient) *client.RPCClient {

		for _, clientInfo := range clients {
			if chatServerID, ok := rpcMsg.Session.Get("chatServerID").(string); ok && clientInfo.ServerConfig.ID == chatServerID {
				return clientInfo
			}
		}

		return nil //if return nil, pine will get one rpc client by random
	})

	app.Start()
}
