package main

import (
	"flag"
	"fmt"
	"osi-replay/pkg/common"
	"osi-replay/pkg/rewriter"
)

func main() {
	var (
		inFile  string
		outFile string
	)
	flag.StringVar(&inFile, "in", "capture.pcap", "Input PCAP file")
	flag.StringVar(&outFile, "out", "rewritten_capture.pcap", "Output PCAP file")
	flag.Parse()

	logger := common.NewLogger("rewriter-cmd")

	cfg := &rewriter.RewriteConfig{
		IPMapSrc: map[string]string{"192.168.1.100": "10.0.0.5"},
		IPMapDst: map[string]string{"192.168.1.200": "10.0.0.10"},
		MACMapSrc: map[string]string{
			"00:11:22:33:44:55": "aa:bb:cc:dd:ee:ff",
		},
		MACMapDst: map[string]string{},
	}

	logger.Info(fmt.Sprintf("Rewriting packets from %s -> %s", inFile, outFile))

	if err := rewriter.Run(cfg, inFile, outFile, logger); err != nil {
		logger.Fatal(err)
	}
	logger.Info("Rewrite complete.")
}
