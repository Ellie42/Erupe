package channelserver

import (
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"sync"
	"time"

	"github.com/Andoryuuta/Erupe/network"
	"github.com/Andoryuuta/Erupe/network/mhfpacket"
	"github.com/Andoryuuta/byteframe"
)

// Session holds state for the channel server connection.
type Session struct {
	sync.Mutex
	server    *Server
	rawConn   net.Conn
	cryptConn *network.CryptConn
}

// NewSession creates a new Session type.
func NewSession(server *Server, conn net.Conn) *Session {
	s := &Session{
		server:    server,
		rawConn:   conn,
		cryptConn: network.NewCryptConn(conn),
	}
	return s
}

// Start starts the session packet read&handle loop.
func (s *Session) Start() {
	go func() {
		fmt.Println("Channel server got connection!")
		// Unlike the sign and entrance server,
		// the client DOES NOT initalize the channel connection with 8 NULL bytes.

		for {
			pkt, err := s.cryptConn.ReadPacket()
			if err != nil {
				fmt.Println(err)
				fmt.Println("Error on channel server readpacket")
				return
			}
			s.handlePacketGroup(pkt)
		}
	}()
}

var loadDataCount int
var clientEnumerateCount int
var getPaperDataCount int
var questListCount int
var setStageBinary []byte

func (s *Session) handlePacketGroup(pktGroup []byte) {
	defer func() {
		if r := recover(); r != nil {
			bf := byteframe.NewByteFrameFromBytes(pktGroup)
			opcode := network.PacketID(bf.ReadUint16())
			fmt.Println("Recovered from panic: ", opcode)
		}
	}()

	bf := byteframe.NewByteFrameFromBytes(pktGroup)
	opcode := network.PacketID(bf.ReadUint16())

	if (opcode != network.MSG_SYS_END) {
		fmt.Printf("Opcode: %s\n", opcode)
		fmt.Printf("Data:\n%s\n", hex.Dump(pktGroup))
	}

	switch opcode {
	case network.MSG_MHF_ENUMERATE_EVENT:
		fallthrough
	case network.MSG_SYS_reserve18B:
		fallthrough
	case network.MSG_MHF_READ_MERCENARY_W:
		fallthrough
	case network.MSG_MHF_GET_ETC_POINTS:
		fallthrough
	case network.MSG_MHF_READ_GUILDCARD:
		fallthrough
	case network.MSG_MHF_CHECK_WEEKLY_STAMP:
		fallthrough
	case network.MSG_MHF_GET_KIJU_INFO:
		fallthrough
	case network.MSG_MHF_GET_KOURYOU_POINT:
		fallthrough
	case network.MSG_MHF_GET_RENGOKU_RANKING_RANK:
		fallthrough
	case network.MSG_MHF_READ_BEAT_LEVEL:
		fallthrough
	case network.MSG_MHF_GET_WEEKLY_SCHEDULE:
		fallthrough
	case network.MSG_MHF_LIST_MEMBER:
		fallthrough
	case network.MSG_MHF_LOAD_PLATE_DATA:
		fallthrough
	case network.MSG_MHF_LOAD_PLATE_BOX:
		fallthrough
	case network.MSG_MHF_LOAD_FAVORITE_QUEST:
		fallthrough
	case network.MSG_MHF_LOAD_DECO_MYSET:
		fallthrough
	case network.MSG_MHF_LOAD_HUNTER_NAVI:
		fallthrough
	case network.MSG_MHF_GET_UD_SCHEDULE:
		fallthrough
	case network.MSG_MHF_GET_UD_INFO:
		fallthrough
	case network.MSG_MHF_GET_UD_MONSTER_POINT:
		fallthrough
	case network.MSG_MHF_GET_RAND_FROM_TABLE:
		fallthrough
	case network.MSG_MHF_LOAD_PLATE_MYSET:
		fallthrough
	case network.MSG_MHF_LOAD_RENGOKU_DATA:
		fallthrough
	case network.MSG_MHF_ENUMERATE_SHOP:
		fallthrough
	case network.MSG_MHF_LOAD_SCENARIO_DATA:
		fallthrough
	case network.MSG_MHF_GET_BOOST_RIGHT:
		fallthrough
	case network.MSG_MHF_GET_REWARD_SONG:
		fallthrough
	case network.MSG_MHF_GET_UD_SELECTED_COLOR_INFO:
		fallthrough
	case network.MSG_MHF_LOAD_MEZFES_DATA:
		fallthrough
	case network.MSG_MHF_GET_GACHA_POINT:
		fallthrough
	case network.MSG_SYS_CREATE_STAGE:
		fallthrough
	case network.MSG_SYS_RESERVE_STAGE:
		fallthrough
	case network.MSG_MHF_POST_BOOST_TIME_QUEST_RETURN:
		fallthrough
	case network.MSG_MHF_UPDATE_USE_TREND_WEAPON_LOG:
		fallthrough
	case network.MSG_MHF_STAMPCARD_STAMP:
		fallthrough
	case network.MSG_MHF_GET_SEIBATTLE:
		fallthrough
	case network.MSG_MHF_GET_GUILD_WEEKLY_BONUS_ACTIVE_COUNT:
		fallthrough
	case network.MSG_MHF_UPDATE_CAFEPOINT:
		// netcafe points
		fallthrough
	case network.MSG_MHF_GET_KEEP_LOGIN_BOOST_STATUS:
		fallthrough
	case network.MSG_MHF_ENUMERATE_PRICE:
		// seems to just be the client grabbing monster bounty prices to pair with quest data
		fallthrough
	case network.MSG_SYS_ISSUE_LOGKEY:
		fallthrough
	case network.MSG_SYS_LOCK_STAGE:
		fallthrough
	case network.MSG_MHF_TRANSFER_ITEM:
		fallthrough
	case network.MSG_MHF_GET_EARTH_STATUS:
		ackHandle := bf.ReadUint32()
		data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp.bin", opcode.String()))
		if err != nil {
			panic(err)
		}

		bfw := byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteBytes(data)
		s.cryptConn.SendPacket(bfw.Data())	
	case network.MSG_MHF_GET_BOOST_TIME_LIMIT:
		ackHandle := bf.ReadUint32()
		data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp.bin", opcode.String()))
		if err != nil {
			panic(err)
		}
		// Actual response 
		bfw := byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteBytes(data)
		s.cryptConn.SendPacket(bfw.Data())
		
		// Second packet with no real content and shared ack, Rasta for partner maybe? No idea on boost
		bfw = byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteUint64(0x0000000000000000)
		bfw.WriteUint16(0x0010)
		s.cryptConn.SendPacket(bfw.Data())	
	case network.MSG_MHF_LOAD_OTOMO_AIROU:
		ackHandle := bf.ReadUint32()
		data, err := ioutil.ReadFile(fmt.Sprintf("save_files/MSG_MHF_SAVE_OTOMO_AIROU.bin"))
		if err != nil {
			panic(err)
		}

		bfw := byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteUint8(0x01)
		bfw.WriteBytes(data)
		bfw.WriteUint16(0x0010)
		s.cryptConn.SendPacket(bfw.Data())
	case network.MSG_MHF_LOAD_PARTNER:
		ackHandle := bf.ReadUint32()
		data, err := ioutil.ReadFile(fmt.Sprintf("save_files/MSG_MHF_SAVE_PARTNER.bin"))
		if err != nil {
			panic(err)
		}

		bfw := byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteUint8(0x01)
		bfw.WriteBytes(data)
		bfw.WriteUint16(0x0010)
		s.cryptConn.SendPacket(bfw.Data())
			
		// Second packet with no real content and shared ack, Rasta for partner maybe? No idea on boost
		bfw = byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteUint64(0x0000000000000000)
		bfw.WriteUint16(0x0010)
		s.cryptConn.SendPacket(bfw.Data())	
	case network.MSG_SYS_GET_FILE:
		// needs to check if scenario or quest file and then deliver the actual request file instead of a pure quest response
		ackHandle := bf.ReadUint32()
		_ = bf.ReadBytes(2)

		var fileNameHex string = string(bf.ReadBytes(7))
		data, err := ioutil.ReadFile(fmt.Sprintf("file_resp/%s.bin", string(fileNameHex)))
		if err != nil {
			panic(err)
		}

		
		var length uint16 = uint16(len(data))
		
		bfw := byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteUint16(0x0100)
		bfw.WriteUint16(length)
		bfw.WriteBytes(data)
		bfw.WriteUint16(0x0010)
		s.cryptConn.SendPacket(bfw.Data())
	case network.MSG_SYS_ENUMERATE_CLIENT:
		if clientEnumerateCount == 0 {
			ackHandle := bf.ReadUint32()
			data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp0.bin", opcode.String()))
			if err != nil {
				panic(err)
			}

			bfw := byteframe.NewByteFrame()
			bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
			bfw.WriteUint32(ackHandle)
			bfw.WriteBytes(data)
			s.cryptConn.SendPacket(bfw.Data())
		} else {
			ackHandle := bf.ReadUint32()
			data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp1.bin", opcode.String()))
			if err != nil {
				panic(err)
			}

			bfw := byteframe.NewByteFrame()
			bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
			bfw.WriteUint32(ackHandle)
			bfw.WriteBytes(data)
			s.cryptConn.SendPacket(bfw.Data())
		}
		
		if clientEnumerateCount < 1 {
			clientEnumerateCount++
		}
	case network.MSG_MHF_ACQUIRE_MONTHLY_REWARD:
		ackHandle := bf.ReadUint32()
		data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp1.bin", opcode.String()))
		if err != nil {
			panic(err)
		}
		bfw := byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteBytes(data)
		s.cryptConn.SendPacket(bfw.Data())
		
		data, err = ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp0.bin", opcode.String()))
		if err != nil {
			panic(err)
		}
		bfw = byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_UPDATE_RIGHT))
		bfw.WriteBytes(data)
		s.cryptConn.SendPacket(bfw.Data())
	case network.MSG_SYS_RIGHTS_RELOAD:
		ackHandle := bf.ReadUint32()
		//combine the sys update and the actual ack packet
		data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp0.bin", opcode.String()))
		if err != nil {
			panic(err)
		}
		bfw := byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_UPDATE_RIGHT))
		bfw.WriteBytes(data)
		s.cryptConn.SendPacket(bfw.Data())
		
		data, err = ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp1.bin", opcode.String()))
		if err != nil {
			panic(err)
		}
		bfw = byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteBytes(data)
		s.cryptConn.SendPacket(bfw.Data())
	case network.MSG_SYS_ENUMERATE_STAGE:
		ackHandle := bf.ReadUint32()
		_ = bf.ReadUint8()
		check := bf.ReadUint8()

		if check >= 15 {
			//returning id
			id := bf.ReadBytes(10)
			
			data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp1.bin", opcode.String()))
			if err != nil {
				panic(err)
			}

			bfw := byteframe.NewByteFrame()
			bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
			bfw.WriteUint32(ackHandle)
			bfw.WriteBytes(data)
			bfw.WriteBytes(id)
			bfw.WriteUint16(0x0010)
			s.cryptConn.SendPacket(bfw.Data())
		} else {
			//fully canned response
		
			data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp.bin", opcode.String()))
			if err != nil {
				panic(err)
			}

			bfw := byteframe.NewByteFrame()
			bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
			bfw.WriteUint32(ackHandle)
			bfw.WriteBytes(data)
			s.cryptConn.SendPacket(bfw.Data())
		}
	case network.MSG_SYS_SET_STAGE_BINARY:
		// uses acks of follow up MSG_SYS_WAIT_STAGE_BINARY for response so just store entire thing
		bf.Seek(0,io.SeekStart)
		setStageBinary = bf.DataFromCurrent()
		bf.Seek(0,io.SeekStart)
		//fmt.Printf("Seek\n")
		var remainingDataAmount uint = uint(8)
		var doneCheck uint16 = uint16(44);
		for true{
			//check for extra packets
			_ = bf.ReadBytes(5)
			remainingDataAmount = uint(bf.ReadUint16())
			//fmt.Printf("Remaining:   %d", remainingDataAmount)
			_ = bf.ReadBytes(15)
			_ = bf.ReadBytes(remainingDataAmount)
			doneCheck = uint16(bf.ReadUint16())
			//fmt.Printf(" Done:%d\n", (doneCheck))
			bf.Seek(-2,io.SeekCurrent)
			
			if (doneCheck != uint16(44)) {
				break

			}
		}
	case network.MSG_SYS_WAIT_STAGE_BINARY:
		// game insists on these responses being combined if multiple
		setStage := byteframe.NewByteFrameFromBytes(setStageBinary)
		bfw := byteframe.NewByteFrame()
		loopCount := 0;
		// get type of packet
		bf.Seek(0,io.SeekStart)
		_ = bf.ReadBytes(7)
//		waitType := uint(bf.ReadUint8())
		bf.Seek(2,io.SeekStart)
		
		var remainingDataAmount uint = uint(8)
		var doneCheck uint16 = uint16(43)
		
		//fmt.Printf("Wait Type:   %d", waitType)
		
//		if (waitType != 23){
			for true{
				bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
				//ack
				ackHandle := bf.ReadUint32()
				//read until expected end of packet
				_ = bf.ReadBytes(22)
				//done with actual MSG_SYS_WAIT_STAGE_BINARY portion
			
				//Deal with MSG_SYS_SET_STAGE_BINARY packet bytes
				//read until length byte and then read length
				
				
				//check for extra packets
				_ = setStage.ReadBytes(5)
				remainingDataAmount = uint(setStage.ReadUint16())
				//fmt.Printf("Remaining:   %d", remainingDataAmount)
				_ = setStage.ReadBytes(15)
				remainingData := setStage.ReadBytes(remainingDataAmount)
				//doneCheck = uint16(setStage.ReadUint16())
				//fmt.Printf(" Done:%d\n", (doneCheck))
				//setStage.Seek(-2,io.SeekCurrent)
				
				//if (doneCheck != uint16(44)) {
				//	break
				//}
				
				//sys ack + ack from wait_stage + 0x00 0x00 + length of data in set_stage + extra data from set stage + 0x00 0x10
				if loopCount < 2{
					bfw.WriteUint32(ackHandle)
					bfw.WriteUint16(uint16(0x0100))
					bfw.WriteUint16(uint16(remainingDataAmount))
					bfw.WriteBytes(remainingData)
				} else {
					bfw.WriteUint32(ackHandle)
					bfw.WriteUint16(uint16(0x0100))
					bfw.WriteUint16(uint16(remainingDataAmount))
					bfw.WriteUint32(0xF906B21A) // ceaseless character ID because bad
					trimRemain := byteframe.NewByteFrameFromBytes(remainingData)
					_ = trimRemain.ReadBytes(4)
					remainingDataTrim := trimRemain.DataFromCurrent()
					bfw.WriteBytes(remainingDataTrim)
				}
				
				
				// back to header for next packet
				//read next potential header and break if end of packets
				doneCheck = bf.ReadUint16()
				//bf.Seek(-2,io.SeekCurrent)
				
				if doneCheck != 43{
					// end packet and send
					bfw.WriteUint16(0x0010)
					s.cryptConn.SendPacket(bfw.Data())
					// update setStage to have second packet in case there's an ungrouped wait packet afterwards
					setStageBinary = setStage.DataFromCurrent()
					
					break
				}
				loopCount++
			}
			/*
		} else {
			// no interactions with a MSG_SYS_WAIT_STAGE_BINARY packet
			ackHandle := bf.ReadUint32()
			data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp.bin", opcode.String()))
			if err != nil {
				panic(err)
			}

			bfw := byteframe.NewByteFrame()
			bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
			bfw.WriteUint32(ackHandle)
			bfw.WriteBytes(data)
			s.cryptConn.SendPacket(bfw.Data())
		}*/
		//fmt.Printf("\n",hex.Dump((bfw.Data())))
	case network.MSG_SYS_CAST_BINARY:
		_ = bf.ReadBytes(6)
		var remainingData uint = uint(bf.ReadUint16())
		//fmt.Println(remainingData)
		bf.ReadBytes(remainingData)
	case network.MSG_SYS_GET_USER_BINARY:
		_ = bf.ReadBytes(9)
	case network.MSG_SYS_SET_OBJECT_BINARY:
		_ = bf.ReadBytes(4)
		var remainingData uint = uint(bf.ReadUint16())
		//fmt.Println(remainingData)
		bf.ReadBytes(remainingData)
	case network.MSG_SYS_RECORD_LOG:
		ackHandle := bf.ReadUint32()
		_ = bf.ReadBytes(6)
		remainingData := bf.ReadUint16()
		_ = bf.ReadBytes(uint(remainingData))
		
		data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp.bin", opcode.String()))
		if err != nil {
			panic(err)
		}


		bfw := byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteBytes(data)
		s.cryptConn.SendPacket(bfw.Data())	
	case network.MSG_MHF_ADD_ACHIEVEMENT:
		_ = bf.ReadBytes(5)
	case network.MSG_MHF_ENUMERATE_RANKING:
		ackHandle := bf.ReadUint32()
		_ = bf.ReadUint32()
		
		data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp.bin", opcode.String()))
		if err != nil {
			panic(err)
		}
//		fmt.Println("\nread ranking bin\n")

		
		bfw := byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteBytes(data)
		bfw.WriteUint32(uint32(time.Now().Unix()))
		bfw.WriteUint64(0x0001000000000010)
		s.cryptConn.SendPacket(bfw.Data())	
		
		

	case network.MSG_SYS_TERMINAL_LOG:
		ackHandle := bf.ReadUint32()
		_ = bf.ReadUint32()
		// number of 36 byte responses before end of terminal log packet
		totalResponses := bf.ReadUint16()
		_ = bf.ReadUint16()
		_ = bf.ReadBytes(uint(totalResponses * 36))
		
		bfw := byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteUint32(0x00000000)
		// any different value?
		bfw.WriteUint32(uint32(time.Now().Unix()))
		bfw.WriteUint16(0x0010)
		s.cryptConn.SendPacket(bfw.Data())
	case network.MSG_MHF_SAVE_FAVORITE_QUEST:
		fallthrough
	case network.MSG_MHF_SAVE_OTOMO_AIROU:
		fallthrough
	case network.MSG_MHF_SAVE_PARTNER:
		ackHandle := bf.ReadUint32()
		remainingData := bf.ReadUint16()
		remainingDataBytes := bf.ReadBytes(uint(remainingData))

		
		// write local savefile for returning in load requests
		go func() {
			savew := byteframe.NewByteFrame()
			// 00 > Length > actual data
			savew.WriteUint8(0x00)
			savew.WriteUint16(remainingData)
			savew.WriteBytes(remainingDataBytes)
			ioutil.WriteFile(fmt.Sprintf("save_files/%s.bin", opcode.String()), []byte(savew.Data()), 0644)
		}()
		
		data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp.bin", opcode.String()))
		if err != nil {
			panic(err)
		}

		bfw := byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteBytes(data)
		s.cryptConn.SendPacket(bfw.Data())
	case network.MSG_MHF_SAVE_DECO_MYSET:
		fallthrough
	case network.MSG_MHF_SAVE_PLATE_MYSET:
		fallthrough
	case network.MSG_MHF_SAVE_MEZFES_DATA:
		ackHandle := bf.ReadUint32()
		_ = bf.ReadBytes(2)
		remainingData := bf.ReadUint16()
		_ = bf.ReadBytes(uint(remainingData))
		data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp.bin", opcode.String()))
		if err != nil {
			panic(err)
		}


		bfw := byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteBytes(data)
		s.cryptConn.SendPacket(bfw.Data())
	case network.MSG_MHF_SAVEDATA:
		ackHandle := bf.ReadUint32()
		_ = bf.ReadBytes(11)
		remainingData := bf.ReadUint16()
		remainingDataBytes := bf.ReadBytes(uint(remainingData))
		
		// write local savefile for returning in load requests
		go func() {
			savew := byteframe.NewByteFrame()
			// 00 > Length > actual data
			savew.WriteUint8(0x00)
			savew.WriteUint16(remainingData)
			savew.WriteBytes(remainingDataBytes)
			ioutil.WriteFile(fmt.Sprintf("save_files/%s.bin", opcode.String()), []byte(savew.Data()), 0644)
		}()
		
		data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp.bin", opcode.String()))
		if err != nil {
			panic(err)
		}


		bfw := byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteBytes(data)
		s.cryptConn.SendPacket(bfw.Data())
	case network.MSG_MHF_SAVE_SCENARIO_DATA:
		ackHandle := bf.ReadUint32()
		_ = bf.ReadBytes(14)
		data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp.bin", opcode.String()))
		if err != nil {
			panic(err)
		}


		bfw := byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteBytes(data)
		s.cryptConn.SendPacket(bfw.Data())
	case network.MSG_SYS_CREATE_OBJECT:
		ackHandle := bf.ReadUint32()
		_ = bf.ReadBytes(16)
		data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp.bin", opcode.String()))
		if err != nil {
			panic(err)
		}


		bfw := byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteBytes(data)
		s.cryptConn.SendPacket(bfw.Data())
	case network.MSG_SYS_SET_USER_BINARY:
		_ = bf.ReadBytes(1)
		var remainingData uint = uint(bf.ReadUint16())
		bf.ReadBytes(remainingData)
	case network.MSG_MHF_SET_ENHANCED_MINIDATA:
		ackHandle := bf.ReadUint32()
		var remainingData uint = uint(bf.ReadUint16())
		bf.ReadBytes(remainingData)
		
		data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp.bin", opcode.String()))
		if err != nil {
			panic(err)
		}

		bfw := byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_CLEANUP_OBJECT))
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteBytes(data)
		s.cryptConn.SendPacket(bfw.Data())
	case network.MSG_MHF_GET_UD_TACTICS_FOLLOWER:
		_ = bf.ReadBytes(4)
	case network.MSG_HEAD:
	// actually lazily handling bad lengths which end up with null bytes instead of a real header
		remainingData := bf.DataFromCurrent()
		bf.Seek(0, io.SeekStart)
		for i := 0;  i<=len(remainingData); i++ {
			check := bf.ReadUint16()
			if check > 0 {
				bf.Seek(int64(i), io.SeekStart)
				break
			} 
        }
	case network.MSG_SYS_BACK_STAGE:
		ackHandle := bf.ReadUint32()
		data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp.bin", opcode.String()))
		if err != nil {
			panic(err)
		}

		bfw := byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_CLEANUP_OBJECT))
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteBytes(data)
		s.cryptConn.SendPacket(bfw.Data())
	case network.MSG_SYS_reserve203:
		// isn't only ever 5 bytes
		ackHandle := bf.ReadUint32()
		_ = bf.ReadBytes(4)
		
		data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp.bin", opcode.String()))
		if err != nil {
			panic(err)
		}

		bfw := byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_CLEANUP_OBJECT))
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteBytes(data)
		s.cryptConn.SendPacket(bfw.Data())
	case network.MSG_MHF_GET_ENHANCED_MINIDATA:
		ackHandle := bf.ReadUint32()
		_ = bf.ReadBytes(5)
		data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp.bin", opcode.String()))
		if err != nil {
			panic(err)
		}

		bfw := byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_CLEANUP_OBJECT))
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteBytes(data)
		s.cryptConn.SendPacket(bfw.Data())	
	case network.MSG_SYS_MOVE_STAGE:
		ackHandle := bf.ReadUint32()
		_ = bf.ReadBytes(1)
		var remainingData uint = uint(bf.ReadUint8())
		bf.ReadBytes(remainingData)
		
		data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp.bin", opcode.String()))
		if err != nil {
			panic(err)
		}

		bfw := byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_CLEANUP_OBJECT))
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteBytes(data)
		s.cryptConn.SendPacket(bfw.Data())
	case network.MSG_SYS_ENTER_STAGE:
		ackHandle := bf.ReadUint32()
		_ = bf.ReadBytes(17)

		data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp.bin", opcode.String()))
		if err != nil {
			panic(err)
		}

		bfw := byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_CLEANUP_OBJECT))
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteBytes(data)
		s.cryptConn.SendPacket(bfw.Data())
		
	case network.MSG_MHF_GET_EARTH_VALUE:

		ackHandle := bf.ReadUint32()
		bf.ReadBytes(11)
		var respType uint = uint(bf.ReadUint8())

		fmt.Printf("MSG_MHF_GET_EARTH_VALUE: %d\n", respType)

		data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp_0%d.bin", opcode.String(), respType))
		if err != nil {
			panic(err)
		}

		bfw := byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteBytes(data)
		s.cryptConn.SendPacket(bfw.Data())
	case network.MSG_MHF_GET_TOWER_INFO:
		ackHandle := bf.ReadUint32()
		bf.ReadBytes(3)
		var respType uint = uint(bf.ReadUint8())

		fmt.Printf("MSG_MHF_GET_TOWER_INFO: %d\n", respType)

		data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp_0%d.bin", opcode.String(), respType))
		if err != nil {
			panic(err)
		}

		bfw := byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteBytes(data)
		s.cryptConn.SendPacket(bfw.Data())
	case network.MSG_MHF_INFO_FESTA:
		ackHandle := bf.ReadUint32()
		_ = bf.ReadUint32()

		data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp.bin", opcode.String()))
		if err != nil {
			panic(err)
		}

		bfw := byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteBytes(data)
		s.cryptConn.SendPacket(bfw.Data())
	case network.MSG_MHF_INFO_GUILD:
		ackHandle := bf.ReadUint32()
		_ = bf.ReadUint32()

		data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp.bin", opcode.String()))
		if err != nil {
			panic(err)
		}

		bfw := byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteBytes(data)
		s.cryptConn.SendPacket(bfw.Data())
	case network.MSG_MHF_ENUMERATE_GUILD_MEMBER:
		ackHandle := bf.ReadUint32()
		_ = bf.ReadUint32()

		data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp.bin", opcode.String()))
		if err != nil {
			panic(err)
		}

		bfw := byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteBytes(data)
		s.cryptConn.SendPacket(bfw.Data())
	case network.MSG_MHF_STATE_FESTA_G:
		ackHandle := bf.ReadUint32()
		_ = bf.ReadBytes(10)

		data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp.bin", opcode.String()))
		if err != nil {
			panic(err)
		}
	
		bfw := byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteBytes(data)
		s.cryptConn.SendPacket(bfw.Data())
	case network.MSG_MHF_LOADDATA:
		ackHandle := bf.ReadUint32()
		if loadDataCount == 0 {
			//initial response
			data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp%d_0.bin", opcode.String(), loadDataCount))
			if err != nil {
				panic(err)
			}

			bfw := byteframe.NewByteFrame()
			bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
			bfw.WriteUint32(ackHandle)
			bfw.WriteBytes(data)
			s.cryptConn.SendPacket(bfw.Data())
						
						
			data, err = ioutil.ReadFile(fmt.Sprintf("bin_resp/MSG_SYS_RIGHTS_RELOAD_resp0.bin"))
			if err != nil {
				panic(err)
			}
			bfw = byteframe.NewByteFrame()
			bfw.WriteUint16(uint16(network.MSG_SYS_UPDATE_RIGHT))
			bfw.WriteBytes(data)
			s.cryptConn.SendPacket(bfw.Data())
			
		} else {
			/*
			data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp%d.bin", opcode.String(), loadDataCount))
			if err != nil {
				panic(err)
			}

			bfw := byteframe.NewByteFrame()
			bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
			bfw.WriteUint32(ackHandle)
			bfw.WriteBytes(data)
			s.cryptConn.SendPacket(bfw.Data())
			*/
			
			data, err := ioutil.ReadFile(fmt.Sprintf("save_files/MSG_MHF_SAVEDATA.bin"))
			if err != nil {
				panic(err)
			}

			bfw := byteframe.NewByteFrame()
			bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
			bfw.WriteUint32(ackHandle)
			bfw.WriteUint8(0x01)
			bfw.WriteBytes(data)
			bfw.WriteUint16(0x0010)
			s.cryptConn.SendPacket(bfw.Data())
		}
		
		loadDataCount++
		if loadDataCount > 1 {
			loadDataCount = 0
		}
	case network.MSG_MHF_GET_PAPER_DATA:
		ackHandle := bf.ReadUint32()

		data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp%d.bin", opcode.String(), getPaperDataCount))
		if err != nil {
			panic(err)
		}

		bfw := byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteBytes(data)
		s.cryptConn.SendPacket(bfw.Data())

		getPaperDataCount++
		if getPaperDataCount > 9 {
			getPaperDataCount = 0
		}
	case network.MSG_MHF_ENUMERATE_QUEST:
		ackHandle := bf.ReadUint32()
		_ = bf.ReadBytes(5)
		questList := bf.ReadUint8()
		_ = bf.ReadBytes(1)
		if questList == 0 {
			questListCount = 0
		} else if (questList <= 44)  {
			questListCount = 1
		} else if questList <= 88 {
			questListCount = 2
		} else if questList <= 126 {
			questListCount = 3
		} else if questList <= 172 {
			questListCount = 4
		} else {
			questListCount = 0
		}
		data, err := ioutil.ReadFile(fmt.Sprintf("bin_resp/%s_resp%d.bin", opcode.String(), questListCount))
		if err != nil {
			panic(err)
		}
		bfw := byteframe.NewByteFrame()
		bfw.WriteUint16(uint16(network.MSG_SYS_ACK))
		bfw.WriteUint32(ackHandle)
		bfw.WriteBytes(data)
		s.cryptConn.SendPacket(bfw.Data())
	case network.MSG_SYS_POSITION_OBJECT:
		_ = bf.ReadBytes(16)
	default:
		// Get the packet parser and handler for this opcode.
		mhfPkt := mhfpacket.FromOpcode(opcode)
		if mhfPkt == nil {
			fmt.Println("Got opcode which we don't know how to parse, can't parse anymore for this group ", opcode)
			return
		}

		// Parse and handle the packet
		mhfPkt.Parse(bf)
		handlerTable[opcode](s, mhfPkt)
		break
	}

	// If there is more data on the stream that the .Parse method didn't read, then read another packet off it.
	
	remainingData := bf.DataFromCurrent()
	readCheck := len(bf.Data()) - len(remainingData)
	
	if len(remainingData) >= 2 && (opcode == network.MSG_SYS_MOVE_STAGE || opcode == network.MSG_MHF_SAVE_PLATE_MYSET || opcode == network.MSG_MHF_SAVE_OTOMO_AIROU || opcode == network.MSG_SYS_SET_OBJECT_BINARY || opcode == network.MSG_MHF_SAVE_DECO_MYSET || opcode == network.MSG_MHF_SAVE_MEZFES_DATA || opcode == network.MSG_SYS_GET_USER_BINARY || opcode == network.MSG_MHF_SAVEDATA || opcode == network.MSG_MHF_ADD_ACHIEVEMENT || opcode == network.MSG_SYS_RECORD_LOG || opcode == network.MSG_SYS_POSITION_OBJECT || opcode == network.MSG_SYS_TERMINAL_LOG || opcode == network.MSG_SYS_CREATE_OBJECT || opcode == network.MSG_SYS_SET_STAGE_BINARY || opcode == network.MSG_MHF_SET_ENHANCED_MINIDATA || opcode == network.MSG_SYS_WAIT_STAGE_BINARY || opcode == network.MSG_MHF_SAVE_PARTNER || opcode == network.MSG_MHF_SAVE_SCENARIO_DATA || opcode == network.MSG_MHF_GET_KEEP_LOGIN_BOOST_STATUS || opcode == network.MSG_SYS_TIME || opcode == network.MSG_SYS_reserve203 /*|| opcode == network.MSG_HEAD*/ ||opcode == network.MSG_MHF_ENUMERATE_QUEST || opcode == network.MSG_MHF_GET_UD_TACTICS_FOLLOWER || opcode == network.MSG_MHF_INFO_FESTA || opcode == network.MSG_SYS_EXTEND_THRESHOLD || opcode == network.MSG_MHF_STATE_FESTA_G || opcode == network.MSG_SYS_CAST_BINARY || opcode == network.MSG_SYS_ENTER_STAGE || opcode == network.MSG_SYS_SET_USER_BINARY) {
			s.handlePacketGroup(remainingData)
	} else if len(remainingData) >= 2 && (readCheck > 6) {
		fmt.Println("Remaining data catch all on ", opcode)
	} 
}