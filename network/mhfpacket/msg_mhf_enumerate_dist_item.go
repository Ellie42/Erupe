package mhfpacket

import (
	"github.com/Andoryuuta/Erupe/network"
	"github.com/Andoryuuta/byteframe"
)

// MsgMhfEnumerateDistItem represents the MSG_MHF_ENUMERATE_DIST_ITEM
type MsgMhfEnumerateDistItem struct{
	AckHandle uint32
	Unk0 uint8
	Unk1 uint16
	Unk2 uint16
}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfEnumerateDistItem) Opcode() network.PacketID {
	return network.MSG_MHF_ENUMERATE_DIST_ITEM
}

// Parse parses the packet from binary
func (m *MsgMhfEnumerateDistItem) Parse(bf *byteframe.ByteFrame) error {
	m.AckHandle = bf.ReadUint32()
	m.Unk0 = bf.ReadUint8()
	m.Unk1 = bf.ReadUint16()
	m.Unk2 = bf.ReadUint16()
	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgMhfEnumerateDistItem) Build(bf *byteframe.ByteFrame) error {
	panic("Not implemented")
}
