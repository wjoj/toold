// +build windows

package toold

/*
#cgo windows LDFLAGS: -L${SRCDIR} -lIPHLPAPI -lws2_32

#include <winsock2.h>
#include <iphlpapi.h>
#include <stdio.h>
#include <stdlib.h>
#include <locale.h>
#include <psdk_inc/_ip_types.h>

#include <winsock.h>
#include <WS2tcpip.h>

#pragma comment(lib, "IPHLPAPI.lib")
#pragma comment(lib, "ws2_32.lib")

PIP_ADAPTER_ADDRESSES initAdapterList(ULONG outBufLen){
	return (IP_ADAPTER_ADDRESSES *)malloc(outBufLen);
}



*/
import "C"
import (
	"unsafe"
)

func GetAdapterAddressList() {
	outBufLen := C.ULONG(15000)
	pAddresses := C.initAdapterList(C.ulong(outBufLen))
	if pAddresses == nil {
		return
	}
	family := C.AF_UNSPEC
	flags := C.GAA_FLAG_INCLUDE_PREFIX | C.GAA_FLAG_INCLUDE_GATEWAYS
	dwRetVal := C.GetAdaptersAddresses(C.ulong(family), C.ulong(flags), nil, pAddresses, &outBufLen)
	if dwRetVal != 0 {
		C.free(unsafe.Pointer(pAddresses))
		pAddresses = nil
		return
	}
	pCurrAddresses := pAddresses
	for {
		if pCurrAddresses == nil {
			break
		}
		pAnycast := pCurrAddresses.FirstAnycastAddress
		for {
			if pAnycast == nil {
				break
			}
			if pAnycast.Address.lpSockaddr != nil {
				continue
			}
			if C.AF_INET == pAnycast.Address.lpSockaddr.sa_family {

			} else if C.AF_INET6 == pAnycast.Address.lpSockaddr.sa_family {

			}
		}

		pCurrAddresses = pCurrAddresses.Next
	}
}
