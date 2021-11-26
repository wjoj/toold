package toold

import (
	"fmt"
	"net"
	"runtime"
	"strings"

	"github.com/thinkeridea/go-extend/exnet"
)

//OSNameType OSNameType
type OSNameType string

//
const (
	OSNameTypeLinux   OSNameType = "linux"
	OSNameTypeWindows            = "windows"
	OSNameTypeDarwin             = "darwin"
)

/*
GetMACAddressFromLocal 获取本机mac地址
*/
func GetMACAddressFromLocal() ([]string, error) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return []string{}, fmt.Errorf("fail to get net interfaces: %v", err)
	}
	macAddrs := []string{}
	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}
		macAddrs = append(macAddrs, macAddr)
	}
	return macAddrs, nil
}

/*
GetIPAddressFromLocal 获取IP地址
*/
func GetIPAddressFromLocal() ([]string, error) {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		return []string{}, fmt.Errorf("fail to get net interface addrs: %v", err)
	}
	ips := []string{}
	for _, address := range interfaceAddr {
		ipNet, isValidIPNet := address.(*net.IPNet)
		if isValidIPNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	return ips, nil
}

//
type IPInfo struct {
	IP       string
	IPBody   []byte
	Type     string
	MaskBody []byte
	Mask     string
}

type IPAdapterInfo struct {
	AdapterName string
	Index       int
	Infos       []IPInfo
}

/*
IsSameNetWorkSegment 是否是相同网段
*/
func IsSameNetWorkSegment(ip1, ip2 IPInfo) bool {
	i1, _ := exnet.IPString2Long(ip1.IP)
	iM1, _ := exnet.IPString2Long(ip1.Mask)
	i2, _ := exnet.IPString2Long(ip2.IP)
	iM2, _ := exnet.IPString2Long(ip2.Mask)
	if i1&iM1 == i2&iM2 {
		return true
	}
	return false
}

/*
GetIPAddressFromLocal 获取IPHemask地址
*/
func GetIPInfoFromLocal() []*IPAdapterInfo {
	ints, err := net.Interfaces()
	if err != nil {
		return []*IPAdapterInfo{}
	}
	list := []*IPAdapterInfo{}
	for _, info := range ints {
		ll, err := info.Addrs()
		if err != nil {
			continue
		}
		llCp := []IPInfo{}
		for _, addr := range ll {
			ipNet, isValidIPNet := addr.(*net.IPNet)
			if !isValidIPNet {
				continue
			}
			ipMask := net.IP(ipNet.Mask)
			if ipNet.IP.To4() != nil {
				llCp = append(llCp, IPInfo{
					IP:       ipNet.IP.String(),
					IPBody:   ipNet.IP,
					Type:     "ipv4",
					MaskBody: ipMask,
					Mask:     ipMask.String(),
				})
			} else {
				llCp = append(llCp, IPInfo{
					IP:       ipNet.IP.String(),
					IPBody:   ipNet.IP,
					Type:     "ipv6",
					MaskBody: ipMask,
					Mask:     ipNet.Mask.String(),
				})
			}
		}
		list = append(list, &IPAdapterInfo{
			Index:       info.Index,
			AdapterName: info.Name,
			Infos:       llCp,
		})
	}
	return list
}

/*
GetIPNetAddressFromLocal 获取IP地址
*/
func GetIPNetAddressFromLocal() ([]*net.IPNet, error) {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		return []*net.IPNet{}, fmt.Errorf("fail to get net interface addrs: %v", err)
	}
	ips := []*net.IPNet{}
	for _, address := range interfaceAddr {
		ipNet, isValidIPNet := address.(*net.IPNet)
		if isValidIPNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet)
			}
		}
	}
	return ips, nil
}

/*
GetIPLocalAddress 获取本地ip
*/
func GetIPLocalAddress() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx], nil
}

/*
GetOSName 获取系统
*/
func GetOSName() OSNameType {
	return OSNameType(runtime.GOOS)
}
