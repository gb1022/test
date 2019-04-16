package capture_robot

import (
	// "data"
	"encoding/binary"

	// "gamenet"
	"gbframe"
	"math/rand"
	"protof"
	"time"

	"github.com/golang/protobuf/proto"
)

const (
	MAPMAX_X = 6 //x轴大小
	MAPMAX_Y = 6 //y轴大小
	// MAXROUNDNUM   = 20 //总回合数
	// MAXPOWERLIMIT = 10 //能量总数
	ADDPOWER = 2 //每回合增加的能量点
)
const (
	COSTPOWER = 5 //每回合动用的能量点
)

type Node struct {
	value *MapInfo
	left  *Node
	right *Node
}

type MapInfo struct {
	x int
	y int
}

var (
	b_x int
	b_y int
	// o_x int
	// o_y int
	// Fight_data_capture *Capture_Fight_data_capture
	paths [][]*MapInfo
)

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

var Fight_data_capture = Capture_fight_data{}

var e_points []*MapInfo

// var m_points []*MapInfo
var best_path []*MapInfo

func RandPoint(m_x, m_y int) (int, int) {
	rand.Seed(time.Now().UnixNano())
	bx := rand.Intn(m_x)
	by := rand.Intn(m_y)
	return bx, by
}

func GameRobotProsses(isFight bool, my_name, other_name string) []*MapInfo {
	e_points = []*MapInfo{}
	Fight_data_capture.MovePoints = []*MapInfo{}

	if Fight_data_capture.OtherBirthPoint.x == -1 && Fight_data_capture.OtherBirthPoint.y == -1 {
		KeepExcavatePointsForNum(COSTPOWER)
		// showMapCaptureGame(isFight, my_name, other_name)
		return e_points
	}

	// if len(best_path) == 0 {
	paths = [][]*MapInfo{}
	GetBestPath()
	// }
	// gbframe.Logger_Info.Print("before move :")
	// showPath(best_path)
	// var ss string
	// gbframe.Logger_Info.Scanln(&ss)
	// gbframe.Logger_Info.Printf(ss)

	AutoMove(isFight, my_name, other_name)
	// gbframe.Logger_Info.Print("after move :")
	// showPath(best_path)

	// gbframe.Logger_Info.Scanln(&ss)
	// gbframe.Logger_Info.Printf(ss)
	// showMapCaptureGame(isFight, my_name, other_name)
	time.Sleep(3 * time.Second)
	return e_points
}

func GetBestPath() {
	x_e := Fight_data_capture.UserAtPoint.x
	y_e := Fight_data_capture.UserAtPoint.y
	x_s := Fight_data_capture.OtherBirthPoint.x
	y_s := Fight_data_capture.OtherBirthPoint.y
	gbframe.Logger_Info.Println("======================x_e:", x_e, "y_e:", y_e, "x_s:", x_s, "y_s", y_s)
	xf := x_e - x_s
	yf := y_e - y_s
	node := CreatBitTree(x_e, y_e, xf, yf, x_s, y_s)
	gbframe.Logger_Info.Println("-=-=-=-=-=-=-=-- root.map:", *node.value)
	GetAllPath(node)
	best_path = FindTheBestPath()
	gbframe.Logger_Info.Println("best path first:", best_path[0], "second:", best_path[1])

}

func CreatBitTree(x, y, x_f, y_f, x_s, y_s int) *Node {
	n := &Node{
		value: &MapInfo{
			x: x,
			y: y,
		},
	}
	// gbframe.Logger_Info.Print(n.value)
	var x_a, y_a = x, y
	if x_f > 0 && x > x_s {
		x_a = x - 1
		n.left = CreatBitTree(x_a, y, x_f, y_f, x_s, y_s)
	} else if x_f < 0 && x < x_s {
		x_a = x + 1
		n.left = CreatBitTree(x_a, y, x_f, y_f, x_s, y_s)
	}
	if y_f > 0 && y > y_s {
		y_a = y - 1
		n.right = CreatBitTree(x, y_a, x_f, y_f, x_s, y_s)
	} else if y_f < 0 && y < y_s {
		y_a = y + 1
		n.right = CreatBitTree(x, y_a, x_f, y_f, x_s, y_s)
	}
	return n
}

func GetAllPath(root *Node) {
	var p = []*MapInfo{}
	// var n *Node
	xianxubianli(root, p)
	// 统计路径(x1, y1)
	// gbframe.Logger_Info.Println("\n paths:", paths)
}

func xianxubianli(node *Node, new_path []*MapInfo) {
	if node == nil {
		return
	}
	if node.left == nil && node.right == nil {
		new_path = append(new_path, node.value)
		paths = append(paths, new_path)
		// gbframe.Logger_Info.Println("new path:", new_path)
		return
	}
	new_path = append(new_path, node.value)
	xianxubianli(node.left, new_path)
	xianxubianli(node.right, new_path)
	return
}

func FindTheBestPath() []*MapInfo {
	m := 0
	// var n_path []MapInfo
	j := 0

	for i, path := range paths {
		n := 0
		for _, p := range path {
			if IsExistList(Fight_data_capture.ExcavatePoits, p.x, p.y) {
				n++
			}
		}
		if m < n {
			m = n
			j = i
		}
	}
	if len(paths) > 0 {
		return paths[j]
	}
	return []*MapInfo{}

}

func AutoMove(isFight bool, my_name, other_name string) {
	i := 0
	for i < len(best_path) {
		p := best_path[i]
		// gbframe.Logger_Info.Println("best_path point :", *p)
		// time.Sleep(3 * time.Second)
		if p.x == Fight_data_capture.UserAtPoint.x && p.y == Fight_data_capture.UserAtPoint.y {
			i++
			continue
		}
		// if Fight_data_capture.Power <= 0 {
		// 	break
		// }

		if !IsExistList(Fight_data_capture.ExcavatePoits, p.x, p.y) {
			point := &MapInfo{
				x: p.x,
				y: p.y,
			}
			Fight_data_capture.ExcavatePoits = append(Fight_data_capture.ExcavatePoits, point)
			e_points = append(e_points, point)
			Fight_data_capture.Power--

			// showMapCaptureGame(isFight, my_name, other_name)
			if Fight_data_capture.Power <= 0 {
				i--
				break
			}
		}

		move := &MapInfo{
			x: p.x,
			y: p.y,
		}
		Fight_data_capture.MovePoints = append(Fight_data_capture.MovePoints, move)
		Fight_data_capture.UserAtPoint.x = p.x
		Fight_data_capture.UserAtPoint.y = p.y
		Fight_data_capture.Power--
		// showMapCaptureGame(isFight, my_name, other_name)
		if Fight_data_capture.Power <= 0 {
			break
		}
		if p.x == Fight_data_capture.OtherBirthPoint.x && p.y == Fight_data_capture.OtherBirthPoint.y {
			// gbframe.Logger_Info.Println("YOU WIN!!!!!!!!!!!!!!!!")
			break
		}
		// time.Sleep(1 * time.Second)
		// var ss string
		// gbframe.Logger_Info.Scanln(&ss)
		// gbframe.Logger_Info.Printf(ss)
		// best_path = best_path[i+1:]
		i++

	}
	best_path = best_path[i+1:]
}

func IsExist(x, y int) bool {
	if x == Fight_data_capture.BirthPoint.x && y == Fight_data_capture.BirthPoint.y {
		gbframe.Logger_Info.Printf("this point[%d,%d] is exist by user at point[%d,%d]!\n", x, y, Fight_data_capture.BirthPoint.x, Fight_data_capture.BirthPoint.y)
		return true
	}
	for _, p := range Fight_data_capture.ExcavatePoits {
		if p.x == x && p.y == y {
			gbframe.Logger_Info.Printf("this point[%d,%d] is exist in excavate points!\n", x, y)
			return true
		}
	}
	return false

}

func KeepExcavatePointsForNum(num int) int {
	j := 0
	for j < num {
		x, y := RandPoint(MAPMAX_X, MAPMAX_Y)
		// time.Sleep(1 * time.Millisecond)
		if IsExist(x, y) {
			continue
		}
		point := &MapInfo{
			x: x,
			y: y,
		}
		Fight_data_capture.ExcavatePoits = append(Fight_data_capture.ExcavatePoits, point)
		e_points = append(e_points, point)
		Fight_data_capture.Power--
		if Fight_data_capture.Power <= 0 {
			gbframe.Logger_Info.Println("Fight_data_capture Power is None!")
			return 2
		}
		// if IsExcavateBirthPoint(point) {
		// 	return 0
		// }
		j++
	}
	return 1
}

//判断是否挖到了对方出生点，这个 在服务器判断 这里用不上
// func IsExcavateBirthPoint(point *MapInfo) bool {
// 	if point.x == o_x && point.y == o_y {
// 		gbframe.Logger_Info.Println("this excavate point is other birth point!\n")
// 		Fight_data_capture.OtherBirthPoint.x = point.x
// 		Fight_data_capture.OtherBirthPoint.y = point.y
// 		return true
// 	}
// 	return false
// }

func showPath(path []*MapInfo) {
	gbframe.Logger_Info.Print("the path is :")
	for _, p := range path {
		gbframe.Logger_Info.Print(*p)
	}
	gbframe.Logger_Info.Println()
}

func RobotFight(sc_msg *protof.Message1) (*protof.Message1, bool) {
	robot_csMsg := &protof.Message1{}
	ScMsgToCaptureData(sc_msg)
	isFight := sc_msg.ScFightDataTunnelCapture.GetIsFight()
	my_name := sc_msg.ScFightDataTunnelCapture.FightStateData.GetName()
	other_name := sc_msg.ScFightDataTunnelCapture.FightStateData.GetOtherName()
	// showMapCaptureGame(isFight, my_name, other_name)
	if sc_msg.ScFightDataTunnelCapture.GetResult() == 1 {
		gbframe.Logger_Info.Println("Tunnel Capture Fight Over!\n  Congratulation！ You Win!!!!!!!!")
		return robot_csMsg, true
	} else if sc_msg.ScFightDataTunnelCapture.GetResult() == 2 {
		gbframe.Logger_Info.Println("Tunnel Capture Fight Over!\nOh I am sorry! You Lose!!!!!!!!!")
		return robot_csMsg, true
	} else if sc_msg.ScFightDataTunnelCapture.GetResult() == 3 {
		gbframe.Logger_Info.Println("Tunnel Capture Fight Over!\n   It is Draw!!!!!!!!!!")
		return robot_csMsg, true
	}

	if isFight := sc_msg.ScFightDataTunnelCapture.GetIsFight(); !isFight {
		return robot_csMsg, false
	}
	// e_points := Fight_data_capture.WriteExcavate()
	// Fight_data_capture.ExcavatePoits = append(Fight_data_capture.ExcavatePoits, e_points...)
	// m_points := Fight_data_capture.WriteMovePoints(len(e_points))
	e_points := GameRobotProsses(isFight, my_name, other_name)
	m_points := Fight_data_capture.MovePoints
	gbframe.Logger_Info.Println("e_points:", e_points)
	gbframe.Logger_Info.Println("m_points:", m_points)
	m_maps_proto := ClientMapsToProtoMaps(m_points)
	e_maps_proto := ClientMapsToProtoMaps(e_points)
	fight_capture_msg := &protof.Message1_CS_FightData_Tunnel_Capture{
		MovePoints:    m_maps_proto,
		CapturePoints: e_maps_proto,
	}
	robot_csMsg.CsFightDataTunnelCapture = fight_capture_msg

	// room_sess = gbframe.MakeSession("0", player.Name)

	// fight_room, ok := fight_tunnel_capture.FightRooms[room_sess]
	// if !ok {
	// 	gbframe.Logger_Error.Println("This fight_room is not exist!room_sess:", room_sess)
	// 	return
	// }
	// w_msg := writeClientMessge(cs_msg, int(protof.Message1_CS_FIGHTDATA_TUNNEL_CAPTURE))
	// var ss string
	// gbframe.Logger_Info.Scanln(&ss)
	// gbframe.Logger_Info.Printf(ss)
	// conn.Write(w_msg)

	return robot_csMsg, false
}

func ScMsgToCaptureData(sc_msg *protof.Message1) {
	gbframe.Logger_Info.Println("fightInput sc_msg:", sc_msg)
	// var movelist = []*MapInfo{}
	// var excavatelist = []*MapInfo{}

	birth_point := &MapInfo{
		x: int(sc_msg.GetScFightDataTunnelCapture().GetFightStateData().GetBirthPoint().GetX()),
		y: int(sc_msg.ScFightDataTunnelCapture.FightStateData.GetBirthPoint().GetY()),
	}

	// gbframe.Logger_Info.Println("birth_point:", birth_point)
	// gbframe.Logger_Info.Printf("%p\n", &birth_point)
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
	gbframe.Logger_Info.Println("------------------- ScMsgToCaptureData Fight_data_capture.userAtpoint :", *Fight_data_capture.UserAtPoint)
}

func IsExistList(points []*MapInfo, x int, y int) bool {
	for _, point := range points {
		p_x := point.x
		p_y := point.y
		if p_x == x && p_y == y {
			return true
		}
	}
	return false
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

func writeClientMessge(msg *protof.Message1, mid int) []byte {
	data, _ := proto.Marshal(msg)
	s_data := MarshalSendMsg(data, mid)
	return s_data
}

func MarshalSendMsg(msg []byte, msgId int) []byte {
	cache := make([]byte, 4)
	var buff []byte
	pctLen := uint32(len(msg) + 2)
	binary.BigEndian.PutUint32(cache, pctLen)
	buff = append(buff, cache...)
	idb := cache[:2]
	binary.BigEndian.PutUint16(cache, uint16(msgId))
	buff = append(buff, idb...)
	buff = append(buff, msg...)
	return buff
}
