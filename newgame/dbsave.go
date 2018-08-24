package main

import (
	"encoding/json"
	"gbframe"
	"strconv"

	"github.com/mediocregopher/radix.v2/redis"
)

type RedisData struct {
	Addr        string
	redisClient *redis.Client
}

func (rd *RedisData) FightRecordSave(fr FightRoom, round int, fight_record FightRecord) {
	fight_time_str := strconv.Itoa(int(fr.FightTime))
	key := FIGHTRECORDKEY + fr.Room + "_" + fight_time_str
	b_fightdata, err1 := json.Marshal(fight_record)
	if err1 != nil {
		gbframe.Logger_Error.Println("FightStart json.Marshal is err:", err1)
		return
	}
	rd.redisClient.Cmd("HSET", key, round, b_fightdata)
}

func (rd *RedisData) PlayerDataSave(player Player) {
	key := "player:" + player.Name
	b_playerData, err := json.Marshal(player)
	if err != nil {
		gbframe.Logger_Error.Println("playerDataSave  json.Marshal is err:", err)
		return
	}
	rd.redisClient.Cmd("SET", key, b_playerData)
}

func (rd *RedisData) RankDataSave(name string, score int) {
	rd.redisClient.Cmd("ZADD", RANKDATAKEY, score, name)
}

func (rd *RedisData) GetAllRankData() []string {
	l, err := rd.redisClient.Cmd("ZREVRANGE", RANKDATAKEY, 0, -1).List()
	if err != nil {
		gbframe.Logger_Error.Println("Get All RankData is err:", err)
	}
	return l
}

func (rd *RedisData) GetNumRankData(num int) []string { //获取排行榜信息，包括积分
	l, err := rd.redisClient.Cmd("ZREVRANGE", RANKDATAKEY, 0, num, "WITHSCORES").List()
	if err != nil {
		gbframe.Logger_Error.Println("Get num RankData is err:", err)
	}
	return l
}

func (rd *RedisData) GetNameRankData(name string) int {
	r, err := rd.redisClient.Cmd("ZREVRANK", RANKDATAKEY, name).Int()
	if err != nil {
		gbframe.Logger_Error.Println("Get Name Rank Data is err:", err)
		return 0
	}
	return r + 1
}

func (rd *RedisData) GetPlayerByName(name string) *Player {
	p, err := rd.redisClient.Cmd("GET", PLAYERDATAKEY+name).Bytes()
	if err != nil {
		gbframe.Logger_Error.Println("GetPlayerByName HGET data is err:", err)
		return nil
	}
	var player_r Player
	err = json.Unmarshal(p, &player_r)
	if err != nil {
		gbframe.Logger_Info.Println("GetPlayerByName json.Unmarshal is err:", err)
		return nil
	}
	gbframe.Logger_Info.Println("GetPlayerByName player:", player_r)
	return &player_r

}
