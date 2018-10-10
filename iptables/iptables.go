package iptables

import (
	"github.com/Bpazy/ssManager/util"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

const (
	iptablesInputAdd  = "iptables -A INPUT -p tcp --dport {}"
	iptablesInputDel  = "iptables -D INPUT -p tcp --dport {}"
	iptablesOutputAdd = "iptables -A OUTPUT -p tcp --sport {}"
	iptablesOutputDel = "iptables -D OUTPUT -p tcp --sport {}"
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
