package client

import (
	"EntryTask/constant"
	"EntryTask/logger"
	"EntryTask/rpc/codec"
	"EntryTask/rpc/network"
	"EntryTask/rpc/pool"
	"EntryTask/rpc/rpcEntity"
)

type RpcClient struct {
	ConnPool *pool.Pool
}

var Client *RpcClient

func MakeClient(addr string) error {
	pool, err := pool.Init(addr)
	if err != nil {
		return err
	}
	Client = &RpcClient{ConnPool: pool}
	return nil
}

// 发起远程过程调用
func (client *RpcClient) Call(methodName string, args interface{}) rpcEntity.RpcResponse {
	conn, err := client.ConnPool.GetConn()
	if err != nil {
		logger.Error("rpcClient.Call GetConn error: " + err.Error())
		return rpcEntity.RpcResponse{ErrCode: constant.ServerError}
	}
	defer client.ConnPool.ReleaseConn(conn)
	request := rpcEntity.RpcRequest{
		MethodName: methodName,
		Args:       args,
	}
	// 编码
	encode, err := codec.Encode(request)
	if err != nil {
		logger.Error("rpcClient.Call encode error: " + err.Error())
		return rpcEntity.RpcResponse{ErrCode: constant.ServerError}
	}
	// 发送
	err = network.Send(conn, encode)
	if err != nil {
		logger.Error("rpcClient.Call send error: " + err.Error())
		return rpcEntity.RpcResponse{ErrCode: constant.ServerError}
	}
	// 接收
	read, err := network.Read(conn)
	if err != nil {
		logger.Error("rpcClient.Call read error: " + err.Error())
		return rpcEntity.RpcResponse{ErrCode: constant.ServerError}
	}
	// 解码
	res, err := codec.ResDecode(read)
	if err != nil {
		logger.Error("rpcClient.Call decode error: " + err.Error())
		return rpcEntity.RpcResponse{ErrCode: constant.ServerError}
	}
	return res
}
