package rpcEntity

import "fmt"

type RpcRequest struct {
	MethodName string
	Args       interface{}
}

func (rreq RpcRequest) ToString() string {
	return fmt.Sprintf("MethodName is %v, Args is %v", rreq.MethodName, rreq.Args)
}
