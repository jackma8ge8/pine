package application

import (
	"github.com/jackma8ge8/pine/rpc/handler/serverhandler"
	"github.com/jackma8ge8/pine/rpc/handler/clienthandler"
)

// RegisteHandler 注册Handler
func (app Application) RegisteHandler(name string, f interface{}) {
	clienthandler.Manager.Register(name, f)
}

// RegisteRemoter 注册RPC Handler
func (app Application) RegisteRemoter(name string, f interface{}) {
	serverhandler.Manager.Register(name, f)
}
