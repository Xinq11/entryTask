package pool

import (
	"EntryTask/config"
	"EntryTask/logger"
	"net"
	"sync"
)

type Pool struct {
	addr     string
	mu       sync.Mutex
	freeConn []net.Conn
	waitConn map[int]chan net.Conn
	waitNum  int // 等待连接数量
	numOpen  int // 打开连接数量
	numFree  int // 空闲连接数量
	numMax   int // 最大连接数
}

// 初始化连接池 懒加载
func Init(addr string) (*Pool, error) {
	return &Pool{
		addr:     addr,
		freeConn: make([]net.Conn, 0, config.ConnNum),
		mu:       sync.Mutex{},
		numFree:  config.ConnNum,
		waitNum:  1,
		numOpen:  0,
		waitConn: make(map[int]chan net.Conn),
		numMax:   config.ConnMax,
	}, nil
}

// 获取连接
func (pool *Pool) GetConn() (net.Conn, error) {
	pool.mu.Lock()
	lenFreeConn := len(pool.freeConn)
	// 如果有空闲连接 则返回
	if lenFreeConn != 0 {
		logger.Info("tcpPool.GetConn get conn from pool...")
		// 取出切片中的第一个连接 截取数组
		conn := pool.freeConn[0]
		pool.freeConn = pool.freeConn[1:lenFreeConn]
		pool.mu.Unlock()
		return conn, nil
	}
	// 如果连接数大于等于最大连接数 阻塞等待可用连接
	if pool.numOpen >= pool.numMax {
		waitChan := make(chan net.Conn, 1)
		pool.waitConn[pool.waitNum] = waitChan
		pool.waitNum++
		pool.mu.Unlock()
		logger.Info("tcpPool.GetConn wait conn release...")
		select {
		case conn := <-waitChan:
			return conn, nil
		}
	}
	// 如果当前既无空闲连接 也未达到最大连接数 则创建新的连接
	logger.Info("tcpPool.GetConn create new conn...")
	pool.numOpen++
	pool.mu.Unlock()
	conn, err := net.Dial("tcp", pool.addr)
	if err != nil {
		logger.Error("rpcClient.MakeClient net dial error: " + err.Error())
		return nil, err
	}
	return conn, nil
}

// 释放连接
func (pool *Pool) ReleaseConn(conn net.Conn) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	// 如果当前存在等待连接 则复用要释放的连接
	if len(pool.waitConn) != 0 {
		logger.Info("tcpPool.ReleaseConn multiplex conn...")
		var num int
		var waitChan chan net.Conn
		// 从等待池中取出一个值
		for num, waitChan = range pool.waitConn {
			break
		}
		delete(pool.waitConn, num)
		waitChan <- conn
		return
	} else if pool.numOpen > pool.numMax || len(pool.freeConn) == pool.numFree {
		// 如果当前连接数大于最大连接数 或空闲连接已满 关闭连接
		logger.Info("tcpPool.ReleaseConn release conn...")
		conn.Close()
		return
	} else {
		// 如果空闲连接未满 则将连接放回空闲池中
		logger.Info("tcpPool.ReleaseConn release conn to pool...")
		pool.freeConn = append(pool.freeConn, conn)
		return
	}
}

// 关闭连接
func (pool *Pool) CloseFreeConn() {
	logger.Info("tcpPool.CloseFreeConn close conn...")
	for _, conn := range pool.freeConn {
		conn.Close()
	}
}
