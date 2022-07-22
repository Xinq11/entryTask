package client

import (
	"EntryTask/rpc/RpcEntity"
	"EntryTask/rpc/codec"
	"EntryTask/rpc/network"
	"github.com/sirupsen/logrus"
	"net"
)

type RpcClient struct {
	connPool chan net.Conn
}

var Client RpcClient

func MakeClient(addr string) error {
	connPool := make(chan net.Conn, 3000)
	for i := 0; i < 3000; i++ {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			logrus.Error("rpcClient.MakeClient net dial error: ", err.Error())
		}
		connPool <- conn
	}
	Client = RpcClient{
		connPool: connPool,
	}
	return nil
}

// 获取连接
func (client RpcClient) getConn() net.Conn {
	select {
	case conn := <-client.connPool:
		return conn
	}
}

// 释放连接
func (client RpcClient) releaseConn(conn net.Conn) {
	select {
	case client.connPool <- conn:
		return
	}
}

func (client RpcClient) Call(methodName string, args interface{}) RpcEntity.RpcResponse {
	conn := client.getConn()
	defer client.releaseConn(conn)
	request := RpcEntity.RpcRequest{
		MethodName: methodName,
		Args:       args,
	}
	// 编码
	encode, err := codec.Encode(request)
	if err != nil {
		logrus.Error("rpcClient.Call encode error: ", err.Error())
	}
	// 发送
	err = network.Send(conn, encode)
	if err != nil {
		logrus.Error("rpcClient.Call send error: ", err.Error())
	}
	// 接收
	read, err := network.Read(conn)
	if err != nil {
		logrus.Error("rpcClient.Call read error: ", err.Error())
	}
	// 解码
	decode, err := codec.ResDecode(read)
	if err != nil {
		logrus.Error("rpcClient.Call decode error: ", err.Error())
	}
	return decode
}
