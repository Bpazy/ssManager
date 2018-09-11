package main

import (
	"bytes"
	"github.com/Bpazy/ssManager/util"
	"log"
	"os/exec"
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

type Port struct {
	Port  int    `json:"port"`
	Alias string `json:"alias"`
	Usage int64  `json:"usage"`
}

func QueryPorts() []Port {
	rows, err := db.Queryx("select * from s_ports order by port;")
	util.ShouldPanic(err)
	defer rows.Close()
	ports := make([]Port, 0)
	for rows.Next() {
		var p Port
		rows.StructScan(&p)
		ports = append(ports, p)
	}

	for i := range ports {
		p := &ports[i]
		if p.Alias == "" {
			p.Alias = "未配置"
		}
		usage := getUsage(p.Port)
		p.Usage = usage
	}
	return ports
}

func getUsage(port int) int64 {
	if runtime.GOOS == "windows" {
		return -1
	}
	i, ok := util.ShouldParseInt64(RunCommand(strings.Replace(iptablesUsage, "{}", strconv.Itoa(port), -1)))
	if !ok {
		return 0
	}
	return i
}

func DeletePort(port string) {
	RunCommand(strings.Replace(iptablesInputDel, "{}", port, -1))
	RunCommand(strings.Replace(iptablesOutputDel, "{}", port, -1))

	_, err := db.Exec("delete from s_ports where port = ?", port)
	util.ShouldPanic(err)
}

func SavePort(p *Port) bool {
	_, err := db.NamedExec("insert into s_ports (port, alias) values (:port, :alias)", p)
	if err != nil {
		return false
	}
	return true
}

func SaveIptables(port int) {
	if runtime.GOOS == "windows" {
		return
	}

	RunCommand(strings.Replace(iptablesInputAdd, "{}", strconv.Itoa(port), -1))
	RunCommand(strings.Replace(iptablesOutputAdd, "{}", strconv.Itoa(port), -1))
}

func RunCommand(c string) string {
	log.Println("command prepare: " + c)
	cmd := exec.Command("bash", "-c", c)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	util.ShouldPanic(err)

	result := strings.TrimSpace(out.String())
	log.Println("command result: " + result)
	return result
}
