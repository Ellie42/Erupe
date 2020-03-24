package channelserver

import (
	"github.com/Andoryuuta/Erupe/network/mhfpacket"
	"github.com/Andoryuuta/byteframe"
	"go.uber.org/zap"
	"math/rand"
)

func handleMsgMhfMercenaryHuntdata(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEnumerateMercenaryLog(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfCreateMercenary(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfCreateMercenary)

	bf := byteframe.NewByteFrame()

	bf.WriteUint32(0x00)          // Unk
	bf.WriteUint32(rand.Uint32()) // Partner ID?

	s.QueueAck(pkt.AckHandle, bf.Data())
}

func handleMsgMhfSaveMercenary(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfSaveMercenary)
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfReadMercenaryW(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfReadMercenaryW)
	var data []byte
	err := s.server.db.QueryRow("SELECT savemercenary FROM characters WHERE id = $1", s.charID).Scan(&data)
	if err != nil {
		s.logger.Fatal("Failed to get savemercenary data from db", zap.Error(err))
	}
	doAckBufSucceed(s, pkt.AckHandle, data)
}

func handleMsgMhfReadMercenaryM(s *Session, p mhfpacket.MHFPacket) {
	// I'm assuming this is just called if your character is male over female but haven't checked
}

func handleMsgMhfContractMercenary(s *Session, p mhfpacket.MHFPacket) {}
