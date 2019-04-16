package fight_tunnel_capture

import (
	"db"
	"encoding/json"
	"gamenet"
	"math"
	"math/rand"

	//	"fmt"
	//	"net"
	//	"encoding/json"
	"capture_robot"
	"data"
	"gbframe"
	"protof"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
)

const (
	MAPMAX_X      = 20 //x轴大小
	MAPMAX_Y      = 20 //y轴大小
	MAXROUNDNUM   = 40 //总回合数
	MAXPOWERLIMIT = 6  //能量总数
	ADDPOWER      = 2  //每回合增加的能量点
)

type FightRoom struct {
	Room         string
	who          bool //表示该谁活动，true表示A,false表示B,第一次是A先
	RoundNum     int  //回合
	FighterA     *FightPlayer
	FighterB     *FightPlayer
	FightRecords map[int]*FightRecord
	FightTime    int
	Result       int //0：没结束，1：A胜利，2：B胜利，3：平局
}

type Map struct {
	x int
	y int
}
type FightPlayer struct {
	Playerdata *data.Player
	Name       string
	Fight_Data *FightData
	Score      int
	ServerSess string
}
type FightData struct {
	BirthPoint          *Map   //出生节点
	ExcavatePoits       []*Map //自己已挖掘的节点
	OtherExcavatePoints []*Map //对方已挖掘的节点
	MovePoints          []*Map //移动
	OtherBirthPoint     *Map   //对方的出生节点
	UserAtPoint         *Map   //当前所在的节点
	OtherUserAtPoint    *Map   //对方所在的节点 ps:（被发现才会有值，未被发现则节点为-1,-1）
	Power               int    //能量点
	OtherPower          int
}

type FightRecord struct { //战斗记录
	RoundNum int
	FighterA *FightPlayer
	FighterB *FightPlayer
}

var MatchPool = map[string]*data.Player{}

var FightRooms = map[string]*FightRoom{}

//创建战斗房间
func CreateFightRoom(playerlist []*data.Player) *FightRoom {
	fight_time := int(time.Now().Unix())
	birth_point_A, birth_point_B := initBirthPoint()

	////////////////////
	// birth_point_A = &Map{
	// 	x: 1,
	// 	y: 1,
	// }
	// birth_point_B = &Map{
	// 	x: 2,
	// 	y: 2,
	// }
	/////////////////////
	fpA := initFightPlayer(playerlist[0], true, birth_point_A, birth_point_B)
	fpB := initFightPlayer(playerlist[1], false, birth_point_B, birth_point_A)
	ids := strconv.Itoa(playerlist[0].MatchId)
	roomid := gbframe.MakeSession(ids, "")
	if playerlist[0].MatchId == 0 {
		roomid = gbframe.MakeSession(ids, playerlist[0].Name)
	}
	fr := &FightRoom{
		Room:         roomid,
		who:          true,
		RoundNum:     0,
		FighterA:     fpA,
		FighterB:     fpB,
		FightRecords: map[int]*FightRecord{},
		FightTime:    fight_time,
		Result:       0,
	}
	return fr
}

func MakeRandPoint() (int, int, int, int) {

	rand.Seed(time.Now().UnixNano())
	x := rand.Intn(MAPMAX_X)
	y := rand.Intn(MAPMAX_Y / 2)
	a := rand.Intn(MAPMAX_X)
	b := rand.Intn(MAPMAX_Y/2) + MAPMAX_Y/2
	return x, y, a, b
}

func initBirthPoint() (*Map, *Map) {
	x_a, y_a, x_b, y_b := MakeRandPoint()
	for {
		if math.Abs(float64(y_a-y_b)) < 10 || math.Abs(float64(x_b-x_a)) < 10 {
			x_a, y_a, x_b, y_b = MakeRandPoint()
		} else {
			break
		}

	}

	var Map_a, Map_b Map
	Map_a.x = x_a
	Map_a.y = y_a
	Map_b.x = x_b
	Map_b.y = y_b
	return &Map_a, &Map_b

}

func initFightData(birth_point, other_point *Map) *FightData {
	var fmaps = []*Map{}
	// o_b_map := &Map{
	// 	x: -1,
	// 	y: -1,
	// }
	o_at_map := &Map{
		x: -1,
		y: -1,
	}
	at_point := &Map{
		x: birth_point.x,
		y: birth_point.y,
	}
	fd := &FightData{
		BirthPoint:          birth_point,   //出生节点
		ExcavatePoits:       fmaps,         //已挖掘的节点
		OtherExcavatePoints: fmaps,         //对方已挖掘的节点
		MovePoints:          fmaps,         //移动
		OtherBirthPoint:     other_point,   //对方的出生节点
		UserAtPoint:         at_point,      //当前所在的节点
		OtherUserAtPoint:    o_at_map,      //对方所在的节点 ps:（被发现才会有值，未被发现则节点为-1,-1）
		Power:               MAXPOWERLIMIT, //能量点
		OtherPower:          MAXPOWERLIMIT, //对手的剩余能量点
	}
	return fd
}

func initFightPlayer(player *data.Player, side bool, birth_point, other_point *Map) *FightPlayer { //side表示是主队还是客队
	// x, y := 1, 1
	// if side {
	// 	x = MAPMAX_X
	// 	y = MAPMAX_Y
	// }
	fight_data := initFightData(birth_point, other_point)
	fp := &FightPlayer{
		Playerdata: player,
		Score:      0,
		Name:       player.Name,
		Fight_Data: fight_data,
		ServerSess: player.Server_sess,
	}
	return fp
}

//战斗准备，创建战斗房间
func FightReady(playerlist []*data.Player) {
	if len(playerlist) != 2 {
		gbframe.Logger_Error.Println("playerlist len is != 2!playerlist:", playerlist)
		return
	}
	ids := strconv.Itoa(playerlist[0].MatchId)
	rsess := gbframe.MakeSession(ids, "")

	if playerlist[0].MatchId == 0 {
		rsess = gbframe.MakeSession(ids, playerlist[0].Name) // 创建与机器人对战的session
	}
	froom, ok := FightRooms[rsess]
	if !ok {
		froom = CreateFightRoom(playerlist)
		FightRooms[rsess] = froom
	} else {
		gbframe.Logger_Error.Println("this fightRoom is exist!fightroom sess:", rsess)
		return
	}
	froom.FightStart()
}

func (fr *FightRoom) FightStart() {
	Aserver := gamenet.NetServers[fr.FighterA.Playerdata.Server_sess]
	robot_flag := false
	Bserver, ok := gamenet.NetServers[fr.FighterB.Playerdata.Server_sess]
	if !ok {
		if fr.FighterB.Playerdata.IsRobot() {
			Bserver = nil
			robot_flag = true
		} else {
			return
		}
	}

	init_Map_Info := &protof.Message1_Map_Info{
		X: proto.Int32(int32(-1)),
		Y: proto.Int32(int32(-1)),
	}
	BirthPoint_A := &protof.Message1_Map_Info{
		X: proto.Int32(int32(fr.FighterA.Fight_Data.BirthPoint.x)),
		Y: proto.Int32(int32(fr.FighterA.Fight_Data.BirthPoint.y)),
	}
	BirthPoint_B := &protof.Message1_Map_Info{
		X: proto.Int32(int32(fr.FighterB.Fight_Data.BirthPoint.x)),
		Y: proto.Int32(int32(fr.FighterB.Fight_Data.BirthPoint.y)),
	}
	fsd_A := &protof.Message1_FightStateData_Capture{
		Name:             &fr.FighterA.Playerdata.Name,
		OtherName:        &fr.FighterB.Playerdata.Name,
		BirthPoint:       BirthPoint_A,
		OtherBirthPoint:  BirthPoint_B,
		UserAtPoint:      BirthPoint_A,
		OtherUserAtPoint: init_Map_Info,
		LastPower:        proto.Int32(MAXPOWERLIMIT),
		OtherPower:       proto.Int32(MAXPOWERLIMIT),
	}
	fsd_B := &protof.Message1_FightStateData_Capture{
		Name:             &fr.FighterB.Playerdata.Name,
		OtherName:        &fr.FighterA.Playerdata.Name,
		BirthPoint:       BirthPoint_B,
		OtherBirthPoint:  BirthPoint_A,
		UserAtPoint:      BirthPoint_B,
		OtherUserAtPoint: init_Map_Info,
		LastPower:        proto.Int32(MAXPOWERLIMIT),
		OtherPower:       proto.Int32(MAXPOWERLIMIT),
	}
	round := 0
	roundnum := proto.Int32(int32(round))
	result := proto.Int32(0)
	mid := int(protof.Message1_SC_FIGHTDATA_TUNNEL_CAPTURE)

	//给A发
	A_sc_fight_data := &protof.Message1_SC_FightData_Tunnel_Capture{
		Code:           proto.Int32(0),
		Round:          roundnum,
		Result:         result,
		FightStateData: fsd_A,
		IsFight:        proto.Bool(true),
	}
	A_sc_msg := &protof.Message1{
		ScFightDataTunnelCapture: A_sc_fight_data,
	}
	B_sc_fight_data := &protof.Message1_SC_FightData_Tunnel_Capture{
		Code:           proto.Int32(0),
		Round:          roundnum,
		Result:         result,
		FightStateData: fsd_B,
		IsFight:        proto.Bool(false),
	}
	B_sc_msg := &protof.Message1{
		ScFightDataTunnelCapture: B_sc_fight_data,
	}
	Aserver.WriteMessage(mid, A_sc_msg)
	if robot_flag {
		capture_robot.RobotFight(B_sc_msg)
		gbframe.Logger_Info.Println("A server RemoteIP:", Aserver.Service.Conn.RemoteAddr().String(), "A name:", fr.FighterA.Playerdata.Name)
		// gbframe.Logger_Info.Println("B server RemoteIP:", Bserver.Service.Conn.RemoteAddr().String(), "B name:", fr.FighterB.Playerdata.Name)
		gbframe.Logger_Info.Println("A_sc_msg:", *A_sc_msg)
		gbframe.Logger_Info.Println("Robot_sc_msg:", *B_sc_msg)
		fr.SaveFightRecordData()
		return
	} else {
		Bserver.WriteMessage(mid, B_sc_msg)
	}

	gbframe.Logger_Info.Println("A server RemoteIP:", Aserver.Service.Conn.RemoteAddr().String(), "A name:", fr.FighterA.Playerdata.Name)
	gbframe.Logger_Info.Println("B server RemoteIP:", Bserver.Service.Conn.RemoteAddr().String(), "B name:", fr.FighterB.Playerdata.Name)
	gbframe.Logger_Info.Println("A_sc_msg:", *A_sc_msg)
	gbframe.Logger_Info.Println("B_sc_msg:", *B_sc_msg)
	fr.SaveFightRecordData()
}

func (fr *FightRoom) SaveFightRecordData() {
	// fight_time_str := strconv.Itoa(int(fr.FightTime))
	fight_record := &FightRecord{
		RoundNum: fr.RoundNum,
		FighterA: fr.FighterA,
		FighterB: fr.FighterB,
	}
	fr.FightRecords[fr.RoundNum] = fight_record
	fightTime := int(fr.FightTime)
	b_fightdata, err1 := json.Marshal(fight_record)
	if err1 != nil {
		gbframe.Logger_Error.Println("SaveFightRecordData json.Marshal is err:", err1)
		return
	}
	db.Gameredis.FightRecordSave(b_fightdata, fr.RoundNum, fr.Room, fightTime)
}

func (fr *FightRoom) FightProsses(cs_msg *protof.Message1, player *data.Player, room_sess string) (*protof.Message1, bool) {
	code := 0
	gbframe.Logger_Info.Println("fight msg:", cs_msg)
	gbframe.Logger_Info.Println("FightRoom Aname:", fr.FighterA.Name, " Bname:", fr.FighterB.Name)
	gbframe.Logger_Info.Println("ServerA:", fr.FighterA.ServerSess, " ServerB:", fr.FighterB.ServerSess)
	// result := 0
	if fr.who && player.Name == fr.FighterA.Name { //这个是主场玩家的消息
		fr.RoundNum += 1
		code = fr.FighterA.UpdateIsFightPlayerInfo(cs_msg)
		if code == 0 {
			fr.FighterB.UpdateNotFightPlayerInfo(fr.FighterA)
			fr.FighterA.GetOtherBirthPointOrAtPoint(*fr.FighterB.Fight_Data.BirthPoint, *fr.FighterB.Fight_Data.UserAtPoint)
		}

	} else if !fr.who && player.Name == fr.FighterB.Name { //这个是客场玩家的消息
		code = fr.FighterB.UpdateIsFightPlayerInfo(cs_msg)
		if code == 0 {
			fr.FighterA.UpdateNotFightPlayerInfo(fr.FighterB)
			fr.FighterB.GetOtherBirthPointOrAtPoint(*fr.FighterA.Fight_Data.BirthPoint, *fr.FighterA.Fight_Data.UserAtPoint)
		} else {
			gbframe.Logger_Error.Println("UpdateIsFightPlayerInfo is error code:", code, "player.Name:", player.Name)
		}
	} else if !fr.who && fr.FighterB.Name == "robot" { //这个是客场机器人的消息
		code = fr.FighterB.UpdateIsFightPlayerInfo(cs_msg)
		gbframe.Logger_Info.Println("111111111111111111111111222222222222223333333333")
		if code == 0 {
			fr.FighterA.UpdateNotFightPlayerInfo(fr.FighterB)
			fr.FighterB.GetOtherBirthPointOrAtPoint(*fr.FighterA.Fight_Data.BirthPoint, *fr.FighterA.Fight_Data.UserAtPoint)
		} else {
			gbframe.Logger_Error.Println("UpdateIsFightPlayerInfo is error code:", code, "player.Name:", player.Name)
		}
	} else {
		gbframe.Logger_Error.Println("fight player is not right!fr.who:", fr.who, "player.Name:", player.Name)
		code = 3
	}

	if fr.who == true {
		result_a, result_b := fr.GetResult(1)
		fr.FighterA.Fight_Data.OtherPower = fr.FighterB.Fight_Data.Power
		fr.FighterB.Fight_Data.OtherPower = fr.FighterA.Fight_Data.Power
		sc_fp_data_a := ServerToProtoPoints(fr.FighterA, fr.FighterB.Name)
		sc_fp_data_b := ServerToProtoPoints(fr.FighterB, fr.FighterA.Name)
		sc_msg_cap_a := &protof.Message1_SC_FightData_Tunnel_Capture{
			Code:           proto.Int32(int32(code)),
			Round:          proto.Int32(int32(fr.RoundNum)),
			Result:         proto.Int32(int32(result_a)),
			FightStateData: sc_fp_data_a,
			IsFight:        proto.Bool(false),
		}
		sc_msg_cap_b := &protof.Message1_SC_FightData_Tunnel_Capture{
			Code:           proto.Int32(int32(code)),
			Round:          proto.Int32(int32(fr.RoundNum)),
			Result:         proto.Int32(int32(result_b)),
			FightStateData: sc_fp_data_b,
			IsFight:        proto.Bool(true),
		}

		sc_msg_a := &protof.Message1{
			ScFightDataTunnelCapture: sc_msg_cap_a,
		}
		sc_msg_b := &protof.Message1{
			ScFightDataTunnelCapture: sc_msg_cap_b,
		}
		fr.who = false
		fr.FighterA.SendFightMsg(sc_msg_a)
		if fr.FighterB.Playerdata.IsRobot() {
			gbframe.Logger_Info.Println("2222222222222222222222222222222222")
			robot_csMsg, ok := capture_robot.RobotFight(sc_msg_b)
			gbframe.Logger_Info.Println("33333333333333333333333333333333333333")
			fr.SaveFightRecordData()
			if ok {
				gbframe.Logger_Info.Println("fight with the robot is over!player:", player)
				// return robot_csMsg, ok
				if fr.Result > 0 {
					fr.CaulFightEnd()
				}
			}
			return robot_csMsg, ok
		} else {
			fr.FighterB.SendFightMsg(sc_msg_b)
		}

	} else {
		result_a, result_b := fr.GetResult(2)
		sc_fp_data_a := ServerToProtoPoints(fr.FighterA, fr.FighterB.Name)
		sc_fp_data_b := ServerToProtoPoints(fr.FighterB, fr.FighterA.Name)
		sc_msg_cap_a := &protof.Message1_SC_FightData_Tunnel_Capture{
			Code:           proto.Int32(int32(code)),
			Round:          proto.Int32(int32(fr.RoundNum)),
			Result:         proto.Int32(int32(result_a)),
			FightStateData: sc_fp_data_a,
			IsFight:        proto.Bool(true),
		}
		sc_msg_cap_b := &protof.Message1_SC_FightData_Tunnel_Capture{
			Code:           proto.Int32(int32(code)),
			Round:          proto.Int32(int32(fr.RoundNum)),
			Result:         proto.Int32(int32(result_b)),
			FightStateData: sc_fp_data_b,
			IsFight:        proto.Bool(false),
		}
		sc_msg_a := &protof.Message1{
			ScFightDataTunnelCapture: sc_msg_cap_a,
		}
		sc_msg_b := &protof.Message1{
			ScFightDataTunnelCapture: sc_msg_cap_b,
		}
		fr.who = true
		fr.FighterA.SendFightMsg(sc_msg_a)
		if !fr.FighterB.Playerdata.IsRobot() {
			fr.FighterB.SendFightMsg(sc_msg_b)
		}

	}
	fr.SaveFightRecordData()
	if fr.Result > 0 {
		fr.CaulFightEnd()
	}
	return nil, true
}

//判断输赢
func (fr *FightRoom) GetResult(side int) (int, int) {
	if fr.RoundNum <= MAXROUNDNUM {
		if side == 1 { //判断A是否赢了
			birth_point := fr.FighterB.Fight_Data.BirthPoint
			num := len(fr.FighterA.Fight_Data.MovePoints)
			if num <= 0 {
				return 0, 0
			}
			move_end_point := fr.FighterA.Fight_Data.MovePoints[num-1]
			if birth_point.x == move_end_point.x && birth_point.y == move_end_point.y {
				gbframe.Logger_Info.Println("GetResult A win!")
				fr.Result = 1
				return 1, 2
			}
		} else if side == 2 { //判断B是否赢了
			birth_point := fr.FighterA.Fight_Data.BirthPoint
			num := len(fr.FighterB.Fight_Data.MovePoints)
			if num <= 0 {
				return 0, 0
			}
			move_end_point := fr.FighterB.Fight_Data.MovePoints[num-1]
			if birth_point.x == move_end_point.x && birth_point.y == move_end_point.y {
				gbframe.Logger_Info.Println("GetResult B win!")
				fr.Result = 2
				return 2, 1
			}
		}
		return 0, 0 //还没有结束
	} else {
		gbframe.Logger_Info.Println("GetResult is draw!")
		fr.Result = 3
		return 3, 3 //平了
	}

}

func (fr *FightRoom) showRecords() {
	for _, data := range fr.FightRecords {
		gbframe.Logger_Info.Println("Record:", data)
	}

}

func (fr *FightRoom) FightDataSave() {
	//保存到redis
	db.Gameredis.PlayerDataSave(fr.FighterA.Playerdata)
	//保存排行榜
	db.Gameredis.RankDataSave(fr.FighterA.Playerdata.Name, fr.FighterA.Playerdata.Score)
	db.Gameredis.PlayerDataSave(fr.FighterB.Playerdata)
	//保存排行榜
	db.Gameredis.RankDataSave(fr.FighterB.Playerdata.Name, fr.FighterB.Playerdata.Score)
}

//战斗结束，删除房间
func (fr *FightRoom) CaulFightEnd() {
	gbframe.Logger_Info.Println("Fight is over!Close Room!roomid:", fr.Room)
	if fr.Result == 1 {
		fr.FighterA.Playerdata.Score += 3
	} else if fr.Result == 2 {
		fr.FighterB.Playerdata.Score += 3
	} else if fr.Result == 3 {
		fr.FighterB.Playerdata.Score += 1
		fr.FighterA.Playerdata.Score += 1
	}
	fr.FightDataSave()
	fr.showRecords()
	delete(FightRooms, fr.Room)
}

func (fp *FightPlayer) SendFightMsg(sc_msg *protof.Message1) {
	server := gamenet.NetServers[fp.ServerSess]
	gbframe.Logger_Info.Println("server sess :", fp.ServerSess)
	mid := int(protof.Message1_SC_FIGHTDATA_TUNNEL_CAPTURE)
	server.WriteMessage(mid, sc_msg)
}

//获取对方的出生点和所在点
func (fp *FightPlayer) GetOtherBirthPointOrAtPoint(birth_point Map, at_point Map) {
	for _, point := range fp.Fight_Data.ExcavatePoits {
		if birth_point.x == point.x && birth_point.y == point.y {
			fp.Fight_Data.OtherBirthPoint.x = point.x
			fp.Fight_Data.OtherBirthPoint.y = point.y
		}
		if at_point.x == point.x && at_point.y == point.y {
			fp.Fight_Data.OtherUserAtPoint.x = point.x
			fp.Fight_Data.OtherUserAtPoint.y = point.y
		}
	}
}

func (fp *FightPlayer) UpdateNotFightPlayerInfo(fight_player *FightPlayer) {
	fp.Fight_Data.ExcavatePoits = fight_player.Fight_Data.OtherExcavatePoints
	fp.Fight_Data.OtherExcavatePoints = fight_player.Fight_Data.ExcavatePoits
}

func (fp *FightPlayer) UpdateIsFightPlayerInfo(cs_msg *protof.Message1) int {
	code := 0
	fp.Fight_Data.Power += ADDPOWER

	new_excavate_points := cs_msg.CsFightDataTunnelCapture.GetCapturePoints()
	new_excavate_points_ := ProtoToServerPoints(new_excavate_points) //新挖掘的节点
	new_excavate_num := len(new_excavate_points_)
	if !IsInMapLimit(new_excavate_points_) {
		gbframe.Logger_Error.Println("UpdateFightPlayerInfo excavate points out of map range!name:", fp.Name)
		return 2
	}
	fp.MakeMapPointsInFP(new_excavate_points_, 1)
	fp.MakeMapPointsInFP(new_excavate_points_, 3)
	move_points := cs_msg.CsFightDataTunnelCapture.GetMovePoints()
	move_points_ := ProtoToServerPoints(move_points) //移动路径
	new_move_num := len(move_points_)
	if !IsInMapLimit(move_points_) {
		gbframe.Logger_Error.Println("UpdateFightPlayerInfo move_points out of map range!name:", fp.Name)
		return 2
	}
	code = fp.MakeMapPointsInFP(move_points_, 2)
	use_power := new_excavate_num + new_move_num
	if use_power <= fp.Fight_Data.Power {
		fp.Fight_Data.Power -= use_power
	} else {
		gbframe.Logger_Error.Println("power is not enough!player name:", fp.Name,
			" player power is ", fp.Fight_Data.Power, " use power is ", use_power)
		code = 4
	}
	return code
}

//处理节点
func (fp *FightPlayer) MakeMapPointsInFP(points []*Map, point_type int) int {
	switch point_type {
	case 1: //添加挖掘点
		fp.Fight_Data.ExcavatePoits = append(fp.Fight_Data.ExcavatePoits, points...)
	case 2: //转换移动点和当前所在点
		if len(points) <= 0 {
			gbframe.Logger_Warning.Println("move points is None!player name:", fp.Name)
			return 0
		}
		if ok := fp.CheckMovePointsAndGetAtPoint(points); !ok {
			return 1
		}
		fp.Fight_Data.MovePoints = points
	case 3: //用自己的挖掘点将对方挖掘点覆盖了
		fp.DelMapPoints(points)
	}
	return 0

}

//检查并确定最终所在点
func (fp *FightPlayer) CheckMovePointsAndGetAtPoint(points []*Map) bool {
	//移动点的合法性，两个节点必须相邻
	tmp_point := fp.Fight_Data.UserAtPoint

	//确定移动点都在挖掘点内
	for _, point := range points {
		if ok := IsInMapInfoList(fp.Fight_Data.ExcavatePoits, *point, *fp.Fight_Data.BirthPoint); !ok {
			gbframe.Logger_Error.Println("CheckMovePointsRight is error:this point ", *point, " is not in ExcavatePoits!Name:", fp.Name)
			return false
		}
	}
	if !IsMovePointsEachForNear(*tmp_point, points) {
		return false
	}
	fp.Fight_Data.UserAtPoint.x = points[len(points)-1].x
	fp.Fight_Data.UserAtPoint.y = points[len(points)-1].y
	return true
}

//将对方的已挖掘节点与自己的挖掘点匹配
func (fp *FightPlayer) DelMapPoints(points []*Map) {
	for _, point := range points {
		for i, other_point := range fp.Fight_Data.OtherExcavatePoints {
			if point.x == other_point.x && point.y == other_point.y {
				fp.Fight_Data.OtherExcavatePoints =
					append(fp.Fight_Data.OtherExcavatePoints[:i], fp.Fight_Data.OtherExcavatePoints[i+1:]...)
			}
		}
	}
}

//转换proto协议转换成服务器的格式
func ProtoToServerPoints(msg_points []*protof.Message1_Map_Info) []*Map {
	var points = []*Map{}
	for _, msg_point := range msg_points {
		x := int(msg_point.GetX())
		y := int(msg_point.GetY())
		point := &Map{
			x: x,
			y: y,
		}
		points = append(points, point)
	}

	return points
}

func ServerToProtoPoints(fp *FightPlayer, other_name string) *protof.Message1_FightStateData_Capture {
	//己方出生点
	point_birth := &protof.Message1_Map_Info{
		X: proto.Int32(int32(fp.Fight_Data.BirthPoint.x)),
		Y: proto.Int32(int32(fp.Fight_Data.BirthPoint.y)),
	}
	//对方出生点
	other_point_birth := &protof.Message1_Map_Info{
		X: proto.Int32(int32(fp.Fight_Data.OtherBirthPoint.x)),
		Y: proto.Int32(int32(fp.Fight_Data.OtherBirthPoint.y)),
	}
	//转换maplist为proto协议的maplist
	//已挖掘的节点
	excavate_points := ServerMapsToProtoMaps(fp.Fight_Data.ExcavatePoits)
	other_excavate_points := ServerMapsToProtoMaps(fp.Fight_Data.OtherExcavatePoints)
	// move_points := ServerMapsToProtoMaps(fp.Fight_Data.MovePoints)

	at_point := &protof.Message1_Map_Info{
		X: proto.Int32(int32(fp.Fight_Data.UserAtPoint.x)),
		Y: proto.Int32(int32(fp.Fight_Data.UserAtPoint.y)),
	}

	other_at_point := &protof.Message1_Map_Info{
		X: proto.Int32(int32(fp.Fight_Data.OtherUserAtPoint.x)),
		Y: proto.Int32(int32(fp.Fight_Data.OtherUserAtPoint.y)),
	}

	last_power := proto.Int32(int32(fp.Fight_Data.Power))
	other_last_power := proto.Int32(int32(fp.Fight_Data.OtherPower))

	fdtc_msg := &protof.Message1_FightStateData_Capture{
		Name:       proto.String(fp.Name),
		OtherName:  proto.String(other_name),
		BirthPoint: point_birth,
		// ExcavatePoits:       excavate_points,
		// OtherExcavatePoints: other_excavate_points,
		// MovePoints:          move_points,
		OtherBirthPoint:  other_point_birth,
		UserAtPoint:      at_point,
		OtherUserAtPoint: other_at_point,
		LastPower:        last_power,
		OtherPower:       other_last_power,
	}
	if len(excavate_points) > 0 {
		fdtc_msg.ExcavatePoits = excavate_points
	}
	if len(other_excavate_points) > 0 {
		fdtc_msg.OtherExcavatePoints = other_excavate_points
	}
	// if len(move_points) > 0 {
	// 	fdtc_msg.MovePoints = move_points
	// }
	return fdtc_msg
}

func ServerMapsToProtoMaps(maplist []*Map) []*protof.Message1_Map_Info {
	var map_infos = []*protof.Message1_Map_Info{}

	for _, data := range maplist {
		map_info := &protof.Message1_Map_Info{
			X: proto.Int32(int32(data.x)),
			Y: proto.Int32(int32(data.y)),
		}
		map_infos = append(map_infos, map_info)
	}
	return map_infos
}

//是否超出地图
func IsInMapLimit(points []*Map) bool {
	for _, point := range points {
		if point.x >= MAPMAX_X || point.y >= MAPMAX_Y || point.x < 0 || point.y < 0 {
			gbframe.Logger_Error.Println("point:", point, " is not in Maplimit,limit_x:", MAPMAX_X, " limit_y:", MAPMAX_Y)
			return false
		}
	}
	return true
}

//是否在挖掘点以外
func IsInMapInfoList(points []*Map, point, birth Map) bool {
	if point.x == birth.x && point.y == birth.y {
		return true
	}
	for _, p := range points {
		if p.x == point.x && p.y == point.y {
			return true
		}
	}
	gbframe.Logger_Error.Println("this point:", point, " is not in Maplist:", points)
	return false
}

func IsMovePointsEachForNear(tmp_point Map, points []*Map) bool {
	//确定移动点是两两相邻的持续移动点
	for _, point := range points {
		tmp_x := int(math.Abs(float64(tmp_point.x - point.x)))
		tmp_y := int(math.Abs(float64(tmp_point.y - point.y)))
		if tmp_x == 0 {
			if tmp_y != 1 {
				gbframe.Logger_Error.Println("tmp_point(", tmp_point.x, tmp_point.y, ") and point(", point.x, point.y, ") is not near!")
				return false
			}
		} else if tmp_y == 0 {
			if tmp_x != 1 {
				gbframe.Logger_Error.Println("tmp_point(", tmp_point.x, tmp_point.y, ") and point(", point.x, point.y, ") is not near!")
				return false
			}
		} else {
			gbframe.Logger_Error.Println("tmp_point(", tmp_point.x, tmp_point.y, ") and point(", point.x, point.y, ") is not near!")
			return false
		}
		tmp_point.x = point.x
		tmp_point.y = point.y
	}
	return true

}
