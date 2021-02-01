package filter

import "github.com/jackma8ge8/pine/rpc/message"

// AfterFilterSlice map
type AfterFilterSlice []func(rpcResp *message.PineMsg) (next bool)

// Register filter
func (slice *AfterFilterSlice) Register(f func(rpcResp *message.PineMsg) (next bool)) {
	*slice = append(*slice, f)
}

// Exec filter
func (slice AfterFilterSlice) Exec(rpcResp *message.PineMsg) (next bool) {
	for _, f := range slice {
		next = f(rpcResp)
		if !next {
			return false
		}
	}
	return true
}

// After filter
var After = make(AfterFilterSlice, 0)
