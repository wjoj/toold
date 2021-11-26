// +build darwin

package toold

import (
	"bytes"
	"fmt"
	"os/exec"
	"syscall"
)

//SysCmdMore SysCmdMore
func SysCmdMore(terminal string, cmds []*CmdInfo) (*SysCmds, error) {
	return nil, nil
}

func SysCmdWinHide(cmd *exec.Cmd) {

}

func GetSysProcAttr(is bool) *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		// HideWindow: is,
		// CreationFlags: syscall.PROCESS_QUERY_INFORMATION | syscall.INFINITE,
	}
}

/*
SysCmdHideDos 系统cmd
return SysCmd,error
*/
func SysCmdHideDos(cmdname string, params []string) (*SysCmds, error) {
	cmd := exec.Command(cmdname, params...)
	// if runtime.GOOS == "windows" {
	// 	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	// }
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
