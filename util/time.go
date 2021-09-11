package util

import (
	"math"
	"strconv"
	"time"

	base_time "adp_backend/infra/time"
)

// 关于时间粒度的格式化单位
const (
	HourP  = "2006-01-02 15"
	DayP   = "2006-01-02"
	MonthP = "2006-01"
	format = "2006-01-02 15:04:05"
)

type WeekDate struct {
	WeekTh    int
	StartTime time.Time
	EndTime   time.Time
}

func String2Time(str string) (time.Time, error) {

	loc, _ := time.LoadLocation("Local")

	the_time, err := time.ParseInLocation(format, str, loc)
	if err != nil {
		return time.Now(), err
	}

	return the_time, nil
}

func Atoll(s string) int64 {
	i64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return i64
}

// 根据前端参数计算实际的时间区间
func GenerateInterval(timing string, now time.Time) (time.Time, time.Time, string) {
	startTm := now
	endTm := now
	var template string
	switch timing {
	case "day":
		template = HourP
		startTm = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
		endTm = startTm.AddDate(0, 0, 1)
	case "week":
		template = DayP
		num := int(time.Now().Weekday())
		if time.Now().Weekday() == 0 {
			num = 7
		}
		startTm = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, 1-num)
		endTm = startTm.AddDate(0, 0, 7)
	case "month":
		template = DayP
		startTm = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
		endTm = startTm.AddDate(0, 1, 0)
	case "year":
		template = MonthP
		startTm = time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, time.Local)
		endTm = startTm.AddDate(1, 0, 0)
	}

	return startTm, endTm, template
}

func InitializeTimeInterval(timeUnit, startTime, entTime, timeTemplate string, now time.Time) (time.Time, time.Time, string, error) {

	startTm, endTm := time.Time{}, time.Time{}

	switch timeTemplate {
	case "hour":
		timeTemplate = HourP
	case "day":
		timeTemplate = DayP
	case "week":
		timeTemplate = "week"
	case "month":
		timeTemplate = MonthP
	}

	if timeUnit == "" {
		startTm, err := base_time.StrToTime(startTime)
		if err != nil {
			return now, now, "", err
		}
		endTm, err := base_time.StrToTime(entTime)
		if err != nil {
			return now, now, "", err
		}

		return startTm, endTm, timeTemplate, nil
	}

	switch timeUnit {
	case "day":
		startTm = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
		endTm = startTm.AddDate(0, 0, 1)
	case "week":
		num := int(time.Now().Weekday())
		if time.Now().Weekday() == 0 {
			num = 7
		}
		startTm = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, 1-num)
		endTm = startTm.AddDate(0, 0, 7)
	case "month":
		startTm = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
		endTm = startTm.AddDate(0, 1, 0)
	case "year":
		startTm = time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, time.Local)
		endTm = startTm.AddDate(1, 0, 0)
	}

	return startTm, endTm, timeTemplate, nil
}

// 返回日期在当前年为第几周
func WeekByDate(t time.Time) int {
	yearDay := t.YearDay()
	yearFirstDay := t.AddDate(0, 0, -yearDay+1)
	firstDayInWeek := int(yearFirstDay.Weekday())

	//今年第一周有几天
	firstWeekDays := 1
	if firstDayInWeek != 0 {
		firstWeekDays = 7 - firstDayInWeek + 1
	}
	var week int
	if yearDay <= firstWeekDays {
		week = 1
	} else {
		week = (yearDay-firstWeekDays)/7 + 2
	}
	return week
}

// 根据时间区间，分割成周为单位
func GroupByWeekDate(startTime, endTime time.Time) []WeekDate {
	weekDate := make([]WeekDate, 0)
	diffDuration := endTime.Sub(startTime)
	days := int(math.Ceil(float64(diffDuration/(time.Hour*24)))) + 1

	currentWeekDate := WeekDate{}
	currentWeekDate.WeekTh = WeekByDate(endTime)
	currentWeekDate.EndTime = endTime
	currentWeekDay := int(endTime.Weekday())
	if currentWeekDay == 0 {
		currentWeekDay = 7
	}
	currentWeekDate.StartTime = endTime.AddDate(0, 0, -currentWeekDay+1)
	nextWeekEndTime := currentWeekDate.StartTime
	weekDate = append(weekDate, currentWeekDate)

	for i := 0; i < (days-currentWeekDay)/7; i++ {
		weekData := WeekDate{}
		weekData.EndTime = nextWeekEndTime
		weekData.StartTime = nextWeekEndTime.AddDate(0, 0, -7)
		weekData.WeekTh = WeekByDate(weekData.StartTime)
		nextWeekEndTime = weekData.StartTime
		weekDate = append(weekDate, weekData)
	}

	if lastDays := (days - currentWeekDay) % 7; lastDays > 0 {
		lastData := WeekDate{}
		lastData.EndTime = nextWeekEndTime
		lastData.StartTime = nextWeekEndTime.AddDate(0, 0, -lastDays)
		lastData.WeekTh = WeekByDate(lastData.StartTime)
		weekDate = append(weekDate, lastData)
	}

	return weekDate
}
