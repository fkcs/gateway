package common

import (
	"github.com/fkcs/gateway/internal/infrastructure/logger"
	"fmt"
	"net"
)

//获取本机ip
func GetLocalIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logger.Logger().Errorf("failed to get local ip!%v", err)
		return "", err
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no local ip")
}

//获取本机Mac
func GetLocalMac() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		logger.Logger().Error("failed to get local mac! %v", err)
		return "", err
	}
	for _, inter := range interfaces {
		mac := inter.HardwareAddr
		if mac.String() != "" {
			return mac.String(), nil
		}
	}
	return "", fmt.Errorf("no local mac")
}
