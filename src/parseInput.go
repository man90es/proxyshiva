package main

import (
	"strings"
	"net"
	"encoding/binary"
	"strconv"
)

func ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

func int2ip(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}

func parseInput(in string) [][]string {
	out := make([][]string, 2)
	for i := range out {
		out[i] = make([]string, 0)
	}

	IPRanges := strings.Split(strings.Split(in, ":")[0], ",")
	portRanges := strings.Split(strings.Split(in, ":")[1], ",")

	for _, IPRange := range IPRanges {
		if strings.Contains(IPRange, "-") {
			r := strings.Split(IPRange, "-")

			for ip := ip2int(net.ParseIP(r[0])); ip < ip2int(net.ParseIP(r[1])); ip++ {
				out[0] = append(out[0], int2ip(ip).String())
			}
		} else {
			out[0] = append(out[0], IPRange)
		}
	}

	for _, portRange := range portRanges {
		if strings.Contains(portRange, "-") {
			r := strings.Split(portRange, "-")
			finish, _ := strconv.Atoi(r[1])

			for port, _ := strconv.Atoi(r[0]); port < finish; port++ {
				out[1] = append(out[1], strconv.Itoa(port))
			}
		} else {
			out[1] = append(out[1], portRange)
		}
	}

	return out
}