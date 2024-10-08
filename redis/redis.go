package redis

import (
	"app/message/pb"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
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

func GetAccessToken(appId string) (string, error) {
	ctx := context.Background()
	key := AccessTokenKey(appId)
	return client.Get(ctx, key).Result()
}

func SetAccessToken(appId string, accessToken string) error {
	ctx := context.Background()
	key := AccessTokenKey(appId)
	return client.Set(ctx, key, accessToken, AccessTokenTTL).Err()
}

func UpdateReport(appId string, report *pb.Report) error {
	ctx := context.Background()
	pip := client.Pipeline()
	keyScore := WeekScoreKey(appId)
	keyScoreMonth := MonthScoreKey(appId)
	keyWinningStreak := WinningStreakKey(appId)
	for _, report := range report.Info {
		pip.ZIncrBy(ctx, keyScore, float64(report.Score), report.OpenId)
		pip.ZIncrBy(ctx, keyScoreMonth, float64(report.Score), report.OpenId)
		if report.IsWin {
			pip.ZIncrBy(ctx, keyWinningStreak, 1, report.OpenId)
		} else {
			pip.ZRem(ctx, keyWinningStreak, report.OpenId)
		}
		keyUserDataCustom := UserDataCustomKey(appId, report.OpenId)
		for key, custom := range report.Custom {
			if customData, err := json.Marshal(custom); err == nil {
				pip.HSet(ctx, keyUserDataCustom, key, customData)
			}
		}
	}
	_, err := pip.Exec(ctx)
	return err
}

func GetScoreRank(appId string, topCount int32) ([]*pb.Audience, error) {
	ctx := context.Background()
	key := WeekScoreKey(appId)
	cmd := client.ZRangeWithScores(ctx, key, 0, int64(topCount))
	result, err := cmd.Result()
	if err != nil {
		fmt.Println("GetScoreRank err", err)
		return nil, err
	}
	var openIdList []string
	for _, info := range result {
		openIdList = append(openIdList, info.Member.(string))
	}
	audienceBasicList := GetAudienceBasicList(appId, openIdList)
	audienceInfoList := GetAudienceInfoList(appId, openIdList)
	if len(openIdList) != len(audienceBasicList) || len(openIdList) != len(audienceInfoList) {
		fmt.Println("GetScoreRank len err", len(openIdList), len(audienceBasicList), len(audienceInfoList))
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
func GetMonthScoreRank(appId string, topCount int32) ([]*pb.Audience, error) {
	ctx := context.Background()
	key := MonthScoreKey(appId)
	cmd := client.ZRangeWithScores(ctx, key, 0, int64(topCount))
	result, err := cmd.Result()
	if err != nil {
		fmt.Println("GetMonthScoreRank err", err)
		return nil, err
	}
	var openIdList []string
	for _, info := range result {
		openIdList = append(openIdList, info.Member.(string))
	}
	audienceBasicList := GetAudienceBasicList(appId, openIdList)
	audienceInfoList := GetAudienceInfoList(appId, openIdList)
	if len(openIdList) != len(audienceBasicList) || len(openIdList) != len(audienceInfoList) {
		fmt.Println("GetMonthScoreRank len err", len(openIdList), len(audienceBasicList), len(audienceInfoList))
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
	ttl, err := client.TTL(ctx, key).Result()
	if err != nil {
		return
	}
	if ttl > AudienceBasicTTLCheck {
		return
	}
	if data == nil {
		return
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}
	cmd := client.Set(ctx, key, string(jsonData), AudienceBasicTTL)
	fmt.Println("[SetAudienceBasic err]", cmd.Err())
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
		Custom: make(map[string]*pb.AudienceCustom),
	}
	pip := client.Pipeline()
	keyScore := WeekScoreKey(appId)
	keyScoreLast := LastWeekScoreKey(appId)
	keyScoreMonth := MonthScoreKey(appId)
	keyScoreLastMonth := LastMonthScoreKey(appId)
	keyWiningStreak := WinningStreakKey(appId)
	keyUserDataCustom := UserDataCustomKey(appId, openId)

	cmdScore := pip.ZScore(ctx, keyScore, openId)
	cmdRank := pip.ZRevRank(ctx, keyScore, openId)
	cmdLastRank := pip.ZRevRank(ctx, keyScoreLast, openId)
	cmdScoreMonth := pip.ZScore(ctx, keyScoreMonth, openId)
	cmdRankMonth := pip.ZRevRank(ctx, keyScoreMonth, openId)
	cmdLastRankMonth := pip.ZRevRank(ctx, keyScoreLastMonth, openId)
	cmdWiningStreak := pip.ZScore(ctx, keyWiningStreak, openId)
	cmdUserDataCustom := pip.HGetAll(ctx, keyUserDataCustom)
	_, _ = pip.Exec(ctx)
	if result, err := cmdScore.Result(); err == nil {
		ret.Score = int32(result)
	}
	if result, err := cmdRank.Result(); err == nil {
		ret.Rank = int32(result) + 1
	}
	if result, err := cmdLastRank.Result(); err == nil {
		ret.LastRank = int32(result) + 1
	}
	if result, err := cmdScoreMonth.Result(); err == nil {
		ret.MonthScore = int32(result)
	}
	if result, err := cmdRankMonth.Result(); err == nil {
		ret.MonthRank = int32(result) + 1
	}
	if result, err := cmdLastRankMonth.Result(); err == nil {
		ret.LastMonthRank = int32(result) + 1
	}
	if result, err := cmdWiningStreak.Result(); err == nil {
		ret.WinningStreak = int32(result)
	}
	if result, err := cmdUserDataCustom.Result(); err == nil {
		for k, v := range result {
			userDataCustom := &pb.AudienceCustom{}
			if err = json.Unmarshal([]byte(v), &userDataCustom); err == nil {
				ret.Custom[k] = userDataCustom
			}
		}
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
	_, _ = pip.Exec(ctx)
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
	keyScore := WeekScoreKey(appId)
	keyScoreLast := LastWeekScoreKey(appId)
	keyScoreMonth := MonthScoreKey(appId)
	keyScoreLastMonth := LastMonthScoreKey(appId)
	keyWiningStreak := WinningStreakKey(appId)

	var cmdScoreList []*redis.FloatCmd
	var cmdRankList []*redis.IntCmd
	var cmdRankLastList []*redis.IntCmd
	var cmdScoreMonthList []*redis.FloatCmd
	var cmdRankMonthList []*redis.IntCmd
	var cmdRankLastMonthList []*redis.IntCmd
	var cmdWiningStreakList []*redis.FloatCmd
	var cmdUserDataCustomList []*redis.StringStringMapCmd
	for _, openId := range openIdList {
		cmdScoreList = append(cmdScoreList, pip.ZScore(ctx, keyScore, openId))
		cmdRankList = append(cmdRankList, pip.ZRevRank(ctx, keyScore, openId))
		cmdRankLastList = append(cmdRankLastList, pip.ZRevRank(ctx, keyScoreLast, openId))
		cmdScoreMonthList = append(cmdScoreMonthList, pip.ZScore(ctx, keyScoreMonth, openId))
		cmdRankMonthList = append(cmdRankMonthList, pip.ZRevRank(ctx, keyScoreMonth, openId))
		cmdRankLastMonthList = append(cmdRankLastMonthList, pip.ZRevRank(ctx, keyScoreLastMonth, openId))
		cmdWiningStreakList = append(cmdWiningStreakList, pip.ZScore(ctx, keyWiningStreak, openId))
		cmdUserDataCustomList = append(cmdUserDataCustomList, pip.HGetAll(ctx, UserDataCustomKey(appId, openId)))
	}
	_, _ = pip.Exec(ctx)
	for i, openId := range openIdList {
		cmdScore := cmdScoreList[i]
		cmdRank := cmdRankList[i]
		cmdRankLast := cmdRankLastList[i]
		cmdScoreMonth := cmdScoreMonthList[i]
		cmdRankMonth := cmdRankMonthList[i]
		cmdRankLastMonth := cmdRankLastMonthList[i]
		cmdWiningStreak := cmdWiningStreakList[i]
		cmdUserDataCustom := cmdUserDataCustomList[i]
		info := &pb.AudienceInfo{
			OpenId: openId,
			Custom: make(map[string]*pb.AudienceCustom),
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
		if result, err := cmdScoreMonth.Result(); err == nil {
			info.MonthScore = int32(result)
		}
		if result, err := cmdRankMonth.Result(); err == nil {
			info.MonthRank = int32(result) + 1
		}
		if result, err := cmdRankLastMonth.Result(); err == nil {
			info.LastMonthRank = int32(result) + 1
		}
		if result, err := cmdWiningStreak.Result(); err == nil {
			info.WinningStreak = int32(result)
		}
		if result, err := cmdUserDataCustom.Result(); err == nil {
			for k, v := range result {
				userDataCustom := &pb.AudienceCustom{}
				if err = json.Unmarshal([]byte(v), &userDataCustom); err == nil {
					info.Custom[k] = userDataCustom
				}
			}
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
