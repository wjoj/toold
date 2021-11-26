package toold

import (
	"fmt"
	"strings"
)

//GetAudiosExt 获取音频扩展名
func GetAudiosExt(ext string) string {
	extL := strings.ToLower(ext)
	switch extL {
	case ".wav", ".mp3", ".cda", ".ape":
		break
	case "wav", "mp3", "cda", "ape":
		ext = fmt.Sprintf(".%v", ext)
		break
	default:
		return ""
	}
	return ext
}
