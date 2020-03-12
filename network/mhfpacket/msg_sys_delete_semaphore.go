package mhfpacket

import (
	"github.com/Andoryuuta/Erupe/network"
	"github.com/Andoryuuta/byteframe"
)

// MsgSysDeleteSemaphore represents the MSG_SYS_DELETE_SEMAPHORE
type MsgSysDeleteSemaphore struct{
	Unk0	uint16
	Unk1	uint16
	Unk2	uint16
	Unk3	uint16
	Unk4	uint16
}

// Opcode returns the ID associated with this packet type.
func (m *MsgSysDeleteSemaphore) Opcode() network.PacketID {
	return network.MSG_SYS_DELETE_SEMAPHORE
}

// Parse parses the packet from binary
func (m *MsgSysDeleteSemaphore) Parse(bf *byteframe.ByteFrame) error {
	m.Unk0 = bf.ReadUint16()
	m.Unk1 = bf.ReadUint16()
	m.Unk2 = bf.ReadUint16()
	m.Unk3 = bf.ReadUint16()
	m.Unk4 = bf.ReadUint16()
	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgSysDeleteSemaphore) Build(bf *byteframe.ByteFrame) error {
	panic("Not implemented")
}
