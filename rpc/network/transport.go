package network

import (
	"encoding/binary"
	"io"
	"net"
)

const headerLen = 10

func Send(conn net.Conn, data []byte) error {
	buf := make([]byte, headerLen+len(data))
	// 按照大端续将数据长度写入报头
	binary.BigEndian.PutUint32(buf[:headerLen], uint32(len(data)))
	// 写入数据
	copy(buf[headerLen:], data)
	if _, err := conn.Write(buf); err != nil {
		return err
	}
	return nil
}

func Read(conn net.Conn) ([]byte, error) {
	header := make([]byte, headerLen)
	if _, err := io.ReadFull(conn, header); err != nil {
		return nil, err
	}
	dataLen := binary.BigEndian.Uint32(header)
	data := make([]byte, dataLen)
	if _, err := io.ReadFull(conn, data); err != nil {
		return nil, err
	}
	return data, nil
}
