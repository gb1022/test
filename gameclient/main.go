package main

import (
	//	"bufio"
	"fmt"
	"net"
	//	"os"
	"encoding/binary"
	"errors"
	"protof"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
)

const addr = "127.0.0.1:9999"

func ReadMessage(msg []byte, n int) (*protof.Message1_SC_LoginMessage, int) {
	fmt.Println("recv msg:", string(msg))
	r_data, msgid, _ := UnmarshalRecMsg(msg[:n])
	sc_msg := &protof.Message1_SC_LoginMessage{}
	err := proto.Unmarshal(r_data, sc_msg)
	if err != nil {
		fmt.Println("proto Unmarshal is error,by:", err)
		return nil, 0
	}
	return sc_msg, msgid

}

func WriteMessge() []byte {
	id := proto.Int32(1)
	name := proto.String("gb")
	opt := proto.Int32(122222)
	cs_msg := &protof.Message1_CS_LoginMessage{
		Id:   id,
		Name: name,
		Opt:  opt,
	}
	fmt.Println("sc_msg:", cs_msg.String())
	data, _ := proto.Marshal(cs_msg)
	//	n := len(data)
	fmt.Println("data:", data)
	mid := int(protof.Message1_CS_LOGINMESSAGE)
	s_data := MarshalSendMsg(data, mid)
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
	fmt.Println("4")
	defer func() {
		fmt.Println("5")
		conn.Close()
	}()
	fmt.Println("6")
	if err != nil {
		fmt.Println("DialTCP error:", err)
		return
	}
	fmt.Println("7")

	for {
		//		inputReader := bufio.NewReader(os.Stdin)
		//		var str string
		//		str, err = inputReader.ReadString('\n')
		msg := WriteMessge()
		conn.Write(msg)
		tmp_msg := &protof.Message1_CS_LoginMessage{}
		proto.Unmarshal(msg, tmp_msg)
		fmt.Println("tmp msg:", tmp_msg.String())
		buf := make([]byte, 2048)
		//		var buf []byte
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("conn is err:", err)
			time.Sleep(5 * time.Second)
			continue
		}
		sc_msg, msgid := ReadMessage(buf, n)
		fmt.Println("server send msg : ", sc_msg.String(), "msgid:", msgid)
		if strings.Contains(string(buf[:n]), "bye bye!") {
			fmt.Println("client will close by 5s!")
			time.Sleep(5 * time.Second)
			return
		}
		time.Sleep(5 * time.Second)
	}

}
