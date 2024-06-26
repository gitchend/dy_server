package redis

import (
	"fmt"
	"time"
)

const (
	AudienceBasicTTL      = 24 * time.Hour
	AudienceBasicTTLCheck = 1 * time.Hour
)

func PublishKey(appId string, roomId string) string {
	return fmt.Sprintf("DataPush_%s_%s", appId, roomId)
}

func UserDataKey(appId string, openId string) string {
	return fmt.Sprintf("UserData_%s_%s", appId, openId)
}
func UserDataCustomKey(appId string, openId string) string {
	return fmt.Sprintf("UserData_%s_%s_Custom", appId, openId)
}

func WinningStreakKey(appId string) string {
	return fmt.Sprintf("Global_%s_WinningStreak", appId)
}

func ThisWeekScoreKey(appId string) string {
	weekStart := getWeekStart(time.Now())
	ret := fmt.Sprintf("Global_%s_Score_%d", appId, weekStart)
	fmt.Println("ThisWeekScoreKey", ret)
	return ret
}

func LastWeekScoreKey(appId string) string {
	weekStart := getWeekStart(time.Now().Add(-time.Hour * 24 * 7))
	ret := fmt.Sprintf("Global_%s_Score_%d", appId, weekStart)
	fmt.Println("LastWeekScoreKey", ret)
	return ret
}

func getWeekStart(now time.Time) int64 {
	gmtTimeLoc := time.FixedZone("UTC+8", 0)
	now = now.In(gmtTimeLoc)
	startOfWeek := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startOfWeek = startOfWeek.AddDate(0, 0, int(-startOfWeek.Weekday()))
	return startOfWeek.Unix()
}

func getWeekEnd(now time.Time) int64 {
	gmtTimeLoc := time.FixedZone("UTC+8", 0)
	now = now.In(gmtTimeLoc)
	endOfWeek := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfWeek = endOfWeek.AddDate(0, 0, int(7-endOfWeek.Weekday()))
	return endOfWeek.Unix()
}
