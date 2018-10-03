package iptables

import (
	"github.com/Bpazy/ssManager/util"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

const (
	iptablesSptUsage  = "iptables -L -nvx | grep spt:{} | awk '{print $2}'"
	iptablesDptUsage  = "iptables -L -nvx | grep dpt:{} | awk '{print $2}'"
	iptablesInputAdd  = "iptables -A INPUT -p tcp --dport {}"
	iptablesInputDel  = "iptables -D INPUT -p tcp --dport {}"
	iptablesOutputAdd = "iptables -A OUTPUT -p tcp --sport {}"
	iptablesOutputDel = "iptables -D OUTPUT -p tcp --sport {}"

	iptablesUsage = "iptables -L -nvx"
	grepSpt       = "grep spt:{}"
)

var (
	sptRegexp = regexp.MustCompile("spt:(\\d+)")
)

func DeleteIptables(port int) {
	if runtime.GOOS == "windows" {
		return
	}

	util.RunCommand(strings.Replace(iptablesInputDel, "{}", strconv.Itoa(port), -1))
	util.RunCommand(strings.Replace(iptablesOutputDel, "{}", strconv.Itoa(port), -1))
}

func SaveIptables(port int) {
	if runtime.GOOS == "windows" {
		return
	}

	util.MustRunCommand(strings.Replace(iptablesInputAdd, "{}", strconv.Itoa(port), -1))
	util.MustRunCommand(strings.Replace(iptablesOutputAdd, "{}", strconv.Itoa(port), -1))
}

func GetSptUsage(port int) int64 {
	if runtime.GOOS == "windows" {
		return 0
	}
	i, ok := util.ShouldParseInt64(util.MustRunCommand(strings.Replace(iptablesSptUsage, "{}", strconv.Itoa(port), -1)))
	if !ok {
		return 0
	}
	return i
}

func GetDptUsage(port int) int64 {
	if runtime.GOOS == "windows" {
		return 0
	}
	i, ok := util.ShouldParseInt64(util.MustRunCommand(strings.Replace(iptablesDptUsage, "{}", strconv.Itoa(port), -1)))
	if !ok {
		return 0
	}
	return i
}

func GetSptUsageMap(ports []int) (m map[int]int64) {
	if runtime.GOOS == "windows" || len(ports) == 0 {
		return
	}

	grepSpts := ""
	for _, port := range ports {
		grepSpts += "|" + strings.Replace(grepSpt, "{}", strconv.Itoa(port), -1)
	}

	lines := strings.Split(util.MustRunCommand(iptablesUsage+"|"+grepSpts[1:]), "\r")
	for _, line := range lines {
		port := sptRegexp.FindString(line)
		m[util.MustParseInt(port)] = util.MustParseInt64(strings.Fields(line)[1])
	}
	return
}
