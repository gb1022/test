package main

import (
	"gbframe"
	"net"
	"protof"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/mediocregopher/radix.v2/redis"
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
const redis_addr = "192.168.1.151:6379"
const redis_db = 14

var NetServers = map[string]*Server{}
var gameredis RedisData

//func ParseConfig(path string){
//	file,err := os.Open(path)
//	if err != nil{
//		gbframe.Logger_Error("ParseConfig is err:",err)
//		return
//	}

//}

func init() {
	conn, err := redis.Dial("tcp", redis_addr)
	gameredis.redisClient = conn
	if err != nil {
		gbframe.Logger_Error.Println("redis connect err:", err)
		return
	}
}

func (s *Server) ReadMessage(msg []byte) (*protof.Message1, int) {
	pMsg, msgId, err := UnmarshalRecMsg(msg)
	if err != nil {
		gbframe.Logger_Error.Println("msg UnmarshalRecMsg is error,by:", err)
		return nil, 0
	}
	cs_msg := &protof.Message1{}
	err = proto.Unmarshal(pMsg, cs_msg)
	if err != nil {
		gbframe.Logger_Error.Println("proto Unmarshal is error,by:", err)
		return nil, 0
	}
	gbframe.Logger_Info.Println("cs_msg:", cs_msg, "msgid:", msgId)
	return cs_msg, msgId

}

func (s *Server) WriteMessage(msgid int, sc_msg *protof.Message1) {
	gbframe.Logger_Info.Println("Send msg,mid:", msgid, "sc_msg:", *sc_msg)
	data, _ := proto.Marshal(sc_msg)
	send_data := MarshalSendMsg(data, msgid)
	s.Service.TranData.OutData <- send_data
}

func (s *Server) GamePosses() {
	gbframe.Logger_Info.Println("......GamePosses start....")
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.sec += 1
			//			gbframe.Info.Println.Println("update,sec:", s.sec, "state:", s.Service.State)
			if s.sec >= 200 || s.Service.TranData.State == false {
				//				s.Service.State = false

				s.logout()
				//				gbframe.Info.Println.Println("Server close !!!!!!!,romoteIp:", s.Service.Conn.RemoteAddr().String())
				s.ConnClose(s.sess)
				return
			}
		case inbyte := <-s.Service.TranData.InData:
			if s.Service.TranData.State == false {
				s.logout()
				//				gbframe.Info.Println.Println("s.Service.TranData.State :", s.Service.TranData.State)
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
	gbframe.Logger_Info.Println("handlerClient...start")
	//	s.wg.Add(1)
	s.GamePosses()

}

func (s *Server) ConnClose(sess string) {
	//	s.Service.Conn.Close()
	//	s.wg.Done()
	//	s.Service.Wg.Done()
	gbframe.Logger_Info.Println("Client close !!!!!!!,client sess:", s.sess)
	delete(NetServers, s.sess)

}

//创建一个服务
func CreateServer(listener *net.TCPListener, id string, sess string) (*Server, error) {
	//	gbframe.Info.Println.Println("222222222222222")
	s, err := gbframe.CreateService(listener, id)
	//	gbframe.Info.Println.Println("3333333333333333333333")
	if err != nil {
		gbframe.Logger_Error.Println("CreateServer is err:", err)
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
	gbframe.Logger_Info.Println("remote ip :", s.TranData.Conn.RemoteAddr().String())
	return ser, nil
}

//开始服务，建立链接，开始接收数据
func Start(addr string, prt string) {
	//	netaddr, _ := net.ResolveTCPAddr("tcp", addr)
	gbframe.Logger_Info.Println("game Server is start...")
	listener, _ := gbframe.ListenTcp("tcp", addr)
	id := 1
	var s *Server
	for {
		id_str := strconv.Itoa(id)
		//		gbframe.Info.Println.Println("1111111111111111111")
		sess := gbframe.MakeSession(id_str, "")
		s, err := CreateServer(listener, id_str, sess)
		//		gbframe.Info.Println.Println("444444444444444444")
		if err != nil {
			return
		}
		//		gbframe.Info.Println.Println("5555555555555555555")

		NetServers[sess] = s
		go s.handlerClient(id_str)
		go s.Service.ServiceProcess()
		//		gbframe.Info.Println.Println("666666666666666666")
		id += 1
	}
	s.Service.Wg.Wait()
}

func (s *Server) login(recMsg *protof.Message1) {
	Login(recMsg, s)
}

func (s *Server) logout() {
	loginout(s)
}

//协议对应接口
func (s *Server) HandleProto(id protof.Message1_ID, recMsg *protof.Message1) {
	//	gbframe.Info.Println.Println("aaaaaaaaaaaaaaaaaa,id:", id)
	player := GetPlayerBySess(s.sess) //从匹配池中查找玩家
	switch id {
	case protof.Message1_PIND:
	//todo
	case protof.Message1_CS_LOGINMESSAGE:
		s.login(recMsg)
		gbframe.Logger_Info.Println("Match pool :", MatchPool)
	case protof.Message1_CS_FIGHTSTART:
		FightReady(recMsg, s)

	case protof.Message1_CS_FIGHTDATA:
		//		gbframe.Info.Println.Println("1111111111111111111")

		//		gbframe.Info.Println.Println("2222222222222222222")
		match_id_str := strconv.Itoa(player.MatchId)
		//		gbframe.Info.Println.Println("3333333333333333")
		room_sess := gbframe.MakeSession(match_id_str, "")
		//		gbframe.Info.Println.Println("4444444444444444")
		fight_room, ok := FightRooms[room_sess]
		if !ok {
			gbframe.Logger_Error.Println("This fight_room is not exist!room_sess:", room_sess)
			return
		}
		//		gbframe.Info.Println.Println("555555555555555555555555")
		fight_room.FightProsses(recMsg, player, room_sess)
		//		gbframe.Info.Println.Println("66666666666666")
	case protof.Message1_CS_GETRANK:
		SCRankData(player.Name, s)
	default:
		gbframe.Logger_Error.Println("No id case! id:", id)
		return

	}
}

func main() {
	Start(addr, "tcp4")
	time.Sleep(10 * time.Second)
}
