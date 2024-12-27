package rewriter

import (
	"fmt"
	"io"
	"net"
	"os"

	"osi-replay/pkg/common"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
)

type RewriteConfig struct {
	IPMapSrc  map[string]string
	IPMapDst  map[string]string
	MACMapSrc map[string]string
	MACMapDst map[string]string
}

func Run(cfg *RewriteConfig, inFile, outFile string, logger *common.Logger) error {
	fIn, err := os.Open(inFile)
	if err != nil {
		return fmt.Errorf("error opening input file %s: %w", inFile, err)
	}
	defer fIn.Close()

	reader, err := pcapgo.NewReader(fIn)
	if err != nil {
		return fmt.Errorf("error creating pcapgo reader: %w", err)
	}

	fOut, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("error creating output pcap file %s: %w", outFile, err)
	}
	defer fOut.Close()

	writer := pcapgo.NewWriter(fOut)
	if err := writer.WriteFileHeader(65536, layers.LinkTypeEthernet); err != nil {
		return fmt.Errorf("error writing pcap header: %w", err)
	}

	var count int
	for {
		data, ci, err := reader.ReadPacketData()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Error(fmt.Errorf("error reading packet data: %w", err))
			continue
		}

		newData, err := RewritePacket(data, cfg)
		if err != nil {
			logger.Error(fmt.Errorf("rewrite error: %w", err))
			continue
		}
		if err := writer.WritePacket(ci, newData); err != nil {
			logger.Error(fmt.Errorf("error writing rewritten packet: %w", err))
			continue
		}
		count++
	}
	logger.Info(fmt.Sprintf("Rewrite complete. %d packets processed.", count))
	return nil
}

func RewritePacket(data []byte, cfg *RewriteConfig) ([]byte, error) {
	packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
	if packet.ErrorLayer() != nil {
		return nil, fmt.Errorf("error decoding packet: %v", packet.ErrorLayer().Error())
	}

	var modified bool

	if ethLayer := packet.Layer(layers.LayerTypeEthernet); ethLayer != nil {
		eth := ethLayer.(*layers.Ethernet)

		if newMAC, ok := cfg.MACMapSrc[eth.SrcMAC.String()]; ok {
			m, err := net.ParseMAC(newMAC)
			if err == nil {
				eth.SrcMAC = m
				modified = true
			}
		}
		if newMAC, ok := cfg.MACMapDst[eth.DstMAC.String()]; ok {
			m, err := net.ParseMAC(newMAC)
			if err == nil {
				eth.DstMAC = m
				modified = true
			}
		}
	}

	if ipv4Layer := packet.Layer(layers.LayerTypeIPv4); ipv4Layer != nil {
		ip4 := ipv4Layer.(*layers.IPv4)
		if newIP, ok := cfg.IPMapSrc[ip4.SrcIP.String()]; ok {
			parsed := net.ParseIP(newIP)
			if parsed != nil {
				ip4.SrcIP = parsed
				modified = true
			}
		}
		if newIP, ok := cfg.IPMapDst[ip4.DstIP.String()]; ok {
			parsed := net.ParseIP(newIP)
			if parsed != nil {
				ip4.DstIP = parsed
				modified = true
			}
		}
	}

	if ipv6Layer := packet.Layer(layers.LayerTypeIPv6); ipv6Layer != nil {
		ip6 := ipv6Layer.(*layers.IPv6)
		if newIP, ok := cfg.IPMapSrc[ip6.SrcIP.String()]; ok {
			parsed := net.ParseIP(newIP)
			if parsed != nil {
				ip6.SrcIP = parsed
				modified = true
			}
		}
		if newIP, ok := cfg.IPMapDst[ip6.DstIP.String()]; ok {
			parsed := net.ParseIP(newIP)
			if parsed != nil {
				ip6.DstIP = parsed
				modified = true
			}
		}
	}

	if !modified {
		return data, nil
	}
	return serializePacket(packet)
}

func serializePacket(packet gopacket.Packet) ([]byte, error) {
	ethLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethLayer == nil {
		return nil, fmt.Errorf("no ethernet layer found")
	}
	eth := ethLayer.(*layers.Ethernet)

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	var (
		arpLayer   *layers.ARP
		ipv4Layer  *layers.IPv4
		ipv6Layer  *layers.IPv6
		tcpLayer   *layers.TCP
		udpLayer   *layers.UDP
		icmp4Layer *layers.ICMPv4
		icmp6Layer *layers.ICMPv6
	)

	if l := packet.Layer(layers.LayerTypeARP); l != nil {
		arpLayer = l.(*layers.ARP)
	}
	if l := packet.Layer(layers.LayerTypeIPv4); l != nil {
		ipv4Layer = l.(*layers.IPv4)
	}
	if l := packet.Layer(layers.LayerTypeIPv6); l != nil {
		ipv6Layer = l.(*layers.IPv6)
	}
	if l := packet.Layer(layers.LayerTypeTCP); l != nil {
		tcpLayer = l.(*layers.TCP)
	}
	if l := packet.Layer(layers.LayerTypeUDP); l != nil {
		udpLayer = l.(*layers.UDP)
	}
	if l := packet.Layer(layers.LayerTypeICMPv4); l != nil {
		icmp4Layer = l.(*layers.ICMPv4)
	}
	if l := packet.Layer(layers.LayerTypeICMPv6); l != nil {
		icmp6Layer = l.(*layers.ICMPv6)
	}

	var layersToSerialize []gopacket.SerializableLayer
	layersToSerialize = append(layersToSerialize, eth)

	if arpLayer != nil {
		layersToSerialize = append(layersToSerialize, arpLayer)
	} else if ipv4Layer != nil {
		layersToSerialize = append(layersToSerialize, ipv4Layer)
	} else if ipv6Layer != nil {
		layersToSerialize = append(layersToSerialize, ipv6Layer)
	}

	switch {
	case tcpLayer != nil:
		layersToSerialize = append(layersToSerialize, tcpLayer)
	case udpLayer != nil:
		layersToSerialize = append(layersToSerialize, udpLayer)
	case icmp4Layer != nil:
		layersToSerialize = append(layersToSerialize, icmp4Layer)
	case icmp6Layer != nil:
		layersToSerialize = append(layersToSerialize, icmp6Layer)
	}

	if err := gopacket.SerializeLayers(buf, opts, layersToSerialize...); err != nil {
		return nil, fmt.Errorf("serialize error: %w", err)
	}
	return buf.Bytes(), nil
}
