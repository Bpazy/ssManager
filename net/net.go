package net

import (
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"time"
)

var (
	ports          = arraylist.New()
	portDownSize   = make(map[int]int)
	portUpDataSize = make(map[int]int)
	portDownSpeed  = make(map[int]float32)
	portUpSpeed    = make(map[int]float32)
)

func AddPorts(_ports []int) {
	for _, port := range _ports {
		ports.Add(port)
	}
}

func RemovePort(port int) {
	ports.Remove(port)
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

			if ports.Contains(int(tcp.DstPort)) {
				portUpDataSize[int(tcp.DstPort)] += len(packet.Data())
				continue
			}
			if ports.Contains(int(tcp.SrcPort)) {
				portDownSize[int(tcp.SrcPort)] += len(packet.Data())
				continue
			}
		}
	}
}

func calculateNetSpeed() {
	for {
		duration := 1
		for port, size := range portUpDataSize {
			portUpSpeed[port] = float32(size) / float32(duration)
			portUpDataSize[port] = 0
		}
		for port, size := range portDownSize {
			portDownSpeed[port] = float32(size) / float32(duration)
			portDownSize[port] = 0
		}

		time.Sleep(time.Duration(duration) * time.Second)
	}
}
