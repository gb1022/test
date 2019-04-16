package match

import (
	"data"
	// "fight"
	"gamenet"
	"gbframe"
	"protof"
	"strconv"
	// "time"
)

type MatchRoom struct {
	Room string
	// who  string //表示首先接收到的是谁的消息 存player的sess值的
	// FighterA     *FightPlayer
	// FighterB     *FightPlayer
	PlayerList []*data.Player
	// FightRecords map[int]*data.FightRecord
	// FightTime    int
}

var MatchRoomPool = map[string]*MatchRoom{}

func FindOrCreateMatchRoom(player *data.Player, msess string) (bool, *MatchRoom) {

	if matchRoom, ok := MatchRoomPool[msess]; !ok {
		mRoom := CreateMatchRoom(player)
		MatchRoomPool[msess] = mRoom
		return false, mRoom
	} else {
		matchRoom.PlayerList = append(matchRoom.PlayerList, player)
		return true, matchRoom
	}

}

func CreateMatchRoom(player *data.Player) *MatchRoom {
	// fight_time := int(time.Now().Unix())

	// fpA := initFightPlayer(player, true)
	ids := strconv.Itoa(player.MatchId)
	roomid := gbframe.MakeSession(ids, "")
	if player.MatchId == 0 {
		roomid = gbframe.MakeSession(ids, player.Name)
	}
	playerlist := []*data.Player{}
	playerlist = append(playerlist, player)
	mr := &MatchRoom{
		Room: roomid,
		// who:          "",
		PlayerList: playerlist,
		// FightRecords: map[int]*FightRecord{},
		// FightTime:    fight_time,
	}

	return mr
}

// func MatchRoom(player *Player) *fight.FightRoom {
// 	//	time.Sleep(10 * time.Second)
// 	ids := strconv.Itoa(player.MatchId)
// 	rsess := gbframe.MakeSession(ids, "")
// 	froom := FightRooms[rsess]
// 	//	for _, p := range MatchPool {
// 	//		if player.MatchId == p.MatchId && player.Name != p.Name {
// 	//			//			froom := initFightRoom(player, p, ids)
// 	//			//			froom := FightRooms[rsess]
// 	//			froom.JoinFightRoom(player)
// 	//			FightRooms[rsess] = froom
// 	//			break
// 	//		}
// 	//	}
// 	froom.JoinFightRoom(player)
// 	FightRooms[rsess] = froom
// 	gbframe.Logger_Info.Println("A name:", froom.FighterA.Playerdata.Name)
// 	gbframe.Logger_Info.Println("B name:", froom.FighterB.Playerdata.Name)
// 	return froom
// 	//	ids := strconv.Itoa(player.MatchId)
// 	//	robotRoomId := MakeSession(ids, "robot")
// 	//	robotPlayer := CreateRobot()                       //创建机器人
// 	//	prfroom := initFightRoom(player, robotPlayer, ids) //创建机器人房间
// 	//	FightRooms[robotRoomId] = prfroom

// 	// func FightReady(mroom *MatchRoom) {
// 	// 	fight.FightReady()

// }

// func MatchRobot(player, s) {
// 	robot := data.CreateRobot()
// 	playerlist := []*Player{}
// 	playerlist = append(playerlist, player)
// 	playerlist = append(playerlist, &robot)

// }

func MatchProcess(player *data.Player, str string) (bool, *MatchRoom) {
	ids := strconv.Itoa(player.MatchId)
	msess := gbframe.MakeSession(ids, str)
	// playerlist := []*data.Player{}
	if res, mroom := FindOrCreateMatchRoom(player, msess); res {
		// FightReady(mroom)
		return true, mroom
	}
	return false, nil

}

func StartMatch(cs_msg *protof.Message1, s *gamenet.Server) (bool, *MatchRoom) {

	isstartA := true
	isstartB := true
	sc_fight_start := &protof.Message1_SC_FightStart{
		IsstartA: &isstartA,
		IsstartB: &isstartB,
	}
	robot := int(cs_msg.CsFightStart.GetTorobot())
	if robot == 2 {
		ok, mroom := matchRobot(cs_msg, s)
		return ok, mroom
	}

	sc_msg := &protof.Message1{
		ScFightStart: sc_fight_start,
	}

	if !cs_msg.CsFightStart.GetIsstart() {
		isstartA = false
	}
	sess := s.Sess
	player := data.GetPlayerBySess(sess)
	if player == nil {
		gbframe.Logger_Info.Println("player is nil,sess:", sess)
		isstartA = false
	}
	if player.MatchFlag {
		gbframe.Logger_Error.Println("this player is matching,do not match again!player name:", player.Name)
		isstartA = false
		isstartB = false
		mid := int(protof.Message1_SC_FIGHTSTART)
		s.WriteMessage(mid, sc_msg)
		return false, nil
	}
	// ids := strconv.Itoa(player.MatchId)
	// rsess := gbframe.MakeSession(ids, "")
	ok, mroom := MatchProcess(player, "")
	if !ok {
		isstartB = false
		mid := int(protof.Message1_SC_FIGHTSTART)
		gbframe.Logger_Info.Println("sc_msg", sc_msg)
		s.WriteMessage(mid, sc_msg)
		return false, nil
	}
	gbframe.Logger_Info.Println("StartMatch sc_msg:", sc_msg)
	for _, p := range mroom.PlayerList {
		ser := gamenet.NetServers[p.Server_sess]
		mid := int(protof.Message1_SC_FIGHTSTART)
		gbframe.Logger_Info.Println("StartMatch server:", ser, " gamenet.NetServers:", gamenet.NetServers, " sess:", p.Server_sess)
		ser.WriteMessage(mid, sc_msg)
		p.MatchFlag = false
	}
	// mid := int(protof.Message1_SC_FIGHTSTART)
	// s.WriteMessage(mid, sc_msg)

	return true, mroom
}

func DeleteMatchRoomPool(mroom *MatchRoom) {

	delete(MatchRoomPool, mroom.Room)
	gbframe.Logger_Info.Println("DeleteMatchRoomPool,MatchRoomPool:", MatchRoomPool)

}

func FindMatchRoomByPlayer(player *data.Player) (bool, *MatchRoom) {
	ids := strconv.Itoa(player.MatchId)
	msess := gbframe.MakeSession(ids, "")
	if ids == "0" {
		msess = gbframe.MakeSession(ids, player.Name)
	}
	if mroom, ok := MatchRoomPool[msess]; ok {
		return true, mroom
	} else {
		return false, nil
	}

}

func matchRobot(cs_msg *protof.Message1, s *gamenet.Server) (bool, *MatchRoom) {
	gbframe.Logger_Info.Println("Match robot begin!player sess :", s.Sess)
	robot := data.CreateRobot()
	isstartA := true
	isstartB := true
	sc_fight_start := &protof.Message1_SC_FightStart{
		IsstartA: &isstartA,
		IsstartB: &isstartB,
	}
	sc_msg := &protof.Message1{
		ScFightStart: sc_fight_start,
	}
	sess := s.Sess
	player := data.GetPlayerBySess(sess)
	if player == nil {
		gbframe.Logger_Info.Println("player is nil,sess:", sess)
		isstartA = false
	}
	if player.MatchFlag {
		gbframe.Logger_Error.Println("this player is matching,do not match again!player name:", player.Name)
		isstartA = false
		isstartB = false
		mid := int(protof.Message1_SC_FIGHTSTART)
		s.WriteMessage(mid, sc_msg)
		return false, nil
	}
	ids := strconv.Itoa(player.MatchId)
	msess := gbframe.MakeSession(ids, player.Name)
	ok, mroom := FindOrCreateMatchRoom(player, msess)
	ok, mroom = FindOrCreateMatchRoom(robot, msess) //将机器人加入匹配房间
	if !ok {
		isstartB = false
		mid := int(protof.Message1_SC_FIGHTSTART)
		gbframe.Logger_Info.Println("sc_msg", sc_msg)
		s.WriteMessage(mid, sc_msg)
		return false, nil
	}

	ser := gamenet.NetServers[player.Server_sess]
	mid := int(protof.Message1_SC_FIGHTSTART)
	gbframe.Logger_Info.Println("StartMatch server:", ser, " gamenet.NetServers:", gamenet.NetServers, " sess:", player.Server_sess)
	ser.WriteMessage(mid, sc_msg)
	player.MatchFlag = false
	return true, mroom
}
