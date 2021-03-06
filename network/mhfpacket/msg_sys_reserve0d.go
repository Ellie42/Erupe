package mhfpacket

import (
	"github.com/Andoryuuta/Erupe/network"
	"github.com/Andoryuuta/byteframe"
)

// MsgSysReserve0D represents the MSG_SYS_reserve0D
type MsgSysReserve0D struct{}

// Opcode returns the ID associated with this packet type.
func (m *MsgSysReserve0D) Opcode() network.PacketID {
	return network.MSG_SYS_reserve0D
}

// Parse parses the packet from binary
func (m *MsgSysReserve0D) Parse(bf *byteframe.ByteFrame) error {
	panic("Not implemented")
}

// Build builds a binary packet from the current data.
func (m *MsgSysReserve0D) Build(bf *byteframe.ByteFrame) error {
	panic("Not implemented")
}
