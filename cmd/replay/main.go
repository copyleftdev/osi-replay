package main

import (
	"flag"
	"fmt"

	"osi-replay/pkg/common"
	"osi-replay/pkg/replay"
)

func main() {
	var (
		iface  string
		inFile string
	)
	flag.StringVar(&iface, "i", "eth0", "Interface to replay on")
	flag.StringVar(&inFile, "f", "capture.pcap", "PCAP file to replay")
	flag.Parse()

	logger := common.NewLogger("replay-cmd")
	logger.Info(fmt.Sprintf("Replaying from %s on interface %s", inFile, iface))

	cfg := &common.CaptureConfig{
		InterfaceName: iface,
		Promiscuous:   false,
		SnapLen:       65535,
		Timeout:       0,
		PcapFile:      inFile,
	}

	if err := replay.ReplayPackets(cfg, logger); err != nil {
		logger.Fatal(err)
	}
	logger.Info("Replay complete.")
}
