package main

import (
	"github.com/Bpazy/ssManager/iptables"
	"github.com/Bpazy/ssManager/util"
)

type PortSorter []Port

func (p PortSorter) Len() int {
	return len(p)
}

func (p PortSorter) Less(i, j int) bool {
	return p[i].Usage < p[j].Usage
}

func (p PortSorter) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

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
		usage := iptables.GetUsage(p.Port)
		p.Usage = usage
	}
	return ports
}

func DeletePort(port int) {
	iptables.DeleteIptables(port)

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

func EditPort(p *Port) bool {
	_, err := db.NamedExec("update s_ports set alias = :alias where port = :port", p)
	if err != nil {
		return false
	}
	return true
}

func ResetPortUsage(port int) {
	iptables.DeleteIptables(port)
	iptables.SaveIptables(port)
}
