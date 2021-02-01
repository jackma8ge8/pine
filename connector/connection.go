package connector

import (
	"fmt"
	"strings"
	"sync"

	"github.com/jackma8ge8/pine/application/config"
	"github.com/jackma8ge8/pine/connector/filter"
	"github.com/jackma8ge8/pine/rpc"
	"github.com/jackma8ge8/pine/rpc/client/clientmanager"
	"github.com/jackma8ge8/pine/rpc/context"
	"github.com/jackma8ge8/pine/rpc/message"
	"github.com/jackma8ge8/pine/rpc/session"
	"github.com/jackma8ge8/pine/serializer"
	"github.com/jackma8ge8/pine/service/compressservice"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// Connection 用户连接信息
type Connection struct {
	uid            string
	conn           *websocket.Conn
	data           map[string]string
	routeRecord    map[string]string
	compressRecord map[string]bool
	mutex          sync.Mutex
}

// Get 从session中查找一个值
func (connection *Connection) Get(key string) string {
	return connection.data[key]
}

// Set 往session中设置一个键值对
func (connection *Connection) Set(key string, v string) {
	connection.data[key] = v
}

// 回复request
func (connection *Connection) response(pineMsg *message.PineMsg) {
	connection.mutex.Lock()
	defer connection.mutex.Unlock()
	err := connection.conn.WriteMessage(message.TypeEnum.BinaryMessage, serializer.ToBytes(pineMsg))

	if err != nil {
		logrus.Error(err)
	}
}

// 主动推送消息
func (connection *Connection) notify(notify *message.PineMsg) {

	connection.mutex.Lock()
	defer connection.mutex.Unlock()

	err := connection.conn.WriteMessage(message.TypeEnum.BinaryMessage, serializer.ToBytes(notify))

	if err != nil {
		logrus.Error(err)
	}
}

// GetSession 获取session
func (connection *Connection) GetSession() *session.Session {
	session := &session.Session{
		UID:  connection.uid,
		CID:  config.GetServerConfig().ID,
		Data: connection.data,
	}
	return session
}

// StartReceiveMsg 开始接收消息
func (connection *Connection) StartReceiveMsg() {

	registerConnectorHandler()

	conn := connection.conn
	connLock := &sync.Mutex{}
	// 开始接收消息
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			conn.CloseHandler()(0, err.Error())
			break
		}
		// 解析消息
		clientMessage := &message.PineMsg{}

		err = proto.Unmarshal(data, clientMessage)

		if err != nil {
			logrus.Error("消息解析失败", err, "Data", data)
			continue
		}

		if clientMessage.Route == "" {
			logrus.Error("Route不能为空", err, "Data", clientMessage)
			continue
		}

		var serverKind string
		var handler string

		routeBytes := []byte(clientMessage.Route)

		if len(routeBytes) == 2 {
			serverKind = compressservice.Server.GetKindByCode(routeBytes[0])
			handler = string(routeBytes[1])
		} else {
			handlerInfos := strings.Split(clientMessage.Route, ".")
			serverKind = handlerInfos[0] // 解析出服务器类型
			handler = handlerInfos[1]    // 真正的handler
		}

		session := connection.GetSession()

		rpcMsg := &message.RPCMsg{
			From:      config.GetServerConfig().ID,
			Handler:   handler,
			Type:      message.RemoterTypeEnum.HANDLER,
			RequestID: clientMessage.RequestID,
			RawData:   clientMessage.Data,
			Session:   session,
		}

		// 获取RPCCLint
		rpcClient := clientmanager.GetClientByRouter(serverKind, rpcMsg, &connection.routeRecord)

		if rpcClient == nil {

			tip := fmt.Sprint("找不到任何", serverKind, "服务器")
			clientMessageResp := &message.PineMsg{
				Route:     clientMessage.Route,
				RequestID: clientMessage.RequestID,
				Data: serializer.ToBytes(&message.PineErrResp{
					Code:    500,
					Message: &tip,
				}),
			}

			connection.response(clientMessageResp)
			continue
		}

		rpcCtx := context.GenRespCtx(conn, rpcMsg, connLock)

		if !filter.Before.Exec(rpcCtx) {
			continue
		}

		if *clientMessage.RequestID == 0 {
			rpc.Notify.ToServer(rpcClient.ServerConfig.ID, rpcMsg)
		} else {
			// 转发Request
			rpc.Request.ToServer(rpcClient.ServerConfig.ID, rpcMsg, func(data []byte) {

				pineMsg := &message.PineMsg{
					RequestID: clientMessage.RequestID,
					Route:     clientMessage.Route,
					Data:      data,
				}

				filter.After.Exec(pineMsg)
				connection.response(pineMsg)
			})
		}
	}
}

// Kick this coonnection
func (connection *Connection) Kick(data []byte) {
	notify := &message.PineMsg{
		Route: string([]byte{
			compressservice.Server.GetCodeByKind(config.GetServerConfig().Kind),
			compressservice.Event.GetCodeByEvent(SysHandlerMap.Kick)}),
		Data: data,
	}
	connection.notify(notify)
	DelConnection(connection.uid)
	connection.conn.Close()
}

func init() {
	compressservice.Event.AddEventRecord(SysHandlerMap.Kick)
}
