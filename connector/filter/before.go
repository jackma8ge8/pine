package filter

import "github.com/jackma8ge8/pine/rpc/context"

// BeforeFilterSlice map
type BeforeFilterSlice []func(rpcCtx *context.RPCCtx) (next bool)

// Register filter
func (slice *BeforeFilterSlice) Register(f func(rpcCtx *context.RPCCtx) (next bool)) {
	*slice = append(*slice, f)
}

// Exec filter(返回true标识继续往下执行)
func (slice BeforeFilterSlice) Exec(rpcCtx *context.RPCCtx) (next bool) {
	for _, f := range slice {
		next = f(rpcCtx)
		if !next {
			return false
		}
	}
	return true
}

// Before filter
var Before = make(BeforeFilterSlice, 0)
