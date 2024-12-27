package transform

import (
	"fmt"
	"io"
	"os"

	"osi-replay/pkg/common"
	"osi-replay/pkg/sanitizer"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
)

// Run reads from inFile, applies sanitizer logic, and writes outFile.
func Run(inFile, outFile string, logger *common.Logger) error {
	fIn, err := os.Open(inFile)
	if err != nil {
		return fmt.Errorf("error opening input file: %w", err)
	}
	defer fIn.Close()

	reader, err := pcapgo.NewReader(fIn)
	if err != nil {
		return fmt.Errorf("error creating pcap reader: %w", err)
	}

	fOut, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer fOut.Close()

	writer := pcapgo.NewWriter(fOut)
	if err := writer.WriteFileHeader(65536, layers.LinkTypeEthernet); err != nil {
		return fmt.Errorf("error writing pcap header: %w", err)
	}

	var total, kept int
	for {
		data, ci, err := reader.ReadPacketData()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Error(fmt.Errorf("error reading packet data: %w", err))
			continue
		}
		total++

		packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
		if packet.ErrorLayer() != nil {
			logger.Warn(fmt.Sprintf("Packet decode error: %v", packet.ErrorLayer().Error()))
			continue
		}

		newPacket, keep := sanitizer.SanitizePacket(packet)
		if !keep || newPacket == nil {
			continue
		}

		if err := writer.WritePacket(ci, newPacket.Data()); err != nil {
			logger.Error(fmt.Errorf("error writing sanitized packet: %w", err))
			continue
		}
		kept++
	}

	logger.Info(fmt.Sprintf("Done. Processed %d packets, kept %d.", total, kept))
	return nil
}
