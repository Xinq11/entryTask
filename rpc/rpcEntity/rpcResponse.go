package rpcEntity

import (
	"EntryTask/constant"
	"fmt"
)

type RpcResponse struct {
	ErrCode constant.ErrCode
	Data    interface{}
}

func (rres RpcResponse) ToString() string {
	return fmt.Sprintf("ErrCode is %v, Data is %v", rres.ErrCode, rres.Data)
}
