package mhfpacket

import (
	"github.com/Andoryuuta/Erupe/network"
	"github.com/Andoryuuta/byteframe"
)

// MsgSysCreateSemaphore represents the MSG_SYS_CREATE_SEMAPHORE
type MsgSysCreateSemaphore struct{
	AckHandle uint32
	Unk0	uint16
	DataSize       uint16
	RawDataPayload []byte
}

// Opcode returns the ID associated with this packet type.
func (m *MsgSysCreateSemaphore) Opcode() network.PacketID {
	return network.MSG_SYS_CREATE_SEMAPHORE
}

// Parse parses the packet from binary
func (m *MsgSysCreateSemaphore) Parse(bf *byteframe.ByteFrame) error {
	m.AckHandle = bf.ReadUint32()
	m.Unk0 = bf.ReadUint16()
	m.DataSize = bf.ReadUint16()
	m.RawDataPayload = bf.ReadBytes(uint(m.DataSize))
	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgSysCreateSemaphore) Build(bf *byteframe.ByteFrame) error {
	panic("Not implemented")
}
