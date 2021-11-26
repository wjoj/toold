package toold

import (
	"fmt"
	"strings"
)

//GetVideoExt GetVideoExt
func GetVideoExt(ext string) string {
	extL := strings.ToLower(ext)
	switch extL {
	case ".rm", ".mp4", ".avi", ".3gp", ".flv", ".wmv":
		break
	case "rm", "mp4", "avi", "3gp", "flv", "wmv":
		ext = fmt.Sprintf(".%v", ext)
		break
	default:
		return ""
	}
	return ext
}

//GetZipExt GetVideoExt
func GetZipExt(ext string) string {
	extL := strings.ToLower(ext)
	switch extL {
	case ".zip":
		break
	case "zip":
		ext = fmt.Sprintf(".%v", ext)
		break
	default:
		return ""
	}
	return ext
}

//GetOExt GetOExt
func GetOExt(ext string) string {
	extL := strings.ToLower(ext)
	switch extL {
	case ".o", ".O":
		break
	case "o", "O":
		ext = fmt.Sprintf(".%v", ext)
		break
	default:
		return ""
	}
	return ext
}

//GetOExt GetOExt
func GetAPKExt(ext string) string {
	extL := strings.ToLower(ext)
	switch extL {
	case ".apk", ".APK":
		break
	case "apk", "APK":
		ext = fmt.Sprintf(".%v", ext)
		break
	default:
		return ""
	}
	return ext
}

//GetDocExt 文档
func GetDocExt(ext string) string {
	extL := strings.ToLower(ext)
	switch extL {
	case ".txt", ".doc", ".docx", ".xls", ".ppt", ".pptx", ".xlsx", `.wps`:
		break
	case "txt", "doc", "docx", "xls", "ppt", "pptx", "xlsx", "wps":
		ext = fmt.Sprintf(".%v", ext)
		break
	default:
		return ""
	}
	return ext
}

//GetCompressExt 压缩
func GetCompressExt(ext string) string {
	extL := strings.ToLower(ext)
	switch extL {
	case ".rar", ".zip", ".cab", ".iso", ".jar", ".ace", ".7z", ".tar", ".gz":
		break
	case "rar", "zip", "cab", "iso", "jar", "ace", "7z", "tar", "gz":
		ext = fmt.Sprintf(".%v", ext)
		break
	default:
		return ""
	}
	return ext
}

//GetOExt GetOExt
func GetFileExt(ext string) string {
	// extL := strings.ToLower(ext)
	// extL = strings.Replace(extL, ".", "", 1)
	// switch extL {
	// case ".rar", ".zip",
	// 	".doc", ".docx", ".pdf",
	// 	".jpg", "jpeg", ".bmp", ".png",
	// 	".rm", ".mp4", ".avi", ".3gp", ".flv", ".wmv":
	// 	break
	// default:
	// 	return ""
	// }
	return ext
}
