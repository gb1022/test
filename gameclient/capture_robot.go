package main

import (
	"fmt"
	"math/rand"
	"time"
)

// const (
// 	MAPMAX_X      = 10  //x轴大小
// 	MAPMAX_Y      = 10  //y轴大小
// 	MAXROUNDNUM   = 300 //总回合数
// 	MAXPOWERLIMIT = 20  //能量总数
// 	ADDPOWER      = 5   //每回合增加的能量点
// )
const (
	COSTPOWER = 5 //每回合动用的能量点
)

type Node struct {
	value *MapInfo
	left  *Node
	right *Node
}

var (
	b_x int
	b_y int
	// o_x int
	// o_y int
	// Fight_data_capture *Capture_Fight_data_capture
	paths [][]*MapInfo
)

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
		showMapCaptureGame(isFight, my_name, other_name)
		return e_points
	}

	if len(best_path) == 0 {
		GetBestPath()
	}
	// fmt.Print("before move :")
	// showPath(best_path)
	// var ss string
	// fmt.Scanln(&ss)
	// fmt.Printf(ss)

	AutoMove(isFight, my_name, other_name)
	// fmt.Print("after move :")
	// showPath(best_path)

	// fmt.Scanln(&ss)
	// fmt.Printf(ss)
	showMapCaptureGame(isFight, my_name, other_name)
	time.Sleep(3 * time.Second)
	return e_points
}

func GetBestPath() {
	x_e := Fight_data_capture.UserAtPoint.x
	y_e := Fight_data_capture.UserAtPoint.y
	x_s := Fight_data_capture.OtherBirthPoint.x
	y_s := Fight_data_capture.OtherBirthPoint.y
	xf := x_e - x_s
	yf := y_e - y_s
	node := CreatBitTree(x_e, y_e, xf, yf, x_s, y_s)
	GetAllPath(node)
	best_path = FindTheBestPath()
	// fmt.Println("best path:", best_path)

}

func CreatBitTree(x, y, x_f, y_f, x_s, y_s int) *Node {
	n := &Node{
		value: &MapInfo{
			x: x,
			y: y,
		},
	}
	// fmt.Print(n.value)
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
	// fmt.Println("\n paths:", paths)
}

func xianxubianli(node *Node, new_path []*MapInfo) {
	if node == nil {
		return
	}
	if node.left == nil && node.right == nil {
		new_path = append(new_path, node.value)
		paths = append(paths, new_path)
		// fmt.Println("new path:", new_path)
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
		// fmt.Println("best_path point :", *p)
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

			showMapCaptureGame(isFight, my_name, other_name)
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
		showMapCaptureGame(isFight, my_name, other_name)
		if Fight_data_capture.Power <= 0 {
			break
		}
		if p.x == Fight_data_capture.OtherBirthPoint.x && p.y == Fight_data_capture.OtherBirthPoint.y {
			// fmt.Println("YOU WIN!!!!!!!!!!!!!!!!")
			break
		}
		// time.Sleep(1 * time.Second)
		// var ss string
		// fmt.Scanln(&ss)
		// fmt.Printf(ss)
		// best_path = best_path[i+1:]
		i++

	}
	best_path = best_path[i+1:]
}

func IsExist(x, y int) bool {
	if x == Fight_data_capture.BirthPoint.x && y == Fight_data_capture.BirthPoint.y {
		fmt.Printf("this point[%d,%d] is exist by user at point[%d,%d]!\n", x, y, Fight_data_capture.BirthPoint.x, Fight_data_capture.BirthPoint.y)
		return true
	}
	for _, p := range Fight_data_capture.ExcavatePoits {
		if p.x == x && p.y == y {
			fmt.Printf("this point[%d,%d] is exist in excavate points!\n", x, y)
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
			fmt.Println("Fight_data_capture Power is None!")
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
// 		fmt.Println("this excavate point is other birth point!\n")
// 		Fight_data_capture.OtherBirthPoint.x = point.x
// 		Fight_data_capture.OtherBirthPoint.y = point.y
// 		return true
// 	}
// 	return false
// }

func showPath(path []*MapInfo) {
	fmt.Print("the path is :")
	for _, p := range path {
		fmt.Print(*p)
	}
	fmt.Println()
}
