package replay

import (
	"fmt"
	"io"
	"os"

	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"

	"osi-replay/pkg/common"
)

// ReplayPackets reads from cfg.PcapFile and writes raw frames to cfg.InterfaceName.
func ReplayPackets(cfg *common.CaptureConfig, logger *common.Logger) error {
	f, err := os.Open(cfg.PcapFile)
	if err != nil {
		return fmt.Errorf("could not open pcap file %s: %w", cfg.PcapFile, err)
	}
	defer f.Close()

	reader, err := pcapgo.NewReader(f)
	if err != nil {
		return fmt.Errorf("could not create pcapgo reader: %w", err)
	}

	handle, err := pcap.OpenLive(cfg.InterfaceName, cfg.SnapLen, cfg.Promiscuous, cfg.Timeout)
	if err != nil {
		return fmt.Errorf("error opening interface %s: %w", cfg.InterfaceName, err)
	}
	defer handle.Close()

	var count int
	for {
		data, _, err := reader.ReadPacketData()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Error(fmt.Errorf("error reading packet: %w", err))
			continue
		}

		if err := handle.WritePacketData(data); err != nil {
			logger.Error(fmt.Errorf("error writing packet: %w", err))
			continue
		}

		count++
		if count % 1000 == 0 {
			logger.Info(fmt.Sprintf("Replayed %d packets so far...", count))
		}
	}

	logger.Info(fmt.Sprintf("Replay complete. Total packets replayed: %d", count))
	return nil
}
