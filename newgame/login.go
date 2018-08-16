package main

import (
	"fmt"
	//	"net"
	"protof"
	"time"

	"github.com/golang/protobuf/proto"
)

func (s *Server) Login(recMsg *protof.Message1) {
	name := recMsg.CsLoginMessage.GetName()
	player := CreatPlayer(name, s.sess, int(recMsg.CsLoginMessage.GetOpt()))
	fmt.Println("player name:", name)
	//	MathPool = append(MathPool, *player)
	s.AddPlayerInPool(player)
	fmt.Println("login play:", player)
	code := proto.Int32(0)
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

func (s *Server) logout() {
	if s.Service.State == false {
		delete(MatchPool, s.sess)
		fmt.Println("logout,MathPool:", MatchPool)
	}

}
