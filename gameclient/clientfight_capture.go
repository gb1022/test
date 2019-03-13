package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"protof"

	"github.com/golang/protobuf/proto"
)

const (
	MAPMAX_X      = 20 //x轴大小
	MAPMAX_Y      = 20 //y轴大小
	MAXROUNDNUM   = 20 //总回合数
	MAXPOWERLIMIT = 20 //能量总数
	ADDPOWER      = 2  //每回合增加的能量点
)

type MapInfo struct {
	x int
	y int
}

type Capture_fight_data struct {
	round               int
	BirthPoint          *MapInfo   //出生节点
	ExcavatePoits       []*MapInfo //自己已挖掘的节点
	OtherExcavatePoints []*MapInfo //对方已挖掘的节点
	MovePoints          []*MapInfo //移动
	OtherBirthPoint     *MapInfo   //对方的出生节点
	UserAtPoint         *MapInfo   //当前所在的节点
	OtherUserAtPoint    *MapInfo   //对方所在的节点 ps:（被发现才会有值，未被发现则节点为-1,-1）
	Power               int        //剩余能量点
	OtherPower          int        //对方剩余能量点
}

func ScMsgToCaptureData(sc_msg *protof.Message1) {
	fmt.Println("fightInput sc_msg:", sc_msg)
	// var movelist = []*MapInfo{}
	// var excavatelist = []*MapInfo{}

	birth_point := &MapInfo{
		x: int(sc_msg.GetScFightDataTunnelCapture().GetFightStateData().GetBirthPoint().GetX()),
		y: int(sc_msg.ScFightDataTunnelCapture.FightStateData.GetBirthPoint().GetY()),
	}

	// fmt.Println("birth_point:", birth_point)
	// fmt.Printf("%p\n", &birth_point)
	// Fight_data_capture.BirthPoint.x = int(sc_msg.GetScFightDataTunnelCapture().GetFightStateData().GetBirthPoint().GetX())
	// Fight_data_capture.BirthPoint.y = int(sc_msg.ScFightDataTunnelCapture.FightStateData.GetBirthPoint().GetY())
	Fight_data_capture.round = int(sc_msg.GetScFightDataTunnelCapture().GetRound())
	Fight_data_capture.BirthPoint = birth_point
	other_birth_point := &MapInfo{
		x: int(*sc_msg.ScFightDataTunnelCapture.FightStateData.GetOtherBirthPoint().X),
		y: int(*sc_msg.ScFightDataTunnelCapture.FightStateData.GetOtherBirthPoint().Y),
	}
	Fight_data_capture.OtherBirthPoint = other_birth_point
	user_at_point := &MapInfo{
		x: int(*sc_msg.ScFightDataTunnelCapture.FightStateData.GetUserAtPoint().X),
		y: int(*sc_msg.ScFightDataTunnelCapture.FightStateData.GetUserAtPoint().Y),
	}
	Fight_data_capture.UserAtPoint = user_at_point
	other_at_point := &MapInfo{
		x: int(*sc_msg.ScFightDataTunnelCapture.FightStateData.GetOtherUserAtPoint().X),
		y: int(*sc_msg.ScFightDataTunnelCapture.FightStateData.GetOtherUserAtPoint().Y),
	}
	Fight_data_capture.OtherUserAtPoint = other_at_point
	Fight_data_capture.Power = int(sc_msg.ScFightDataTunnelCapture.FightStateData.GetLastPower())
	Fight_data_capture.OtherPower = int(sc_msg.ScFightDataTunnelCapture.FightStateData.GetOtherPower())

	Fight_data_capture.ExcavatePoits = ProtoToClientData(sc_msg.ScFightDataTunnelCapture.FightStateData.GetExcavatePoits())
	Fight_data_capture.OtherExcavatePoints = ProtoToClientData(sc_msg.ScFightDataTunnelCapture.FightStateData.GetOtherExcavatePoints())
	// Fight_data_capture.MovePoints = ProtoToClientData(sc_msg.ScFightDataTunnelCapture.FightStateData.GetMovePoints())
	// Fight_data_capture.ExcavatePoits = append(Fight_data_capture.ExcavatePoits, excavatePoits...)
	// Fight_data_capture.OtherExcavatePoints = append(Fight_data_capture.OtherExcavatePoints, otherExcavatePoints...)
	// return movelist, excavatelist
}

func (fdc *Capture_fight_data) GetExcavatePoints() []*MapInfo {

	var excavate_maplist = []*MapInfo{}
	if fdc.Power > 3 {
		return excavate_maplist
	}
	return excavate_maplist

}

func (fdc *Capture_fight_data) WriteExcavate() []*MapInfo {
	fmt.Println("please Input your excavate points:")
	var excavate_maplist = []*MapInfo{}
	for i := 0; i < fdc.Power; i++ {
		var x int
		var y int
		fmt.Print("x:")
		fmt.Scanln(&x)
		fmt.Print("y:")
		fmt.Scanln(&y)
		if x == -1 && y == -1 {
			return excavate_maplist
		} else if InTheseMaps(x, y, fdc.ExcavatePoits) {
			fmt.Println("continue,because this point is exist!point:", x, y)
			continue
		}
		e_mapinfo := &MapInfo{
			x: x,
			y: y,
		}
		excavate_maplist = append(excavate_maplist, e_mapinfo)
	}

	return excavate_maplist

}

func InTheseMaps(x, y int, map_list []*MapInfo) bool {
	fmt.Println("map_list:", map_list)
	for _, point := range map_list {
		if point.x == x && point.y == y {
			return true
		}
	}
	return false
}

func (fdc *Capture_fight_data) WriteMovePoints(num int) []*MapInfo {
	fmt.Println("please Input your move points:")
	var move_maplist = []*MapInfo{}
	last_power := fdc.Power - num
	for i := 0; i < last_power; i++ {
		var x int
		var y int
		fmt.Print("x:")
		fmt.Scanln(&x)
		fmt.Print("y:")
		fmt.Scanln(&y)
		if x == -1 && y == -1 {
			return move_maplist
		} else if !InTheseMaps(x, y, fdc.ExcavatePoits) {
			fmt.Println("continue,because this point not in ExcavatePoints!point:", x, y)
			continue
		}
		m_mapinfo := &MapInfo{
			x: x,
			y: y,
		}
		move_maplist = append(move_maplist, m_mapinfo)
	}
	return move_maplist
}

func (fdc *Capture_fight_data) AutoWriteGamePoint() ([]*MapInfo, []*MapInfo) {
	var movelist = []*MapInfo{}
	var excavatelist = []*MapInfo{}
	// e_point := &MapInfo{
	// 	x: fdc.ExcavatePoits()
	// }
	return excavatelist, movelist

}

func ProtoToClientData(points []*protof.Message1_Map_Info) []*MapInfo {
	var maplist = []*MapInfo{}
	for _, point := range points {
		x := int(point.GetX())
		y := int(point.GetY())
		map_ := &MapInfo{
			x: x,
			y: y,
		}
		maplist = append(maplist, map_)
	}
	return maplist
}

func ClientMapsToProtoMaps(maplist []*MapInfo) []*protof.Message1_Map_Info {
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

func showMapCaptureGame(sc_msg *protof.Message1) {
	round := int(sc_msg.ScFightDataTunnelCapture.GetRound())
	fmt.Println("====================MAP================================")
	m_b_x := int(sc_msg.ScFightDataTunnelCapture.FightStateData.BirthPoint.GetX())
	m_b_y := int(sc_msg.ScFightDataTunnelCapture.FightStateData.BirthPoint.GetY())
	o_b_x := int(sc_msg.ScFightDataTunnelCapture.FightStateData.OtherBirthPoint.GetX())
	o_b_y := int(sc_msg.ScFightDataTunnelCapture.FightStateData.OtherBirthPoint.GetY())
	my_cap_points := sc_msg.ScFightDataTunnelCapture.FightStateData.GetExcavatePoits()
	other_cap_points := sc_msg.ScFightDataTunnelCapture.FightStateData.GetOtherExcavatePoints()
	m_at_x := int(sc_msg.ScFightDataTunnelCapture.FightStateData.UserAtPoint.GetX())
	m_at_y := int(sc_msg.ScFightDataTunnelCapture.FightStateData.UserAtPoint.GetY())
	o_at_x := int(sc_msg.ScFightDataTunnelCapture.FightStateData.OtherUserAtPoint.GetX())
	o_at_y := int(sc_msg.ScFightDataTunnelCapture.FightStateData.OtherUserAtPoint.GetY())

	fmt.Println("Fight Round:", round, "my birth:[", m_b_x, ",", m_b_y, "] other birth:[", o_b_x, ",", o_b_y, "]")
	fmt.Printf("my at:[%d,%d],other at:[%d,%d]\n", m_at_x, m_at_y, o_at_x, o_at_y)
	fmt.Println("*:自己的位置，+:对方的位置，%:自己的出生点，#:对方的出生点，@:自己挖的，&:对方挖的")
	isFight := sc_msg.ScFightDataTunnelCapture.GetIsFight()

	for i := 0; i < (MAP_Y + 1); i++ {
		if i > 0 {
			fmt.Printf("%d", i-1)
		}
		for j := 0; j < MAP_X; j++ {
			if i == 0 {
				if j == 0 {
					fmt.Printf("   %d", j)
				} else if j > 0 && j < 10 {
					fmt.Printf("  %d", j)
				} else {
					fmt.Printf(" %d", j)
				}

			} else {
				if i >= 11 && j == 0 {
					if j == (m_at_x) && i == (m_at_y+1) {
						fmt.Print(" *")
					} else if i == o_at_y+1 && j == o_at_x {
						fmt.Print(" +")
					} else if i == m_b_y+1 && j == m_b_x {
						fmt.Print(" %")
					} else if j == (o_b_x) && i == (o_b_y+1) {
						fmt.Print(" #")
					} else if IsExistList(my_cap_points, j, i-1) {
						fmt.Print(" @")
					} else if IsExistList(other_cap_points, j, i-1) {
						fmt.Print(" &")
					} else {
						fmt.Print(" .")
					}
				} else {
					if j == (m_at_x) && i == (m_at_y+1) {
						fmt.Print("  *")
					} else if i == o_at_y+1 && j == o_at_x {
						fmt.Print("  +")
					} else if i == m_b_y+1 && j == m_b_x {
						fmt.Print("  %")
					} else if j == (o_b_x) && i == (o_b_y+1) {
						fmt.Print("  #")
					} else if IsExistList(my_cap_points, j, i-1) {
						fmt.Print("  @")
					} else if IsExistList(other_cap_points, j, i-1) {
						fmt.Print("  &")
					} else {
						fmt.Print("  .")
					}
				}

			}

		}
		fmt.Println("")
	}
	fmt.Println("====================================================")
	fmt.Println("Are you fighting! :", isFight)
	power := int(sc_msg.ScFightDataTunnelCapture.FightStateData.GetLastPower())
	other_power := int(sc_msg.ScFightDataTunnelCapture.FightStateData.GetOtherPower())
	my_name := sc_msg.ScFightDataTunnelCapture.FightStateData.GetName()
	other_name := sc_msg.ScFightDataTunnelCapture.FightStateData.GetOtherName()
	fmt.Println("My Name:", my_name, "Power:", power, "| Other Name:", other_name, "Other Power:", other_power)
	fmt.Println("====================Fight Show End================================")

}

func fight_tunnel_capture(sc_msg *protof.Message1, conn net.Conn) bool {
	showMapCaptureGame(sc_msg)
	if sc_msg.ScFightDataTunnelCapture.GetResult() == 1 {
		fmt.Println("Tunnel Capture Fight Over!\n  Congratulation！ You Win!!!!!!!!")
		return true
	} else if sc_msg.ScFightDataTunnelCapture.GetResult() == 2 {
		fmt.Println("Tunnel Capture Fight Over!\nOh I am sorry! You Lose!!!!!!!!!")
		return true
	} else if sc_msg.ScFightDataTunnelCapture.GetResult() == 3 {
		fmt.Println("Tunnel Capture Fight Over!\n   It is Draw!!!!!!!!!!")
		return true
	}
	ScMsgToCaptureData(sc_msg)
	if isFight := sc_msg.ScFightDataTunnelCapture.GetIsFight(); !isFight {
		return false
	}
	e_points := Fight_data_capture.WriteExcavate()
	Fight_data_capture.ExcavatePoits = append(Fight_data_capture.ExcavatePoits, e_points...)
	m_points := Fight_data_capture.WriteMovePoints(len(e_points))

	fmt.Println("m_points:", e_points)
	fmt.Println("e_points:", m_points)
	m_maps_proto := ClientMapsToProtoMaps(m_points)
	e_maps_proto := ClientMapsToProtoMaps(e_points)
	fight_capture_msg := &protof.Message1_CS_FightData_Tunnel_Capture{
		MovePoints:    m_maps_proto,
		CapturePoints: e_maps_proto,
	}
	cs_msg := &protof.Message1{
		CsFightDataTunnelCapture: fight_capture_msg,
	}
	w_msg := WriteMessge(cs_msg, int(protof.Message1_CS_FIGHTDATA_TUNNEL_CAPTURE))
	var ss string
	fmt.Scanln(&ss)
	fmt.Printf(ss)
	conn.Write(w_msg)
	return false
}

func IsExistList(points []*protof.Message1_Map_Info, x int, y int) bool {
	for _, point := range points {
		p_x := int(point.GetX())
		p_y := int(point.GetY())
		if p_x == x && p_y == y {
			return true
		}
	}
	return false
}

func Flush屏幕() {
	cmd := exec.Command("cmd", "/C", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
