package server

import (
	"EntryTask/constant"
	"EntryTask/rpc/RpcEntity"
	"EntryTask/rpc/codec"
	"EntryTask/rpc/network"
	"EntryTask/rpc/service"
	"github.com/sirupsen/logrus"
	"net"
	"strings"
	"sync"
)

type RpcServer struct {
	mu         sync.Mutex
	serviceMap map[string]*service.RPCService
}

func MakeServer() RpcServer {
	return RpcServer{
		mu:         sync.Mutex{},
		serviceMap: make(map[string]*service.RPCService),
	}
}

func (server RpcServer) Register(svc *service.RPCService) {
	server.mu.Lock()
	defer server.mu.Unlock()
	server.serviceMap[svc.ServiceName] = svc

}

func (server RpcServer) Accept(addr string) {
	// 监听端口
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.Panic("rpcServer.accept error: ", err.Error())
	}
	for {
		accept, _ := listen.Accept()
		// 开启协程 处理RPC
		go func(conn net.Conn) {
			for {
				// 读取数据
				data, err := network.Read(conn)
				if err != nil {
					logrus.Error("rpcServer.Accept read error: ", err.Error())
				}
				// 解码
				req, err := codec.ReqDecode(data)
				if err != nil {
					logrus.Error("rpcServer.Accept decode error: ", err.Error())
				}
				server.mu.Lock()
				index := strings.LastIndex(req.MethodName, ".")
				serviceName := req.MethodName[:index]
				methodName := req.MethodName[index+1:]
				svc, ok := server.serviceMap[serviceName]
				server.mu.Unlock()
				res := RpcEntity.RpcResponse{}
				// 调用处理函数
				if !ok {
					logrus.Error("rpcServer.Accept error: ", err.Error())
					res.Err_code = constant.ServerError
				} else {
					res = svc.RpcHandler(methodName, req)
				}
				// 编码
				data, err = codec.Encode(res)
				if err != nil {
					logrus.Error("rpcServer.Accept encode error: ", err.Error())
				}
				// 发送数据
				network.Send(conn, data)
			}
		}(accept)
	}
}
