package toold

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
)

//SysCmds  系统cmd
type SysCmds struct {
	Cmd []string
	byt string
	Out bytes.Buffer
}

var commands = map[string]string{
	"windows": "cmd /c start",
	"darwin":  "sh",
	"linux":   "xdg-open",
}

/*
OpenBrowserURL 打开浏览器
*/
func OpenBrowserURL(params []string) {
	run, _ := commands[runtime.GOOS]
	err := exec.Command(run, params...).Start()
	if err != nil {
		fmt.Printf("error:%v", err)
	}
}

/*
SysCmdAway 系统cmd
return SysCmd,error
*/
func SysCmdAway(cmdname string, params []string) (*SysCmds, error) {
	cmd := exec.Command(cmdname, params...)
	if runtime.GOOS == "windows" {
		SysCmdWinHide(cmd)
	}
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	whoami, err := cmd.Output()
	if err != nil {
		err2 := cmd.Start()
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

/*
SysCmd 系统cmd
return SysCmd,error
*/
func SysCmd(cmdname string, params []string) (*SysCmds, error) {
	cmd := exec.Command(cmdname, params...)
	if runtime.GOOS == "windows" {
		SysCmdWinHide(cmd)
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

//Cmd Cmd
func Cmd() {
	// cmd := exec.Command(cmdname, params...)
}

//CmdInfo CmdInfo
type CmdInfo struct {
	Cmdname     string
	Params      []string
	ParamsBytes [][]byte
}

// func SysCmdMore(terminal string, cmds []*CmdInfo) (*SysCmds, error) {
// 	cmd := exec.Command(terminal)
// 	if runtime.GOOS == "windows" {
// 		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
// 	}
// 	in := bytes.NewBuffer(nil)
// 	cmd.Stdin = in
// 	var out bytes.Buffer
// 	cmd.Stdout = &out
// 	var outErr bytes.Buffer
// 	cmd.Stderr = &outErr
// 	go func() {
// 		for _, info := range cmds {
// 			run := info.Cmdname
// 			for _, com := range info.Params {
// 				run = fmt.Sprintf("%v %v", run, com)
// 			}
// 			fmt.Printf("\nyun:%v", run)
// 			// in.Write([]byte(run))
// 			// in.WriteByte('\n')
// 			// in.WriteString(fmt.Sprintf("\n"))
// 			// ss := ConvertToString(run, "gbk", "utf8")
// 			in.WriteString(fmt.Sprintf("%v\n", run))
// 		}
// 	}()
// 	err := cmd.Start()
// 	if err != nil {
// 		return &SysCmds{
// 			Cmd: cmd.Args,
// 			byt: "nil",
// 			Out: out,
// 		}, fmt.Errorf("%v info:%v %v", err.Error(), outErr.String(), out.String())
// 	}
// 	err = cmd.Wait()
// 	if err != nil {
// 		return &SysCmds{
// 			Cmd: cmd.Args,
// 			byt: "nil",
// 			Out: out,
// 		}, fmt.Errorf("%v info:%v %v", err.Error(), outErr.String(), out.String())
// 	}
// 	errs := outErr.String()
// 	if len(errs) != 0 {
// 		return &SysCmds{
// 			Cmd: cmd.Args,
// 			byt: "nil",
// 			Out: out,
// 		}, fmt.Errorf("runing error:%v %v", outErr.String(), out.String())
// 	}
// 	FileWrite("ssss", out.String())
// 	// fmt.Printf("\n sdds:%v", out.String())
// 	return &SysCmds{
// 		Cmd: cmd.Args,
// 		byt: out.String(),
// 		Out: out,
// 	}, nil
// }
