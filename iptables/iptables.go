package iptables

import (
	"github.com/Bpazy/ssManager/util"
	"runtime"
	"strconv"
	"strings"
)

const (
	iptablesUsage     = "iptables -L -nvx | grep spt:{} | awk '{print $2}'"
	iptablesInputAdd  = "iptables -A INPUT -p tcp --dport {}"
	iptablesInputDel  = "iptables -D INPUT -p tcp --dport {}"
	iptablesOutputAdd = "iptables -A OUTPUT -p tcp --sport {}"
	iptablesOutputDel = "iptables -D OUTPUT -p tcp --sport {}"
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

func GetUsage(port int) int64 {
	if runtime.GOOS == "windows" {
		return 0
	}
	i, ok := util.ShouldParseInt64(util.MustRunCommand(strings.Replace(iptablesUsage, "{}", strconv.Itoa(port), -1)))
	if !ok {
		return 0
	}
	return i
}
