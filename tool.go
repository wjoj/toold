package toold

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

/*
MD5 md5加密
*/
func MD5(con string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(con))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

//GenValidateCode GenValidateCode
func GenValidateCode(width int) string {
	numeric := [9]byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

//Code Code
func Code() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06v", rnd.Int31n(999999))
}

/*
IsDateFromString 日期格式
*/
func IsDateFromString(val string) bool {
	k, _ := regexp.MatchString(`^\d{4}\-(0?[1-9]|[1][012])\-(0?[1-9]|[12][0-9]|3[01])$`, val)
	return k
}

/*
IsMondayFromString 月日期格式
*/
func IsMondayFromString(val string) bool {
	k, _ := regexp.MatchString(`^(0?[1-9]|[1][012])\-(0?[1-9]|[12][0-9]|3[01])$`, val)
	return k
}

/*
IsMondayFromString 小时:分格式
*/
func IsHourMinFromString(val string) bool {
	k, _ := regexp.MatchString(`^(20|21|22|23|[0-1]\d):([0-5]\d)$`, val)
	return k
}

/*
IsTimeFromString 时间格式
*/
func IsTimeFromString(val string) bool {
	k, _ := regexp.MatchString(`^\d{4}[\-](0?[1-9]|1[012])[\-](0?[1-9]|[12][0-9]|3[01])(\s+(0?[0-9]|1[0-9]|2[0-3])\:(0?[0-9]|[1-5][0-9])\:(0?[0-9]|[1-5][0-9]))?$`, val)
	return k
}

/*
IsPhoneFromString 手机号验证
*/
func IsPhoneFromString(val string) bool {
	k, _ := regexp.MatchString("^1[3456789]\\d{9}$", val)
	return k
}

/*
IsNumberFromString 字符串是否是存数字
*/
func IsNumberFromString(val string) bool {
	k, _ := regexp.MatchString("^[0-9]*$", val)
	return k
}

/*
IsFloatFromString 字符串是否是存数字
*/
func IsFloatFromString(val string) bool {
	k, _ := regexp.MatchString("^[0-9]+.[0-9]+$", val)
	return k
}

/*
IsMinusNumberFromString 字符串是否是存数字
*/
func IsMinusNumberFromString(val string) bool {
	k, _ := regexp.MatchString("^-[0-9]*$", val)
	return k
}

/*
IsPercentageFromString 是否是百分数
*/
func IsPercentageFromString(val string) bool {
	k, _ := regexp.MatchString("^-[0-9]*%$", val)
	return k
}

/*
IsIPv4 IPv4
*/
func IsIPv4(val string) bool {
	k, _ := regexp.MatchString(`^(((\d{1,2})|(1\d{1,2})|(2[0-4]\d)|(25[0-5]))\.){3}((\d{1,2})|(1\d{1,2})|(2[0-4]\d)|(25[0-5]))$`, val)
	return k
}

//IsUserName IsUserName
func IsUserName(val string) bool {
	k, _ := regexp.MatchString("^[a-z0-9_-]{1,}$", val)
	return k
}

//IsPassword 大小写字母与数字组成，可以是全数字或全字母
func IsPassword(val string) bool {
	k, _ := regexp.MatchString("^([a-z]|[A-Z]|\\d){6,30}$", val)
	return k
}

//IsEnAndNumber IsEnAndNumber
func IsEnAndNumber(val string) bool {
	k, _ := regexp.MatchString("^[a-zA-Z0-9]+$", val)
	return k
}

/*
IsEmailFromString 邮箱验证
*/
func IsEmailFromString(val string) bool {
	k, _ := regexp.MatchString("^\\w+@[a-z0-9]+\\.[a-z]{2,4}$", val)
	return k
}

//IsCutApartFromString IsCutApartFromString
func IsCutApartFromString(vals, sep string, val string) bool {
	arr := strings.Split(vals, sep)
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

/*
DNSAnalysis 域名解析
*/
func DNSAnalysis(domain string) ([]string, error) {
	ns, err := net.LookupHost(domain)
	if err != nil {
		return []string{}, fmt.Errorf("DNS:%v", err)
	}
	ips := []string{}
	for _, n := range ns {
		ips = append(ips, n)
		fmt.Fprintf(os.Stdout, "--%s\n", n)
	}
	return ips, nil
}

/*
IsTimeRange 验证两个时间是否 start<end
*/
func IsTimeRange(start interface{}, end interface{}) bool {
	fotmat := func(times interface{}) int64 {
		switch times.(type) {
		case string:
			return TimeOftenStringConversionUninx(fmt.Sprintf("%v", times))
		case time.Time:
			return times.(time.Time).Local().Unix()
		default:
			break
		}
		return 0
	}
	return fotmat(start) < fotmat(end)
}
