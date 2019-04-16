package main

import (
	// "os"
	// "os/exec"

	//	"bufio"
	"fmt"
	"net"

	//	"os"
	"encoding/binary"
	"errors"
	"protof"
	"strconv"

	//	"strings"
	"time"

	"github.com/golang/protobuf/proto"
)

const addr = "127.0.0.1:9999"

//const addr = "192.168.1.151:9999"
const (
	MAP_X = 20
	MAP_Y = 20
)

var Fight_data_capture = Capture_fight_data{}

//type ClientData struct{
//	Conn net.Conn

//}

func ReadMessage(msg []byte) (*protof.Message1, int) {
	// fmt.Println("recv msg:", string(msg[:n]), "msg byte:", msg[:n])
	//	r_data, msgid, _ := UnmarshalRecMsg(msg[:n])
	//	sc_msg := &protof.Message1_SC_LoginMessage{}
	//	err := proto.Unmarshal(r_data, sc_msg)
	//	if err != nil {
	//		fmt.Println("proto Unmarshal is error,by:", err)
	//		return nil, 0
	//	}
	//	return sc_msg, msgid
	pMsg, msgId, err := UnmarshalRecMsg(msg)
	if err != nil {
		return nil, 0
	}
	// fmt.Println("pMsg:", pMsg)
	fmt.Println("msgId:", msgId)
	sc_msg := &protof.Message1{}
	err = proto.Unmarshal(pMsg, sc_msg)
	if err != nil {
		fmt.Println("proto Unmarshal is error,by:", err)
		return nil, 0
	}
	fmt.Println("sc_msg:", sc_msg, "msgid:", msgId)
	return sc_msg, msgId

}

func WriteMessge(msg *protof.Message1, mid int) []byte {
	//	id := proto.Int32(1)
	//	name := proto.String("gb")
	//	opt := proto.Int32(122222)
	//	cs_msg := &protof.Message1_CS_LoginMessage{
	//		Id:   id,
	//		Name: name,
	//		Opt:  opt,
	//	}
	//	msg := &protof.Message1{
	//		CsLoginMessage: cs_msg,
	//	}
	//	fmt.Println("sc_msg:", msg.String())
	data, _ := proto.Marshal(msg)
	//	n := len(data)
	// fmt.Println("data:", data)
	//	mid := int(mid)
	s_data := MarshalSendMsg(data, mid)
	// fmt.Println("s_data:", s_data)
	return s_data
}

//解析接收的消息
func UnmarshalRecMsg(msg []byte) ([]byte, int, error) {
	msgLen := binary.BigEndian.Uint32(msg[0:4])
	msgId := binary.BigEndian.Uint16(msg[4:6])
	if msgLen != (uint32(len(msg)) - uint32(4)) {
		fmt.Println("UnmalRecMsg is error,Msg lenth is wrong!,msgLen:", msgLen, "len(msg):", len(msg))
		// fmt.Println("UnmalRecMsg msg is :", msg)
		return nil, 0, errors.New("Msg lenth is wrong")
	}
	rmsg := msg[6:]
	return rmsg, int(msgId), nil
}

//构建发送的消息
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

func login(conn net.Conn) {
	var ids string
	var name string
	var opt int
	fmt.Print("please input name:")
	fmt.Scanln(&name)
	fmt.Print("please input id:")
	fmt.Scanln(&ids)
	id, err := strconv.Atoi(ids)
	if err != nil {
		fmt.Println("login is wrong,err:", err)
	}
	fmt.Print("please input opt:")
	fmt.Scanln(&opt)
	pid := proto.Int32(int32(id))
	pname := proto.String(name)
	opt_ := proto.Int32(int32(opt))
	cs_msg := &protof.Message1_CS_LoginMessage{
		Id:   pid,
		Name: pname,
		Opt:  opt_,
	}
	msg := &protof.Message1{
		CsLoginMessage: cs_msg,
	}
	w_msg := WriteMessge(msg, int(protof.Message1_CS_LOGINMESSAGE))
	conn.Write(w_msg)
}

func showMapAndPlayerState(sc_msg *protof.Message1) { //显示战斗界面和战斗数据
	round := int(*sc_msg.ScFightData.Round)
	fmt.Println("====================Fight Show================================")
	m_x := int(*sc_msg.ScFightData.MySide.MapX)
	m_y := int(*sc_msg.ScFightData.MySide.MapY)
	o_x := int(*sc_msg.ScFightData.OtherSide.MapX)
	o_y := int(*sc_msg.ScFightData.OtherSide.MapY)
	fmt.Println("Fight Round:", round, "my:", m_x, m_y, "other:", o_x, o_y)
	for i := 0; i < MAP_X; i++ {
		for j := 0; j < MAP_Y; j++ {
			if i == m_x-1 && j == m_y-1 {
				fmt.Print("*")
			} else if i == (o_x-1) && j == (o_y-1) {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println("")
	}
	m_life := *sc_msg.ScFightData.MySide.Life
	o_life := *sc_msg.ScFightData.OtherSide.Life
	fmt.Println("My Name:", *sc_msg.ScFightData.MySide.Name, "My life:", m_life, "| Other Name:", *sc_msg.ScFightData.OtherSide.Name, "-Other life:", o_life)
	fmt.Println("====================Fight Show End================================")
}

func readyFight(conn net.Conn, game_type, robot int) {
	read := true
	cs_fightstart := &protof.Message1_CS_FightStart{
		Isstart:  &read,
		Gametype: proto.Int32(int32(game_type)),
		Torobot:  proto.Int32(int32(robot)),
	}
	msg := &protof.Message1{
		CsFightStart: cs_fightstart,
	}
	w_msg := WriteMessge(msg, int(protof.Message1_CS_FIGHTSTART))
	conn.Write(w_msg)
	// SendPingMsg(conn, 0)
}

func fight(sc_msg *protof.Message1, conn net.Conn) bool {
	showMapAndPlayerState(sc_msg)
	if sc_msg.ScFightData.GetResult() == 1 {
		fmt.Println("Fight Over!\n  Congratulation！ You Win!!!!!!!!")
		return true
	} else if sc_msg.ScFightData.GetResult() == 2 {
		fmt.Println("Fight Over!\nOh I am sorry! You Lose!!!!!!!!!")
		return true
	} else if sc_msg.ScFightData.GetResult() == 3 {
		fmt.Println("Fight Over!\n   It is Draw!!!!!!!!!!")
		return true
	}
	//	var x int //位移 x
	//	var y int //位移 y
	//	var s int //速度
	//	var a int //攻击
	//	fmt.Print("please input x:")
	//	fmt.Scanln(&x)
	//	fmt.Print("please input y:")
	//	fmt.Scanln(&y)
	//	fmt.Print("please input speed:")
	//	fmt.Scanln(&s)
	//	fmt.Print("please input attack:")
	//	fmt.Scanln(&a)
	x, y, s, a := fightInput(sc_msg)
	m_x := int32(x)
	m_y := int32(y)
	attack := int32(a)
	speed := int32(s)
	fight_msg := &protof.Message1_CS_FightData{
		Speed:  &speed,
		Attack: &attack,
		MoveX:  &m_x,
		MoveY:  &m_y,
	}
	cs_msg := &protof.Message1{
		CsFightData: fight_msg,
	}
	w_msg := WriteMessge(cs_msg, int(protof.Message1_CS_FIGHTDATA))
	conn.Write(w_msg)
	return false

}

func SendRankMsg(conn net.Conn) {
	fmt.Println("111111111111")
	cs_msg := &protof.Message1{
		CsGetRank: &protof.Message1_CS_GetRank{
			Code: proto.Int32(0),
		},
	}
	mid := int(protof.Message1_CS_GETRANK)
	w_msg := WriteMessge(cs_msg, mid)
	conn.Write(w_msg)
	fmt.Println("Send Rank Requist Successful!")

}

func ShowRank(sc_msg *protof.Message1) {
	if *sc_msg.ScGetRank.Code != 0 {
		fmt.Println("Rank is Error!")
		return
	}

	fmt.Println("================Rank List========================")
	fmt.Println("排名        名字         分数")
	for k, v := range sc_msg.ScGetRank.RankData {
		fmt.Println(k+1, "        ", *v.Name, "         ", *v.Score)
	}
	fmt.Println("----------------------------------------------")
	fmt.Println("我的排名:", sc_msg.ScGetRank.MyRanking.GetRanking(), "我的分数:", sc_msg.ScGetRank.MyRanking.GetScore())
	fmt.Println("=================================================")
}

func SendPingMsg(conn net.Conn, code int32) {
	cs_msg := &protof.Message1{
		CsPing: &protof.Message1_CS_Ping{
			Code: proto.Int32(code),
		},
	}
	mid := int(protof.Message1_CSPIND)
	w_msg := WriteMessge(cs_msg, mid)
	conn.Write(w_msg)
	fmt.Println("fight ready time:", time.Now().Second())
}

func main() {
	fmt.Println("1")
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	fmt.Println("2")
	if err != nil {
		fmt.Println("ResolveIPAddr error:", err)
		return
	}
	fmt.Println("3")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("conn is err:", err)
		return
	}
	fmt.Println("4")

	fmt.Println("5")
	if err != nil {
		fmt.Println("DialTCP error:", err)
		return
	}
	fmt.Println("7")
	login(conn) //登陆
	//开始战斗：
	ticker := time.NewTicker(time.Second * 1)
	defer func() {
		ticker.Stop()
		conn.Close()
		fmt.Println("The server connect is over!")
	}()
	var server_time float32
	var match_flag int
	for {
		//		inputReader := bufio.NewReader(os.Stdin)
		//		var str string
		//		str, err = inputReader.ReadString('\n')
		//		msg := WriteMessge()
		//		conn.Write(msg)
		//		tmp_msg := &protof.Message1_CS_LoginMessage{}
		//		proto.Unmarshal(msg, tmp_msg)
		//		fmt.Println("tmp msg:", tmp_msg.String())

		buf := make([]byte, 65535)
		//		var buf []byte
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("conn Read is err:", err)
			break
		}

		sc_msg, msgid := ReadMessage(buf[:n])
		switch msgid {
		case int(protof.Message1_SCPIND):
			server_time = float32(sc_msg.ScPing.GetTime())
			fmt.Println("Server time is ", server_time)
			if match_flag == 1 {
				SendPingMsg(conn, 1)
			}
		case int(protof.Message1_SC_LOGINMESSAGE):
			if sc_msg.ScLoginMessage.GetCode() == 0 {
				fmt.Println("client login success!")
			} else {
				fmt.Println("client login fail!please reset game!")
				login(conn)
				continue
			}
		case int(protof.Message1_SC_FIGHTSTART):
			fmt.Println("sc msg :", sc_msg)
			if sc_msg.ScFightStart.GetIsstartA() && sc_msg.ScFightStart.GetIsstartB() {
				fmt.Println("Fight is begin now!")
				continue
			} else {
				fmt.Println("somebody is not ready! please wait a moment...")
				// SendPingMsg(conn, 1)
				// match_flag = 1
				continue
			}
		case int(protof.Message1_SC_FIGHTDATA):
			Flush屏幕()
			isover := fight(sc_msg, conn)
			if isover {
				fmt.Println("######################################")
			} else {
				continue
			}
		case int(protof.Message1_SC_FIGHTDATA_TUNNEL_CAPTURE):
			Flush屏幕()
			isover := fight_tunnel_capture(sc_msg, conn)
			if isover {
				fmt.Println("######################################")
			} else {
				continue
			}
		case int(protof.Message1_SC_GETRANK):
			ShowRank(sc_msg)
		}
		for {
			var choice int
			fmt.Print("1:Are you ready fighting!\n")
			fmt.Print("2:Show the Rank!\n")
			fmt.Print("3:Do you quit!\n")
			fmt.Print("Which do you choice?\n")
			fmt.Scanln(&choice)
			if choice == 1 {
				var ch int
				var robot int
				fmt.Print("1 map game; \n2 tunnel_capture_game\n")
				fmt.Println("Which game do you choose?")
				fmt.Scanln(&ch)
				fmt.Print("1 with player?\n2 with robot?\n")
				fmt.Println("Which type?")
				fmt.Scanln(&robot)
				fmt.Println("You are ready fight!")

				readyFight(conn, ch, robot)
				break
			} else if choice == 2 {
				SendRankMsg(conn)
				break
			} else if choice == 3 {
				fmt.Println("bye bye!")
				return
			}
		}

		//		fmt.Println("server send msg : ", sc_msg.String(), "msgid:", msgid)
		//		if strings.Contains(string(buf[:n]), "bye bye!") {
		//			fmt.Println("client will close by 5s!")
		//			time.Sleep(10 * time.Second)
		//			return
		//		}
		//		time.Sleep(5 * time.Second)
	}

}
