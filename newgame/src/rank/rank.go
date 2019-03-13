package rank

import (
	"db"
	"gamenet"
	"gbframe"
	"protof"
	"strconv"

	"github.com/golang/protobuf/proto"
)

const (
	RANK_MAX_NUM = 20
)

type RankData struct {
	Name    string
	Score   int
	Ranking int
}

func SCRankData(name string, s *gamenet.Server) {
	gbframe.Logger_Info.Println("Get Rank Data begin!")
	rlist := GetRankData()
	myr := GetPlayerRankData(name, rlist)
	code := 0
	var proto_rlist []*protof.Message1_RankData
	for _, v := range rlist {
		protof_rank_data := &protof.Message1_RankData{
			Name:    proto.String(v.Name),
			Score:   proto.Int32(int32(v.Score)),
			Ranking: proto.Int32(int32(v.Ranking)),
		}
		proto_rlist = append(proto_rlist, protof_rank_data)
	}
	my_rank := &protof.Message1_RankData{
		Name:    &myr.Name,
		Score:   proto.Int32(int32(myr.Score)),
		Ranking: proto.Int32(int32(myr.Ranking)),
	}
	sc_msg := &protof.Message1{
		ScGetRank: &protof.Message1_SC_GetRank{
			Code:      proto.Int32(int32(code)),
			RankData:  proto_rlist,
			MyRanking: my_rank,
		},
	}
	mid := int(protof.Message1_SC_GETRANK)
	s.WriteMessage(mid, sc_msg)
}

func GetRankData() []*RankData {
	l := db.Gameredis.GetNumRankData(RANK_MAX_NUM)
	var rank_list []*RankData
	for i := 0; i < len(l); {
		var rd RankData
		rd.Name = l[i]
		rd.Score, _ = strconv.Atoi(l[i+1])
		rd.Ranking = i/2 + 1
		rank_list = append(rank_list, &rd)
		i = i + 2
		gbframe.Logger_Info.Println("====", rd)
	}
	gbframe.Logger_Info.Println("rank_list:", rank_list)
	return rank_list
}

func GetPlayerRankData(name string, rank_list []*RankData) *RankData {
	//	var my_rank RankData
	for _, v := range rank_list {
		if v.Name == name {
			return v
		}
	}
	//如果此玩家没有在排行榜内，则初始化他的信息
	my_rank := &RankData{
		Name:    name,
		Score:   0,
		Ranking: -1,
	}
	gbframe.Logger_Error.Println("Get Player Rank data is wrong,this name is not exist,name:", name)
	return my_rank
}
