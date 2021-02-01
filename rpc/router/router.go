package router

import (
	"github.com/jackma8ge8/pine/rpc/client"
	"github.com/jackma8ge8/pine/rpc/message"
)

// Map 存储自定义的路由
type Map map[string]func(rpcMsg *message.RPCMsg, clients []*client.RPCClient) *client.RPCClient

// Register 注册一个路由函数
func (routeMap Map) Register(serverKind string, route func(rpcMsg *message.RPCMsg, clients []*client.RPCClient) *client.RPCClient) {
	routeMap[serverKind] = route
}

// Get 获取一个路由函数
func (routeMap Map) Get(serverKind string) func(rpcMsg *message.RPCMsg, clients []*client.RPCClient) *client.RPCClient {
	return routeMap[serverKind]
}

// Manager 管理router
var Manager = make(Map)
