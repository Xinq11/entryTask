package client

import (
	"EntryTask/config"
	"EntryTask/constant"
	"EntryTask/logger"
	"EntryTask/rpc/codec"
	"EntryTask/rpc/network"
	"EntryTask/rpc/rpcEntity"
	"net"
)

type RpcClient struct {
	connPool chan net.Conn
}

var Client *RpcClient

func MakeClient(addr string) {
	connPool := make(chan net.Conn, config.ConnNum)
	for i := 0; i < config.ConnNum; i++ {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			logger.Error("rpcClient.MakeClient net dial error: " + err.Error())
		}
		connPool <- conn
	}
	Client = &RpcClient{
		connPool: connPool,
	}
}

// 获取连接
func (client *RpcClient) getConn() net.Conn {
	select {
	case conn := <-client.connPool:
		return conn
	}
}

// 释放连接
func (client *RpcClient) releaseConn(conn net.Conn) {
	select {
	case client.connPool <- conn:
		return
	}
}

// 发起远程过程调用
func (client *RpcClient) Call(methodName string, args interface{}) rpcEntity.RpcResponse {
	conn := client.getConn()
	defer client.releaseConn(conn)
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
