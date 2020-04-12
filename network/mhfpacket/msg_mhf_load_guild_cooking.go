package mhfpacket

import (
	"github.com/Andoryuuta/Erupe/network"
	"github.com/Andoryuuta/byteframe"
)

// MsgMhfLoadGuildCooking represents the MSG_MHF_LOAD_GUILD_COOKING
type MsgMhfLoadGuildCooking struct{
	AckHandle   uint32
	Unk0      		uint8
}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfLoadGuildCooking) Opcode() network.PacketID {
	return network.MSG_MHF_LOAD_GUILD_COOKING
}

// Parse parses the packet from binary
func (m *MsgMhfLoadGuildCooking) Parse(bf *byteframe.ByteFrame) error {
	m.AckHandle = bf.ReadUint32()
	m.Unk0 = bf.ReadUint8()
	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgMhfLoadGuildCooking) Build(bf *byteframe.ByteFrame) error {
	panic("Not implemented")
}
