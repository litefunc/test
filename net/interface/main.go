package main

import (
	"VodoPlay/logger"
	"errors"
	"fmt"
	"net"
)

func f() {

	ifaces, err := net.Interfaces()
	if err != nil {
		logger.Error(err)
		return
	}
	// handle err
	for i, v := range ifaces {
		logger.Info(i, v)
		addrs, err := v.Addrs()
		if err != nil {
			logger.Error(err)
			return
		}
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				logger.Warn("net")
				ip = v.IP
			case *net.IPAddr:
				logger.Warn("addr")
				ip = v.IP
			}
			ip.To4()
			logger.Debug(i, "ip:", net.ParseIP(ip.String()), "addr:", addr)
		}
	}

}

func main() {

	logger.Debug(getIP("wlp1s0", "eth0", "docker0"))

}

func getIP(names ...string) (net.IP, error) {

	ifaces, err := net.Interfaces()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if _, err := getInterfaceByNames(ifaces, names...); err != nil {
		return nil, err
	}

	return getIpByInterfacesAndNames(ifaces, names...)
}

func getIpByInterfacesAndNames(ifaces []net.Interface, names ...string) (net.IP, error) {

	for _, name := range names {
		for _, v := range ifaces {
			if v.Name == name {
				ip, err := getIPByInterface(v)
				if err != nil {
					break
				}
				return ip, nil
			}
		}
	}

	return nil, errors.New("ip not found")
}

func getInterfaceByNames(ifaces []net.Interface, names ...string) (net.Interface, error) {

	for _, name := range names {
		for _, v := range ifaces {
			if v.Name == name {
				return v, nil
			}
		}
	}
	err := errors.New("interface not found")
	logger.Error(err)
	return net.Interface{}, err
}

func getIPByInterface(iface net.Interface) (net.IP, error) {
	addrs, err := iface.Addrs()
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	ip, ok := findIPByAddrs(addrs)
	if !ok {
		err := fmt.Errorf(`can not find ip of %s`, iface.Name)
		logger.Error(err)
		return nil, err
	}
	return ip, nil
}

func findIPByAddrs(addrs []net.Addr) (net.IP, bool) {
	for _, addr := range addrs {
		ipn, ok := addr.(*net.IPNet)
		if ok {
			if ipn.IP.To4() != nil {
				return ipn.IP.To4(), true
			}
		}
	}

	return nil, false
}
