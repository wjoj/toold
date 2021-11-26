package toold

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
	//	"strings"
)

// Files 文件
type Files struct {
	File   *os.File
	Errors error
}

// IsFile 判断文件是否存在
func IsFile(file string) bool {
	_, err := os.Stat(file)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

//IsFileFolder 判断是否是文件夹
func IsFileFolder(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	return fi.IsDir()
}

//FileRemove FileRemove
func FileRemove(path string) error {
	return os.Remove(path)
}

//FileRemoveAll FileRemoveAll
func FileRemoveAll(path string) error {
	return os.RemoveAll(path)
}

func FileReplaceName(old, new string) error {
	return os.Rename(old, new)
}

// Midk 生成文件夹
func Midk(name string) error {
	err := os.Mkdir(name, os.ModeDir|os.ModePerm)
	return err
}

// MidkAllPath 指定文件生成所有文件夹
func MidkAllPath(path string) string {
	os.MkdirAll(path, os.ModeDir|os.ModePerm)
	//}
	return path
}

// MidkAllPath 指定文件生成所有文件夹
func MidkAllPathV2(path string) error {

	//}
	return os.MkdirAll(path, os.ModeDir|os.ModePerm)
}

// NewFile 生成新文件
func NewFile(fileName string) *Files {
	var dstFile *os.File
	var err error
	if IsFile(fileName) == true {
		dstFile, err = os.OpenFile(fileName, os.O_RDWR, os.ModePerm)
	} else {
		dstFile, err = os.Create(fileName)
	}
	return &Files{
		File:   dstFile,
		Errors: err,
	}
}

// NewFile 生成新文件
func NewFileTrunc(fileName string) *Files {
	var dstFile *os.File
	var err error
	if IsFile(fileName) == true {
		dstFile, err = os.OpenFile(fileName, os.O_RDWR|os.O_TRUNC, os.ModePerm)
	} else {
		dstFile, err = os.Create(fileName)
	}
	return &Files{
		File:   dstFile,
		Errors: err,
	}
}

//NewFileAppEnd 生成新文件 继续添加
func NewFileAppEnd(fileName string) *Files {
	var dstFile *os.File
	var err error
	if IsFile(fileName) == true {
		dstFile, err = os.OpenFile(fileName, os.O_RDWR|os.O_APPEND, os.ModePerm)
	} else {
		dstFile, err = os.Create(fileName)
	}
	return &Files{
		File:   dstFile,
		Errors: err,
	}
}

// WriteMap 像文件写对象
func (f *Files) WriteMap(maps map[string]string) (int, error) {
	f.File.Seek(0, 2)
	str, err := json.Marshal(maps)
	f.Errors = err
	return f.File.Write(str)
}

// WriteString 向文件写字符串
func (f *Files) WriteString(str string) (int, error) {
	f.File.Seek(0, 2)
	return f.File.WriteString(str)
}

// WriteInterface 向文件写interface
func (f *Files) WriteInterface(v interface{}) (int, error) {
	f.File.Seek(0, 2)
	str, _ := json.Marshal(v)
	return f.File.Write(str)
}

func (f *Files) Read() ([]byte, error) {
	return ioutil.ReadAll(f.File)
}

// Close 关闭对象
func (f *Files) Close() {
	f.File.Close()
}

// FileReadName 读文件
func FileReadName(pathname string) ([]byte, error) {
	b, errsd := ioutil.ReadFile(pathname)
	return b, errsd
}

// FileAfreshWrite 写重新文件
func FileAfreshWrite(fileName string, v interface{}) (err error) {
	logFile := NewFileTrunc(fileName)
	if logFile.Errors != nil {
		err = logFile.Errors
		return
	}

	switch v.(type) {
	case string:
		logFile.WriteString(v.(string))
		break
	case []byte:
		logFile.File.Write(v.([]byte))
		break
	default:
		logFile.WriteInterface(v)
		break
	}
	logFile.Close()
	return
}

// FileWrite 写文件
func FileWrite(fileName string, v interface{}) {
	logFile := NewFile(fileName)
	switch v.(type) {
	case string:
		logFile.WriteString(v.(string))
		break
	default:
		logFile.WriteInterface(v)
		break
	}
	logFile.Close()
}

// FileWrites 写文件
func FileWrites(fileName string, v ...interface{}) {
	logFile := NewFile(fileName)
	for _, arg := range v {
		btcopy, _ := json.Marshal(arg)
		logFile.File.Seek(0, 2)
		logFile.File.Write(btcopy)
	}
	logFile.Close()
}

/*
FileWritesEnter 文件写入
*/
func FileWritesEnter(fileName string, v ...interface{}) {
	by := []byte("\n")
	logFile := NewFile(fileName)
	for _, arg := range v {
		btcopy, _ := json.Marshal(arg)
		logFile.File.Seek(0, 2)
		logFile.File.Write(btcopy)
	}
	logFile.File.Seek(0, 2)
	logFile.File.Write(by)
	logFile.Close()
}

/*
FileWritesEnterTime 文件打印
*/
func FileWritesEnterTime(fileName string, v ...interface{}) error {
	logFile := NewFile(fileName)

	timestr := "\n| -- " + Times() + " -- |\n"
	bys := []byte(timestr)
	logFile.File.Seek(0, 2)
	_, err := logFile.File.Write(bys)
	if err != nil {
		return err
	}
	by := []byte("\n")
	for _, arg := range v {
		btcopy, _ := json.Marshal(arg)
		logFile.File.Seek(0, 2)
		logFile.File.Write(btcopy)
	}
	logFile.File.Seek(0, 2)
	logFile.File.Write(by)
	logFile.Close()
	return nil
}

/*
FileSysPath 文件路径
*/
func FileSysPath() string {
	strP, _ := exec.LookPath(os.Args[0])
	return strP
}

/*
GetFilePathIdentifier 获取文件分级符
*/
func GetFilePathIdentifier() string {
	path := FileSysPath()
	if strings.Contains(path, "\\") {
		return "\\"
	}
	return "/"
}

/*
IsFilePathIdentifier 获取文件分级符
*/
func IsFilePathIdentifier() bool {
	path := FileSysPath()
	if strings.Contains(path, "\\") {
		return false
	}
	return true
}

/*
FileSysRootDirectory 获取根目录
*/
func FileSysRootDirectory() string {
	inv := "/"
	if GetOSName() == OSNameTypeWindows {
		inv = "\\"
	}
	ps := strings.Split(FileSysPath(), inv)
	path := ""
	end := len(ps) - 1
	for _, s := range ps[0:end] {
		if len(path) == 0 {
			path = fmt.Sprintf("%v", s)
		} else {
			path += fmt.Sprintf("%v%v", inv, s)
		}
	}
	return path
}

var sym sync.Mutex

/*
FileWriteByte 写入文件
*/
func FileWriteByte(filename string, data []byte) error {
	sym.Lock()
	defer sym.Unlock()
	files := NewFileAppEnd(filename)
	if files.Errors != nil {
		return files.Errors
	}
	defer files.File.Close()
	_, err := files.File.Write(data)
	if err != nil {
		return err
	}
	return nil
}

/*
FileCopy src:源文件 dst:新文件
*/
func FileCopy(srcName, dstName string) (written int64, err error) {
	src, err := os.OpenFile(srcName, os.O_RDWR|os.O_APPEND, os.ModePerm)
	if err != nil {
		return
	}
	defer src.Close()
	var dst *os.File
	if IsFile(dstName) == true {
		dst, err = os.OpenFile(dstName, os.O_RDWR|os.O_APPEND, os.ModePerm)
	} else {
		dst, err = os.Create(dstName)
	}
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

/*
GetRootDirectory 获取当前跟目录
*/
func GetRootDirectory() string {
	path, _ := os.Getwd()
	return path
}
