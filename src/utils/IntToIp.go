package utils

import (
	"net"
	"encoding/binary"
)

func IntToIp(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}