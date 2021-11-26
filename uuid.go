package toold

import (
	"fmt"
	"strings"
)

//GetLinuxDiskCurrentUUID GetLinuxDiskCurrentUUID
func GetLinuxDiskCurrentUUID(path string) (uuid string, err error) {
	log, err := SysCmd("df", []string{"-kh", path})
	if err != nil {
		return
	}
	con := log.Out.String()
	cons := strings.Split(con, "\n")
	if len(cons) < 1 {
		err = fmt.Errorf("获取根错误")
		return
	}
	cons = strings.Split(cons[1], " ")
	if len(cons) < 1 {
		err = fmt.Errorf("获取根错误为空")
		return
	}
	disk := cons[0]
	if len(disk) == 0 {
		err = fmt.Errorf("获取根错误为空")
		return
	}
	log, err = SysCmd("blkid", []string{"-s", "UUID", disk})
	if err != nil {
		err = fmt.Errorf("获取根错误为空")
		return
	}
	con = log.Out.String()
	if len(con) == 0 {
		err = fmt.Errorf("获取id错误")
		return
	}
	cons = strings.Split(con, "=")
	if len(cons) < 0 {
		err = fmt.Errorf("获取id错误为空")
		return
	}
	uuid = strings.ReplaceAll(strings.ReplaceAll(cons[1], "\"", ""), "\n", "")
	return
}

//GetWinDiskCurrentUUID GetWinDiskCurrentUUID
func GetWinDiskCurrentUUID() (uuid string, err error) {
	log, err := SysCmd("wmic", []string{"csproduct", "list", "full"})
	if err != nil {
		return
	}
	con := log.Out.String()
	cons := strings.Split(con, "\r\n")
	for _, c := range cons {
		uuid = strings.ReplaceAll(c, "UUID=", "")
		if c != uuid {
			uuid = strings.ReplaceAll(uuid, "\r", "")
			return uuid, nil
		}
	}

	return
}

//GetDiskCurrentUUID GetDiskCurrentUUID
func GetDiskCurrentUUID(pwd string) (uuid string, err error) {
	oss := GetOSName()
	if oss == OSNameTypeWindows {
		return GetWinDiskCurrentUUID()
	} else if oss == OSNameTypeLinux {
		return GetLinuxDiskCurrentUUID(pwd)
	}
	return
}

//GetWinDiskCurrentUUID GetWinDiskCurrentUUID
func GetWinDiskCurrentSeriaNumber() (uuid string, err error) {
	log, err := SysCmd("wmic", []string{"diskdrive", "get", "serialnumber"})
	if err != nil {
		return
	}
	con := log.Out.String()
	cons := strings.Split(con, "\r\n")
	for _, c := range cons {
		c = strings.ReplaceAll(c, "\t", "")
		c = strings.ReplaceAll(c, "\n", "")
		c = strings.ReplaceAll(c, "\r", "")
		c = strings.TrimSpace(c)
		if c == "SerialNumber" || len(c) == 0 {

		} else {
			uuid += c
		}
	}
	return
}

//GetWinDiskCurrentUUID GetWinDiskCurrentUUID
func GetWinDiskCurrentSeriaNumbers() (uuids []string, err error) {
	log, err := SysCmd("wmic", []string{"diskdrive", "get", "serialnumber"})
	if err != nil {
		return
	}
	con := log.Out.String()
	cons := strings.Split(con, "\r\n")
	for _, c := range cons {
		c = strings.ReplaceAll(c, "\t", "")
		c = strings.ReplaceAll(c, "\n", "")
		c = strings.ReplaceAll(c, "\r", "")
		c = strings.TrimSpace(c)
		if c == "SerialNumber" || len(c) == 0 {

		} else {
			uuids = append(uuids, c)
		}
	}
	return
}

//GetDiskCurrentSeriaNumber 获取磁盘序列号
func GetDiskCurrentSeriaNumber() (uuid string, err error) {
	oss := GetOSName()
	if oss == OSNameTypeWindows {
		ll, _ := GetWinDiskCurrentSeriaNumber()
		if len(ll) == 0 {
			return GetDiskCurrentUUID("")
		}
		return ll, nil
	} else if oss == OSNameTypeLinux {
		return GetLinuxDiskCurrentUUID("./")
	} else if oss == OSNameTypeDarwin {
		return GetLinuxDiskCurrentUUID("")
	}
	return
}

//GetDiskCurrentSeriaNumber 获取磁盘序列号
func GetDiskCurrentSeriaNumbers() (uuid []string, err error) {
	oss := GetOSName()
	if oss == OSNameTypeWindows {
		ll, _ := GetWinDiskCurrentSeriaNumbers()
		if len(ll) == 0 {
			kk, err := GetDiskCurrentUUID("")
			if err != nil {
				return []string{}, err
			}
			return []string{kk}, nil
		}
		return ll, nil
	} else if oss == OSNameTypeLinux {
		return []string{}, fmt.Errorf("获取uuid该系统暂无开发")
	}
	return
}
