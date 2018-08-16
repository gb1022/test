package main

import (
	"fmt"
	"net"
	"strconv"
	//	"strings"
	//	"encoding/json"
	"gbframe"
	"protof"
	//	"sync"

	"time"

	"github.com/golang/protobuf/proto"
)

type GameData struct {
	FrameTranData *gbframe.TransportData
	GamerId       string
	RecMsg        string
	SendMsg       string
}

type Server struct {
	Service  *gbframe.Service
	GameChan chan GameData
	//	Conn     net.Conn
	GamerId string
	RecMsg  string
	SendMsg string
	sec     int
	sess    string
	//	wg      *sync.WaitGroup
}

//const ip = "127.0.0.1"
const addr = ":9999"

var NetServers = map[string]*Server{}

func (s *Server) ReadMessage(msg []byte) (*protof.Message1, int) {
	pMsg, msgId, err := UnmarshalRecMsg(msg)
	if err != nil {
		return nil, 0
	}
	cs_msg := &protof.Message1{}
	err = proto.Unmarshal(pMsg, cs_msg)
	if err != nil {
		fmt.Println("proto Unmarshal is error,by:", err)
		return nil, 0
	}
	fmt.Println("cs_msg:", cs_msg, "msgid:", msgId)
	return cs_msg, msgId

}

func (s *Server) WriteMessage(msgid int, sc_msg *protof.Message1) {

	data, _ := proto.Marshal(sc_msg)
	send_data := MarshalSendMsg(data, msgid)
	fmt.Println("send_data:", send_data)
	s.Service.TranData.OutData <- send_data
}

func (s *Server) GamePosses() {
	fmt.Println("......GamePosses start....")
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.sec += 1
			fmt.Println("update,sec:", s.sec, "state:", s.Service.State)
			if s.sec >= 60 || s.Service.TranData.State == false {
				//				s.Service.State = false

				s.logout()
				//				fmt.Println("Server close !!!!!!!,romoteIp:", s.Service.Conn.RemoteAddr().String())
				s.ConnClose(s.sess)
				return
			}
		case inbyte := <-s.Service.TranData.InData:
			if s.Service.TranData.State == false {
				s.logout()
				//				fmt.Println("s.Service.TranData.State :", s.Service.TranData.State)
				break
			}
			recMsg, msgid := s.ReadMessage(inbyte)
			mid := protof.Message1_ID(msgid)
			s.HandleProto(mid, recMsg)
			//			s.WriteMessage(msgid)
			s.sec = 0
		}
	}

}

func (s *Server) CreateGameData(t *gbframe.TransportData, gid string) *GameData {
	g := &GameData{
		FrameTranData: t,
		GamerId:       gid,
		RecMsg:        "",
		SendMsg:       "",
	}
	return g
}

func (s *Server) handlerClient(id string) {
	fmt.Println("handlerClient...start")
	//	s.wg.Add(1)
	s.GamePosses()

}

func (s *Server) ConnClose(sess string) {
	//	s.Service.Conn.Close()
	//	s.wg.Done()
	//	s.Service.Wg.Done()
	fmt.Println("Client close !!!!!!!,client sess:", s.sess)
	delete(NetServers, s.sess)

}

//创建一个服务
func CreateServer(listener *net.TCPListener, id string, sess string) (*Server, error) {
	//	fmt.Println("222222222222222")
	s, err := gbframe.CreateService(listener, id)
	//	fmt.Println("3333333333333333333333")
	if err != nil {
		fmt.Println("CreateServer is err:", err)
		return nil, err
	}
	ser := &Server{
		Service:  s,
		GameChan: make(chan GameData),
		GamerId:  id,
		RecMsg:   "",
		SendMsg:  "",
		sess:     sess,
		//		wg:       s.Wg,
	}
	fmt.Println("remote ip :", s.TranData.Conn.RemoteAddr().String())
	return ser, nil
}

//开始服务，建立链接，开始接收数据
func Start(addr string, prt string) {
	//	netaddr, _ := net.ResolveTCPAddr("tcp", addr)
	fmt.Println("game Server is start...")
	listener, _ := gbframe.ListenTcp("tcp", addr)
	id := 1
	var s *Server
	for {
		id_str := strconv.Itoa(id)
		//		fmt.Println("1111111111111111111")
		sess := MakeSession(id_str, "")
		s, err := CreateServer(listener, id_str, sess)
		//		fmt.Println("444444444444444444")
		if err != nil {
			return
		}
		//		fmt.Println("5555555555555555555")

		NetServers[sess] = s
		go s.handlerClient(id_str)
		go s.Service.ServiceProcess()
		//		fmt.Println("666666666666666666")
		id += 1
	}
	s.Service.Wg.Wait()
}

//协议对应接口
func (s *Server) HandleProto(id protof.Message1_ID, recMsg *protof.Message1) {
	switch id {
	case protof.Message1_PIND:
	//todo
	case protof.Message1_CS_LOGINMESSAGE:
		s.Login(recMsg)
		fmt.Println("Match pool :", MatchPool)
	case protof.Message1_CS_FIGHTSTART:
		s.FightReady(recMsg)

	case protof.Message1_CS_FIGHTDATA:

	default:
		return

	}
}

func main() {
	Start(addr, "tcp4")
	time.Sleep(10 * time.Second)
}
