package gnet

import (
	"errors"
	"github.com/go-ll/ginx/giface"
	"sync"
)

type ConnManager struct {
	conns    map[uint32]giface.Conn
	connLock sync.RWMutex
}

// NewConnManager 创建一个连接管理
func NewConnManager() *ConnManager {
	return &ConnManager{
		conns: make(map[uint32]giface.Conn),
	}
}

func (cm *ConnManager) Add(conn giface.Conn) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	cm.conns[conn.GetConnId()] = conn

}
func (cm *ConnManager) Remove(conn giface.Conn) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	delete(cm.conns, conn.GetConnId())
}
func (cm *ConnManager) Get(connId uint32) (giface.Conn, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	if conn, ok := cm.conns[connId]; ok {
		return conn, nil
	}
	return nil, errors.New("connection not found")
}
func (cm *ConnManager) Count() int {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()
	return len(cm.conns)
}
func (cm *ConnManager) ClearConn() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	for connId, conn := range cm.conns {
		// 停止
		conn.Stop()
		// 删除
		delete(cm.conns, connId)
	}
}
