package main

import (
	"fmt"
	//	"net"
	"protof"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
)

const (
	MAPMAX_X  = 20
	MAPMAX_Y  = 20
	LIFEMAX   = 100
	SPEEDMAX  = 100
	ACTTAKMAX = 100
)

type FightRoom struct {
	Room         string
	FighterA     *FightPlayer
	FighterB     *FightPlayer
	FightRecords map[int]*FightRecord
}

type FightPlayer struct {
	Playerdata *Player
	MapX       int
	MapY       int
	Life       int
	Acttak     int
	Speed      int
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
	fpA := initFightPlayer(player, true)
	ids := strconv.Itoa(player.MatchId)
	roomid := MakeSession(ids, "")
	fr := &FightRoom{
		Room:         roomid,
		FighterA:     fpA,
		FightRecords: map[int]*FightRecord{},
	}

	return fr
}

func MatchRoom(player *Player) *FightRoom {
	//	time.Sleep(10 * time.Second)
	ids := strconv.Itoa(player.MatchId)
	rsess := MakeSession(ids, "")
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
	fmt.Println("A name:", froom.FighterA.Playerdata.Name)
	fmt.Println("B name:", froom.FighterB.Playerdata.Name)
	return froom
	//	ids := strconv.Itoa(player.MatchId)
	//	robotRoomId := MakeSession(ids, "robot")
	//	robotPlayer := CreateRobot()                       //创建机器人
	//	prfroom := initFightRoom(player, robotPlayer, ids) //创建机器人房间
	//	FightRooms[robotRoomId] = prfroom
}

func (s *Server) FightReady(cs_msg *protof.Message1) { //这里面的A,B 并不是FightRoom中的AB，只是区别己方和对方
	isstartA := true
	isstartB := true
	if !cs_msg.CsFightStart.GetIsstart() {
		isstartA = false
	}
	sess := s.sess
	player := GetPlayerBySess(sess)
	if player == nil {
		fmt.Println("player is nil,sess:", sess)
		isstartA = false
	}
	ids := strconv.Itoa(player.MatchId)
	rsess := MakeSession(ids, "")
	//	ids := strconv.Itoa(player.MatchId)
	//	rsess := MakeSession(ids, "")
	froom, ok := FightRooms[rsess]
	if !ok {
		froom = CreateFightRoom(player)
		FightRooms[rsess] = froom
		isstartB = false
	} else {
		//		c := *player.Conn
		fmt.Println("player sess:", sess)
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

func (r *FightRoom) FightStart() {
	Aserver := NetServers[r.FighterA.Playerdata.Server_sess]
	Amx := proto.Int32(int32(r.FighterA.MapX))
	Amy := proto.Int32(int32(r.FighterA.MapY))
	Alife := proto.Int32(int32(r.FighterA.Life))
	Bserver := NetServers[r.FighterB.Playerdata.Server_sess]
	Bmx := proto.Int32(int32(r.FighterB.MapX))
	Bmy := proto.Int32(int32(r.FighterB.MapY))
	Blife := proto.Int32(int32(r.FighterB.Life))
	fsd_A := &protof.Message1_FightStateData{
		Name: &r.FighterA.Playerdata.Name,
		MapX: Amx,
		MapY: Amy,
		Life: Alife,
	}
	fsd_B := &protof.Message1_FightStateData{
		Name: &r.FighterB.Playerdata.Name,
		MapX: Bmx,
		MapY: Bmy,
		Life: Blife,
	}
	round := 0
	roundnum := proto.Int32(0)
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
	fmt.Println("A server RemoteIP:", Aserver.Service.Conn.RemoteAddr().String(), "A name:", r.FighterA.Playerdata.Name)
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
	fmt.Println("B server RemoteIP:", Bserver.Service.Conn.RemoteAddr().String(), "B name:", r.FighterB.Playerdata.Name)
	fight_record := &FightRecord{
		RoundNum: round,
		FighterA: *r.FighterA,
		FighterB: *r.FighterB,
	}
	r.FightRecords[round] = fight_record
}
