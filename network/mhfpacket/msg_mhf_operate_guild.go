package mhfpacket

import (
	"github.com/Andoryuuta/Erupe/network"
	"github.com/Andoryuuta/byteframe"
)

type GuildOperateAction uint8

const (
	GUILD_OPERATE_DISBAND = 0x01
	GUILD_OPERATE_LEAVE   = 0x08
	GUILD_OPERATE_UPGRADE = 0x0a
)

// MsgMhfOperateGuild represents the MSG_MHF_OPERATE_GUILD
type MsgMhfOperateGuild struct {
	AckHandle uint32
	GuildID   uint32
	Action    GuildOperateAction
	UnkData   []byte
}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfOperateGuild) Opcode() network.PacketID {
	return network.MSG_MHF_OPERATE_GUILD
}

// Parse parses the packet from binary
func (m *MsgMhfOperateGuild) Parse(bf *byteframe.ByteFrame) error {
	m.AckHandle = bf.ReadUint32()
	m.GuildID = bf.ReadUint32()
	m.Action = GuildOperateAction(bf.ReadUint8())

	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgMhfOperateGuild) Build(bf *byteframe.ByteFrame) error {
	panic("Not implemented")
}
