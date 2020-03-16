package channelserver

import (
	"encoding/binary"
	"fmt"
	"github.com/Andoryuuta/Erupe/network/mhfpacket"
	"github.com/Andoryuuta/byteframe"
)

func handleMsgMhfCreateGuild(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfCreateGuild)

	guildId, err := CreateGuild(s, stripNullTerminator(pkt.Name))

	if err != nil {
		bf := byteframe.NewByteFrameFromBytes([]byte{0x00, 0x00, 0x00, 0x00})

		// No reasoning behind these values other than they cause a 'failed to create'
		// style message, it's better than nothing for now.
		bf.WriteUint32(0x01010101)

		ack := &mhfpacket.MsgSysAck{AckHandle: pkt.AckHandle, AckData: bf.Data()}

		s.QueueSendMHF(ack)
		return
	}

	bf := byteframe.NewByteFrameFromBytes([]byte{0x00, 0x00, 0x00, 0x00})

	bf.WriteUint32(uint32(guildId))

	ack := &mhfpacket.MsgSysAck{AckHandle: pkt.AckHandle, AckData: bf.Data()}

	s.QueueSendMHF(ack)
}

func handleMsgMhfOperateGuild(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfOperateGuild)

	guild, err := GetGuildInfoByID(s, pkt.GuildID)

	if err != nil {
		return
	}

	// Only leader can manage guild for now
	if guild.Leader.CharID != s.charID {
		s.logger.Warn(fmt.Sprintf("character '%d' is attempting to manage guild '%d' without permission", s.charID, guild.ID))
		return
	}

	switch pkt.Action {
	case mhfpacket.GUILD_OPERATE_DISBAND:
		err = DisbandGuild(s, guild.ID)
	default:
		panic(fmt.Sprintf("unhandled operate guild action '%d'", pkt.Action))
	}

	response := 0x01

	if err != nil {
		response = 0x00
	}

	bf := byteframe.NewByteFrameFromBytes([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, byte(response)})

	ackMessage := &mhfpacket.MsgSysAck{
		AckHandle: pkt.AckHandle,
		AckData:   bf.Data(),
	}

	s.QueueSendMHF(ackMessage)
}

func handleMsgMhfOperateGuildMember(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfInfoGuild(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfInfoGuild)

	var guild *Guild
	var err error

	if pkt.GuildID > 0 {
		guild, err = GetGuildInfoByID(s, pkt.GuildID)
	} else {
		guild, err = GetGuildInfoByCharacterId(s, s.charID)
	}

	if err == nil && guild != nil {
		characterGuildData, err := GetCharacterGuildData(s, s.charID)

		characterJoinedAt := uint32(0)

		if characterGuildData != nil {
			characterJoinedAt = uint32(characterGuildData.JoinedAt.Unix())
		}

		if err != nil {
			resp := byteframe.NewByteFrame()
			resp.WriteUint32(0) // Count
			resp.WriteUint8(0)  // Unk, read if count == 0.

			doSizedAckResp(s, pkt.AckHandle, resp.Data())
			return
		}

		bf := byteframe.NewByteFrame()

		bf.WriteUint32(uint32(guild.ID))
		bf.WriteUint32(guild.Leader.CharID)
		bf.WriteUint16(0x0) // Guild festival ranking (I think)
		bf.WriteUint16(guild.MemberCount)

		// Unk appears to be static
		bf.WriteBytes([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

		bf.WriteUint16(0x01) // Unk appears to be static

		guildMainMotto := fmt.Sprintf("%s\x00", guild.MainMotto)

		bf.WriteUint32(uint32(guild.CreatedAt.Unix()))
		bf.WriteUint32(characterJoinedAt)
		bf.WriteUint8(uint8(len(guild.Name)))
		bf.WriteUint8(uint8(len(guildMainMotto)))
		bf.WriteUint8(uint8(5)) // Length of unknown string below
		bf.WriteUint8(uint8(len(guild.Leader.Name)))
		bf.WriteBytes([]byte(guild.Name))
		bf.WriteBytes([]byte(guildMainMotto))
		bf.WriteBytes([]byte{0xFF, 0x00, 0x00, 0x00, 0x00}) // Unk string
		bf.WriteBytes([]byte(guild.Leader.Name))
		bf.WriteBytes([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x02, 0x00, 0x00, 0x00, 0x00}) // Unk

		// Here there are always 3 null terminated names, not sure what they relate to though
		// Having all three as null bytes is perfectly valid
		for i := 0; i < 3; i++ {
			bf.WriteUint8(0x1) // Name Length - 1 minimum due to null byte
			bf.WriteUint8(0x0) // Name string
		}

		// Unk
		bf.WriteBytes([]byte{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1E, 0x00,
			0x00, 0xD6, 0xD8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		})

		// Unk-ish Indicates an alliance with 0x0a
		// When using 0x1b there is a lot more data expected after
		// 0x0 = no alliance
		bf.WriteUint8(0x0)

		// TODO add alliance parts here

		bf.WriteBytes([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

		doSizedAckResp(s, pkt.AckHandle, bf.Data())
	} else {
		//// REALLY large/complex format... stubbing it out here for simplicity.
		//resp := byteframe.NewByteFrame()
		//resp.WriteUint32(0) // Count
		//resp.WriteUint8(0)  // Unk, read if count == 0.

		doSizedAckResp(s, pkt.AckHandle, make([]byte, 8))
	}
}

func handleMsgMhfEnumerateGuild(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfEnumerateGuild)

	var guilds []*Guild
	var err error

	switch pkt.Type {
	case mhfpacket.ENUMERATE_GUILD_TYPE_NAME:
		// I have no idea if is really little endian, but it seems too weird to have a random static
		// 0x00 before the string
		searchTermLength := binary.LittleEndian.Uint16(pkt.RawDataPayload[9:11])
		searchTerm := pkt.RawDataPayload[11 : 11+searchTermLength]

		guilds, err = FindGuildsByName(s, stripNullTerminator(string(searchTerm)))
	default:
		panic(fmt.Sprintf("no handler for guild search type '%d'", pkt.Type))
	}

	if err != nil || guilds == nil {
		stubEnumerateNoResults(s, pkt.AckHandle)
		return
	}

	bf := byteframe.NewByteFrame()
	bf.WriteUint16(uint16(len(guilds)))

	for _, guild := range guilds {
		bf.WriteUint8(0x00) // Unk
		bf.WriteUint32(guild.ID)
		bf.WriteUint32(guild.Leader.CharID)
		bf.WriteUint16(guild.MemberCount)
		bf.WriteUint8(0x00)  // Unk
		bf.WriteUint8(0x00)  // Unk
		bf.WriteUint16(0x00) // Rank
		bf.WriteUint32(uint32(guild.CreatedAt.Unix()))

		guildName := fmt.Sprintf("%s\x00", guild.Name)
		leaderName := fmt.Sprintf("%s\x00", guild.Leader.Name)

		bf.WriteUint8(uint8(len(guildName)))
		bf.WriteBytes([]byte(guildName))
		bf.WriteUint8(uint8(len(leaderName)))
		bf.WriteBytes([]byte(leaderName))
		bf.WriteUint8(0x01)
	}

	bf.WriteUint8(0x01)
	bf.WriteUint8(0x00)

	doSizedAckResp(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfUpdateGuild(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfArrangeGuildMember(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEnumerateGuildMember(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfEnumerateGuildMember)
	guild, err := GetGuildInfoByCharacterId(s, s.charID)

	if err != nil {
		s.logger.Warn("failed to retrieve guild sending no result message")
		doSizedAckResp(s, pkt.AckHandle, make([]byte, 2))
		return
	} else if guild == nil {
		doSizedAckResp(s, pkt.AckHandle, make([]byte, 2))
		return
	}

	guildMembers, err := GetGuildMembers(s, guild.ID, guild.MemberCount)

	if err != nil {
		s.logger.Error("failed to retrieve guild")
		return
	}

	bf := byteframe.NewByteFrame()

	bf.WriteUint16(guild.MemberCount)

	for _, member := range guildMembers {
		bf.WriteUint32(member.CharID)
		bf.WriteBytes([]byte{0x00, 0x63, 0x00, 0x00, 0x3A, 0xE9, 0x06, 0x00, 0x01}) // Unk
		bf.WriteUint16(uint16(len(member.Name)))
		bf.WriteBytes([]byte(member.Name))
	}

	for _, member := range guildMembers {
		bf.WriteUint32(uint32(member.JoinedAt.Unix()))
	}

	bf.WriteBytes([]byte{0x00, 0x00}) // Unk, might be to do with alliance, 0x00 == no alliance

	for range guildMembers {
		bf.WriteUint32(0x00) // Unk
	}

	doSizedAckResp(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfPostGuildScout(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfCancelGuildScout(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfAnswerGuildScout(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetGuildScoutList(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetGuildScoutList)

	// No scouting allowed
	doSizedAckResp(s, pkt.AckHandle, make([]byte, 4))
}

func handleMsgMhfGetGuildManageRight(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetGuildManageRight)

	guild, err := GetGuildInfoByCharacterId(s, s.charID)

	if err != nil {
		s.logger.Warn("failed to respond to manage rights message")
		return
	} else if guild == nil {
		bf := byteframe.NewByteFrame()
		bf.WriteUint16(0x00) // Unk
		bf.WriteUint16(0x00) // Member count

		doSizedAckResp(s, pkt.AckHandle, bf.Data())
		return
	}

	bf := byteframe.NewByteFrame()

	bf.WriteUint16(0x00) // Unk
	bf.WriteUint16(guild.MemberCount)

	members, err := GetGuildMembers(s, guild.ID, guild.MemberCount)

	for _, member := range members {
		bf.WriteUint32(member.CharID)
		bf.WriteUint32(0x0)
	}

	doSizedAckResp(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfSetGuildManageRight(s *Session, p mhfpacket.MHFPacket) {}
