package rpcEntity

import "EntryTask/constant"

type RpcResponse struct {
	ErrCode constant.ErrCode
	Data    interface{}
}
