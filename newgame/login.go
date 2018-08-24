package main

import (
	//	"fmt"
	//	"net"
	"gbframe"
	"protof"
	"time"

	"github.com/golang/protobuf/proto"
)

func Login(recMsg *protof.Message1, s *Server) {
	name := recMsg.CsLoginMessage.GetName()
	code := proto.Int32(0)
	gbframe.Logger_Info.Println("player name:", name)

	player := gameredis.GetPlayerByName(name)
	if player == nil {
		player = CreatPlayer(name, s.sess, int(recMsg.CsLoginMessage.GetOpt()))
	}
	gbframe.Logger_Info.Println("login play:", player)
	if player.IsPlayerExist() {
		code = proto.Int32(1)
	} else {
		AddPlayerInPool(player, s.sess)
		gameredis.redisClient.Cmd("SADD", PLAYERNAMEKEY, name)
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

func loginout(s *Server) {
	if s.Service.State == false {
		delete(MatchPool, s.sess)
		gbframe.Logger_Info.Println("logout,MathPool:", MatchPool)
	}

}
