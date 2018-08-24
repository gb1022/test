package main

import (
	//	"fmt"
	//	"net"
	//	"encoding/json"
	"gbframe"
	"protof"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
)

const (
	MAPMAX_X    = 20
	MAPMAX_Y    = 20
	LIFEMAX     = 100
	SPEEDMAX    = 100
	ACTTAKMAX   = 100
	MAXROUNDNUM = 60
)

type FightRoom struct {
	Room         string
	who          string //表示首先接收到的是谁的消息 存player的sess值的
	FighterA     *FightPlayer
	FighterB     *FightPlayer
	FightRecords map[int]*FightRecord
	FightTime    int
}

type FightPlayer struct {
	Playerdata *Player
	MapX       int
	MapY       int
	Life       int
	Acttak     int
	Speed      int
	Result     int
}

type FightRecord struct {
	RoundNum int
	FighterA FightPlayer
	FighterB FightPlayer
}

var MatchPool = map[string]*Player{}

var FightRooms = map[string]*FightRoom{}

func initFightPlayer(player *Player, side bool) *FightPlayer {
	x, y := 1, 1
	if side {
		x = MAPMAX_X
		y = MAPMAX_Y
	}
	fp := &FightPlayer{
		Playerdata: player,
		MapX:       x,
		MapY:       y,
		Life:       LIFEMAX,
		Acttak:     0,
		Speed:      0,
		Result:     0,
	}
	return fp
}

func (fr *FightRoom) JoinFightRoom(playerB *Player) {
	//	fpA := initFightPlayer(playerA, true)
	//	ids := strconv.Itoa(player.MatchId)
	fpB := initFightPlayer(playerB, false)
	fr.FighterB = fpB
}

func CreateFightRoom(player *Player) *FightRoom {
	fight_time := int(time.Now().Unix())

	fpA := initFightPlayer(player, true)
	ids := strconv.Itoa(player.MatchId)
	roomid := gbframe.MakeSession(ids, "")
	fr := &FightRoom{
		Room:         roomid,
		who:          "",
		FighterA:     fpA,
		FightRecords: map[int]*FightRecord{},
		FightTime:    fight_time,
	}

	return fr
}

func MatchRoom(player *Player) *FightRoom {
	//	time.Sleep(10 * time.Second)
	ids := strconv.Itoa(player.MatchId)
	rsess := gbframe.MakeSession(ids, "")
	froom := FightRooms[rsess]
	//	for _, p := range MatchPool {
	//		if player.MatchId == p.MatchId && player.Name != p.Name {
	//			//			froom := initFightRoom(player, p, ids)
	//			//			froom := FightRooms[rsess]
	//			froom.JoinFightRoom(player)
	//			FightRooms[rsess] = froom
	//			break
	//		}
	//	}
	froom.JoinFightRoom(player)
	FightRooms[rsess] = froom
	gbframe.Logger_Info.Println("A name:", froom.FighterA.Playerdata.Name)
	gbframe.Logger_Info.Println("B name:", froom.FighterB.Playerdata.Name)
	return froom
	//	ids := strconv.Itoa(player.MatchId)
	//	robotRoomId := MakeSession(ids, "robot")
	//	robotPlayer := CreateRobot()                       //创建机器人
	//	prfroom := initFightRoom(player, robotPlayer, ids) //创建机器人房间
	//	FightRooms[robotRoomId] = prfroom
}

func FightReady(cs_msg *protof.Message1, s *Server) { //这里面的A,B 并不是FightRoom中的AB，只是区别己方和对方
	isstartA := true
	isstartB := true
	if !cs_msg.CsFightStart.GetIsstart() {
		isstartA = false
	}
	sess := s.sess
	player := GetPlayerBySess(sess)
	if player == nil {
		gbframe.Logger_Info.Println("player is nil,sess:", sess)
		isstartA = false
	}
	ids := strconv.Itoa(player.MatchId)
	rsess := gbframe.MakeSession(ids, "")
	//	ids := strconv.Itoa(player.MatchId)
	//	rsess := MakeSession(ids, "")
	froom, ok := FightRooms[rsess]
	if !ok {
		froom = CreateFightRoom(player)
		FightRooms[rsess] = froom
		isstartB = false
	} else {
		//		c := *player.Conn
		gbframe.Logger_Info.Println("player sess:", sess)
		froom = MatchRoom(player)
	}

	sc_fight_start := &protof.Message1_SC_FightStart{
		IsstartA: &isstartA,
		IsstartB: &isstartB,
	}
	sc_msg := &protof.Message1{
		ScFightStart: sc_fight_start,
	}
	mid := int(protof.Message1_SC_FIGHTSTART)
	s.WriteMessage(mid, sc_msg)

	if isstartB && isstartA {
		var sess_B string
		if sess == froom.FighterA.Playerdata.Server_sess {
			sess_B = froom.FighterB.Playerdata.Server_sess

		} else {
			sess_B = froom.FighterA.Playerdata.Server_sess
		}
		serverB := NetServers[sess_B]
		serverB.WriteMessage(mid, sc_msg)
		time.Sleep(10 * time.Second)
		froom.FightStart()
	}

}

func (fr *FightRoom) FightStart() {
	Aserver := NetServers[fr.FighterA.Playerdata.Server_sess]
	Amx := proto.Int32(int32(fr.FighterA.MapX))
	Amy := proto.Int32(int32(fr.FighterA.MapY))
	Alife := proto.Int32(int32(fr.FighterA.Life))
	Bserver := NetServers[fr.FighterB.Playerdata.Server_sess]
	Bmx := proto.Int32(int32(fr.FighterB.MapX))
	Bmy := proto.Int32(int32(fr.FighterB.MapY))
	Blife := proto.Int32(int32(fr.FighterB.Life))
	fsd_A := &protof.Message1_FightStateData{
		Name: &fr.FighterA.Playerdata.Name,
		MapX: Amx,
		MapY: Amy,
		Life: Alife,
	}
	fsd_B := &protof.Message1_FightStateData{
		Name: &fr.FighterB.Playerdata.Name,
		MapX: Bmx,
		MapY: Bmy,
		Life: Blife,
	}
	round := 1
	roundnum := proto.Int32(int32(round))
	result := proto.Int32(0)
	mid := int(protof.Message1_SC_FIGHTDATA)
	//给A发
	A_sc_fight_data := &protof.Message1_SC_FightData{
		Round:     roundnum,
		Result:    result,
		MySide:    fsd_A,
		OtherSide: fsd_B,
	}
	A_sc_msg := &protof.Message1{
		ScFightData: A_sc_fight_data,
	}
	Aserver.WriteMessage(mid, A_sc_msg)
	gbframe.Logger_Info.Println("A server RemoteIP:", Aserver.Service.Conn.RemoteAddr().String(), "A name:", fr.FighterA.Playerdata.Name)
	//给B发
	B_sc_fight_data := &protof.Message1_SC_FightData{
		Round:     roundnum,
		Result:    result,
		MySide:    fsd_B,
		OtherSide: fsd_A,
	}
	B_sc_msg := &protof.Message1{
		ScFightData: B_sc_fight_data,
	}
	Bserver.WriteMessage(mid, B_sc_msg)
	gbframe.Logger_Info.Println("B server RemoteIP:", Bserver.Service.Conn.RemoteAddr().String(), "B name:", fr.FighterB.Playerdata.Name)
	fight_record := &FightRecord{
		RoundNum: round,
		FighterA: *fr.FighterA,
		FighterB: *fr.FighterB,
	}
	fr.FightRecords[round] = fight_record

	//保存战斗到redis
	//	fight_time_str := strconv.Itoa(int(fr.FightTime))
	//	key := FIGHTRECORDKEY + "_" + fr.Room + "_" + fight_time_str
	//	b_fightdata, err1 := json.Marshal(fight_record)
	//	if err1 != nil {
	//		gbframe.Logger_Error.Println("FightStart json.Marshal is err:", err1)
	//		return
	//	}
	//	gameredis.redisClient.Cmd("HSET", key, round, b_fightdata)
	gameredis.FightRecordSave(*fr, round, *fight_record)
}

func (fr *FightRoom) FightProsses(cs_msg *protof.Message1, player *Player, room_sess string) {
	gbframe.Logger_Info.Println("fight msg:", cs_msg)
	if fr.who == "" {
		fr.who = player.Server_sess
		if fr.FighterA.Playerdata.Server_sess == fr.who {
			fr.FighterA.Acttak = int(*cs_msg.CsFightData.Attack)
			fr.FighterA.MapX = Move(fr.FighterA.MapX, int(*cs_msg.CsFightData.MoveX))
			fr.FighterA.MapY = Move(fr.FighterA.MapY, int(*cs_msg.CsFightData.MoveY))
			fr.FighterA.Speed = int(*cs_msg.CsFightData.Speed)
			//			fr.who = player.Server_sess

		} else if fr.FighterB.Playerdata.Server_sess == fr.who {
			fr.FighterB.Acttak = int(*cs_msg.CsFightData.Attack)
			fr.FighterB.MapX = Move(fr.FighterB.MapX, int(*cs_msg.CsFightData.MoveX))
			fr.FighterB.MapY = Move(fr.FighterB.MapY, int(*cs_msg.CsFightData.MoveY))
			fr.FighterB.Speed = int(*cs_msg.CsFightData.Speed)
			//			fr.who = player.Server_sess
		} else {
			gbframe.Logger_Info.Println("player server sess is err!sess:", player.Server_sess)
			return
		}
	} else if fr.who != player.Server_sess && fr.who != "" {
		fr.who = ""
		if fr.FighterA.Playerdata.Server_sess == player.Server_sess {
			fr.FighterA.Acttak = int(*cs_msg.CsFightData.Attack)
			fr.FighterA.MapX = Move(fr.FighterA.MapX, int(*cs_msg.CsFightData.MoveX))
			fr.FighterA.MapY = Move(fr.FighterA.MapY, int(*cs_msg.CsFightData.MoveY))
			fr.FighterA.Speed = int(*cs_msg.CsFightData.Speed)
			//			fr.who = player.Server_sess

		} else if fr.FighterB.Playerdata.Server_sess == player.Server_sess {
			fr.FighterB.Acttak = int(*cs_msg.CsFightData.Attack)
			fr.FighterB.MapX = Move(fr.FighterB.MapX, int(*cs_msg.CsFightData.MoveX))
			fr.FighterB.MapY = Move(fr.FighterB.MapY, int(*cs_msg.CsFightData.MoveY))
			fr.FighterB.Speed = int(*cs_msg.CsFightData.Speed)
			//			fr.who = player.Server_sess
		} else {

			gbframe.Logger_Error.Println("fr.who is err!who:", fr.who)
			return
		}
		fr.FightLogic()
	} else if fr.who == player.Server_sess {
		return
	}
	FightRooms[room_sess] = fr
	if fr.FighterA.Result != 0 && fr.FighterB.Result != 0 { //如果得出结果了，那么战斗结束，删除房间
		fr.FightOver()
	}
}

func (fr *FightRoom) FightLogic() {
	fighterA := fr.FighterA
	fighterB := fr.FighterB
	round := len(fr.FightRecords) + 1
	//	result := 0
	if fighterA.Acttak >= fighterB.Life {
		fighterB.Life = 0
	} else {
		fighterB.Life = fighterB.Life - fighterA.Acttak
	}
	if fighterA.Life <= fighterB.Acttak {
		fighterA.Life = 0
	} else {
		fighterA.Life = fighterA.Life - fighterB.Acttak
	}
	if fighterA.Life == 0 && fighterB.Life == 0 {
		fighterA.Result = 3
		fighterB.Result = 3
	} else if fighterA.Life == 0 && fighterB.Life > 0 {
		fighterA.Result = 2
		fighterB.Result = 1
	} else if fighterA.Life > 0 && fighterB.Life == 0 {
		fighterA.Result = 1
		fighterB.Result = 2
	} else {
		fighterA.Result = 0
		fighterB.Result = 0
		if round >= MAXROUNDNUM {
			fighterA.Result = 3
			fighterB.Result = 3
		}
	}

	fr.WriteFightDataMsg(round) //发送战斗消息给客户端
	//添加记录
	record := &FightRecord{
		RoundNum: round,
		FighterA: *fighterA,
		FighterB: *fighterB,
	}
	fr.FightRecords[round] = record
	gameredis.FightRecordSave(*fr, round, *record)
}

func (fr *FightRoom) WriteFightDataMsg(round int) {
	//	gbframe.Logger_Info.Println("")
	mid := int(protof.Message1_SC_FIGHTDATA)
	//先处理A的消息
	my_reslut := proto.Int32(int32(fr.FighterA.Result))
	Amx := proto.Int32(int32(fr.FighterA.MapX))
	Amy := proto.Int32(int32(fr.FighterA.MapY))
	Alife := proto.Int32(int32(fr.FighterA.Life))
	//	Aspeed := proto.Int32(int32(fr.FighterA.Speed))
	Bmx := proto.Int32(int32(fr.FighterB.MapX))
	Bmy := proto.Int32(int32(fr.FighterB.MapY))
	Blife := proto.Int32(int32(fr.FighterB.Life))
	my_fight_state_data := &protof.Message1_FightStateData{
		Name: &fr.FighterA.Playerdata.Name,
		MapX: Amx,
		MapY: Amy,
		Life: Alife,
	}
	other_fight_state_data := &protof.Message1_FightStateData{
		Name: &fr.FighterB.Playerdata.Name,
		MapX: Bmx,
		MapY: Bmy,
		Life: Blife,
	}
	A_sc_fightdata := &protof.Message1_SC_FightData{
		Round:     proto.Int32(int32(round)),
		Result:    my_reslut,
		MySide:    my_fight_state_data,
		OtherSide: other_fight_state_data,
	}
	A_sc_msg := &protof.Message1{
		ScFightData: A_sc_fightdata,
	}
	Aserver := NetServers[fr.FighterA.Playerdata.Server_sess]
	Aserver.WriteMessage(mid, A_sc_msg)

	//先处理B的消息
	my_reslut = proto.Int32(int32(fr.FighterB.Result))
	Amx = proto.Int32(int32(fr.FighterA.MapX))
	Amy = proto.Int32(int32(fr.FighterA.MapY))
	Alife = proto.Int32(int32(fr.FighterA.Life))
	//	Aspeed := proto.Int32(int32(fr.FighterA.Speed))
	Bmx = proto.Int32(int32(fr.FighterB.MapX))
	Bmy = proto.Int32(int32(fr.FighterB.MapY))
	Blife = proto.Int32(int32(fr.FighterB.Life))
	other_fight_state_data = &protof.Message1_FightStateData{
		Name: &fr.FighterA.Playerdata.Name,
		MapX: Amx,
		MapY: Amy,
		Life: Alife,
	}
	my_fight_state_data = &protof.Message1_FightStateData{
		Name: &fr.FighterB.Playerdata.Name,
		MapX: Bmx,
		MapY: Bmy,
		Life: Blife,
	}
	B_sc_fightdata := &protof.Message1_SC_FightData{
		Round:     proto.Int32(int32(round)),
		Result:    my_reslut,
		MySide:    my_fight_state_data,
		OtherSide: other_fight_state_data,
	}
	B_sc_msg := &protof.Message1{
		ScFightData: B_sc_fightdata,
	}
	Bserver := NetServers[fr.FighterB.Playerdata.Server_sess]
	Bserver.WriteMessage(mid, B_sc_msg)
}

//移动
func Move(firstpoint, movedance int) int {
	point := firstpoint + movedance
	if point >= 1 && point <= 20 {
		return point
	} else if point < 1 {
		point = 1
	} else if point > 20 {
		point = 20
	}
	return point
}

//战斗结束，删除房间
func (fr *FightRoom) FightOver() {
	gbframe.Logger_Info.Println("Fight is over!Close Room!roomid:", fr.Room)
	if fr.FighterA.Result == 1 {
		fr.FighterA.Playerdata.Score += 3
	} else if fr.FighterA.Result == 3 {
		fr.FighterA.Playerdata.Score += 1
	}
	if fr.FighterB.Result == 1 {
		fr.FighterB.Playerdata.Score += 3
	} else if fr.FighterB.Result == 3 {
		fr.FighterB.Playerdata.Score += 1
	}
	//保存到redis
	gameredis.PlayerDataSave(*fr.FighterA.Playerdata)
	fr.FighterA.Playerdata.SaveScoreInRank() //保存排行榜
	gameredis.PlayerDataSave(*fr.FighterB.Playerdata)
	fr.FighterB.Playerdata.SaveScoreInRank() //保存排行榜
	fr.showRecords()
	delete(FightRooms, fr.Room)
}

func (fr *FightRoom) showRecords() {
	for _, data := range fr.FightRecords {
		gbframe.Logger_Info.Println("Record:", data)
	}

}
