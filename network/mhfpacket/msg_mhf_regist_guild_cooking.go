package mhfpacket

import (
	"github.com/Andoryuuta/Erupe/network"
	"github.com/Andoryuuta/byteframe"
)

// MsgMhfRegistGuildCooking represents the MSG_MHF_REGIST_GUILD_COOKING
type MsgMhfRegistGuildCooking struct{
	AckHandle      uint32
	Unk0           uint32
	Unk1           uint16
	Unk2           uint8
}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfRegistGuildCooking) Opcode() network.PacketID {
	return network.MSG_MHF_REGIST_GUILD_COOKING
}

// Parse parses the packet from binary
func (m *MsgMhfRegistGuildCooking) Parse(bf *byteframe.ByteFrame) error {
	m.AckHandle = bf.ReadUint32()
	m.Unk0 = bf.ReadUint32()
	m.Unk1 = bf.ReadUint16()
	m.Unk2 = bf.ReadUint8()
	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgMhfRegistGuildCooking) Build(bf *byteframe.ByteFrame) error {
	panic("Not implemented")
}
