package replay_test

import (
	"os"
	"path/filepath"
	"testing"

	"osi-replay/pkg/common"
	"osi-replay/pkg/replay"
)

// Test missing PCAP file
func TestReplayPackets_NoSuchFile(t *testing.T) {
	logger := common.NewLogger("test-replay")
	cfg := &common.CaptureConfig{
		InterfaceName: "eth999", // likely non-existent
		SnapLen:       65535,
		Promiscuous:   false,
		PcapFile:      "no_such_file.pcap",
	}

	err := replay.ReplayPackets(cfg, logger)
	if err == nil {
		t.Errorf("Expected error when PCAP file does not exist")
	}
}

// Optional real replay test (skipped by default)
func TestReplayPackets_Basic(t *testing.T) {
	t.Skip("Skipping real replay test in CI or restricted environment")

	// Provide a small real PCAP in a testdata/ folder:
	pcapFile := filepath.Join("testdata", "small_test.pcap")
	if _, err := os.Stat(pcapFile); os.IsNotExist(err) {
		t.Skipf("Skipping: testdata file %v not found", pcapFile)
	}

	logger := common.NewLogger("test-replay")
	cfg := &common.CaptureConfig{
		InterfaceName: "lo", // or eth0, if valid
		SnapLen:       65535,
		Promiscuous:   false,
		PcapFile:      pcapFile,
	}

	err := replay.ReplayPackets(cfg, logger)
	if err != nil {
		t.Errorf("Unexpected error replaying pcap: %v", err)
	}
}
