package toold

import (
	"math/rand"
	"strings"
	"time"
)

//CodeInvitation CodeInvitation
type CodeInvitation struct {
	Base    string //进制的包含字符, string类型
	deciMal int64
	Pad     string //补位字符,若生成的code小于最小长度,则补位+随机字符, 补位字符不能在进制字符中
	Len     int    //code最小长度
}

//IDTransformationCode IdTransformationCode
func (c *CodeInvitation) IdTransformationCode(id int64) string {
	c.deciMal = int64(len(c.Base))
	mod := int64(0)
	res := ""
	for id != 0 {
		mod = id % c.deciMal
		id = id / c.deciMal
		res += string(c.Base[mod])
	}
	resLen := len(res)
	if resLen < c.Len {
		res += c.Pad
		for i := 0; i < c.Len-resLen-1; i++ {
			rand.Seed(time.Now().UnixNano())
			res += string(c.Base[rand.Intn(int(c.deciMal))])
		}
	}
	return res
}

//CodeTransformationID CodeTransformationID
func (c *CodeInvitation) CodeTransformationID(code string) int64 {
	res := int64(0)
	lenCode := len(code)
	c.deciMal = int64(len(c.Base))
	//var baseArr [] byte = []byte(c.base)
	baseArr := []byte(c.Base)     // 字符串进制转换为byte数组
	baseRev := make(map[byte]int) // 进制数据键值转换为map
	for k, v := range baseArr {
		baseRev[v] = k
	}

	// 查找补位字符的位置
	isPad := strings.Index(code, c.Pad)
	if isPad != -1 {
		lenCode = isPad
	}

	r := 0
	for i := 0; i < lenCode; i++ {
		// 补充字符直接跳过
		if string(code[i]) == c.Pad {
			continue
		}
		index := baseRev[code[i]]
		b := int64(1)
		for j := 0; j < r; j++ {
			b *= c.deciMal
		}
		// pow 类型为 float64 , 类型转换太麻烦, 所以自己循环实现pow的功能
		//res += float64(index) * math.Pow(float64(32), float64(2))
		res += int64(index) * b
		r++
	}
	return res
}
