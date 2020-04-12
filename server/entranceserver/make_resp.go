package entranceserver

import (
	"encoding/binary"
	"github.com/Andoryuuta/Erupe/common/stringsupport"
	"net"
	"time"

	"github.com/Andoryuuta/Erupe/config"
	"github.com/Andoryuuta/byteframe"
)

func paddedString(x string, size uint) []byte {
	out := make([]byte, size)
	copy(out, x)

	// Null terminate it.
	out[len(out)-1] = 0
	return out
}

func encodeServerInfo(serverInfos []config.EntranceServerInfo) []byte {
	bf := byteframe.NewByteFrame()

	for serverIdx, si := range serverInfos {
		bf.WriteUint32(binary.LittleEndian.Uint32(net.ParseIP(si.IP).To4()))
		bf.WriteUint16(16 + uint16(serverIdx))
		bf.WriteUint16(si.Unk2)
		bf.WriteUint16(uint16(len(si.Channels)))
		bf.WriteUint8(si.Type)
		bf.WriteUint8(si.Season)
		bf.WriteUint8(si.Unk6)
		bf.WriteBytes(paddedString(stringsupport.MustConvertUTF8ToShiftJIS(si.Name), 66))
		bf.WriteUint32(si.AllowedClientFlags)

		for channelIdx, ci := range si.Channels {
			bf.WriteUint16(ci.Port)
			bf.WriteUint16(16 + uint16(channelIdx))
			bf.WriteUint16(ci.MaxPlayers)
			bf.WriteUint16(ci.CurrentPlayers)
			bf.WriteUint16(ci.Unk4)
			bf.WriteUint16(ci.Unk5)
			bf.WriteUint16(ci.Unk6)
			bf.WriteUint16(ci.Unk7)
			bf.WriteUint16(ci.Unk8)
			bf.WriteUint16(ci.Unk9)
			bf.WriteUint16(ci.Unk10)
			bf.WriteUint16(ci.Unk11)
			bf.WriteUint16(ci.Unk12)
			bf.WriteUint16(ci.Unk13)
		}
	}
	bf.WriteUint32(uint32(time.Now().In(time.FixedZone("UTC+9", 9*60*60)).Unix()))
	bf.WriteUint32(0x0000003C)
	return bf.Data()
}

func makeHeader(data []byte, respType string, entryCount uint16, key byte) []byte {
	bf := byteframe.NewByteFrame()
	bf.WriteBytes([]byte(respType))
	bf.WriteUint16(entryCount)
	bf.WriteUint16(uint16(len(data)))
	if len(data) > 0 {
		bf.WriteUint32(CalcSum32(data))
		bf.WriteBytes(data)
	}

	dataToEncrypt := bf.Data()

	bf = byteframe.NewByteFrame()
	bf.WriteUint8(key)
	bf.WriteBytes(EncryptBin8(dataToEncrypt, key))
	return bf.Data()
}

func makeSv2Resp(servers []config.EntranceServerInfo) []byte {
	rawServerData := encodeServerInfo(servers)
	bf := byteframe.NewByteFrame()
	bf.WriteBytes(makeHeader(rawServerData, "SV2", uint16(len(servers)), 0x00))
	return bf.Data()

}

func makeUsrResp(pkt []byte) []byte {
	// TODO(Andoryuuta): Figure out what this user data is.
	// Is it for the friends list at the world selection screen?
	// If so, how does it work without the entrance server connection being authenticated?

	// uint16 for number of requested ids
	// uint32 for each id
	// response seems to be server number starting from 10 10 00 00 for server 1 channel 1?
	bf := byteframe.NewByteFrameFromBytes(pkt)
	_ = bf.ReadUint32() // ALL+
	_ = bf.ReadUint8()  // 0x00

	userEntries := bf.ReadUint16()
	// actual process will be reading all ids and returning real server, just returning all in server 1 for now
	bf = byteframe.NewByteFrame()
	for i := 0; i < int(userEntries); i++ {
		bf.WriteBytes([]byte{0x10, 0x10, 0x00, 0x00})
	}
	return makeHeader(bf.Data(), "USR", userEntries, 0x00)

}
