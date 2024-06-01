package redis

import (
	"app/message/pb"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
)

var client *redis.Client

func Init() {
	redisAddr := os.Getenv("REDIS_ADDRESS")
	redisUserName := os.Getenv("REDIS_USERNAME")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	client = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Username: redisUserName,
		Password: redisPassword,
		DB:       0, // use default DB
	})
	_, err := client.Ping(context.TODO()).Result()
	if err != nil {
		fmt.Printf("redisClient init error. err %s", err)
		panic(fmt.Sprintf("redis init failed. err %s\n", err))
	}
}

func UpdateReport(appId string, report *pb.Report) error {
	ctx := context.Background()
	pip := client.Pipeline()
	keyScore := ThisWeekScoreKey(appId)
	WinningStreak := WinningStreakKey(appId)
	for _, report := range report.Info {
		if report.Score > 0 {
			pip.ZIncrBy(ctx, keyScore, float64(report.Score), report.OpenId)
		}
		if report.IsWin {
			pip.ZIncrBy(ctx, WinningStreak, 1, report.OpenId)
		} else {
			pip.ZRem(ctx, WinningStreak, report.OpenId)
		}
	}
	_, err := pip.Exec(ctx)
	return err
}

func GetRank(appId string, topCount int32) ([]*pb.Audience, error) {
	ctx := context.Background()
	key := ThisWeekScoreKey(appId)
	cmd := client.ZRangeWithScores(ctx, key, 0, int64(topCount))
	result, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	var openIdList []string
	for _, info := range result {
		openIdList = append(openIdList, info.Member.(string))
	}
	audienceBasicList := GetAudienceBasicList(appId, openIdList)
	audienceInfoList := GetAudienceInfoList(appId, openIdList)
	if len(openIdList) != len(audienceBasicList) || len(openIdList) != len(audienceInfoList) {
		return nil, fmt.Errorf("GetRank len err: %d %d %d", len(openIdList), len(audienceBasicList), len(audienceInfoList))
	}
	var ret []*pb.Audience
	for i := range openIdList {
		ret = append(ret, &pb.Audience{
			AudienceBasic: audienceBasicList[i],
			AudienceInfo:  audienceInfoList[i],
		})
	}
	return ret, nil
}

func SetAudienceBasic(appId string, data *pb.AudienceBasic) {
	fmt.Println("[SetAudienceBasic]", data.OpenId, data)
	ctx := context.Background()
	key := UserDataKey(appId, data.OpenId)
	exist, err := client.Exists(ctx, key).Result()
	if err != nil {
		fmt.Println("[SetAudienceBasic err]", err)
		return
	}
	if exist == 1 {
		return
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}
	client.Set(ctx, key, string(jsonData), -1)
}

func GetAudience(appId string, openId string) *pb.Audience {
	return &pb.Audience{
		AudienceBasic: GetAudienceBasic(appId, openId),
		AudienceInfo:  GetAudienceInfo(appId, openId),
	}
}

func GetAudienceBasic(appId string, openId string) *pb.AudienceBasic {
	ctx := context.Background()
	key := UserDataKey(appId, openId)
	result, err := client.Get(ctx, key).Result()
	if err != nil {
		return nil
	}
	data := &pb.AudienceBasic{}
	err = json.Unmarshal([]byte(result), data)
	if err != nil {
		return nil
	}
	return data
}

func GetAudienceInfo(appId string, openId string) *pb.AudienceInfo {
	ctx := context.Background()
	ret := &pb.AudienceInfo{
		OpenId: openId,
	}
	pip := client.Pipeline()
	keyScore := ThisWeekScoreKey(appId)
	keyScoreLast := LastWeekScoreKey(appId)
	keyWiningStreak := WinningStreakKey(appId)

	cmdScore := pip.ZScore(ctx, keyScore, openId)
	cmdRank := pip.ZRank(ctx, keyScore, openId)
	cmdRankLast := pip.ZRank(ctx, keyScoreLast, openId)
	cmdWiningStreak := pip.ZScore(ctx, keyWiningStreak, openId)
	_, err := pip.Exec(ctx)
	if err != nil {
		return ret
	}
	if result, err := cmdScore.Result(); err == nil {
		ret.Score = int32(result)
	}
	if result, err := cmdRank.Result(); err == nil {
		ret.Rank = int32(result) + 1
	}
	if result, err := cmdRankLast.Result(); err == nil {
		ret.LastRank = int32(result) + 1
	}
	if result, err := cmdWiningStreak.Result(); err == nil {
		ret.WinningStreak = int32(result)
	}
	return ret
}

func GetAudienceBasicList(appId string, openIdList []string) (ret []*pb.AudienceBasic) {
	ctx := context.Background()
	pip := client.Pipeline()
	var cmdList []*redis.StringCmd
	for _, openId := range openIdList {
		key := UserDataKey(appId, openId)
		cmdList = append(cmdList, pip.Get(ctx, key))
	}
	_, err := pip.Exec(ctx)
	if err != nil {
		return nil
	}
	for _, cmd := range cmdList {
		data := &pb.AudienceBasic{}
		if result, err := cmd.Result(); err == nil {
			_ = json.Unmarshal([]byte(result), data)
		}
		ret = append(ret, data)
	}
	return
}

func GetAudienceInfoList(appId string, openIdList []string) (ret []*pb.AudienceInfo) {
	ctx := context.Background()
	pip := client.Pipeline()
	keyScore := ThisWeekScoreKey(appId)
	keyScoreLast := LastWeekScoreKey(appId)
	keyWiningStreak := WinningStreakKey(appId)

	var cmdScoreList []*redis.FloatCmd
	var cmdRankList []*redis.IntCmd
	var cmdRankLastList []*redis.IntCmd
	var cmdWiningStreakList []*redis.FloatCmd
	for _, openId := range openIdList {
		cmdScoreList = append(cmdScoreList, pip.ZScore(ctx, keyScore, openId))
		cmdRankList = append(cmdRankList, pip.ZRank(ctx, keyScore, openId))
		cmdRankLastList = append(cmdRankLastList, pip.ZRank(ctx, keyScoreLast, openId))
		cmdWiningStreakList = append(cmdWiningStreakList, pip.ZScore(ctx, keyWiningStreak, openId))
	}
	_, _ = pip.Exec(ctx)
	for i, openId := range openIdList {
		cmdScore := cmdScoreList[i]
		cmdRank := cmdRankList[i]
		cmdRankLast := cmdRankLastList[i]
		cmdWiningStreak := cmdWiningStreakList[i]
		info := &pb.AudienceInfo{
			OpenId: openId,
		}
		if result, err := cmdScore.Result(); err == nil {
			info.Score = int32(result)
		}
		if result, err := cmdRank.Result(); err == nil {
			info.Rank = int32(result) + 1
		}
		if result, err := cmdRankLast.Result(); err == nil {
			info.LastRank = int32(result) + 1
		}
		if result, err := cmdWiningStreak.Result(); err == nil {
			info.WinningStreak = int32(result)
		}
		ret = append(ret, info)
	}
	return ret
}

func Subscribe(appId string, roomId string, pushFunc func(string), pushCloseChan chan struct{}) {
	ps := client.Subscribe(context.Background(), PublishKey(appId, roomId))
	go func() {
		for {
			msg, err := ps.ReceiveMessage(context.Background())
			if err != nil {
				return
			}
			pushFunc(msg.Payload)
		}
	}()
	go func() {
		<-pushCloseChan
		_ = ps.Unsubscribe(context.Background(), PublishKey(appId, roomId))
	}()
}

func Publish(appId string, roomId string, data *pb.NotifyAudienceAction) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}
	client.Publish(context.Background(), PublishKey(appId, roomId), string(jsonData))
}