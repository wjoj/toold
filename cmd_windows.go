// +build windows

package toold

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"syscall"

	"golang.org/x/text/encoding/simplifiedchinese"
)

func Uint16ToBytes(n uint16) []byte {
	return []byte{
		byte(n),
		byte(n >> 8),
	}
}

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

func ConvertByte2String(byte []byte, charset Charset) string {

	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}

	return str
}

func SysCmdWinHide(cmd *exec.Cmd) {
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}
}

func SysCmdMore(terminal string, cmds []*CmdInfo) (*SysCmds, error) {
	cmd := exec.Command(terminal)
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}
	in := bytes.NewBuffer(nil)
	cmd.Stdin = in
	var out bytes.Buffer
	cmd.Stdout = &out
	var outErr bytes.Buffer
	cmd.Stderr = &outErr
	go func() {
		for _, info := range cmds {
			run := info.Cmdname
			for _, com := range info.Params {
				run += " " + syscall.EscapeArg(com)
			}
			in.WriteString(fmt.Sprintf("%v \n", run))
		}
	}()
	err := cmd.Start()
	if err != nil {
		return &SysCmds{
			Cmd: cmd.Args,
			byt: "nil",
			Out: out,
		}, fmt.Errorf("%v info:%v %v", err.Error(), outErr.String(), out.String())
	}
	err = cmd.Wait()
	defer func() {
		out.Reset()
		in.Reset()
	}()
	if err != nil {
		return &SysCmds{
			Cmd: cmd.Args,
			byt: "nil",
			Out: out,
		}, fmt.Errorf("%v info:%v %v", err.Error(), outErr.String(), out.String())
	}
	errs := outErr.String()
	if len(errs) != 0 {
		return &SysCmds{
			Cmd: cmd.Args,
			byt: "nil",
			Out: out,
		}, fmt.Errorf("runing error:%v %v", outErr.String(), out.String())
	}
	return &SysCmds{
		Cmd: cmd.Args,
		byt: out.String(),
		Out: out,
	}, nil
}

/*
SysCmdHideDos 系统cmd
return SysCmd,error
*/
func SysCmdHideDos(cmdname string, params []string) (*SysCmds, error) {
	cmd := exec.Command(cmdname, params...)
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	whoami, err := cmd.Output()
	if err != nil {
		err2 := cmd.Run()
		if err2 != nil {
			return &SysCmds{
				Cmd: cmd.Args,
				byt: "nil",
				Out: out,
			}, fmt.Errorf("%v info:%v %v", err2.Error(), stderr.String(), out.String())
		}
		if len(stderr.String()) == 0 {
			return &SysCmds{
				Cmd: cmd.Args,
				byt: "nil",
				Out: out,
			}, nil
		}
		return &SysCmds{
			Cmd: cmd.Args,
			byt: "nil",
			Out: out,
		}, fmt.Errorf("info:%v %v", stderr.String(), out.String())
	}
	return &SysCmds{
		Cmd: cmd.Args,
		byt: string(whoami),
		Out: out,
	}, nil
}

func GetSysProcAttr(is bool) *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		HideWindow: is,
		// CreationFlags: syscall.PROCESS_QUERY_INFORMATION | syscall.INFINITE,
	}
}
