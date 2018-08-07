package main

import (
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
