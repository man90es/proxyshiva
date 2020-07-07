package main

import (
	"strings"
	"net"
	"strconv"
	"utils"
)

func scheduleScheme(k KRequest, schemes []string, taskQueue chan KRequest) int {
	for _, scheme := range schemes {
		k.Scheme = scheme
		taskQueue <- k
	}

	return len(schemes)	
}

func schedulePort(k KRequest, in *string, schemes []string, taskQueue chan KRequest) int {
	portRanges := strings.Split(strings.Split(*in, ":")[1], ",")
	resultCount := 0

	for _, portRange := range portRanges {
		if strings.Contains(portRange, "-") {
			p := strings.Split(portRange, "-")
			fin, _ := strconv.Atoi(p[1])

			for port, _ := strconv.Atoi(p[0]); port < fin; port++ {
				k.Port = strconv.Itoa(port)

				resultCount += scheduleScheme(k, schemes, taskQueue)
			}
		} else {
			k.Port = portRange

			resultCount += scheduleScheme(k, schemes, taskQueue)
		}
	}

	return resultCount
}

func scheduleAddress(k KRequest, in *string, schemes []string, taskQueue chan KRequest) int {
	IPRanges := strings.Split(strings.Split(*in, ":")[0], ",")
	resultCount := 0

	for _, IPRange := range IPRanges {
		if strings.Contains(IPRange, "-") {
			r := strings.Split(IPRange, "-")
			passes := utils.IpToInt(net.ParseIP(r[1])) - utils.IpToInt(net.ParseIP(r[0]))

			for i := uint32(0); i < passes; i++ {
				for ip := utils.IpToInt(net.ParseIP(r[0])); ip < utils.IpToInt(net.ParseIP(r[1])); ip += passes {
					k.Address = utils.IntToIp(ip).String()

					resultCount += schedulePort(k, in, schemes, taskQueue)
				}
			}
		} else {
			k.Address = IPRange

			resultCount += schedulePort(k, in, schemes, taskQueue)
		}
	}

	return resultCount
}

func schedule(in string, schemes []string, taskQueue chan KRequest) int {
	k := KRequest{}
	return scheduleAddress(k, &in, schemes, taskQueue)
}