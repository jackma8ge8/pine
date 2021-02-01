package serverhandler

import (
	"github.com/jackma8ge8/pine/rpc/handler"
)

// ServerHandler ServerHandler
type ServerHandler struct {
	*handler.Handler
}

// Manager return RPCHandler
var Manager = &ServerHandler{
	Handler: &handler.Handler{
		Map: make(handler.Map),
	},
}
