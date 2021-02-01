package application

import (
	"github.com/jackma8ge8/pine/connector/filter"
	"github.com/jackma8ge8/pine/rpc/context"
	"github.com/jackma8ge8/pine/rpc/message"
)

// RegisteHandlerBeforeFilter 注册before filter
func (app Application) RegisteHandlerBeforeFilter(f func(rpcCtx *context.RPCCtx) (next bool)) {
	filter.Before.Register(f)
}

// RegisteHandlerAfterFilter 注册after filter
func (app Application) RegisteHandlerAfterFilter(f func(rpcResp *message.PineMsg) (next bool)) {
	filter.After.Register(f)
}
