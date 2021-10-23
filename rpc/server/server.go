package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/jackma8ge8/pine/application/config"
	"github.com/jackma8ge8/pine/connector"
	"github.com/jackma8ge8/pine/rpc"
	"github.com/jackma8ge8/pine/rpc/context"
	"github.com/jackma8ge8/pine/rpc/handler/clienthandler"
	"github.com/jackma8ge8/pine/rpc/handler/serverhandler"
	"github.com/jackma8ge8/pine/rpc/message"
	"github.com/jackma8ge8/pine/rpc/zookeeper"
	"github.com/jackma8ge8/pine/service/compressservice"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

// Start rpc server
func Start() {
	registerProtoHandler()
	// 注册到zookeeper
	go registToZk()

	// 获取服务器配置
	serverConfig := config.GetServerConfig()
	// RPC server启动
	logrus.Info("Rpc server started ws://" + serverConfig.Host + ":" + fmt.Sprint(serverConfig.Port))
	http.HandleFunc("/rpc", webSocketHandler)

	// 对客户端暴露的ws接口
	if serverConfig.IsConnector {
		http.HandleFunc("/", connector.WebSocketHandler)
	}
	// 开启并监听
	err := http.ListenAndServe(":"+fmt.Sprint(serverConfig.Port), nil)
	logrus.Error("Rpc server start fail: ", err.Error())
}

// WebSocketHandler deal with ws request
func webSocketHandler(w http.ResponseWriter, r *http.Request) {

	// 建立连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Error("连接失败", err.Error())
		return
	}

	// 断开连接自动清除连接信息
	conn.SetCloseHandler(func(code int, text string) error {
		conn.Close()
		return nil
	})

	// 用户认证
	token := r.URL.Query().Get("token")

	// token校验
	if token != config.GetServerConfig().Token {
		logrus.Error("用户校验失败!!!")
		conn.CloseHandler()(0, "认证失败")
		return
	}
	connLock := &sync.Mutex{}
	// 开始接收消息
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			conn.CloseHandler()(0, err.Error())
			break
		}
		// 解析消息
		rpcMsg := &message.RPCMsg{}
		err = proto.Unmarshal(data, rpcMsg)

		rpcCtx := context.GenRespCtx(conn, rpcMsg, connLock)

		if err != nil {
			logrus.Error(err)
			continue
		}

		if rpcMsg.Type == message.RemoterTypeEnum.HANDLER {
			ok := clienthandler.Manager.Exec(rpcCtx)
			if !ok {
				if rpcCtx.GetRequestID() == 0 {
					logrus.Warn(fmt.Sprintf("NotifyHandler(%v)不存在", rpcCtx.GetHandler()))
				} else {
					logrus.Warn(fmt.Sprintf("Handler(%v)不存在", rpcCtx.GetHandler()))
				}
			}

		} else if rpcMsg.Type == message.RemoterTypeEnum.REMOTER {
			ok := serverhandler.Manager.Exec(rpcCtx)
			if !ok {
				if rpcCtx.GetRequestID() == 0 {
					logrus.Warn(fmt.Sprintf("NotifyRemoter(%v)不存在", rpcCtx.GetHandler()))
				} else {
					logrus.Warn(fmt.Sprintf("Remoter(%v)不存在", rpcCtx.GetHandler()))
				}
			}
		} else {
			logrus.Panic("无效的消息类型")
		}
	}
}

// 注册到zookeeper
func registToZk() {
	zookeeper.Start()
}

func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func registerProtoHandler() {
	var serverProtoCentent []byte

	// 获取数据压缩元数据
	clienthandler.Manager.Register("__CompressMetadata__", func(rpcCtx *context.RPCCtx, hash string) {
		pwd, _ := os.Getwd()

		serverProto := path.Join(pwd, "/proto/server.proto")

		var result = map[string]interface{}{}

		// server proto
		if serverProtoCentent == nil && checkFileIsExist(serverProto) {
			var err error
			serverProtoCentent, err = ioutil.ReadFile(serverProto)

			if err != nil {
				logrus.Error(err)
				return
			}
		}
		result["proto"] = string(serverProtoCentent)

		// handlers
		handlers := compressservice.Handler.GetHandlers()
		result["handlers"] = handlers

		// events
		result["events"] = compressservice.Event.GetEvents()

		// serverKind
		result["serverKind"] = config.GetServerConfig().Kind

		rpcMsg := &message.RPCMsg{
			Handler: connector.SysHandlerMap.ServerCode,
		}

		rpc.Request.ToServer(rpcCtx.From, rpcMsg, func(serverCode byte) {
			// serverCode
			result["serverCode"] = serverCode
			rpcCtx.Response(result)
		})
	})
}
