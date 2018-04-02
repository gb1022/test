package main

import (
	//	"bufio"
	"fmt"
	"net"
	//	"os"
	"protof"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
)

const addr = "127.0.0.1:9999"

func ReadMessage(msg []byte, n int) *protof.Message1_SC_LoginMessage {
	fmt.Println("recv msg:", string(msg))
	sc_msg := &protof.Message1_SC_LoginMessage{}
	err := proto.Unmarshal(msg[:n], sc_msg)
	if err != nil {
		fmt.Println("proto Unmarshal is error,by:", err)
		return nil
	}
	return sc_msg

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
	return data
}

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		fmt.Println("ResolveIPAddr error:", err)
		return
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	defer conn.Close()
	if err != nil {
		fmt.Println("DialTCP error:", err)
		return
	}

	for {
		//		inputReader := bufio.NewReader(os.Stdin)
		//		var str string
		//		str, err = inputReader.ReadString('\n')
		msg := WriteMessge()
		conn.Write(msg)
		tmp_msg := &protof.Message1_CS_LoginMessage{}
		proto.Unmarshal(msg, tmp_msg)
		fmt.Println("tmp msg:", tmp_msg.String())
		buf := make([]byte, 20)
		//		var buf []byte
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("conn is err:", err)
			time.Sleep(5 * time.Second)
			continue
		}
		sc_msg := ReadMessage(buf, n)
		fmt.Println("server send msg : ", sc_msg.String())
		if strings.Contains(string(buf[:n]), "bye bye!") {
			fmt.Println("client will close by 5s!")
			time.Sleep(5 * time.Second)
			return
		}
		time.Sleep(5 * time.Second)
	}

}
