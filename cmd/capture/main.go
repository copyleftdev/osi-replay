package main

import (
	"flag"
	"fmt"

	"osi-replay/pkg/capture"
	"osi-replay/pkg/common"
)

func main() {
	var (
		iface   string
		outFile string
	)
	flag.StringVar(&iface, "i", "eth0", "Network interface to capture on")
	flag.StringVar(&outFile, "o", "capture.pcap", "Output PCAP file name")
	flag.Parse()

	logger := common.NewLogger("capture-cmd")

	// We do use fmt.Sprintf below, so keep "fmt"
	logger.Info(fmt.Sprintf("Starting capture on interface: %s", iface))

	cfg := &common.CaptureConfig{
		InterfaceName: iface,
		Promiscuous:   true,
		SnapLen:       65535,
		PcapFile:      outFile,
		Timeout:       0,
	}

	if err := capture.CapturePackets(cfg, logger); err != nil {
		logger.Fatal(err)
	}
	logger.Info("Capture complete.")
}
