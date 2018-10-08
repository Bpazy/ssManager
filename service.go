package main

import (
	"github.com/Bpazy/ssManager/db"
	"github.com/Bpazy/ssManager/iptables"
)

type PortStructSorter []PortStruct

func (p PortStructSorter) Len() int {
	return len(p)
}

func (p PortStructSorter) Less(i, j int) bool {
	return p[i].DownstreamUsage < p[j].DownstreamUsage
}

func (p PortStructSorter) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type PortStruct struct {
	Port            int    `json:"port"`
	Alias           string `json:"alias"`
	UpstreamUsage   int64  `json:"upstreamUsage"`
	DownstreamUsage int64  `json:"downstreamUsage"`
}

func QueryPorts() []int {
	portStructs := QueryPortStructs()
	ports := make([]int, len(portStructs))

	for _, portStruct := range portStructs {
		ports = append(ports, portStruct.Port)
	}
	return ports
}

func QueryPortStructs() []PortStruct {
	rows, err := db.Ins.Query("select port, alias from s_ports order by port")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	ports := make([]PortStruct, 0)
	for rows.Next() {
		var p PortStruct
		rows.Scan(&p.Port, &p.Alias)
		ports = append(ports, p)
	}

	return ports
}

func QueryPortStructsWithUsage() []PortStruct {
	ports := QueryPortStructs()

	for i := range ports {
		p := &ports[i]
		if p.Alias == "" {
			p.Alias = "未配置"
		}
		p.UpstreamUsage = iptables.GetDptUsage(p.Port)
		p.DownstreamUsage = iptables.GetSptUsage(p.Port)
	}
	return ports
}

func DeletePort(port int) {
	iptables.DeleteIptables(port)

	_, err := db.Ins.Exec("delete from s_ports where port = ?", port)
	if err != nil {
		panic(err)
	}
}

func SavePort(p *PortStruct) bool {
	_, err := db.Ins.Exec("insert into s_ports (port, alias) values (?, ?)", p.Port, p.Alias)
	if err != nil {
		return false
	}
	return true
}

func EditPort(p *PortStruct) bool {
	_, err := db.Ins.Exec("update s_ports set alias = ? where port = ?", p.Alias, p.Port)
	if err != nil {
		return false
	}
	return true
}

func ResetPortUsage(port int) {
	iptables.DeleteIptables(port)
	iptables.SaveIptables(port)
}

type User struct {
	UserId       string `json:"userId" db:"user_id"`
	Username     string `json:"username" db:"username"`
	Nickname     string `json:"nickname" db:"nickname"`
	EmailAddress string `json:"emailAddress" db:"email_address"`
}

func SaveUser(u *User, password string) {
	_, err := db.Ins.Exec("insert into s_user (user_id, username, nickname, email_address) values (?, ?, ?, ?)",
		u.UserId, u.Username, u.Nickname, u.EmailAddress)
	if err != nil {
		panic(err)
	}
	_, err = db.Ins.Exec("insert into s_user_password (user_id, password) values (?, ?)", u.UserId, password)
	if err != nil {
		panic(err)
	}
}

func FindUser(userId string) *User {
	row := db.Ins.QueryRow("select user_id, username, nickname, email_address from s_user where userId = ?", userId)
	u := User{}
	row.Scan(&u.UserId, &u.Username, &u.Nickname, &u.EmailAddress)
	return &u
}

func FindUserByAuth(username, password string) *User {
	row := db.Ins.QueryRow("select A.user_id, A.username, A.nickname, A.email_address from s_user A "+
		"join s_user_password B on A.user_id = B.user_id where A.username = ? and B.password = ?", username, password)
	u := User{}
	err := row.Scan(&u.UserId, &u.Username, &u.Nickname, &u.EmailAddress)
	if err != nil {
		return nil
	}
	return &u
}

func AddPortPassword(port, password string) {
	sc.AddPortPassword(port, password)
}
