package redis

import (
	"fmt"
	"time"
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
	return fmt.Sprintf("Global_%s_Score_%d", appId, weekStart)
}

func LastWeekScoreKey(appId string) string {
	weekStart := getWeekStart(time.Now().Add(-time.Hour * 24 * 7))
	return fmt.Sprintf("Global_%s_Score_%d", appId, weekStart)
}

func getWeekStart(now time.Time) int64 {
	startOfWeek := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startOfWeek = startOfWeek.AddDate(0, 0, int(-startOfWeek.Weekday()))
	return startOfWeek.Unix()
}
