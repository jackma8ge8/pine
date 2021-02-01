package zookeeper

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/jackma8ge8/pine/application/config"
	"github.com/jackma8ge8/pine/rpc/client/clientmanager"
	"github.com/jackma8ge8/pine/serializer"
	"github.com/jackma8ge8/pine/service/compressservice"
	"github.com/sirupsen/logrus"

	"github.com/samuel/go-zookeeper/zk"
)

var zkClient *ZkClient

// zkSessionTimeout Session timeout of zookeeper connection
var zkSessionTimeout = time.Second * 3

func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

// Start zookeeper
func Start() {

	// 读取配置文件
	zkConfig := config.GetZkConfig()

	// 建立连接
	conn, _, err := zk.Connect([]string{zkConfig.Host + ":" + fmt.Sprint(zkConfig.Port)}, zkSessionTimeout)
	zkClient = &ZkClient{conn: conn}
	if err != nil {
		logrus.Panic(err)
	}

	// 初始化节点
	initNode()

	// 监听节点变化
	watch()
}

// 初始化节点
func initNode() {

	// 服务器配置
	serverConfig := config.GetServerConfig()

	// 检查根节点是否存在，不存在则创建
	rootPath := fmt.Sprint("/", serverConfig.SystemName)

	if !zkClient.exists(rootPath) {
		zkClient.create(rootPath, []byte{}, 0, zk.WorldACL(zk.PermAll))
	}

	// 检查服务器数据节点是否存在，不存在则创建
	nodePath := fmt.Sprint(rootPath, "/", serverConfig.ID)

	tryTimes := 0
	// 最大尝试次数
	maxTryTimes := int(50 + zkSessionTimeout/100/time.Millisecond)

	for tryTimes = 0; tryTimes < maxTryTimes; tryTimes++ {
		// 不存在则跳出循环，创建节点
		if !zkClient.exists(nodePath) {
			break
		}
		// node 存在则休眠100毫秒
		time.Sleep(time.Millisecond * 100)
	}

	if tryTimes >= maxTryTimes {
		// 操过最大尝试次数则报错
		logrus.Panic(fmt.Sprint("Duplicated server."))
	}

	zkClient.create(nodePath, serializer.ToBytes(serverConfig), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	zkClient.serverID = serverConfig.ID
	logrus.Info("Node created:", nodePath)

	go recreatedNode(nodePath, serverConfig)
}

func watch() {
	// 服务器配置
	serverConfig := config.GetServerConfig()
	zkpath := fmt.Sprint("/", serverConfig.SystemName)
	for {
		// 遍历所有的serverID
		serverIDs, _, eventChan, err := zkClient.conn.ChildrenW(zkpath)
		if err != nil {
			// if err == zk.ErrNoServer {
			// 	time.Sleep(time.Millisecond * 100)
			// 	continue
			// } else if err == zk.ErrConnectionClosed {
			// 	return
			// } else {
			logrus.Error(err)
			return
			// }
		}
		// 监听每个server的情况
		for _, serverID := range serverIDs {

			if clientmanager.GetClientByID(serverID) != nil {
				// 如果已经建立果监听则跳过
				continue
			}

			func(serverID string) {

				for i := 0; i < 30; i++ {
					zkpath := fmt.Sprint(zkpath, "/", serverID)
					isExists, _, err := zkClient.conn.Exists(zkpath)
					if err != nil {
						logrus.Error(err)
						continue
					}

					if !isExists {
						time.Sleep(time.Millisecond * 100)
						continue
					}

					// 监听服务器变化
					data, _, err := zkClient.conn.Get(zkpath)
					if err != nil {
						clientmanager.DelClientByID(serverID)
						logrus.Error(err)
						continue
					}
					// 解析服务器信息
					serverConfig := &config.ServerConfig{}
					err = json.Unmarshal(data, serverConfig)
					if err != nil {
						logrus.Error(err)
						continue
					}
					// 创建客户端，并与该服务器连接
					clientmanager.CreateClient(serverConfig, zkSessionTimeout)

					if config.GetServerConfig().IsConnector {
						compressservice.Server.AddRecord(serverConfig.Kind)
					}

					break
				}
			}(serverID)
		}
		// 没有新事件，则阻塞
		<-eventChan
	}
}

var isCreatingNode = false

func recreatedNode(nodePath string, serverConfig *config.ServerConfig) {

	defer func() {
		isCreatingNode = false
	}()

	if isCreatingNode {
		return
	}

	isCreatingNode = true

	for {
		ok, _, err := zkClient.conn.Exists(nodePath)
		if err != nil {
			if err == zk.ErrNoServer || err == zk.ErrConnectionClosed || err == zk.ErrSessionExpired {
				zkClient.conn.Close()
				go Start()
				return
			}
			logrus.Panic(err)
		}

		if ok {
			time.Sleep(time.Second * 3)
		} else if zkClient.conn.State() == zk.StateConnected {
			zkClient.create(nodePath, serializer.ToBytes(serverConfig), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
			logrus.Info("Node recreated:", nodePath)
		} else {
			time.Sleep(time.Millisecond * 100)
		}
	}
}
