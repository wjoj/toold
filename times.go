package toold

import (
	"strconv"
	"time"
)

// Times 获取服务当前时间
func Times() string {
	now := time.Now()
	year, mon, day := now.Local().Date()
	hour, min, sec := now.Local().Clock()
	zone, _ := now.Local().Zone()
	return strconv.Itoa(year) + "-" + strconv.Itoa(int(mon)) + "-" + strconv.Itoa(day) + " " + strconv.Itoa(hour) + ":" + strconv.Itoa(min) + ":" + strconv.Itoa(sec) + " " + zone
}

// TimeNow 获取当前时间
func TimeNow() string {
	year := time.Now()
	str := strconv.Itoa(year.Year()) + "-" + strconv.Itoa(int(year.Month())) +
		"-" + strconv.Itoa(year.Day()) + " " + strconv.Itoa(year.Hour()) + ":" + strconv.Itoa(year.Minute()) + ":" + strconv.Itoa(year.Second())
	return str
}

// TimeLocalNow 获取服务当前时间
func TimeLocalNow() string {
	times := time.Now()
	return times.Local().Format("2006-01-02 15:04:05")
}

//TimeLocalNowDay 获取服务当天
func TimeLocalNowDay() string {
	times := time.Now()
	return times.Local().Format("2006-01-02")
}

//TimeLocalNowFormat 获取服务当天
func TimeLocalNowFormat(format string) string {
	times := time.Now()
	return times.Local().Format(format)
}

// TimeLocalNowArea 包含时区
func TimeLocalNowArea() time.Time {
	times := time.Now()
	return times.Local()
}

// TimeNowFromUTC 获取当前时间
func TimeNowFromUTC() string {
	times := time.Now()
	return times.UTC().Format("2006-01-02 15:04:05")
}

/*
TimeNowFromUTCArea s
*/
func TimeNowFromUTCArea() time.Time {
	times := time.Now()
	return times.UTC()
}

/*
TimeNowStampFromUnixNano 获取时间戳 纳秒 us
*/
func TimeNowStampFromUnixNano() int64 {
	return time.Now().UnixNano()
}

/*
TimeNowStampFromUnixNanoMs 获取时间戳 微秒 ms
*/
func TimeNowStampFromUnixWeiMs() int64 {
	times := time.Now()
	return times.UnixNano() / 1e3
}

/*
TimeNowStampFromUnixNanoMs 获取时间戳 毫秒 ms
*/
func TimeNowStampFromUnixNanoMs() int64 {
	times := time.Now()
	return times.UnixNano() / 1e6
}

/*
TimeNowStampFromUnix 获取当前时间戳  s
*/
func TimeNowStampFromUnix() int64 {
	times := time.Now()
	return times.Unix()
}

//TimeTimeStringConversionObjFomat TimeTimeStringConversionObjFomat
func TimeTimeStringConversionObjFomat(fomat string, times string) time.Time {
	loc, _ := time.LoadLocation("Local")
	tm2, _ := time.ParseInLocation(fomat, times, loc)
	return tm2
}

//TimeTimeStringConversionObj TimeTimeStringConversionObj
func TimeTimeStringConversionObj(times string) time.Time {
	loc, _ := time.LoadLocation("Local")
	tm2, _ := time.ParseInLocation("2006-01-02 15:04:05", times, loc)
	return tm2
}

/*
TimeStringConversionUninx 转换时间戳
*/
func TimeStringConversionUninx(fomat string, times string) int64 {
	loc, _ := time.LoadLocation("Local")
	tm2, _ := time.ParseInLocation(fomat, times, loc)
	return tm2.Unix()
}

/*
TimeOftenStringConversionUninx 时间格式2006-01-02 15:04:05
*/
func TimeOftenStringConversionUninx(times string) int64 {
	return TimeStringConversionUninx("2006-01-02 15:04:05", times)
}

/*
TimeZoneStringConversionUninx 时间格式2006-01-02 15:04:05
2006-01-02T15:04:05.999999999Z07:00
*/
func TimeZoneStringConversionUninx(times string) int64 {
	return TimeStringConversionUninx("2006-01-02T15:04:05.999999999Z07:00", times)
}

//TimeZone0StringConversionUninx TimeZone0StringConversionUninx
func TimeZone0StringConversionUninx(times string) int64 {
	return TimeStringConversionUninx("2006-01-02T15:04:05Z", times)
}

/*
TimeConversionTime 传时间
*/
func TimeConversionTime(strtime string) time.Time {
	loc, _ := time.LoadLocation("Local")
	tm2, _ := time.ParseInLocation("2006-01-02 15:04:05", strtime, loc)
	return tm2
}

/*
TimeConversionTime 传时间
*/
func TimeConversionFromDate(strdate string) time.Time {
	loc, _ := time.LoadLocation("Local")
	tm2, _ := time.ParseInLocation("2006-01-02", strdate, loc)
	return tm2
}

/*
TimeConversionTimeFomat 传时间
*/
func TimeConversionTimeFomat(fm string, strtime string) time.Time {
	loc, _ := time.LoadLocation("Local")
	tm2, _ := time.ParseInLocation(fm, strtime, loc)
	return tm2
}

/*
TimeConversionTimeFomat 传时间
*/
func TimeConversionTimeFomat2(fm string, strtime string) (time.Time, error) {
	loc, _ := time.LoadLocation("Local")
	tm2, err := time.ParseInLocation(fm, strtime, loc)
	return tm2, err
}

/*
TimeConversionString 时间戳转时间字符串
*/
func TimeConversionString(times int64) string {
	return time.Unix(times, 0).Format("2006-01-02 15:04:05")
}

/*
TimeConversionString 时间戳转时间
*/
func TimeConversionObj(times int64) time.Time {
	return time.Unix(times, 0)
}

/*
TimeConversionStringFormat 时间戳转时间字符串2006-01-02 15:04:05
*/
func TimeConversionStringFormat(times int64, format string) string {
	return time.Unix(times, 0).Format(format)

}

//TimeLocalFromTime 时间转字符串
func TimeLocalFromTime(t time.Time) string {
	return t.Local().Format("2006-01-02 15:04:05")
}

//TimeLocalFromTimeFormat 时间转字符串
func TimeLocalFromTimeFormat(t time.Time, format string) string {
	return t.Local().Format(format)
}

//TimeSubDayNumber 获取时间差天数
func TimeSubDayNumber(tStart, tEnd string) float64 {
	tS := TimeConversionTimeFomat("2006-01-02", tStart[0:10])
	tE := TimeConversionTimeFomat("2006-01-02", tEnd[0:10])
	a := tE.Sub(tS)
	return a.Hours() / float64(24)
}

//TimeSubDays TimeSubDays
func TimeSubDays(tStart, tEnd string) []string {
	days := TimeSubDayNumber(tStart, tEnd)
	if days == 0 {
		return []string{tStart[0:10]}
	}
	ds := []string{}
	for i := 0; i <= int(days); i++ {
		tS := TimeConversionTimeFomat("2006-01-02", tStart[0:10])
		t := tS.Add(24 * time.Hour * time.Duration(i))
		ds = append(ds, TimeLocalFromTimeFormat(t, "2006-01-02"))
	}
	return ds
}

func TimeSubTimes(tStart, tEnd string) [][]string {
	days := TimeSubDays(tStart, tEnd)
	daysCp := [][]string{}
	if len(days) == 1 {
		cps := []string{days[0] + tStart[10:], days[0] + tEnd[10:]}
		return append(daysCp, cps)
	}

	for i, day := range days {
		if i == 0 {
			cps := []string{day + tStart[10:], day + " 23:59:59"}
			daysCp = append(daysCp, cps)
		} else if i == len(days)-1 {
			cps := []string{day + " 00:00:00", day + tEnd[10:]}
			daysCp = append(daysCp, cps)
		} else {
			cps := []string{day + " 00:00:00", day + " 23:59:59"}
			daysCp = append(daysCp, cps)
		}
	}
	return daysCp
}

//GetZeroTime 获取00点
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

//GetFirstDateOfMonth 获取这个d 的第一天
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

//GetLastDateOfMonth 获取这个d 的最后一天
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

func GetMonthDayFistAndLast() (fTime, eTime string) {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	fTime = firstOfMonth.Format("2006-01-02 15:04:05")
	eTime = lastOfMonth.Format("2006-01-02 15:04:05")
	return
}
