package toold

import (
	"math/rand"
	"time"
)

//RandType RandType
type RandType int

//
const (
	KC_RAND_KIND_NUM   RandType = 0 // 纯数字
	KC_RAND_KIND_LOWER          = 1 // 小写字母
	KC_RAND_KIND_UPPER          = 2 // 大写字母
	KC_RAND_KIND_ALL            = 3 // 数字、大小写字母
)

//Krand 随机字符串
func Krand(size int, kind RandType) []byte {
	ikind, kinds, result := int(kind), [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	isAll := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if isAll { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}

//RandString keys已存在的key
func RandString(lng int, keys map[string]bool) string {
	lls := string(Krand(lng, 3))
	if keys[lls] {
		lls = string(Krand(lng, 3))
	}
	return lls
}

//RandNumber keys已存在的key
func RandNumber(lng int, keys map[string]bool) string {
	lls := string(Krand(lng, 0))
	if keys[lls] {
		lls = string(Krand(lng, 0))
	}
	return lls
}

func RandNumberAndCapital(width int) string {
	numeric := []string{
		"1", "2", "3", "4", "5", "6", "7", "8", "9",
		"A", "B", "C", "D", "E", "F", "G",
		"H", "I", "J", "K", "L", "M", "N",
		"O", "P", "Q", "R", "S", "T", "U",
		"V", "W", "X", "Y", "Z",
	}
	r := len(numeric)
	sdsd := ""
	for i := 0; i < width; i++ {
		sdsd += numeric[rand.Intn(r)]
	}
	return sdsd
}
