package channelserver

import (
	"github.com/Andoryuuta/Erupe/common/stringsupport"
	"github.com/Andoryuuta/Erupe/network/mhfpacket"
	"github.com/Andoryuuta/byteframe"
)

func handleMsgMhfSendMail(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfReadMail(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfReadMail)

	mailId := s.mailList[pkt.AccIndex]

	if mailId == 0 {
		doAckBufFail(s, pkt.AckHandle, make([]byte, 4))
		panic("attempting to read mail that doesn't exist in session map")
	}

	mail, err := GetMailByID(s, mailId)

	if err != nil {
		doAckBufFail(s, pkt.AckHandle, make([]byte, 4))
		panic(err)
	}

	_ = mail.MarkRead(s)

	bodyBytes := []byte(stringsupport.MustConvertUTF8ToShiftJIS(mail.Body) + "\x00")

	doAckBufSucceed(s, pkt.AckHandle, bodyBytes)
}

func handleMsgMhfListMail(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfListMail)

	mail, err := GetMailListForCharacter(s, s.charID)

	if err != nil {
		doAckBufFail(s, pkt.AckHandle, make([]byte, 4))
		panic(err)
	}

	if s.mailList == nil {
		s.mailList = make([]int, 256)
	}

	msg := byteframe.NewByteFrame()

	msg.WriteUint32(uint32(len(mail)))

	startIndex := s.mailAccIndex

	for i, m := range mail {
		accIndex := startIndex + uint8(i)
		s.mailList[accIndex] = m.ID
		s.mailAccIndex++

		itemAttached := m.AttachedItemID != nil
		subjectBytes := []byte(stringsupport.MustConvertUTF8ToShiftJIS(m.Subject) + "\x00")
		senderNameBytes := []byte(stringsupport.MustConvertUTF8ToShiftJIS(m.SenderName) + "\x00")

		msg.WriteUint32(m.SenderID)
		msg.WriteUint32(uint32(m.CreatedAt.Unix()))

		msg.WriteUint8(uint8(accIndex))
		msg.WriteUint8(uint8(i))

		flags := uint8(0x00)

		if m.Read {
			flags |= 0x01
		}

		if m.AttachedItemReceived {
			flags |= 0x08
		}

		if m.IsGuildInvite {
			// Guild Invite
			flags |= 0x10

			// System message?
			flags |= 0x04
		}

		msg.WriteUint8(flags)
		msg.WriteBool(itemAttached)
		msg.WriteUint8(uint8(len(subjectBytes)))
		msg.WriteUint8(uint8(len(senderNameBytes)))
		msg.WriteBytes(subjectBytes)
		msg.WriteBytes(senderNameBytes)

		if itemAttached {
			msg.WriteInt16(m.AttachedItemAmount)
			msg.WriteUint16(*m.AttachedItemID)
		}
	}

	doAckBufSucceed(s, pkt.AckHandle, msg.Data())
}

func handleMsgMhfOprtMail(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfOprtMail)

	mail, err := GetMailByID(s, s.mailList[pkt.AccIndex])

	if err != nil {
		doAckSimpleFail(s, pkt.AckHandle, nil)
		panic(err)
	}

	switch mhfpacket.OperateMailOperation(pkt.Operation) {
	case mhfpacket.OperateMailOperationDelete:
		err = mail.MarkDeleted(s)

		if err != nil {
			doAckSimpleFail(s, pkt.AckHandle, nil)
			panic(err)
		}
	}

	doAckSimpleSucceed(s, pkt.AckHandle, nil)
}
