package channelserver

import (
	"github.com/Andoryuuta/Erupe/network"
	"github.com/Andoryuuta/Erupe/network/mhfpacket"
	"github.com/Andoryuuta/byteframe"
	"time"
)

func handleMsgHead(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve01(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve02(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve03(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve04(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve05(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve06(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve07(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysAddObject(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysDelObject(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysDispObject(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysHideObject(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve0C(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve0D(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve0E(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysExtendThreshold(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysEnd(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysNop(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysAck(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysTerminalLog(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysLogin(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysLogin)

	bf := byteframe.NewByteFrame()
	bf.WriteUint16(uint16(network.MSG_SYS_ACK))
	bf.WriteUint32(pkt.AckHandle)
	bf.WriteUint32(0x00000000)
	bf.WriteUint32(uint32(time.Now().Unix()))
	bf.WriteUint16(0x0010)
	s.cryptConn.SendPacket(bf.Data())
}

func handleMsgSysLogout(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysSetStatus(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysPing(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysPing)

	bf := byteframe.NewByteFrame()
	bf.WriteUint16(uint16(network.MSG_SYS_ACK))
	ack := mhfpacket.MsgSysAck{
		AckHandle: pkt.AckHandle,
		Unk0:      0,
		Unk1:      0,
	}
	ack.Build(bf)
	bf.WriteUint16(0x0010)
	s.cryptConn.SendPacket(bf.Data())
}

func handleMsgSysCastBinary(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysHideClient(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysTime(s *Session, p mhfpacket.MHFPacket) {
	//pkt := p.(*mhfpacket.MsgSysTime)

	bf := byteframe.NewByteFrame()
	bf.WriteUint16(uint16(network.MSG_SYS_TIME))
	resp := mhfpacket.MsgSysTime{
		Unk0:      0,
		Timestamp: uint32(time.Now().Unix()),
	}
	resp.Build(bf)
	bf.WriteUint16(0x0010)
	s.cryptConn.SendPacket(bf.Data())
}

func handleMsgSysCastedBinary(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysGetFile(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysIssueLogkey(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysRecordLog(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysEcho(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysCreateStage(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysStageDestruct(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysEnterStage(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysBackStage(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysMoveStage(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysLeaveStage(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysLockStage(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysUnlockStage(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserveStage(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysUnreserveStage(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysSetStagePass(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysWaitStageBinary(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysSetStageBinary(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysGetStageBinary(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysEnumerateClient(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysEnumerateStage(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysCreateMutex(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysCreateOpenMutex(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysDeleteMutex(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysOpenMutex(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysCloseMutex(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysCreateSemaphore(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysCreateAcquireSemaphore(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysDeleteSemaphore(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysAcquireSemaphore(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReleaseSemaphore(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysLockGlobalSema(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysUnlockGlobalSema(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysCheckSemaphore(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysOperateRegister(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysLoadRegister(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysNotifyRegister(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysCreateObject(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysDeleteObject(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysPositionObject(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysRotateObject(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysDuplicateObject(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysSetObjectBinary(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysGetObjectBinary(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysGetObjectOwner(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysUpdateObjectBinary(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysCleanupObject(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve4A(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve4B(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve4C(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve4D(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve4E(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve4F(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysInsertUser(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysDeleteUser(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysSetUserBinary(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysGetUserBinary(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysNotifyUserBinary(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve55(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve56(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve57(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysUpdateRight(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysAuthQuery(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysAuthData(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysAuthTerminal(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve5C(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysRightsReload(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve5E(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve5F(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSavedata(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfLoaddata(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfListMember(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfOprMember(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEnumerateDistItem(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfApplyDistItem(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfAcquireDistItem(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetDistDescription(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSendMail(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfReadMail(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfListMail(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfOprtMail(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfLoadFavoriteQuest(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSaveFavoriteQuest(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfRegisterEvent(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfReleaseEvent(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfTransitMessage(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve71(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve72(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve73(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve74(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve75(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve76(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve77(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve78(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve79(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve7A(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve7B(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve7C(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgCaExchangeItem(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve7E(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfPresentBox(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfServerCommand(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfShutClient(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfAnnounce(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSetLoginwindow(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysTransBinary(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysCollectBinary(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysGetState(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysSerialize(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysEnumlobby(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysEnumuser(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysInfokyserver(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetCaUniqueID(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSetCaAchievement(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfCaravanMyScore(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfCaravanRanking(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfCaravanMyRank(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfCreateGuild(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfOperateGuild(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfOperateGuildMember(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfInfoGuild(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEnumerateGuild(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfUpdateGuild(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfArrangeGuildMember(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEnumerateGuildMember(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEnumerateCampaign(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfStateCampaign(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfApplyCampaign(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEnumerateItem(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfAcquireItem(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfTransferItem(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfMercenaryHuntdata(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEntryRookieGuild(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEnumerateQuest(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEnumerateEvent(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEnumeratePrice(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEnumerateRanking(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEnumerateOrder(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEnumerateShop(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetExtraInfo(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfUpdateInterior(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEnumerateHouse(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfUpdateHouse(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfLoadHouse(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfOperateWarehouse(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEnumerateWarehouse(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfUpdateWarehouse(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfAcquireTitle(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEnumerateTitle(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEnumerateGuildItem(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfUpdateGuildItem(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEnumerateUnionItem(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfUpdateUnionItem(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfCreateJoint(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfOperateJoint(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfInfoJoint(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfUpdateGuildIcon(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfInfoFesta(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEntryFesta(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfChargeFesta(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfAcquireFesta(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfStateFestaU(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfStateFestaG(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEnumerateFestaMember(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfVoteFesta(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfAcquireCafeItem(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfUpdateCafepoint(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfCheckDailyCafepoint(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetCogInfo(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfCheckMonthlyItem(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfAcquireMonthlyItem(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfCheckWeeklyStamp(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfExchangeWeeklyStamp(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfCreateMercenary(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSaveMercenary(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfReadMercenaryW(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfReadMercenaryM(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfContractMercenary(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEnumerateMercenaryLog(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEnumerateGuacot(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfUpdateGuacot(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfInfoTournament(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEntryTournament(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEnterTournamentQuest(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfAcquireTournament(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetAchievement(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfResetAchievement(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfAddAchievement(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfPaymentAchievement(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfDisplayedAchievement(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfInfoScenarioCounter(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSaveScenarioData(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfLoadScenarioData(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetBbsSnsStatus(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfApplyBbsArticle(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetEtcPoints(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfUpdateEtcPoint(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetMyhouseInfo(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfUpdateMyhouseInfo(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetWeeklySchedule(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEnumerateInvGuild(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfOperationInvGuild(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfStampcardStamp(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfStampcardPrize(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfUnreserveSrg(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfLoadPlateData(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSavePlateData(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfLoadPlateBox(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSavePlateBox(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfReadGuildcard(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfUpdateGuildcard(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfReadBeatLevel(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfUpdateBeatLevel(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfReadBeatLevelAllRanking(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfReadBeatLevelMyRanking(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfReadLastWeekBeatRanking(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfAcceptReadReward(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetAdditionalBeatReward(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetFixedSeibatuRankingTable(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetBbsUserStatus(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfKickExportForce(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetBreakSeibatuLevelReward(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetWeeklySeibatuRankingReward(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetEarthStatus(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfLoadPartner(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSavePartner(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetGuildMissionList(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetGuildMissionRecord(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfAddGuildMissionCount(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSetGuildMissionTarget(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfCancelGuildMissionTarget(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfLoadOtomoAirou(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSaveOtomoAirou(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEnumerateGuildTresure(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEnumerateAiroulist(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfRegistGuildTresure(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfAcquireGuildTresure(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfOperateGuildTresureReport(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetGuildTresureSouvenir(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfAcquireGuildTresureSouvenir(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEnumerateFestaIntermediatePrize(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfAcquireFestaIntermediatePrize(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfLoadDecoMyset(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSaveDecoMyset(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfReserve010F(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfLoadGuildCooking(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfRegistGuildCooking(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfLoadGuildAdventure(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfRegistGuildAdventure(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfAcquireGuildAdventure(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfChargeGuildAdventure(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfLoadLegendDispatch(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfLoadHunterNavi(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSaveHunterNavi(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfRegistSpabiTime(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetGuildWeeklyBonusMaster(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetGuildWeeklyBonusActiveCount(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfAddGuildWeeklyBonusExceptionalUser(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetTowerInfo(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfPostTowerInfo(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetGemInfo(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfPostGemInfo(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetEarthValue(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfDebugPostValue(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetPaperData(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetNotice(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfPostNotice(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetBoostTime(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfPostBoostTime(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetBoostTimeLimit(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfPostBoostTimeLimit(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEnumerateFestaPersonalPrize(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfAcquireFestaPersonalPrize(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetRandFromTable(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetCafeDuration(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetCafeDurationBonusInfo(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfReceiveCafeDurationBonus(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfPostCafeDurationBonusReceived(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetGachaPoint(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfUseGachaPoint(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfExchangeFpoint2Item(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfExchangeItem2Fpoint(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetFpointExchangeList(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfPlayStepupGacha(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfReceiveGachaItem(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetStepupStatus(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfPlayFreeGacha(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetTinyBin(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfPostTinyBin(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetSenyuDailyCount(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetGuildTargetMemberNum(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetBoostRight(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfStartBoostTime(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfPostBoostTimeQuestReturn(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetBoxGachaInfo(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfPlayBoxGacha(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfResetBoxGachaInfo(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetSeibattle(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfPostSeibattle(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetRyoudama(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfPostRyoudama(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetTenrouirai(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfPostTenrouirai(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfPostGuildScout(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfCancelGuildScout(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfAnswerGuildScout(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetGuildScoutList(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetGuildManageRight(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSetGuildManageRight(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfPlayNormalGacha(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetDailyMissionMaster(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetDailyMissionPersonal(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSetDailyMissionPersonal(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetGachaPlayHistory(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetRejectGuildScout(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSetRejectGuildScout(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetCaAchievementHist(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSetCaAchievementHist(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetKeepLoginBoostStatus(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfUseKeepLoginBoost(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetUdSchedule(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetUdInfo(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetKijuInfo(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSetKiju(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfAddUdPoint(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetUdMyPoint(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetUdTotalPointInfo(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetUdBonusQuestInfo(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetUdSelectedColorInfo(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetUdMonsterPoint(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetUdDailyPresentList(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetUdNormaPresentList(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetUdRankingRewardList(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfAcquireUdItem(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetRewardSong(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfUseRewardSong(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfAddRewardSongCount(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetUdRanking(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetUdMyRanking(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfAcquireMonthlyReward(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetUdGuildMapInfo(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGenerateUdGuildMap(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetUdTacticsPoint(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfAddUdTacticsPoint(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetUdTacticsRanking(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetUdTacticsRewardList(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetUdTacticsLog(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetEquipSkinHist(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfUpdateEquipSkinHist(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetUdTacticsFollower(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSetUdTacticsFollower(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetUdShopCoin(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfUseUdShopCoin(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetEnhancedMinidata(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSetEnhancedMinidata(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSexChanger(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetLobbyCrowd(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve180(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGuildHuntdata(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfAddKouryouPoint(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetKouryouPoint(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfExchangeKouryouPoint(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetUdTacticsBonusQuest(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetUdTacticsFirstQuestBonus(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetUdTacticsRemainingPoint(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve188(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfLoadPlateMyset(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSavePlateMyset(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve18B(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetRestrictionEvent(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSetRestrictionEvent(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve18E(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve18F(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetTrendWeapon(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfUpdateUseTrendWeaponLog(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve192(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve193(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve194(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSaveRengokuData(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfLoadRengokuData(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetRengokuBinary(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEnumerateRengokuRanking(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetRengokuRankingRank(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetRengokuRankingRank)

	bf := byteframe.NewByteFrame()
	bf.WriteUint16(uint16(network.MSG_SYS_ACK))
	bf.WriteUint32(pkt.AckHandle)
	bf.WriteBytes([]byte{0x01, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	bf.WriteUint16(0x0010)
	s.cryptConn.SendPacket(bf.Data())
}

func handleMsgMhfAcquireExchangeShop(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve19B(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSaveMezfesData(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfLoadMezfesData(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve19E(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve19F(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfUpdateForceGuildRank(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfResetTitle(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve202(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve203(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve204(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve205(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve206(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve207(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve208(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve209(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve20A(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve20B(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve20C(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve20D(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve20E(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReserve20F(s *Session, p mhfpacket.MHFPacket) {}
