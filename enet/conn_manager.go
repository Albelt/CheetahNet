package enet

import (
	"EagleNet/eiface"
	"fmt"
	"sync"
)

//IConnManager接口实现
type ConnManager struct {
	//connID->Connection
	connMap map[uint32]eiface.IConnection
	//connMap的锁
	mutex sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connMap: make(map[uint32]eiface.IConnection),
	}
}

func (m *ConnManager) Add(conn eiface.IConnection) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.connMap[conn.GetConnID()] = conn
}

func (m *ConnManager) Remove(connID uint32) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.connMap, connID)
}

func (m *ConnManager) Get(connID uint32) (eiface.IConnection, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if conn, ok := m.connMap[connID]; ok {
		return conn, nil
	} else {
		return nil, fmt.Errorf("Connection(id=%d) not found", connID)
	}
}

func (m *ConnManager) Len() int {
	return len(m.connMap)
}

func (m *ConnManager) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	//停止所有Connection并删除
	for _, conn := range m.connMap {
		conn.Stop()
		delete(m.connMap, conn.GetConnID())
	}

	fmt.Printf("ConnManager clear all connections.\n")
}
