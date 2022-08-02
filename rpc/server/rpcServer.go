package server

import (
	"EntryTask/constant"
	"EntryTask/logger"
	"EntryTask/rpc/codec"
	"EntryTask/rpc/network"
	"EntryTask/rpc/rpcEntity"
	"EntryTask/rpc/service"
	"net"
	"strings"
	"sync"
)

type RpcServer struct {
	mu         sync.Mutex
	serviceMap map[string]*service.RPCService
}

// 初始化RPC服务端
func MakeServer() *RpcServer {
	return &RpcServer{
		mu:         sync.Mutex{},
		serviceMap: make(map[string]*service.RPCService),
	}
}

func (server *RpcServer) Register(svc *service.RPCService) {
	server.mu.Lock()
	defer server.mu.Unlock()
	server.serviceMap[svc.ServiceName] = svc

}

func (server *RpcServer) Accept(addr string) {
	// 监听端口
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Panic("rpcServer.accept listen error: " + err.Error())
	}
	for {
		accept, _ := listen.Accept()
		// 开启协程 处理RPC
		go func(conn net.Conn) {
			for {
				// 读取数据
				data, err := network.Read(conn)
				if err != nil {
					//Read方法返回EOF错误，表示本端感知到对端已经关闭连接（本端已接收到对端发送的FIN）。此后如果本端不调用Close方法，只释放本端的连接对象，则连接处于非完全关闭状态（CLOSE_WAIT）。即文件描述符发生泄漏。
					if err.Error() == "EOF" {
						conn.Close()
						break
					}
					logger.Error("rpcServer.Accept read error: " + err.Error())
				}
				// 解码
				req, err := codec.ReqDecode(data)
				if err != nil {
					logger.Error("rpcServer.Accept decode error: " + err.Error())
				}
				ok := false
				var methodName string
				var svc *service.RPCService
				// 判断rpc请求的服务是否存在
				if req.MethodName != "" && strings.Contains(req.MethodName, ".") {
					server.mu.Lock()
					index := strings.LastIndex(req.MethodName, ".")
					serviceName := req.MethodName[:index]
					methodName = req.MethodName[index+1:]
					svc, ok = server.serviceMap[serviceName]
					server.mu.Unlock()
				}
				res := rpcEntity.RpcResponse{}
				// 如果rpc请求的服务存在，则反射调用本地方法
				if !ok {
					logger.Error("rpcServer.Accept handler error: " + err.Error())
					res.ErrCode = constant.ServerError
				} else {
					res = svc.RpcHandler(methodName, req)
				}
				// 编码
				data, err = codec.Encode(res)
				if err != nil {
					logger.Error("rpcServer.Accept encode error: " + err.Error())
				}
				// 发送数据
				err = network.Send(conn, data)
				if err != nil {
					logger.Error("rpcServer.Accept send error: " + err.Error())
				}
			}
		}(accept)
	}
}
