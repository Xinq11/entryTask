package service

import (
	"EntryTask/constant"
	"EntryTask/rpc/rpcEntity"
	"errors"
	"github.com/sirupsen/logrus"
	"reflect"
)

type RPCService struct {
	ServiceName string
	svcValue    reflect.Value
	methodMap   map[string]reflect.Method
}

func MakeService(svcStruct interface{}) *RPCService {
	svc := &RPCService{}
	svcType := reflect.TypeOf(svcStruct)
	svc.svcValue = reflect.ValueOf(svcStruct)
	// reflect.Indirect: ptr ->ptr.Value 获取service名称
	svc.ServiceName = reflect.Indirect(svc.svcValue).Type().Name()
	svc.methodMap = make(map[string]reflect.Method)
	// 遍历service方法，统计到map中
	for i := 0; i < svcType.NumMethod(); i++ {
		method := svcType.Method(i)
		methodName := method.Name
		svc.methodMap[methodName] = method
	}
	return svc
}

// 反射调用本地方法
func (svc *RPCService) RpcHandler(methodName string, req rpcEntity.RpcRequest) rpcEntity.RpcResponse {
	// 判断service是否存在请求方法
	if method, ok := svc.methodMap[methodName]; ok {
		function := method.Func
		// 反射调用本地方法
		res := function.Call([]reflect.Value{svc.svcValue, reflect.ValueOf(req.Args)})
		return res[0].Interface().(rpcEntity.RpcResponse)
	} else {
		logrus.Error("rpcService.RpcHandler error: ", errors.New("UnKnown method").Error())
		reply := rpcEntity.RpcResponse{
			ErrCode: constant.ServerError,
		}
		return reply
	}
}
