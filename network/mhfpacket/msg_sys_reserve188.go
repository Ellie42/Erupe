package mhfpacket

import (
	"github.com/Andoryuuta/Erupe/network"
	"github.com/Andoryuuta/byteframe"
)

// MsgSysReserve188 represents the MSG_SYS_reserve188
type MsgSysReserve188 struct {
	AckHandle uint32
}

// Opcode returns the ID associated with this packet type.
func (m *MsgSysReserve188) Opcode() network.PacketID {
	return network.MSG_SYS_reserve188
}

// Parse parses the packet from binary
func (m *MsgSysReserve188) Parse(bf *byteframe.ByteFrame) error {
	m.AckHandle = bf.ReadUint32()
	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgSysReserve188) Build(bf *byteframe.ByteFrame) error {
	panic("Not implemented")
}
