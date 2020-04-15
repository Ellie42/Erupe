package channelserver

import (
	"github.com/Andoryuuta/Erupe/common/stringsupport"
	"github.com/Andoryuuta/Erupe/network/mhfpacket"
	"github.com/Andoryuuta/byteframe"
)

func handleMsgMhfSendMail(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfReadMail(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfListMail(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfListMail)

	mail, err := GetMailListForCharacter(s, s.charID)

	if err != nil {
		doAckBufFail(s, pkt.AckHandle, make([]byte, 4))
		panic(err)
	}

	msg := byteframe.NewByteFrame()

	msg.WriteUint32(uint32(len(mail)))

	for i, m := range mail {
		itemAttached := m.AttachedItemID != nil
		subjectBytes := []byte(stringsupport.MustConvertUTF8ToShiftJIS(m.Subject) + "\x00")
		senderNameBytes := []byte(stringsupport.MustConvertUTF8ToShiftJIS(m.SenderName) + "\x00")

		msg.WriteUint32(m.SenderID)
		msg.WriteUint32(uint32(m.CreatedAt.Unix()))

		// Index with an inconsistent start number, increments normally within a message though
		msg.WriteUint8(uint8(i))

		// Seems to increment sometimes, a little odd
		msg.WriteUint8(uint8(i))
		msg.WriteUint8(0x00)
		msg.WriteBool(itemAttached)
		msg.WriteUint8(uint8(len(subjectBytes)))
		msg.WriteUint8(uint8(len(senderNameBytes)))
		msg.WriteBytes(subjectBytes)
		msg.WriteBytes(senderNameBytes)

		if itemAttached {
			msg.WriteUint16(m.AttachedItemAmount)
			msg.WriteUint16(*m.AttachedItemID)
		}
	}

	doAckBufSucceed(s, pkt.AckHandle, msg.Data())
}

func handleMsgMhfOprtMail(s *Session, p mhfpacket.MHFPacket) {}
