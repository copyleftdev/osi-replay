package main

import (
	"flag"
	"fmt"

	"osi-replay/pkg/common"
	"osi-replay/pkg/transform"
)

func main() {
	var (
		inFile  string
		outFile string
	)
	flag.StringVar(&inFile, "in", "capture.pcap", "Input PCAP file")
	flag.StringVar(&outFile, "out", "sanitized_capture.pcap", "Output PCAP file")
	flag.Parse()

	logger := common.NewLogger("transform-cmd")
	logger.Info(fmt.Sprintf("Transforming %s -> %s", inFile, outFile))

	if err := transform.Run(inFile, outFile, logger); err != nil {
		logger.Fatal(err)
	}
	logger.Info("Transformation complete.")
}
