package channelserver

import (
	"github.com/Andoryuuta/Erupe/network/mhfpacket"
	"github.com/Andoryuuta/byteframe"
	"go.uber.org/zap"
	"math/rand"
	"fmt"
)

func handleMsgMhfMercenaryHuntdata(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfMercenaryHuntdata)
	doAckBufSucceed(s, pkt.AckHandle, make([]byte, 0x0A))
}

func handleMsgMhfEnumerateMercenaryLog(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfCreateMercenary(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfCreateMercenary)

	bf := byteframe.NewByteFrame()

	bf.WriteUint32(0x00)          // Unk
	bf.WriteUint32(rand.Uint32()) // Partner ID?

  doAckSimpleSucceed(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfSaveMercenary(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfSaveMercenary)
	bf := byteframe.NewByteFrameFromBytes(pkt.RawDataPayload)
	GCPValue := bf.ReadUint32()
	_ = bf.ReadUint32() // unk
	MercDataSize := bf.ReadUint32()
	MercData := bf.ReadBytes(uint(MercDataSize))
	_ = bf.ReadUint32() // unk

	if MercDataSize > 0{
	// the save packet has an extra null byte after its size
	_, err := s.server.db.Exec("UPDATE characters SET savemercenary=$1 WHERE id=$2", MercData[:MercDataSize], s.charID)
		if err != nil {
			s.logger.Fatal("Failed to update savemercenary and gcp in db", zap.Error(err))
		}
	}
	// gcp value is always present regardless
	_, err := s.server.db.Exec("UPDATE characters SET gcp=$1 WHERE id=$2", GCPValue, s.charID)
	if err != nil {
		s.logger.Fatal("Failed to update savemercenary and gcp in db", zap.Error(err))
	}
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfReadMercenaryW(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfReadMercenaryW)
	var data []byte
	var gcp uint32
	// still has issues
	err := s.server.db.QueryRow("SELECT savemercenary FROM characters WHERE id = $1", s.charID).Scan(&data)
	if err != nil {
		s.logger.Fatal("Failed to get savemercenary data from db", zap.Error(err))
	}

	err = s.server.db.QueryRow("SELECT COALESCE(gcp, 0) FROM characters WHERE id = $1", s.charID).Scan(&gcp)
	if err != nil {
		panic(err)
	}
	if len(data) == 0{
		data = []byte{0x00}
	}

	resp := byteframe.NewByteFrame()
	resp.WriteBytes(data)
	resp.WriteUint16(0)
	resp.WriteUint32(gcp)
	fmt.Printf("% x", resp.Data())
	doAckBufSucceed(s, pkt.AckHandle, resp.Data())
}

func handleMsgMhfReadMercenaryM(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfReadMercenaryM)
	// accessing actual rasta data of someone else still unsure of the formatting of this
	doAckBufSucceed(s, pkt.AckHandle,  []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfContractMercenary(s *Session, p mhfpacket.MHFPacket) {}
