package main

import (
	"fight_tunnel_capture"
	"gbframe"
	"net"
	"protof"
	"rank"
	"strconv"
	"time"

	"data"
	"db"
	"fight"
	"gamenet"
	"login"
	"match"

	// "github.com/golang/protobuf/proto"
	"github.com/mediocregopher/radix.v2/redis"
)

//const ip = "127.0.0.1"
const addr = ":9999"
const redis_addr = "192.168.1.152:6379"
const redis_db = 14
const (
	MAP_FIGHT          = 1
	TUNNEL_CPTURE_GAME = 2
)

//func ParseConfig(path string){
//	file,err := os.Open(path)
//	if err != nil{
//		gbframe.Logger_Error("ParseConfig is err:",err)
//		return
//	}

//}

func init() {
	conn, err := redis.Dial("tcp", redis_addr)
	if err != nil {
		gbframe.Logger_Error.Println("redis connect err:", err)
		return
	}
	db.Gameredis.RedisClient = conn
}

//创建一个服务
func CreateServer(listener *net.TCPListener, id string, sess string) (*gamenet.Server, error) {
	//	gbframe.Info.Println.Println("222222222222222")
	s, err := gbframe.CreateService(listener, id)
	//	gbframe.Info.Println.Println("3333333333333333333333")
	if err != nil {
		gbframe.Logger_Error.Println("CreateServer is err:", err)
		return nil, err
	}
	ser := &gamenet.Server{
		Service:  s,
		GameChan: make(chan gamenet.GameData),
		GamerId:  id,
		RecMsg:   "",
		SendMsg:  "",
		Sess:     sess,
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
	var s *gamenet.Server
	for {
		id_str := strconv.Itoa(id)
		gbframe.Logger_Info.Println("1111111111111111111,id:", id_str)
		sess := gbframe.MakeSession(id_str, "")
		gbframe.Logger_Info.Println("1111111111111111111,sess:", sess)
		s, err := CreateServer(listener, id_str, sess)
		//		gbframe.Info.Println.Println("444444444444444444")
		if err != nil {
			return
		}
		//		gbframe.Info.Println.Println("5555555555555555555")

		gamenet.NetServers[sess] = s
		gbframe.Logger_Info.Println("gamenet.NetServers:", gamenet.NetServers)
		s.Service.Wg.Add(1)
		go handlerClient(id_str, s)
		go s.Service.ServiceProcess()
		//		gbframe.Info.Println.Println("666666666666666666")
		id += 1
	}
	s.Service.Wg.Wait()
}

func handlerClient(id string, s *gamenet.Server) {
	gbframe.Logger_Info.Println("handlerClient...start")
	//	s.wg.Add(1)
	GamePosses(s)

}

func GamePosses(s *gamenet.Server) {
	gbframe.Logger_Info.Println("......GamePosses start....")

	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.Sec += 1
			//			gbframe.Info.Println.Println("update,sec:", s.sec, "state:", s.Service.State)
			if s.Sec >= 600 || s.Service.State == false {
				//				s.Service.State = false
				player := data.GetPlayerBySess(s.Sess)
				if player != nil {
					if ok, mroom := match.FindMatchRoomByPlayer(player); ok {
						match.DeleteMatchRoomPool(mroom)
					}
				} else {
					gbframe.Logger_Warning.Println("player is nil;s.Sess:", s.Sess)
					// s.ConnClose(s.Sess)
					s.Service.Wg.Done()
					return
				}
				s.Service.Wg.Done()
				login.Logout(s, player)
				//				gbframe.Info.Println.Println("Server close !!!!!!!,romoteIp:", s.Service.Conn.RemoteAddr().String())

				return
			} else {
				player := data.GetPlayerBySess(s.Sess) //从匹配池中查找玩家
				update(s, player)
			}
		case inbyte := <-s.Service.TranData.InData:
			if s.Service.State == false {
				player := data.GetPlayerBySess(s.Sess)
				if player != nil {
					if ok, mroom := match.FindMatchRoomByPlayer(player); ok {
						match.DeleteMatchRoomPool(mroom)
					}

				}

				login.Logout(s, player)
				//				gbframe.Info.Println.Println("s.Service.TranData.State :", s.Service.TranData.State)
				break
			}
			gbframe.Logger_Info.Println("11111111111111111111 inbyte:", inbyte)
			recMsg, msgid := s.ReadMessage(inbyte)
			mid := protof.Message1_ID(msgid)
			HandleProto(mid, recMsg, s)
			//			s.WriteMessage(msgid)
			s.Sec = 0
		}
	}

}

//协议对应接口
func HandleProto(id protof.Message1_ID, recMsg *protof.Message1, s *gamenet.Server) {
	//	gbframe.Info.Println.Println("aaaaaaaaaaaaaaaaaa,id:", id)
	player := data.GetPlayerBySess(s.Sess)
	switch id {
	case protof.Message1_CSPIND:
		code := recMsg.CsPing.GetCode()
		if code != 1 {
			s.Outtime += 1
		} else {
			s.Outtime = 0
		}
		s.SendPingMsg()

	case protof.Message1_CS_LOGINMESSAGE:
		login.Login(recMsg, s)
		gbframe.Logger_Info.Println("Match pool :", data.MatchPool, "MatchRoom pool:", match.MatchRoomPool)
	case protof.Message1_CS_FIGHTSTART:
		ok, mroom := match.StartMatch(recMsg, s)
		if !ok {
			return
		}
		time.Sleep(10 * time.Millisecond) //暂停1秒，免得消息粘连在一起
		game_type := int(recMsg.CsFightStart.GetGametype())
		switch game_type {
		case MAP_FIGHT:
			fight.FightReady(mroom.PlayerList)
		case TUNNEL_CPTURE_GAME:
			fight_tunnel_capture.FightReady(mroom.PlayerList)
		}

		match.DeleteMatchRoomPool(mroom)

	case protof.Message1_CS_FIGHTDATA:
		match_id_str := strconv.Itoa(player.MatchId)
		room_sess := gbframe.MakeSession(match_id_str, "")
		fight_room, ok := fight.FightRooms[room_sess]
		if !ok {
			gbframe.Logger_Error.Println("This fight_room is not exist!room_sess:", room_sess)
			return
		}
		//		gbframe.Info.Println.Println("555555555555555555555555")
		fight_room.FightProsses(recMsg, player, room_sess)
		//		gbframe.Info.Println.Println("66666666666666")
	case protof.Message1_CS_FIGHTDATA_TUNNEL_CAPTURE:
		time.Sleep(1 * time.Second) //暂停1下，免得消息粘连在一起
		match_id_str := strconv.Itoa(player.MatchId)
		room_sess := gbframe.MakeSession(match_id_str, "")
		fight_room, ok := fight_tunnel_capture.FightRooms[room_sess]
		if !ok {
			gbframe.Logger_Error.Println("This fight_room is not exist!room_sess:", room_sess)
			return
		}
		fight_room.FightProsses(recMsg, player, room_sess)
	case protof.Message1_CS_GETRANK:
		rank.SCRankData(player.Name, s)
	default:
		gbframe.Logger_Error.Println("No id case! id:", id)
		return

	}
}

func update(s *gamenet.Server, player *data.Player) {
	if player == nil {
		gbframe.Logger_Error.Println("this player is nil!sess is ", s.Sess)
		return
	} else {
		if !player.MatchFlag {
			return
		}
		player.MatchTime++
		// if player.MatchTime > 20 {
		// 	match.MatchRobot(player, s)
		// }
	}

}

func main() {
	Start(addr, "tcp4")
	time.Sleep(10 * time.Second)
}
