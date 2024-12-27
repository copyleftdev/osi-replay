package capture

import (
	"fmt"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"

	"osi-replay/pkg/common"
)

func CapturePackets(cfg *common.CaptureConfig, logger *common.Logger) error {
	handle, err := pcap.OpenLive(cfg.InterfaceName, cfg.SnapLen, cfg.Promiscuous, cfg.Timeout)
	if err != nil {
		return fmt.Errorf("failed to open device %s: %w", cfg.InterfaceName, err)
	}
	defer handle.Close()

	f, err := os.Create(cfg.PcapFile)
	if err != nil {
		return fmt.Errorf("failed to create pcap file %s: %w", cfg.PcapFile, err)
	}
	defer f.Close()

	writer := pcapgo.NewWriter(f)
	if err := writer.WriteFileHeader(uint32(cfg.SnapLen), handle.LinkType()); err != nil {
		return fmt.Errorf("failed to write file header: %w", err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	logger.Info("Capture started. Press Ctrl+C to stop...")

	for packet := range packetSource.Packets() {
		ci := packet.Metadata().CaptureInfo
		if err := writer.WritePacket(ci, packet.Data()); err != nil {
			logger.Error(fmt.Errorf("failed to write packet: %w", err))
		}
	}

	logger.Info("No more packets to read. Capture done.")
	return nil
}
