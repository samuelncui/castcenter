package castcenter

import (
	"fmt"
	"net"
	"strings"
)

var (
	// ErrLocalIPNotFound local ip not found
	ErrLocalIPNotFound = fmt.Errorf("castcenter: local ip not found")

	netinfo *NetworkInfo
)

func init() {
	var err error
	netinfo, err = getNetworkInfo()
	if err != nil {
		panic(err)
	}
}

// NetworkInfo .
type NetworkInfo struct {
	HardwareName string
	IP           string
	SubNet       string
}

func getNetworkInfo() (*NetworkInfo, error) {
	info := NetworkInfo{}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				info.IP = ipnet.IP.String()
				break
			}
		}
	}

	if info.IP == "" {
		return nil, ErrLocalIPNotFound
	}

	interfaces, _ := net.Interfaces()
	for _, interf := range interfaces {
		if addrs, err := interf.Addrs(); err == nil {
			for _, addr := range addrs {
				if strings.Contains(addr.String(), info.IP) {
					info.HardwareName = interf.Name

					_, net, err := net.ParseCIDR(addr.String())
					if err != nil {
						return nil, err
					}

					info.SubNet = net.String()
					return &info, nil
				}
			}
		}
	}

	return nil, ErrLocalIPNotFound
}
