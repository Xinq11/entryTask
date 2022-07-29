package client

import (
	"EntryTask/internal/entity"
	"EntryTask/rpc/rpcEntity"
	"EntryTask/rpc/server"
	"sync"

	"EntryTask/rpc/service"
	"testing"
	"time"
)

type test struct {
}

func (t test) Add(a entity.UserDTO) rpcEntity.RpcResponse {
	return rpcEntity.RpcResponse{
		ErrCode: 0,
		Data:    a.Username + "rpc",
	}
}

func TestCall(t *testing.T) {
	rpcServer := server.MakeServer()
	test := test{}
	rpcService := service.MakeService(test)
	rpcServer.Register(rpcService)
	go rpcServer.Accept("127.0.0.1:20001")
	time.Sleep(5 * time.Second)
	MakeClient("127.0.0.1:20001")
	wg := sync.WaitGroup{}
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			call := Client.Call("test.Add", entity.UserDTO{Username: "xq"})
			if call.ErrCode == 0 {
				t.Log("success", call.Data)
			} else {
				t.Error("fail", call.ErrCode.GetErrMsgByCode())
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
