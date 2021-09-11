package time

import (
	"testing"
)

func TestTime(t *testing.T) {
	t.Log(CurMSecond())
	t.Log(CurSecond())
	t.Log(CurTime())
	t.Log(TimeToDate())
	t.Log(TimeToDBSufix())
	t.Log(TimeToString())
	t.Log(TimeAddDate(7))

	tm, _ := StrToTime("2017-11-07 21:05:38")
	t.Log(tm.Unix())
}
