package main

import (
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
const (
	MAP_X = 20
	MAP_Y = 20
)

//type ClientData struct{
//	Conn net.Conn

//}

func ReadMessage(msg []byte) (*protof.Message1, int) {
	//	fmt.Println("recv msg:", string(msg[:n]), "msg byte:", msg[:n])
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
	fmt.Println("data:", data)
	//	mid := int(mid)
	s_data := MarshalSendMsg(data, mid)
	fmt.Println("s_data:", s_data)
	return s_data
}

//解析接收的消息
func UnmarshalRecMsg(msg []byte) ([]byte, int, error) {
	msgLen := binary.BigEndian.Uint32(msg[0:4])
	msgId := binary.BigEndian.Uint16(msg[4:6])
	if msgLen != (uint32(len(msg)) - uint32(4)) {
		fmt.Println("UnmalRecMsg is error,Msg lenth is wrong!,msgLen:", msgLen, "len(msg):", len(msg))
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
	fmt.Println("====================================================")
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

}

func readFight(conn net.Conn) {
	read := true
	cs_fightstart := &protof.Message1_CS_FightStart{
		Isstart: &read,
	}
	msg := &protof.Message1{
		CsFightStart: cs_fightstart,
	}
	w_msg := WriteMessge(msg, int(protof.Message1_CS_FIGHTSTART))
	conn.Write(w_msg)
}

func fight(sc_msg *protof.Message1, conn net.Conn) {
	showMapAndPlayerState(sc_msg)
	var x int //位移 x
	var y int //位移 y
	var s int //速度
	var a int //攻击
	fmt.Print("please input x:")
	fmt.Scanln(&x)
	fmt.Print("please input y:")
	fmt.Scanln(&y)
	fmt.Print("please input speed:")
	fmt.Scanln(&s)
	fmt.Print("please input attack:")
	fmt.Scanln(&a)
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
	w_msg := WriteMessge(cs_msg, int(protof.Message1_SC_FIGHTDATA))
	conn.Write(w_msg)

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
	for {
		//		inputReader := bufio.NewReader(os.Stdin)
		//		var str string
		//		str, err = inputReader.ReadString('\n')
		//		msg := WriteMessge()
		//		conn.Write(msg)
		//		tmp_msg := &protof.Message1_CS_LoginMessage{}
		//		proto.Unmarshal(msg, tmp_msg)
		//		fmt.Println("tmp msg:", tmp_msg.String())

		buf := make([]byte, 2048)
		//		var buf []byte
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("conn is err:", err)
			//			time.Sleep(5 * time.Second)

			break
		}

		sc_msg, msgid := ReadMessage(buf[:n])
		if msgid == int(protof.Message1_SC_LOGINMESSAGE) {
			if sc_msg.ScLoginMessage.GetCode() == 0 {
				fmt.Println("client login success!")
			} else {
				fmt.Println("client login fail!please restart game!")
				break
			}
		} else if msgid == int(protof.Message1_SC_FIGHTSTART) {
			if sc_msg.ScFightStart.GetIsstartA() && sc_msg.ScFightStart.GetIsstartB() {
				fmt.Println("Fight is begin now!")
				continue
			} else {
				fmt.Println("somebody is not ready! please wait a moment...")
				continue
			}

		} else if msgid == int(protof.Message1_SC_FIGHTDATA) {
			fight(sc_msg, conn)
			continue
		}
		for {
			var choice int
			fmt.Print("1:Are you ready fight!")
			fmt.Print("2:Are you quit!")
			fmt.Print("Which do you choice?\n")
			fmt.Scanln(&choice)
			if choice == 1 {
				fmt.Println("You are ready fight!")
				readFight(conn)
				break
			} else if choice == 2 {
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
