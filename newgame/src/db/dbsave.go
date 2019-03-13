package db

import (
	"encoding/json"
	"gbframe"

	"strconv"

	"data"

	"github.com/mediocregopher/radix.v2/redis"
)

type RedisData struct {
	Addr        string
	RedisClient *redis.Client
}

const (
	PLAYERNAMEKEY  = "player_name"
	FIGHTRECORDKEY = "fight_record:"
	PLAYERDATAKEY  = "player:"
	RANKDATAKEY    = "rank_data"
)

var Gameredis RedisData

func (rd *RedisData) FightRecordSave(b_fightdata []byte, round int, roomName string, fightTime int) {
	fight_time_str := strconv.Itoa(fightTime)
	key := FIGHTRECORDKEY + roomName + "_" + fight_time_str
	// b_fightdata, err1 := json.Marshal(fight_record)
	// if err1 != nil {
	// 	gbframe.Logger_Error.Println("FightStart json.Marshal is err:", err1)
	// 	return
	// }
	rd.RedisClient.Cmd("HSET", key, round, b_fightdata)
}

func (rd *RedisData) PlayerDataSave(player *data.Player) {
	key := "player:" + player.Name
	b_playerData, err := json.Marshal(player)
	if err != nil {
		gbframe.Logger_Error.Println("playerDataSave  json.Marshal is err:", err)
		return
	}
	rd.RedisClient.Cmd("SET", key, b_playerData)
}

func (rd *RedisData) RankDataSave(name string, score int) {
	rd.RedisClient.Cmd("ZADD", RANKDATAKEY, score, name)
}

func (rd *RedisData) GetAllRankData() []string {
	l, err := rd.RedisClient.Cmd("ZREVRANGE", RANKDATAKEY, 0, -1).List()
	if err != nil {
		gbframe.Logger_Error.Println("Get All RankData is err:", err)
	}
	return l
}

func (rd *RedisData) GetNumRankData(num int) []string { //获取排行榜信息，包括积分
	l, err := rd.RedisClient.Cmd("ZREVRANGE", RANKDATAKEY, 0, num, "WITHSCORES").List()
	if err != nil {
		gbframe.Logger_Error.Println("Get num RankData is err:", err)
	}
	return l
}

func (rd *RedisData) GetNameRankData(name string) int {
	r, err := rd.RedisClient.Cmd("ZREVRANK", RANKDATAKEY, name).Int()
	if err != nil {
		gbframe.Logger_Error.Println("Get Name Rank Data is err:", err)
		return 0
	}
	return r + 1
}

func (rd *RedisData) GetPlayerByName(name string) *data.Player {
	p, err := rd.RedisClient.Cmd("GET", PLAYERDATAKEY+name).Bytes()
	if err != nil {
		gbframe.Logger_Error.Println("GetPlayerByName HGET data is err:", err)
		return nil
	}
	var player_r data.Player
	err = json.Unmarshal(p, &player_r)
	if err != nil {
		gbframe.Logger_Info.Println("GetPlayerByName json.Unmarshal is err:", err)
		return nil
	}
	gbframe.Logger_Info.Println("GetPlayerByName player:", player_r)
	return &player_r

}
