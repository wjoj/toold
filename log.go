package toold

import (
	"fmt"
	"strings"
)

/*
LogFileSysPath 读文件
*/
func LogFileSysPath(filename string) string {
	path := FileSysPath()
	if strings.Contains(path, "/T/go-build") || strings.Contains(path, "\\T\\go-build") {
		path = ""
	} else {
		var array []string
		var sign = 0
		if strings.Contains(path, "\\") == true {
			array = strings.Split(path, "\\")
			sign = 1
		} else {
			array = strings.Split(path, "/")
		}
		var paths string
		for i := 0; i < len(array)-1; i++ {
			if i == 0 {
				paths += array[i]
			} else {
				paths += "/" + array[i]
			}
		}
		paths += "/" + filename
		if sign == 1 {
			paths = strings.Replace(paths, "/", "\\", -1)
		}
		path = paths
	}
	return path
}

/*
LogsPrintf 打印
*/
func LogsPrintf(IsFileLog bool, v ...interface{}) error {
	if IsFileLog == false {
		fmt.Println(Times(), v)
	} else {
		path := LogFileSysPath("log")
		err := FileWritesEnterTime(path, v)
		if err != nil {
			fmt.Println(err, v)
			return err
		}
	}
	return nil
}

/*
LogPrintln 普通打印
*/
func LogPrintln(v ...interface{}) {
	fmt.Println(Times(), v)
}
