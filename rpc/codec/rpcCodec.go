package codec

import (
	"EntryTask/internal/entity"
	"EntryTask/rpc/rpcEntity"
	"bytes"
	"encoding/gob"
)

// 序列化
func Encode(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	gob.Register(rpcEntity.RpcResponse{})
	gob.Register(entity.UserDTO{})
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// RpcRequest反序列化
func ReqDecode(data []byte) (rpcEntity.RpcRequest, error) {
	gob.Register(entity.UserDTO{})
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	var req rpcEntity.RpcRequest
	err := decoder.Decode(&req)
	return req, err
}

// RpcResponse反序列化
func ResDecode(data []byte) (rpcEntity.RpcResponse, error) {
	gob.Register(entity.UserDTO{})
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	var reply rpcEntity.RpcResponse
	err := decoder.Decode(&reply)
	return reply, err
}
