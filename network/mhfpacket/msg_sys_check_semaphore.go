package mhfpacket

import (
	"github.com/Andoryuuta/Erupe/network"
	"github.com/Andoryuuta/byteframe"
)

// MsgSysCheckSemaphore represents the MSG_SYS_CHECK_SEMAPHORE
type MsgSysCheckSemaphore struct{
	AckHandle uint32
	DataSize uint8
	DataBuf []byte
}

// Opcode returns the ID associated with this packet type.
func (m *MsgSysCheckSemaphore) Opcode() network.PacketID {
	return network.MSG_SYS_CHECK_SEMAPHORE
}

// Parse parses the packet from binary
func (m *MsgSysCheckSemaphore) Parse(bf *byteframe.ByteFrame) error {
	m.AckHandle = bf.ReadUint32()
	m.DataSize = bf.ReadUint8()
	m.DataBuf = bf.ReadBytes(uint(m.DataSize))
	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgSysCheckSemaphore) Build(bf *byteframe.ByteFrame) error {
	panic("Not implemented")
}
