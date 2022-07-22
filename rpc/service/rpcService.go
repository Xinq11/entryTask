package service

import (
	"EntryTask/constant"
	"EntryTask/rpc/RpcEntity"
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
	svc.ServiceName = reflect.Indirect(svc.svcValue).Type().Name()
	svc.methodMap = make(map[string]reflect.Method)
	for i := 0; i < svcType.NumMethod(); i++ {
		method := svcType.Method(i)
		methodName := method.Name
		svc.methodMap[methodName] = method
	}
	return svc
}

// 反射调用本地方法
func (svc *RPCService) RpcHandler(methodName string, req RpcEntity.RpcRequest) RpcEntity.RpcResponse {
	if method, ok := svc.methodMap[methodName]; ok {
		function := method.Func
		res := function.Call([]reflect.Value{svc.svcValue, reflect.ValueOf(req.Args)})
		return res[0].Interface().(RpcEntity.RpcResponse)
	} else {
		reply := RpcEntity.RpcResponse{
			Err_code: constant.ServerError,
		}
		return reply
	}
}
