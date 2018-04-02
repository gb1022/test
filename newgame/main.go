package main

import (
	"fmt"
	"net"
	//	"strings"
	//	"encoding/json"
	"protof"
	"time"

	"github.com/golang/protobuf/proto"
)

//const ip = "127.0.0.1"
const addr = ":9999"

var pro = 0

func ReadMessage(msg []byte, n int) *protof.Message1_CS_LoginMessage {
	fmt.Println("recv msg:", string(msg))
	cs_msg := &protof.Message1_CS_LoginMessage{}
	err := proto.Unmarshal(msg[:n], cs_msg)
	if err != nil {
		fmt.Println("proto Unmarshal is error,by:", err)
		return nil
	}
	return cs_msg

}

func WriteMessage() []byte {
	code := proto.Int32(0)
	name := proto.String("gb_test")
	nowtime := time.Now()
	ss := int32(nowtime.Unix())
	logintime := proto.Int32(ss)
	sc_msg := &protof.Message1_SC_LoginMessage{
		Code:      code,
		Name:      name,
		LoginTime: logintime,
	}
	fmt.Println("sc_msg:", sc_msg.String())
	data, _ := proto.Marshal(sc_msg)
	//	n := len(data)
	fmt.Println("data:", data)
	return data
}

func handlerClient(conn net.Conn) {
	defer func() {
		conn.Close()
	}()
	for {
		fmt.Println("remoteAddr:", conn.RemoteAddr().String())
		buf := make([]byte, 2048)
		//		var buf []byte
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("conn.Read err:%v", err)
			break
		}
		//		msg := string(buf[:n-2])
		fmt.Println("Client send message len :", n)
		cs_msg := ReadMessage(buf, n)
		//		msg_json, err := json.Marshal(cs_msg)
		fmt.Println("Client send message :", cs_msg.String())
		//		if msg == "exit" {
		//			fmt.Println("this conn is close!msg:", msg)
		//			conn.Write([]byte("bye bye!"))
		//			return
		//		} else {
		//			conn.Write([]byte("aaaaaaaaaaaaaaa"))
		//		}
		sc_buf := WriteMessage()
		conn.Write(sc_buf)

	}
}

func Start() {
	pro += 1
	fmt.Println("pro:", pro)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		fmt.Println("ResolveIPAddr error:", err)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("ListenTCP error:", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept(),error:", err)
			continue
		}
		go handlerClient(conn)
	}

}

func main() {
	go Start()
	time.Sleep(10 * time.Second)
}
