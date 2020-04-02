package mhfpacket

import (
	"github.com/Andoryuuta/Erupe/network"
	"github.com/Andoryuuta/byteframe"
)

// MsgMhfMercenaryHuntdata represents the MSG_MHF_MERCENARY_HUNTDATA
type MsgMhfMercenaryHuntdata struct{
		AckHandle      uint32
		Unk0       uint8
}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfMercenaryHuntdata) Opcode() network.PacketID {
	return network.MSG_MHF_MERCENARY_HUNTDATA
}

// Parse parses the packet from binary
func (m *MsgMhfMercenaryHuntdata) Parse(bf *byteframe.ByteFrame) error {
	m.AckHandle = bf.ReadUint32()
	m.Unk0 = bf.ReadUint8()
	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgMhfMercenaryHuntdata) Build(bf *byteframe.ByteFrame) error {
	panic("Not implemented")
}
