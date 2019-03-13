package gamenet

import (
	"gbframe"
	// "net"
	"protof"
	// "strconv"
	"time"

	"data"

	"github.com/golang/protobuf/proto"
	// "github.com/mediocregopher/radix.v2/redis"
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
	Sec     int
	Sess    string
	//	wg      *sync.WaitGroup
	Outtime int
}

var NetServers = map[string]*Server{}

func (s *Server) ReadMessage(msg []byte) (*protof.Message1, int) {
	pMsg, msgId, err := data.UnmarshalRecMsg(msg)
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
	datas, _ := proto.Marshal(sc_msg)
	send_data := data.MarshalSendMsg(datas, msgid)
	s.Service.TranData.OutData <- send_data
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

func (s *Server) ConnClose(sess string) {
	//	s.Service.Conn.Close()
	//	s.wg.Done()
	//	s.Service.Wg.Done()
	gbframe.Logger_Info.Println("Client close !!!!!!!,client sess:", s.Sess)
	//	s.logout()
	// s.Service.Conn.Close()
	s.Service.ConnClose()

	delete(NetServers, s.Sess)
	gbframe.Logger_Info.Println("ConnClose !!!!!!!,NetServers:", NetServers)

}

// func (s *Server) login(recMsg *protof.Message1) {
// 	login.Login(recMsg, s)
// }

// func (s *Server) logout() {
// 	login.loginout(s)
// }

func (s *Server) SendPingMsg() {
	nowtime := float32(time.Now().Unix())
	sc_msg := &protof.Message1{
		ScPing: &protof.Message1_SC_Ping{
			Time: proto.Float32(nowtime),
		},
	}
	mid := int(protof.Message1_SCPIND)
	s.WriteMessage(mid, sc_msg)
}
