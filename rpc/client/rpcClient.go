package client

import (
	"EntryTask/rpc/RpcEntity"
	"EntryTask/rpc/codec"
	"EntryTask/rpc/network"
	"github.com/sirupsen/logrus"
	"net"
)

type RpcClient struct {
	conn net.Conn
}

var Client RpcClient

func MakeClient(addr string) error {
	conn, err := net.Dial("tcp", addr)
	Client = RpcClient{}
	if err != nil {
		return err
	} else {
		Client.conn = conn
		return nil
	}
}

func (client RpcClient) Call(methodName string, args interface{}) RpcEntity.RpcResponse {
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
	err = network.Send(client.conn, encode)
	if err != nil {
		logrus.Error("rpcClient.Call send error: ", err.Error())
	}
	// 接收
	read, err := network.Read(client.conn)
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
