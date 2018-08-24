package main

import (
	"encoding/binary"
	"errors"

	"gbframe"
	//	"net"

	//	"net"
)

type Player struct {
	Name string
	//	Room         string
	//	MyFightDatas map[int]*FightData
	//	BodyData     *SelfData
	//	Conn      *net.Conn
	Server_sess string
	MatchTime   int
	MatchId     int
	Score       int //赢一场 得3分，平一场 得1分，输了 不得分
}

const (
	PLAYERNAMEKEY  = "player_name"
	FIGHTRECORDKEY = "fight_record:"
	PLAYERDATAKEY  = "player:"
	RANKDATAKEY    = "rank_data"
)

func (p *Player) IsPlayerExist() bool {
	for _, player := range MatchPool {
		if p.Name == player.Name {
			return true
		}
	}
	return false

}

func (p *Player) SaveScoreInRank() {
	gameredis.RankDataSave(p.Name, p.Score)
}

func AddPlayerInPool(player *Player, sess string) {
	MatchPool[sess] = player
}

func GetPlayerBySess(sess string) *Player {
	for k, p := range MatchPool {
		if sess == k {
			return p
		}
	}
	return nil

}

//解析接收的消息
func UnmarshalRecMsg(msg []byte) ([]byte, int, error) {
	msgLen := binary.BigEndian.Uint32(msg[0:4])
	msgId := binary.BigEndian.Uint16(msg[4:6])
	if msgLen != (uint32(len(msg)) - uint32(4)) {
		gbframe.Logger_Error.Println("UnmalRecMsg is error,Msg lenth is wrong!,msgLen:", msgLen, "len(msg):", len(msg))
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

func CreatPlayer(name string, sess string, opt int) *Player {
	gbframe.Logger_Info.Println("CreatPlayer, sess:", sess)
	player := &Player{
		Name: name,
		//		Conn:      conn,
		Server_sess: sess,
		MatchTime:   0,
		MatchId:     opt,
		Score:       0,
	}
	return player

}

func CreateRobot() *Player { //创建机器人
	rp := &Player{
		Name:        "robot",
		Server_sess: "robot1",
		MatchTime:   0,
		MatchId:     0,
	}
	return rp
}
