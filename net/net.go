package net

import (
	"fmt"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/rickb777/date"
	"time"
)

var (
	ports = arraylist.New()

	perUnitTimePortDownSize = make(map[int]int)
	perUnitTimePortUpSize   = make(map[int]int)
	portDownSpeed           = make(map[int]float32)
	portUpSpeed             = make(map[int]float32)

	datePortUsage = make(map[date.Date]map[int]Usage) // Date => map[port]Usage
)

type Usage struct {
	Down int
	Up   int
}

func GetTodayUsage() map[int]Usage {

}

func AddPorts(_ports []int) {
	for _, port := range _ports {
		ports.Add(port)
	}
}

func DetectNet(deviceName string) *pcap.Handle {
	handle, err := pcap.OpenLive(deviceName, 1024, true, 30*time.Second)
	if err != nil {
		panic(err)
	}

	go calculateNetSpeed()
	go capture(handle)

	return handle
}

func GetSpeed(port int) (float32, float32) {
	return portDownSpeed[port], portUpSpeed[port]
}

func capture(handle *pcap.Handle) {
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer != nil {
			tcp := tcpLayer.(*layers.TCP)

			today := date.Today()
			dataSize := len(packet.Data())
			if ports.Contains(int(tcp.DstPort)) {
				perUnitTimePortUpSize[int(tcp.DstPort)] += dataSize
				datePortUsage[today][int(tcp.DstPort)].Down += dataSize
				continue
			}
			if ports.Contains(int(tcp.SrcPort)) {
				perUnitTimePortDownSize[int(tcp.SrcPort)] += dataSize
				datePortUpSize[today] += dataSize
				continue
			}
		}
	}
}

func calculateNetSpeed() {
	for {
		duration := 1
		for port, size := range perUnitTimePortUpSize {
			portUpSpeed[port] = float32(size) / float32(duration)
			perUnitTimePortUpSize[port] = 0
		}
		for port, size := range perUnitTimePortDownSize {
			portDownSpeed[port] = float32(size) / float32(duration)
			perUnitTimePortDownSize[port] = 0
		}

		time.Sleep(time.Duration(duration) * time.Second)
	}
}
