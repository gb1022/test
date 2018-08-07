package main

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type Player struct {
	Name         string
	Room         string
	MyFightDatas map[int]*FightData
	BodyData     *SelfData
}

type SelfData struct {
	Life int
}

type FightData struct {
	MapX int
	MapY int
}

//解析接收的消息
func UnmarshalRecMsg(msg []byte) ([]byte, int, error) {
	msgLen := binary.BigEndian.Uint32(msg[0:4])
	msgId := binary.BigEndian.Uint16(msg[4:6])
	if msgLen != (uint32(len(msg)) - uint32(4)) {
		fmt.Println("UnmalRecMsg is error,Msg lenth is wrong!,msgLen:", msgLen, "len(msg):", len(msg))
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
