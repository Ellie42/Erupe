package mhfpacket

import (
	"github.com/Andoryuuta/Erupe/network"
	"github.com/Andoryuuta/byteframe"
)

// MsgSysReleaseSemaphore represents the MSG_SYS_RELEASE_SEMAPHORE
type MsgSysReleaseSemaphore struct{
	Unk0	uint32
}

// Opcode returns the ID associated with this packet type.
func (m *MsgSysReleaseSemaphore) Opcode() network.PacketID {
	return network.MSG_SYS_RELEASE_SEMAPHORE
}

// Parse parses the packet from binary
func (m *MsgSysReleaseSemaphore) Parse(bf *byteframe.ByteFrame) error {
	m.Unk0 = bf.ReadUint32()
	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgSysReleaseSemaphore) Build(bf *byteframe.ByteFrame) error {
	panic("Not implemented")
}
