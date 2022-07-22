package client

import (
	"EntryTask/internal/entity"
	"EntryTask/rpc/RpcEntity"
	"EntryTask/rpc/server"
	"sync"

	"EntryTask/rpc/service"
	"testing"
	"time"
)

type test struct {
	res int
}

func (t test) Add(a entity.UserDTO) RpcEntity.RpcResponse {

	return RpcEntity.RpcResponse{
		Err_code: 0,
		Data:     a.Username,
	}
}

func TestCall(t *testing.T) {
	rpcServer := server.MakeServer()
	test := test{}
	rpcService := service.MakeService(test)
	rpcServer.Register(rpcService)
	go rpcServer.Accept("127.0.0.1:20000")
	time.Sleep(5 * time.Second)
	err := MakeClient("127.0.0.1:20000")
	if err != nil {
		t.Error("fail", err.Error())
	}
	wg := sync.WaitGroup{}
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			call := Client.Call("test.Add", entity.UserDTO{Username: "xq"})
			if call.Err_code == 0 {
				t.Log("success", call.Data)
			} else {
				t.Error("fail", call.Err_code.GetErrMsgByCode())
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
