package rewriter_test

import (
	"net"
	"testing"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"

	"osi-replay/pkg/rewriter"
)

func TestRewritePacket_Basic(t *testing.T) {
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	eth := &layers.Ethernet{
		SrcMAC:       net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
		DstMAC:       net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
		EthernetType: layers.EthernetTypeIPv4,
	}
	ip4 := &layers.IPv4{
		SrcIP:    net.ParseIP("192.168.1.100"),
		DstIP:    net.ParseIP("192.168.1.200"),
		Version:  4,
		IHL:      5,
		Protocol: layers.IPProtocolTCP,
	}

	if err := gopacket.SerializeLayers(buf, opts, eth, ip4); err != nil {
		t.Fatalf("Error serializing layers: %v", err)
	}

	original := buf.Bytes()

	cfg := &rewriter.RewriteConfig{
		IPMapSrc: map[string]string{"192.168.1.100": "10.0.0.5"},
		IPMapDst: map[string]string{"192.168.1.200": "10.0.0.10"},
		MACMapSrc: map[string]string{
			"00:11:22:33:44:55": "66:77:88:99:aa:bb",
		},
		MACMapDst: map[string]string{
			"aa:bb:cc:dd:ee:ff": "11:22:33:44:55:66",
		},
	}

	newData, err := rewriter.RewritePacket(original, cfg)
	if err != nil {
		t.Fatalf("RewritePacket returned error: %v", err)
	}

	newPacket := gopacket.NewPacket(newData, layers.LayerTypeEthernet, gopacket.Default)
	ethLayer := newPacket.Layer(layers.LayerTypeEthernet)
	ipLayer := newPacket.Layer(layers.LayerTypeIPv4)

	if ethLayer == nil || ipLayer == nil {
		t.Fatalf("Missing expected layers in rewritten packet")
	}

	ethOut, _ := ethLayer.(*layers.Ethernet)
	ip4Out, _ := ipLayer.(*layers.IPv4)

	// Verify MAC addresses
	if ethOut.SrcMAC.String() != "66:77:88:99:aa:bb" {
		t.Errorf("Expected MAC src 66:77:88:99:aa:bb, got %s", ethOut.SrcMAC)
	}
	if ethOut.DstMAC.String() != "11:22:33:44:55:66" {
		t.Errorf("Expected MAC dst 11:22:33:44:55:66, got %s", ethOut.DstMAC)
	}

	// Verify IP addresses
	if ip4Out.SrcIP.String() != "10.0.0.5" {
		t.Errorf("Expected IP src 10.0.0.5, got %s", ip4Out.SrcIP)
	}
	if ip4Out.DstIP.String() != "10.0.0.10" {
		t.Errorf("Expected IP dst 10.0.0.10, got %s", ip4Out.DstIP)
	}
}
