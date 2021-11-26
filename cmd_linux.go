// +build linux

package toold

import (
	"bytes"
	"fmt"
	"os/exec"
	"syscall"
)

//SysCmdMore SysCmdMore
func SysCmdMore(terminal string, cmds []*CmdInfo) (*SysCmds, error) {
	cmd := exec.Command(terminal)
	in := bytes.NewBuffer(nil)
	cmd.Stdin = in
	var out bytes.Buffer
	cmd.Stdout = &out
	var outErr bytes.Buffer
	cmd.Stderr = &outErr
	go func() {
		// in.WriteString("chcp 65001")
		// in.WriteString("netsh interface ip show config\n")
		for _, info := range cmds {
			run := info.Cmdname
			for _, com := range info.Params {
				run += " " + com
			}

			in.WriteString(fmt.Sprintf("%v\n", run))
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

func SysCmdWinHide(cmd *exec.Cmd) {

}

func GetSysProcAttr(is bool) *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		// HideWindow: is,
		// CreationFlags: syscall.PROCESS_QUERY_INFORMATION | syscall.INFINITE,
	}
}
