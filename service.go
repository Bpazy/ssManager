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

type User struct {
	UserId       string `json:"userId" db:"USER_ID"`
	Username     string `json:"username" db:"USERNAME"`
	Nickname     string `json:"nickname" db:"NICKNAME"`
	EmailAddress string `json:"emailAddress" db:"EMAIL_ADDRESS"`
}

func SaveUser(u *User, password string) {
	_, err := db.NamedExec("insert into s_user (user_id, username, nickname, email_address) values (:userId, :username, :nickname, :emailAddress)", u)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("insert into s_user_password (user_id, password) values (?, ?)", u.UserId, password)
	if err != nil {
		panic(err)
	}
}

func FindUser(userId string) *User {
	row := db.QueryRowx("select user_id, username, nickname, email_address from s_user where userId = ?", userId)
	u := User{}
	row.StructScan(&u)
	return &u
}

func FindUserByAuth(username, password string) *User {
	row := db.QueryRowx("select A.user_id, A.username, A.nickname, A.email_address from s_user A "+
		"join s_user_password B on A.user_id = B.user_id where A.username = ? and B.password = ?", username, password)
	u := User{}
	err := row.StructScan(&u)
	if err != nil {
		return nil
	}
	return &u
}
