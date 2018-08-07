package main

import (
	"fmt"
)

func Login(name string, room string) {
	CreatPlayer(name, room)
}

func CreatPlayer(name string, room string) *Player {
	player := &Player{
		Name: name,
		Room: room,
	}
	return player

}
