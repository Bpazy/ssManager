package main

import (
	"bytes"
	"github.com/Bpazy/ssManager/util"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type Port struct {
	Port  int    `json:"port"`
	Alias string `json:"alias"`
	Usage int    `json:"usage"`
}

func QueryPorts() []Port {
	rows, err := db.Queryx("select * from s_ports;")
	util.ShouldPanic(err)
	defer rows.Close()
	ports := make([]Port, 0)
	for rows.Next() {
		var p Port
		rows.StructScan(&p)
		ports = append(ports, p)
	}

	for _, p := range ports {
		p.Usage = getUsage(p.Port)
	}
	return ports
}

func getUsage(port int) int {
	if runtime.GOOS == "windows" {
		return 0
	}
	replace := strings.Replace(usage, "{}", strconv.Itoa(port), -1)
	cmd := exec.Command("bash", "-c", replace)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	util.ShouldPanic(err)
	return util.MustParseInt(strings.TrimSpace(out.String()))
}

func SavePort(p *Port) {
	_, err := db.NamedExec("insert into s_ports (port, alias) values (:port, :alias)", p)
	util.ShouldPanic(err)
}
