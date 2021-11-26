package toold

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
)

//IsSysCmd 默认true
var IsSysCmd = true

var muIsSysCmd sync.RWMutex

//SetIsSysCmd SetIsSysCmd
func SetIsSysCmd(is bool) {
	muIsSysCmd.Lock()
	defer muIsSysCmd.Unlock()
	IsSysCmd = is
}

//GetIsSysCmd GetIsSysCmd
func GetIsSysCmd() bool {
	muIsSysCmd.Lock()
	defer muIsSysCmd.Unlock()
	return IsSysCmd
}

//IPInfoConfig ip信息
type IPInfoConfig struct {
	NetName        string `form:"net_name"`
	NetNameByte    []byte
	IP             string `form:"ip"`
	DNS            string `form:"dns"`
	DNS2           string `form:"dns2"`
	DefaultGateway string `form:"default_gateway"`
	MaskIP         string `form:"mask_ip"`
	Status         int
	Types          int //1无线
}

//WlanSSID WlanSSID
type WlanSSID struct {
	Signal         int    //信号强度
	SSID           int    //ssid
	Status         int    //连接状态 2连接成功 1正在连接 -1连接失败 -2未知
	StatusMsg      string //
	IsPassword     int    //0无 1已设置
	Name           string //名称
	Authentication string //身份验证
	NetworkType    string //网络类型
	Encryption     string //加密协议
}

//Wlan Wlan
type Wlan struct {
	InterfaceName string
	WlanSSIDs     []WlanSSID
}

//WlanAfresh wlan刷新
type WlanAfresh struct {
	InterfaceName string `form:"interface"`
	Name          string `form:"name"`
}

//WlanStatus WlanStatus
type WlanStatus struct {
	Msg    string
	Status int
}

//WlanDisconnect wlan刷新
type WlanDisconnect struct {
	InterfaceName string `form:"interface"`
}

//WlanConnect wlan刷新
type WlanConnect struct {
	InterfaceName string `form:"interface"`
	Name          string `form:"name"`
	Password      string `form:"password"`
}

//WlanStatusMap WlanStatusMap
var WlanStatusMap = map[string]*WlanStatus{
	"disconnected": &WlanStatus{
		Msg:    "连接断开",
		Status: -1,
	},
	"authenticating": &WlanStatus{
		Msg:    "连接中",
		Status: 1,
	},
	"connected": &WlanStatus{
		Msg:    "连接成功",
		Status: 2,
	},
}

//ArrayWlanSSID ArrayWlanSSID
type ArrayWlanSSID []WlanSSID

func (a ArrayWlanSSID) Len() int {
	return len(a)
}

func (a ArrayWlanSSID) Less(i, j int) bool {
	return a[i].Signal > a[j].Signal
}
func (a ArrayWlanSSID) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

//WlanCmdBase WlanCmdBase
func WlanCmdBase(cmd *CmdInfo) (*SysCmds, error) {
	info, err := SysCmdMore("cmd", []*CmdInfo{
		&CmdInfo{
			Cmdname: "chcp",
			Params:  []string{"65001"}, //65001
		},
		cmd,
	})
	return info, err
}

//WlanSysCmdRun WlanSysCmdRun
func WlanSysCmdRun(cmd string, p []string) (*SysCmds, error) {
	var info *SysCmds
	var err error
	if GetIsSysCmd() {
		info, err = SysCmd(cmd, p)
	} else {
		info, err = WlanCmdBase(&CmdInfo{
			Cmdname: cmd,
			Params:  p,
		})
	}
	return info, err
}

//WlanGetIPInfo 获取ip信息
func WlanGetIPInfo() ([]*IPInfoConfig, error) {
	configs := []*IPInfoConfig{}
	if GetOSName() == OSNameTypeWindows {
		// var info *SysCmds
		// var err error
		// if GetIsSysCmd() {
		// 	info, err = SysCmd("netsh", []string{"interface", "ip", "show", "config"})
		// } else {
		// 	info, err = WlanCmdBase(&CmdInfo{
		// 		Cmdname: "netsh",
		// 		Params:  []string{"interface", "ip", "show", "config"},
		// 	})
		// }
		info, err := WlanSysCmdRun("netsh", []string{"interface", "ip", "show", "config"})
		//

		if err != nil {
			return nil, err
		}
		cmdOut := info.Out.String()
		maionCMD := regexp.MustCompile("Configuration for interface +([.\n]*)+")
		coms := maionCMD.Split(cmdOut, -1)
		for _, info := range coms {
			info = fmt.Sprintf("Configuration for interface %v", info)
			config := &IPInfoConfig{}
			regIP := regexp.MustCompile("IP Address: +([0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3})")
			com := regIP.FindStringSubmatch(info)
			if len(com) == 2 && len(com[1]) != 0 {
				config.IP = com[1]
			}
			if len(config.IP) == 0 {
				continue
			}
			regDNS := regexp.MustCompile("(?:(?:Statically Configured DNS Servers)|(?:DNS servers configured through DHCP)): +([0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3})\r\n(?: *([0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3})\r\n){0,1}")
			com = regDNS.FindStringSubmatch(info)
			if (len(com) == 2 || len(com) == 3) && len(com[1]) != 0 {
				config.DNS = com[1]
				if len(com) == 3 && len(com[2]) != 0 {
					config.DNS2 = com[2]
				}
			}

			regMaskIP := regexp.MustCompile("Subnet Prefix:.*mask ([0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3})")
			com = regMaskIP.FindStringSubmatch(info)
			if len(com) == 2 && len(com[1]) != 0 {
				config.MaskIP = com[1]
			}

			regDefaultGateway := regexp.MustCompile("Default Gateway: +([0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3})")
			com = regDefaultGateway.FindStringSubmatch(info)
			if len(com) == 2 && len(com[1]) != 0 {
				config.DefaultGateway = com[1]
			}

			regNetName := regexp.MustCompile("Configuration for interface *\"(.*)\"")
			com = regNetName.FindStringSubmatch(info)
			// lls := regNetName.FindSubmatch([]byte(info))
			if len(com) == 2 && len(com[1]) != 0 {
				config.NetName = com[1]
				// config.NetNameByte = lls[1]
			}
			if len(strings.Replace(config.NetName, "\r\n", "", -1)) == 0 {
				continue
			}
			status, err := WlanCheckNetwork(config.NetName)
			if err != nil {
				continue
			}
			config.Status = status
			configs = append(configs, config)
		}
	}
	return configs, nil
}

/*
WlanCheckNetwork 根据网络名称获取状态
*/
func WlanCheckNetwork(name string) (int, error) {
	if len(name) == 0 {
		return -1, fmt.Errorf("网络名称不能为空")
	}
	// ss := ConvertToString(name, "gbk", "utf-8")
	var info *SysCmds
	var err error
	if GetIsSysCmd() {
		info, err = SysCmd("netsh", []string{"interface", "show", "interface", fmt.Sprintf("name=%v", name)})
	} else {
		info, err = WlanCmdBase(&CmdInfo{
			Cmdname: "netsh",
			Params:  []string{"interface", "show", "interface", fmt.Sprintf("name=%v", name)},
		})
	}
	if err != nil {
		return -1, err
	}
	cmd := info.Out.String()
	regStatus := regexp.MustCompile("Connect state: +([a-zA-Z]+) *\r")
	coms := regStatus.FindStringSubmatch(cmd)
	if len(coms) == 2 {
		if coms[1] == "Connected" || coms[1] == "connected" {
			return 2, nil
		} else {
			return 0, nil
		}
	}
	return -1, fmt.Errorf("无获取")
}

/*
WlanCheckNetworkBytes 根据网络名称获取状态
*/
func WlanCheckNetworkBytes(name []byte) (int, error) {
	if len(name) == 0 {
		return -1, fmt.Errorf("网络名称不能为空")
	}
	// ss := ConvertToString(name, "gbk", "utf-8")
	var info *SysCmds
	var err error
	if GetIsSysCmd() {
		info, err = SysCmd("netsh", []string{"interface", "show", "interface", fmt.Sprintf("name=%v", name)})
	} else {
		lls := []byte("name=\"")
		lls = append(lls, name...)
		lls = append(lls, byte('"'))
		info, err = WlanCmdBase(&CmdInfo{
			Cmdname: "netsh",
			// Params:      []string{"interface", "show", "interface", fmt.Sprintf("name=\"%v\"", name)},
			ParamsBytes: [][]byte{[]byte("interface"), []byte("show"), []byte("interface"), lls},
		})
	}
	//
	if err != nil {
		return -1, err
	}
	cmd := info.Out.String()
	regStatus := regexp.MustCompile("Connect state: +([a-zA-Z]+) *\r")
	coms := regStatus.FindStringSubmatch(cmd)
	if len(coms) == 2 {
		if coms[1] == "Connected" || coms[1] == "connected" {
			return 2, nil
		} else {
			return 0, nil
		}
	}
	return -1, fmt.Errorf("无获取")
}

//WlanGetIPInfoFromInerfaceName WlanGetIPInfoFromInerfaceName
func WlanGetIPInfoFromInerfaceName(name string) (*IPInfoConfig, error) {
	// info, err := SysCmd("netsh", []string{"interface", "ip", "show", "config", fmt.Sprintf("name=%v", name)})
	// info, err := WlanCmdBase(&CmdInfo{
	// 	Cmdname: "netsh",
	// 	Params:  []string{"interface", "ip", "show", "config", fmt.Sprintf("name=\"%v\"", name)},
	// })
	var info *SysCmds
	var err error
	if GetIsSysCmd() {
		info, err = SysCmd("netsh", []string{"interface", "ip", "show", "config", fmt.Sprintf("name=%v", name)})
	} else {
		info, err = WlanSysCmdRun("netsh", []string{"interface", "ip", "show", "config", fmt.Sprintf("name=\"%v\"", name)})
	}
	//
	// if err != nil {
	// 	return nil, err
	// }
	// info, err := WlanSysCmdRun("netsh", []string{"interface", "ip", "show", "config", fmt.Sprintf("name=\"%v\"", name)})
	if err != nil {
		return nil, err
	}
	config := &IPInfoConfig{}
	cmdOut := info.Out.String()
	regIP := regexp.MustCompile("IP Address: +([0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3})")
	com := regIP.FindStringSubmatch(cmdOut)
	if len(com) == 2 && len(com[1]) != 0 {
		config.IP = com[1]
	}
	regDNS := regexp.MustCompile("(?:(?:Statically Configured DNS Servers)|(?:DNS servers configured through DHCP)): +([0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3})\r\n(?: *([0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3})\r\n){0,1}")
	com = regDNS.FindStringSubmatch(cmdOut)
	if (len(com) == 2 || len(com) == 3) && len(com[1]) != 0 {
		config.DNS = com[1]
		if len(com) == 3 && len(com[2]) != 0 {
			config.DNS2 = com[2]
		}
	}

	regMaskIP := regexp.MustCompile("Subnet Prefix:.*mask ([0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3})")
	com = regMaskIP.FindStringSubmatch(cmdOut)
	if len(com) == 2 && len(com[1]) != 0 {
		config.MaskIP = com[1]
	}

	regDefaultGateway := regexp.MustCompile("Default Gateway: +([0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3})")
	com = regDefaultGateway.FindStringSubmatch(cmdOut)
	if len(com) == 2 && len(com[1]) != 0 {
		config.DefaultGateway = com[1]
	}

	regNetName := regexp.MustCompile("Configuration for interface *\"(.*)\"")
	com = regNetName.FindStringSubmatch(cmdOut)
	if len(com) == 2 && len(com[1]) != 0 {
		config.NetName = com[1]
	}

	return config, nil
}

//WlanCheckIP WlanCheckIP
func WlanCheckIP(ip string) error {
	//
	var info *SysCmds
	// var err error
	if GetIsSysCmd() {
		info, _ = SysCmd("ping", []string{"-n", "1", "-w", "5000", ip})
	} else {
		info, _ = WlanCmdBase(&CmdInfo{
			Cmdname: "ping",
			Params:  []string{"-n", "1", "-w", "5000", ip},
		})
	}
	// info, _ := WlanSysCmdRun("ping", []string{"-n", "1", "-w", "5000", ip})
	// if err != nil {
	// 	return fmt.Errorf("该IP无法使用，请检查该IP是否正确或是否在该局域网内的频段")
	// }
	cmdOut := info.Out.String()
	regNetName := regexp.MustCompile("Reply from .*: (.*)\r\n")
	com := regNetName.FindStringSubmatch(cmdOut)
	if len(com) == 2 && com[1] == "Destination host unreachable." {
		return nil
	}
	if len(com) == 2 {
		regTTLS := regexp.MustCompile("TTL=([0-9]+).*")
		ttls := regTTLS.FindStringSubmatch(com[1])
		if len(ttls) == 2 {
			return fmt.Errorf("该IP以被占用")
		}
	}
	return nil
}

/*
WlanChangeIPInfo 更改ip信息
*/
func WlanChangeIPInfo(netName, ip, dns, dns2, maskip, defaultgateway string) (err error) {
	if !IsIPv4(ip) || !IsIPv4(dns) || (len(dns2) != 0 && !IsIPv4(dns2)) || !IsIPv4(maskip) || !IsIPv4(defaultgateway) {
		return fmt.Errorf("上传的网络地址格式错误")
	}
	if len(netName) == 0 {
		return fmt.Errorf("网络名称错误")
	}
	infoconfig, _ := WlanGetIPInfoFromInerfaceName(netName)
	if infoconfig.IP != ip {
		err = WlanCheckIP(ip)
		if err != nil {
			return
		}
	}

	if GetOSName() == OSNameTypeWindows {
		//
		// var info *SysCmds
		var err error
		if GetIsSysCmd() {
			_, err = SysCmd("netsh", []string{"interface", "ip", "set", "address", fmt.Sprintf("%v", netName), "static", ip, maskip, defaultgateway})
		} else {
			_, err = WlanCmdBase(&CmdInfo{
				Cmdname: "netsh",
				Params:  []string{"interface", "ip", "set", "address", fmt.Sprintf("%v", netName), "static", ip, maskip, defaultgateway},
			})
		}
		// _, err := WlanSysCmdRun("netsh", []string{"interface", "ip", "set", "address", fmt.Sprintf("%v", netName), "static", ip, maskip, defaultgateway})
		if err != nil {
			return err
		}
		if GetIsSysCmd() {
			_, err = SysCmd("netsh", []string{"interface", "ip", "delete", "dns", fmt.Sprintf("%v", netName), "all"})
		} else {
			_, err = WlanCmdBase(&CmdInfo{
				Cmdname: "netsh",
				Params:  []string{"interface", "ip", "delete", "dns", fmt.Sprintf("%v", netName), "all"},
			})
		}

		// _, err = WlanSysCmdRun("netsh", []string{"interface", "ip", "delete", "dns", fmt.Sprintf("%v", netName), "all"})
		if err != nil {
			return err
		}
		if GetIsSysCmd() {
			_, err = SysCmd("netsh", []string{"interface", "ip", "add", "dns", fmt.Sprintf("%v", netName), dns})
		} else {
			_, err = WlanCmdBase(&CmdInfo{
				Cmdname: "netsh",
				Params:  []string{"interface", "ip", "add", "dns", fmt.Sprintf("%v", netName), dns},
			})
		}

		// _, err = WlanSysCmdRun("netsh", []string{"interface", "ip", "add", "dns", fmt.Sprintf("%v", netName), dns})
		if err != nil {
			return err
		}
		if len(dns2) != 0 {
			if GetIsSysCmd() {
				_, err = SysCmd("netsh", []string{"interface", "ip", "add", "dns", fmt.Sprintf("%v", netName), dns2, "index=2"})
			} else {
				_, err = WlanCmdBase(&CmdInfo{
					Cmdname: "netsh",
					Params:  []string{"interface", "ip", "add", "dns", fmt.Sprintf("%v", netName), dns2, "index=2"},
				})
			}

			// _, err = WlanSysCmdRun("netsh", []string{"interface", "ip", "add", "dns", fmt.Sprintf("%v", netName), dns2, "index=2"})
			if err != nil {
				return err
			}
		}
	} else {
		err = fmt.Errorf("该系统暂时不支持该功能")
		return
	}
	return
}

/*
WlanSatusConnectNetWorks 获取网络状态
*/
func WlanSatusConnectNetWorks(wlaninfo WlanAfresh) (*WlanStatus, error) {
	var info *SysCmds
	var err error
	if GetIsSysCmd() {
		info, err = SysCmd("netsh", []string{"wlan", "show", "interfaces"})
	} else {
		info, err = WlanCmdBase(&CmdInfo{
			Cmdname: "netsh",
			Params:  []string{"wlan", "show", "interfaces"},
		})
	}

	if err != nil {
		return nil, err
	}
	cmdOut := info.Out.String()

	regName := regexp.MustCompile("Name *:([.\n]*)")
	coms := regName.Split(cmdOut, -1)

	wlan := &WlanStatus{
		Msg:    "未连接",
		Status: 0,
	}
	for _, info := range coms {
		info = fmt.Sprintf("Name :%v", info)
		regWlanName := regexp.MustCompile("Name *: *(.+)\r")
		name := regWlanName.FindStringSubmatch(info)
		if len(name) < 2 {
			continue
		}
		if name[1] != wlaninfo.InterfaceName {
			continue
		}

		regState := regexp.MustCompile("State *: *(.+)\r")
		state := regState.FindStringSubmatch(info) // disconnected authenticating
		if len(state) < 2 {
			continue
		}
		wlan = WlanStatusMap[state[1]]
		if wlan == nil {
			wlan = &WlanStatus{
				Msg:    state[1],
				Status: -2,
			}
		}
		regSSID := regexp.MustCompile("SSID *: *(.+)\r")
		ssid := regSSID.FindStringSubmatch(info)
		if len(ssid) < 2 {
			continue
		}
		if ssid[1] != wlaninfo.Name {
			wlan = &WlanStatus{
				Msg:    "未连接",
				Status: 0,
			}
			continue
		}
	}
	return wlan, nil
}

/*
WlanGetWlanNetWorks 获取无线网络
*/
func WlanGetWlanNetWorks(path string) ([]*Wlan, error) {
	wlans := []*Wlan{}
	if GetOSName() == OSNameTypeWindows {
		var info *SysCmds
		var err error
		if GetIsSysCmd() {
			info, err = SysCmd("netsh", []string{"wlan", "show", "networks", "mode=bssid"})
		} else {
			info, err = WlanCmdBase(&CmdInfo{
				Cmdname: "netsh",
				Params:  []string{"wlan", "show", "networks", "mode=bssid"},
			})
		}

		if err != nil {
			return []*Wlan{}, err
		}
		cmdOut := info.Out.String()
		regIP := regexp.MustCompile("Interface name : ([.\n]*)+")
		com := regIP.Split(cmdOut, -1)

		for _, wlan := range com {
			isNull := strings.NewReplacer(" ", "", "\n", "", "\t", "", "\r", "").Replace(wlan)
			if len(isNull) == 0 {
				continue
			}
			wlanObj := &Wlan{}
			interfaces := strings.Replace(fmt.Sprintf("Interface name : %v", wlan), "BSSID", "BCCID", -1)
			regWlanName := regexp.MustCompile("Interface name : *(.+)* \r")
			coms := regWlanName.FindStringSubmatch(interfaces)
			if len(coms) == 2 {
				wlanObj.InterfaceName = strings.Replace(coms[1], "\r", "", -1)
			} else {
				continue
			}
			regSSID := regexp.MustCompile("SSID *([.\n]*)+")
			ssids := regSSID.Split(interfaces, -1)
			wlanSSIDs := ArrayWlanSSID{}
			wlanSSIDLines := ArrayWlanSSID{}
			for i, ssid := range ssids {
				if i == 0 {
					continue
				}
				ssidObj := WlanSSID{}
				ssid = fmt.Sprintf("SSID %v", ssid)
				regSSIDName := regexp.MustCompile("SSID ([0-9]+) : *(.+)[\n]+")
				ssidNames := regSSIDName.FindStringSubmatch(ssid)
				if len(ssidNames) == 3 {
					name := strings.Replace(ssidNames[2], "\r", "", -1)
					if len(name) == 0 {
						continue
					}
					ssidObj.Name = name
					ssidObj.SSID = ConversionToInt(ssidNames[1])
					obj, _ := WlanSatusConnectNetWorks(WlanAfresh{
						InterfaceName: wlanObj.InterfaceName,
						Name:          name,
					})
					if obj != nil && obj.Status == 2 {
						ssidObj.Status = obj.Status
						ssidObj.StatusMsg = obj.Msg
					} else {
						ssidObj.Status = 0
						ssidObj.StatusMsg = "未连接"
					}

					if IsWlanProfile(name, path) {
						ssidObj.IsPassword = 1
					}
				}

				regNetType := regexp.MustCompile("Network type +: *(.+)\r")
				netTypes := regNetType.FindStringSubmatch(ssid)
				if len(netTypes) == 2 {
					ssidObj.NetworkType = strings.Replace(netTypes[1], "\r", "", -1)
				}

				regAuthentication := regexp.MustCompile("Authentication +: *(.+) *\r")
				authentications := regAuthentication.FindStringSubmatch(ssid)
				if len(authentications) == 2 {
					ssidObj.Authentication = strings.Replace(authentications[1], "\r", "", -1)
				}

				regEncryption := regexp.MustCompile("Encryption +: *(.+) \r")
				encryptions := regEncryption.FindStringSubmatch(ssid)
				if len(encryptions) == 2 {
					ssidObj.Encryption = strings.Replace(encryptions[1], "\r", "", -1)
				}

				regSignal := regexp.MustCompile("Signal +: *([0-9]+).*\r")
				signals := regSignal.FindStringSubmatch(ssid)
				if len(signals) == 2 {
					ssidObj.Signal = ConversionToInt(strings.Replace(signals[1], "\r", "", -1))
				}
				if ssidObj.Status == 2 {
					wlanSSIDLines = append(wlanSSIDLines, ssidObj)
				} else {
					wlanSSIDs = append(wlanSSIDs, ssidObj)
				}
			}
			sort.Sort(wlanSSIDs)
			for _, info := range wlanSSIDs {
				wlanSSIDLines = append(wlanSSIDLines, info)
			}
			wlanObj.WlanSSIDs = wlanSSIDLines
			wlans = append(wlans, wlanObj)
		}
	} else {
		return []*Wlan{}, fmt.Errorf("该系统暂时不支持该功能")
	}
	return wlans, nil
}

//WlanAfreshWlanNetWorks netsh wlan add profile filename="WLAN-@BHRGZN.CN.xml" 重新连接
func WlanAfreshWlanNetWorks(wlaninfo WlanAfresh, path string, funs func() error) error {
	// filename := interfaceName + "-" + name + ".xml"
	name := wlaninfo.Name
	interfaceName := wlaninfo.InterfaceName
	paths := MidkAllPath(path)
	if !IsWlanProfile(name, paths) {
		return fmt.Errorf("请重新输入密码")
	}
	// var info *SysCmds
	var err error
	if GetIsSysCmd() {
		_, err = SysCmd("netsh", []string{"wlan", "connect", fmt.Sprintf("name=%v", name), fmt.Sprintf("ssid=%v", name), fmt.Sprintf("interface=%v", interfaceName)})
	} else {
		_, err = WlanCmdBase(&CmdInfo{
			Cmdname: "netsh",
			Params:  []string{"wlan", "connect", fmt.Sprintf("name=\"%v\"", name), fmt.Sprintf("ssid=\"%v\"", name), fmt.Sprintf("interface=\"%v\"", interfaceName)},
		})
	}

	if err != nil {
		return err
	}
	return nil
}

/*
WlanDisconnectnNetWorks 断开连接
*/
func WlanDisconnectnNetWorks(form WlanDisconnect) error {
	var info *SysCmds
	var err error
	if GetIsSysCmd() {
		info, err = SysCmd("netsh", []string{"wlan", "disconnect", fmt.Sprintf("interface=%v", form.InterfaceName)})
	} else {
		info, err = WlanCmdBase(&CmdInfo{
			Cmdname: "netsh",
			Params:  []string{"wlan", "disconnect", fmt.Sprintf("interface=\"%v\"", form.InterfaceName)},
		})
	}

	if err != nil {
		return err
	}
	if strings.Contains(info.Out.String(), "There is no such wireless interface on the system.") {
		return fmt.Errorf("系统上没有这样的无线接口")
	}
	return nil
}

/*
WlanSetConnectNetWorks 设置网络
*/
func WlanSetConnectNetWorks(form WlanConnect) error {
	if len(form.Password) < 8 {
		return fmt.Errorf("密码长度至少8位")
	}
	// var info *SysCmds
	var err error
	if GetIsSysCmd() {
		_, err = SysCmd("netsh", []string{"wlan", "set", "profileparameter", fmt.Sprintf("name=%v", form.Name), fmt.Sprintf("keyMaterial=%v", form.Password)})
	} else {
		_, err = WlanCmdBase(&CmdInfo{
			Cmdname: "netsh",
			Params:  []string{"wlan", "set", "profileparameter", fmt.Sprintf("name=\"%v\"", form.Name), fmt.Sprintf("keyMaterial=\"%v\"", form.Password)},
		})
	}

	if err != nil {
		if strings.Contains(err.Error(), "not found on any interface") {
			name := form.Name
			path := MidkAllPath("./data/wlan")
			isfile, path, err := WlanProfile(name, path)
			if err != nil {
				return err
			}
			if isfile == 1 {
				// var info *SysCmds
				var err error
				if GetIsSysCmd() {
					_, err = SysCmd("netsh", []string{"wlan", "add", "profile", fmt.Sprintf("filename=%v", path)})
				} else {
					_, err = WlanCmdBase(&CmdInfo{
						Cmdname: "netsh",
						Params:  []string{"wlan", "add", "profile", fmt.Sprintf("filename=\"%v\"", path)},
					})
				}

				if err != nil {
					return err
				}
				if GetIsSysCmd() {
					_, err = SysCmd("netsh", []string{"wlan", "set", "profileparameter", fmt.Sprintf("name=%v", form.Name), fmt.Sprintf("SSIDname=%v", form.Name)})
				} else {
					_, err = WlanCmdBase(&CmdInfo{
						Cmdname: "netsh",
						Params:  []string{"wlan", "set", "profileparameter", fmt.Sprintf("name=\"%v\"", form.Name), fmt.Sprintf("SSIDname=\"%v\"", form.Name)},
					})
				}

				if err != nil {
					return err
				}
				if GetIsSysCmd() {
					_, err = SysCmd("netsh", []string{"wlan", "set", "profileparameter", fmt.Sprintf("name=%v", form.Name), fmt.Sprintf("keyMaterial=%v", form.Password)})
				} else {
					_, err = WlanCmdBase(&CmdInfo{
						Cmdname: "netsh",
						Params:  []string{"wlan", "set", "profileparameter", fmt.Sprintf("name=\"%v\"", form.Name), fmt.Sprintf("keyMaterial=\"%v\"", form.Password)},
					})
				}

				if err != nil {
					return err
				}
			}
		}
	}
	if GetIsSysCmd() {
		_, err = SysCmd("netsh", []string{"wlan", "connect", fmt.Sprintf("name=%v", form.Name), fmt.Sprintf("ssid=%v", form.Name), fmt.Sprintf("interface=%v", form.InterfaceName)})
	} else {
		_, err = WlanCmdBase(&CmdInfo{
			Cmdname: "netsh",
			Params:  []string{"wlan", "connect", fmt.Sprintf("name=\"%v\"", form.Name), fmt.Sprintf("ssid=%v", form.Name), fmt.Sprintf("interface=%v", form.InterfaceName)},
		})
	}
	if err != nil {
		return err
	}
	return nil
}
