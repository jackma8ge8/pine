package channelservice

import (
	"github.com/jackma8ge8/pine/application/config"
	"github.com/jackma8ge8/pine/connector"
	"github.com/jackma8ge8/pine/rpc"
	"github.com/jackma8ge8/pine/rpc/message"
	"github.com/jackma8ge8/pine/rpc/session"
	"github.com/jackma8ge8/pine/serializer"
	"github.com/jackma8ge8/pine/service/compressservice"
)

// Channel ChannelService
type Channel map[string]*session.Session

func compressEvent(event string) string {
	code := compressservice.Event.GetCodeByEvent(event)
	if code != 0 {
		return string(code)
	}
	return config.GetServerConfig().Kind + "." + event
}

// PushMessage 推送消息给所有人
func (channel Channel) PushMessage(event string, data interface{}) {

	for _, session := range channel {
		PushMessageBySession(session, event, data)
	}
}

// PushMessageToOthers 推送消息给其他人
func (channel Channel) PushMessageToOthers(uids []string, event string, data interface{}) {

	for uid := range channel {
		findIndex := -1
		for index, value := range uids {
			if uid == value {
				findIndex = index
				break
			}
		}
		if findIndex == -1 {
			channel.PushMessageToUser(uid, event, data)
		}
	}
}

// PushMessageToUsers 推送消息给指定玩家
func (channel Channel) PushMessageToUsers(uids []string, event string, data interface{}) {

	for _, uid := range uids {
		channel.PushMessageToUser(uid, event, data)
	}

}

// PushMessageToUser 推送消息给指定玩家
func (channel Channel) PushMessageToUser(uid string, event string, data interface{}) {

	session, ok := channel[uid]
	if !ok {
		return
	}

	PushMessageBySession(session, event, data)

}

// Add 推送消息给其他人
func (channel Channel) Add(uid string, session *session.Session) {
	channel[uid] = session
}

// PushMessageBySession 通过session推送消息
func PushMessageBySession(session *session.Session, event string, data interface{}) {

	notify := &message.PineMsg{
		Route: compressEvent(event),
		Data:  serializer.ToBytes(data),
	}

	rpcMsg := &message.RPCMsg{
		Handler: connector.SysHandlerMap.PushMessage,
		RawData: serializer.ToBytes(notify),
		Session: session,
	}
	rpc.Notify.ToServer(session.CID, rpcMsg)
}

// BroadCast 广播
func BroadCast(event string, data interface{}) {

	pineMsg := &message.PineMsg{
		Route: compressEvent(event),
		Data:  serializer.ToBytes(data),
	}

	rpcMsg := &message.RPCMsg{
		Handler: connector.SysHandlerMap.BroadCast,
		RawData: serializer.ToBytes(pineMsg),
	}

	rpc.BroadCast(rpcMsg)
}
