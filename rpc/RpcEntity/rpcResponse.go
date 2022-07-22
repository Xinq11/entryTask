package RpcEntity

import "EntryTask/constant"

type RpcResponse struct {
	Err_code constant.ErrCode
	Data     interface{}
}
