package mhfpacket

import (
	"github.com/Andoryuuta/Erupe/network"
	"github.com/Andoryuuta/byteframe"
)

// MsgMhfApplyDistItem represents the MSG_MHF_APPLY_DIST_ITEM
type MsgMhfApplyDistItem struct{
	AckHandle uint32
	Unk0 uint8
	RequestType uint32
	Unk2 uint32
	Unk3 uint32
}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfApplyDistItem) Opcode() network.PacketID {
	return network.MSG_MHF_APPLY_DIST_ITEM
}

// Parse parses the packet from binary
func (m *MsgMhfApplyDistItem) Parse(bf *byteframe.ByteFrame) error {
	m.AckHandle = bf.ReadUint32()
	m.Unk0 = bf.ReadUint8()
	m.RequestType = bf.ReadUint32()
	m.Unk2 = bf.ReadUint32()
	m.Unk3 = bf.ReadUint32()
	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgMhfApplyDistItem) Build(bf *byteframe.ByteFrame) error {
	panic("Not implemented")
}
