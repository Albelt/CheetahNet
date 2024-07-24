package enet

//IMessage实现
type Message struct {
	ID   uint32
	Len  uint32
	Data []byte
}

func (m *Message) GetID() uint32 {
	return m.ID
}

func (m *Message) GetLen() uint32 {
	return m.Len
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetID(id uint32) {
	m.ID = id
}

func (m *Message) SetLen(len uint32) {
	m.Len = len
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}

