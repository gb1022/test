package login

import (
	//	"fmt"
	//	"net"
	"data"
	"db"
	"gamenet"
	"gbframe"
	"protof"
	"time"

	"github.com/golang/protobuf/proto"
)

func Login(recMsg *protof.Message1, s *gamenet.Server) {
	name := recMsg.CsLoginMessage.GetName()
	code := proto.Int32(0)
	gbframe.Logger_Info.Println("player name:", name, "s.Sess:", s.Sess)

	player := db.Gameredis.GetPlayerByName(name)
	if player == nil {
		player = data.CreatPlayer(name, s.Sess, int(recMsg.CsLoginMessage.GetOpt()))
	} else {
		player.InitPlayer(s.Sess, int(recMsg.CsLoginMessage.GetOpt()))
	}
	gbframe.Logger_Info.Println("login play:", player)
	if player.IsPlayerExist() {
		gbframe.Logger_Error.Println("this player is exist!")
		code = proto.Int32(1)
	} else {
		data.AddPlayerInPool(player, s.Sess)
		db.Gameredis.RedisClient.Cmd("SADD", db.PLAYERNAMEKEY, name)
	}
	sc_name := proto.String(name)
	nowtime := time.Now()
	ss := int32(nowtime.Unix())
	logintime := proto.Int32(ss)

	sc_login := &protof.Message1_SC_LoginMessage{
		Code:      code,
		Name:      sc_name,
		LoginTime: logintime,
	}
	sc_msg := &protof.Message1{
		ScLoginMessage: sc_login,
	}
	mid := int(protof.Message1_SC_LOGINMESSAGE)
	s.WriteMessage(mid, sc_msg)
	//	fmt.Println("sc_msg:", sc_msg.String())

}

func Logout(s *gamenet.Server, player *data.Player) {
	//	if s.Service.State == false {
	player.IsOnline = false
	gbframe.Logger_Info.Println("logout,s.sess:", s.Sess)
	delete(data.MatchPool, s.Sess)
	gbframe.Logger_Info.Println("logout,MathPool:", data.MatchPool)
	s.ConnClose(s.Sess)
	//	}

}
