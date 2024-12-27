package capture_test

import (
	"os"
	"testing"

	"osi-replay/pkg/capture"
	"osi-replay/pkg/common"
)

// Test invalid interface name
func TestCapturePackets_InvalidInterface(t *testing.T) {
	logger := common.NewLogger("test-capture")
	cfg := &common.CaptureConfig{
		InterfaceName: "does-not-exist",
		Promiscuous:   true,
		SnapLen:       65535,
		PcapFile:      "test_capture.pcap",
	}

	err := capture.CapturePackets(cfg, logger)
	if err == nil {
		t.Fatalf("Expected an error for a non-existent interface, got nil")
	}
	// Clean up if the file got created
	_ = os.Remove("test_capture.pcap")
}

// Test file creation error by providing an invalid path
func TestCapturePackets_FileCreationError(t *testing.T) {
	logger := common.NewLogger("test-capture")
	cfg := &common.CaptureConfig{
		InterfaceName: "lo", // or "eth0" depending on your system
		Promiscuous:   true,
		SnapLen:       65535,
		PcapFile:      "/invalid-dir/test_capture.pcap",
	}

	err := capture.CapturePackets(cfg, logger)
	if err == nil {
		t.Fatal("Expected an error due to invalid file path, got nil")
	}
}

// Optional real capture test (skipped by default)
func TestCapturePackets_Simple(t *testing.T) {
	t.Skip("Skipping real capture test in CI or restricted environment")

	logger := common.NewLogger("test-capture")
	cfg := &common.CaptureConfig{
		InterfaceName: "lo", // or another valid interface
		Promiscuous:   true,
		SnapLen:       65535,
		PcapFile:      "test_capture.pcap",
	}

	err := capture.CapturePackets(cfg, logger)
	if err != nil {
		t.Errorf("Unexpected error capturing on interface 'lo': %v", err)
	}
	_ = os.Remove("test_capture.pcap")
}
