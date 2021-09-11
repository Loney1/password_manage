package time

import (
	"fmt"
	"time"
)

// CurMSecond 当前毫秒数.
func CurMSecond() int64 {
	tm := time.Now()
	return tm.UnixNano() / 1e6
}

// 当前时间 是否 +8 与之前保持一致待定.
func CurTime() time.Time {
	return time.Now()
}

// CurSecond 当前秒数.
func CurSecond() int64 {
	tm := time.Now()
	return tm.Unix()
}

// TimeToDate 返回 yyyy-mm-dd
func TimeToDate() string {
	tm := time.Now()
	return fmt.Sprintf("%d-%02d-%02d", tm.Year(), tm.Month(), tm.Day())
}

// TimeToDBSufix 返回yyyymmdd
func TimeToDBSufix() string {
	tm := time.Now()
	return fmt.Sprintf("%d%02d%02d", tm.Year(), tm.Month(), tm.Day())
}

// TimeToString 返回yyyy-mm-dd hh:mm:ss
func TimeToString() string {
	tm := time.Now()
	return fmt.Sprintf("%02d-%02d-%02d %02d:%02d:%02d",
		tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second())
}

// StrToTime "yyyy-mm-dd hh:mm:ss"
func StrToTime(strTime string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04:05", strTime, time.Local)
}

// TimeToString 返回yyyy-mm-dd hh:mm:ss
func TimeFormat(tm time.Time, format string) string {
	if format == "" {
		format = "%02d-%02d-%02d %02d:%02d:%02d"
	}

	return fmt.Sprintf(format, tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second())
}

//
func TimeAddDate(days int) time.Time {
	cTime := time.Now()
	return cTime.AddDate(0, 0, days)
}

func Str2TimeStamp(str string) (error, int64) {
	formatTime, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return err, 0
	}
	return nil, formatTime.Unix()
}

//file time 时间 转换成time.time
func FileTime2Time(input int64) time.Time {
	t := time.Date(1601, 1, 1, 0, 0, 0, 0, time.Local)
	d := time.Duration(input)
	for i := 0; i < 100; i++ {
		t = t.Add(d)
	}
	return t
}
