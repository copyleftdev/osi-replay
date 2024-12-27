package sanitizer_test

import (
	"net"
	"testing"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"

	"osi-replay/pkg/sanitizer"
)

func TestSanitizePacket_BlockedIP(t *testing.T) {
	// Synthetic IPv4 packet with src=10.0.0.1
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	ip := &layers.IPv4{
		Version: 4,
		IHL:     5,
		SrcIP:   net.ParseIP("10.0.0.1"),
		DstIP:   net.ParseIP("192.168.1.10"),
	}

	if err := gopacket.SerializeLayers(buf, opts, ip); err != nil {
		t.Fatalf("Error serializing IP packet: %v", err)
	}

	packet := gopacket.NewPacket(buf.Bytes(), layers.LayerTypeIPv4, gopacket.Default)
	newPacket, keep := sanitizer.SanitizePacket(packet)
	if keep || newPacket != nil {
		t.Errorf("Expected packet to be dropped for blocked IP 10.0.0.1")
	}
}

func TestSanitizePacket_AllowedIP(t *testing.T) {
	// Synthetic IPv4 packet with src=192.168.100.5, not in blockedIPs
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	ip := &layers.IPv4{
		Version: 4,
		IHL:     5,
		SrcIP:   net.ParseIP("192.168.100.5"),
		DstIP:   net.ParseIP("192.168.100.10"),
	}

	if err := gopacket.SerializeLayers(buf, opts, ip); err != nil {
		t.Fatalf("Error serializing IP packet: %v", err)
	}

	packet := gopacket.NewPacket(buf.Bytes(), layers.LayerTypeIPv4, gopacket.Default)
	newPacket, keep := sanitizer.SanitizePacket(packet)
	if !keep || newPacket == nil {
		t.Errorf("Expected packet to pass sanitizer, but was dropped")
	}
}
