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
	monthStart := getMonthStart(time.Now())
	ret := fmt.Sprintf("Global_%s_WinningStreak_%d", appId, monthStart)
	fmt.Println("WinningStreakKey", ret)
	return ret
}

func WeekScoreKey(appId string) string {
	weekStart := getWeekStart(time.Now())
	ret := fmt.Sprintf("Global_%s_Score_%d", appId, weekStart)
	fmt.Println("WeekScoreKey", ret)
	return ret
}

func LastWeekScoreKey(appId string) string {
	weekStart := getWeekStart(time.Now().Add(-time.Hour * 24 * 7))
	ret := fmt.Sprintf("Global_%s_Score_%d", appId, weekStart)
	fmt.Println("LastWeekScoreKey", ret)
	return ret
}

func MonthScoreKey(appId string) string {
	monthStart := getMonthStart(time.Now())
	ret := fmt.Sprintf("Global_%s_Month_Score_%d", appId, monthStart)
	fmt.Println("MonthScoreKey", ret)
	return ret
}

func LastMonthScoreKey(appId string) string {
	monthStart := getLastMonthStart(time.Now())
	ret := fmt.Sprintf("Global_%s_Month_Score_%d", appId, monthStart)
	fmt.Println("MonthScoreKey", ret)
	return ret
}

func getWeekStart(now time.Time) int64 {
	gmtTimeLoc := time.FixedZone("UTC+8", 0)
	now = now.In(gmtTimeLoc)
	startOfWeek := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startOfWeek = startOfWeek.AddDate(0, 0, int(-startOfWeek.Weekday()))
	return startOfWeek.Unix()
}

func getMonthStart(now time.Time) int64 {
	gmtTimeLoc := time.FixedZone("UTC+8", 0)
	now = now.In(gmtTimeLoc)
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	return startOfMonth.Unix()
}

func getLastMonthStart(now time.Time) int64 {
	gmtTimeLoc := time.FixedZone("UTC+8", 0)
	now = now.In(gmtTimeLoc)
	year := now.Year()
	month := now.Month()
	if month == 1 {
		year = year - 1
		month = 12
	} else {
		month = month - 1
	}
	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
	return startOfMonth.Unix()
}
