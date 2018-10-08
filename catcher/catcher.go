package catcher

import (
	"github.com/Bpazy/ssManager/db"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/rickb777/date"
	"time"
)

const (
	saveTodayUsageInterval = 10 * time.Second
)

type Usage struct {
	Port      int
	Time      time.Time
	DownUsage int
	UpUsage   int
}

func (u Usage) GetDate() date.Date {
	return date.NewAt(u.Time)
}

type Catcher struct {
	DatePortUsage map[date.Date]map[int]*Usage
	Ports         *arraylist.List

	perUnitTimePortDownSize map[int]int
	perUnitTimePortUpSize   map[int]int
	portDownSpeed           map[int]float32
	portUpSpeed             map[int]float32
}

func New(deviceName string, ports []int) *Catcher {
	c := Catcher{}
	c.DatePortUsage = make(map[date.Date]map[int]*Usage)
	c.Ports = arraylist.New()
	c.perUnitTimePortDownSize = make(map[int]int)
	c.perUnitTimePortUpSize = make(map[int]int)
	c.portDownSpeed = make(map[int]float32)
	c.portUpSpeed = make(map[int]float32)

	c.AddPorts(ports)
	c.DatePortUsage = readTodayUsageFromDb()

	go c.detectNet(deviceName)
	go c.saveTodayUsage()

	return &c
}

func readTodayUsageFromDb() map[date.Date]map[int]*Usage {
	rows, err := db.Ins.Query("select port, `date`, downUsage, upUsage from s_usage where date(`date`) = date('now')")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var l []*Usage

	for rows.Next() {
		u := Usage{}
		rows.Scan(&u.Port, &u.Time, &u.DownUsage, &u.UpUsage)
		l = append(l, &u)
	}

	m := map[int]*Usage{}
	for _, u := range l {
		m[u.Port] = u
	}

	r := make(map[date.Date]map[int]*Usage)
	r[date.Today()] = m
	return r
}

func (c *Catcher) AddPorts(ports []int) {
	for _, port := range ports {
		c.Ports.Add(port)
	}
}

func (c Catcher) GetTodayUsage() map[int]*Usage {
	return c.DatePortUsage[date.Today()]
}

func (c *Catcher) detectNet(deviceName string) *pcap.Handle {
	handle, err := pcap.OpenLive(deviceName, 1024, true, 30*time.Second)
	if err != nil {
		panic(err)
	}

	go c.calculateNetSpeed()
	go c.capture(handle)

	return handle
}

func (c *Catcher) GetSpeed(port int) (float32, float32) {
	return c.portDownSpeed[port], c.portUpSpeed[port]
}

func (c *Catcher) capture(handle *pcap.Handle) {
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer != nil {
			tcp := tcpLayer.(*layers.TCP)

			today := date.Today()
			dataSize := len(packet.Data())
			if c.Ports.Contains(int(tcp.DstPort)) {
				c.perUnitTimePortUpSize[int(tcp.DstPort)] += dataSize

				_, ok := c.DatePortUsage[today]
				if !ok {
					c.DatePortUsage[today] = map[int]*Usage{}
				}
				_, ok = c.DatePortUsage[today][int(tcp.DstPort)]
				if !ok {
					c.DatePortUsage[today][int(tcp.DstPort)] = &Usage{Time: time.Now()}
				}
				c.DatePortUsage[today][int(tcp.DstPort)].DownUsage += dataSize
				continue
			}
			if c.Ports.Contains(int(tcp.SrcPort)) {
				c.perUnitTimePortDownSize[int(tcp.SrcPort)] += dataSize

				_, ok := c.DatePortUsage[today]
				if !ok {
					c.DatePortUsage[today] = map[int]*Usage{}
				}
				_, ok = c.DatePortUsage[today][int(tcp.SrcPort)]
				if !ok {
					c.DatePortUsage[today][int(tcp.SrcPort)] = &Usage{Time: time.Now()}
				}
				c.DatePortUsage[today][int(tcp.SrcPort)].UpUsage += dataSize
				continue
			}
		}
	}
}

func (c *Catcher) calculateNetSpeed() {
	for {
		duration := 1
		for port, size := range c.perUnitTimePortUpSize {
			c.portUpSpeed[port] = float32(size) / float32(duration)
			c.perUnitTimePortUpSize[port] = 0
		}
		for port, size := range c.perUnitTimePortDownSize {
			c.portDownSpeed[port] = float32(size) / float32(duration)
			c.perUnitTimePortDownSize[port] = 0
		}

		time.Sleep(time.Duration(duration) * time.Second)
	}
}

func (c *Catcher) saveTodayUsage() {
	for {
		todayUsage := c.GetTodayUsage()

		for p, u := range todayUsage {
			row := db.Ins.QueryRow("select port, `date`, downUsage, upUsage from s_usage")
			u2 := Usage{}
			err := row.Scan(&u2.Port, &u2.Time, &u2.DownUsage, &u2.UpUsage)
			if err != nil {
				_, err = db.Ins.Exec("insert into s_usage (port, `date`, downUsage, upUsage) VALUES (?,date(?),?,?)", p, time.Now(), u.DownUsage, u.UpUsage)
				if err != nil {
					panic(err)
				}
				continue
			}
			_, err = db.Ins.Exec("update s_usage set downUsage = ?, upUsage = ? where port = ? and `date` = date(?)",
				u.DownUsage, u.UpUsage, p, u.Time)
		}

		time.Sleep(saveTodayUsageInterval)
	}
}
