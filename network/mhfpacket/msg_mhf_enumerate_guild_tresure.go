package mhfpacket

import (
	"github.com/Andoryuuta/Erupe/network"
	"github.com/Andoryuuta/byteframe"
)

// MsgMhfEnumerateGuildTresure represents the MSG_MHF_ENUMERATE_GUILD_TRESURE
type MsgMhfEnumerateGuildTresure struct {
	AckHandle uint32
}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfEnumerateGuildTresure) Opcode() network.PacketID {
	return network.MSG_MHF_ENUMERATE_GUILD_TRESURE
}

// Parse parses the packet from binary
func (m *MsgMhfEnumerateGuildTresure) Parse(bf *byteframe.ByteFrame) error {
	m.AckHandle = bf.ReadUint32()

	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgMhfEnumerateGuildTresure) Build(bf *byteframe.ByteFrame) error {
	panic("Not implemented")
}
