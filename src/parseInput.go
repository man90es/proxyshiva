package main

import (
	"strings"
	"net"
	"strconv"
	"utils"
)

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

			for ip := utils.IpToInt(net.ParseIP(r[0])); ip < utils.IpToInt(net.ParseIP(r[1])); ip++ {
				out[0] = append(out[0], utils.IntToIp(ip).String())
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