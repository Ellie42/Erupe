package mhfpacket

import (
	"github.com/Andoryuuta/Erupe/network"
	"github.com/Andoryuuta/byteframe"
)

// MsgSysLockGlobalSema represents the MSG_SYS_LOCK_GLOBAL_SEMA
type MsgSysLockGlobalSema struct {
	AckHandle uint32
	StageName string
	ServerID  string
}

// Opcode returns the ID associated with this packet type.
func (m *MsgSysLockGlobalSema) Opcode() network.PacketID {
	return network.MSG_SYS_LOCK_GLOBAL_SEMA
}

// Parse parses the packet from binary
func (m *MsgSysLockGlobalSema) Parse(bf *byteframe.ByteFrame) error {
	m.AckHandle = bf.ReadUint32()

	stageNameLength := bf.ReadUint16()
	serverIDLength := bf.ReadUint16()

	m.StageName = string(bf.ReadBytes(uint(stageNameLength)))
	m.ServerID = string(bf.ReadBytes(uint(serverIDLength)))

	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgSysLockGlobalSema) Build(bf *byteframe.ByteFrame) error {
	panic("Not implemented")
}
