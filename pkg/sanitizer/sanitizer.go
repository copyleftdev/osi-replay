package sanitizer

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// Example simple map for blocked IP addresses
var blockedIPs = map[string]bool{
	"10.0.0.1": true,
}

// SanitizePacket returns (nil, false) if the packet should be dropped,
// or (packet, true) otherwise. You can extend logic for more filtering.
func SanitizePacket(packet gopacket.Packet) (gopacket.Packet, bool) {
	if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
		ip4, _ := ipLayer.(*layers.IPv4)
		if blockedIPs[ip4.SrcIP.String()] || blockedIPs[ip4.DstIP.String()] {
			// drop the packet
			return nil, false
		}
	}
	return packet, true
}
