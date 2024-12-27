package transform_test

import (
	"os"
	"path/filepath"
	"testing"

	"osi-replay/pkg/common"
	"osi-replay/pkg/transform"
)

// Test error on nonexistent input file
func TestRun_NoSuchFile(t *testing.T) {
	logger := common.NewLogger("test-transform")
	err := transform.Run("no_such_file.pcap", "out.pcap", logger)
	if err == nil {
		t.Errorf("Expected error with nonexistent input file, got nil")
	}
	_ = os.Remove("out.pcap")
}

// Optional real transform test
func TestRun_Basic(t *testing.T) {
	t.Skip("Skipping real transform test if no testdata pcap available")

	pcapIn := filepath.Join("testdata", "example_in.pcap")
	pcapOut := filepath.Join("testdata", "example_out.pcap")

	if _, err := os.Stat(pcapIn); os.IsNotExist(err) {
		t.Skipf("Skipping test, no input pcap at %s", pcapIn)
	}

	logger := common.NewLogger("test-transform")
	_ = os.Remove(pcapOut)

	err := transform.Run(pcapIn, pcapOut, logger)
	if err != nil {
		t.Errorf("Unexpected error transforming pcap: %v", err)
	}

	// Check if output got created
	if _, err := os.Stat(pcapOut); os.IsNotExist(err) {
		t.Errorf("Output file %s not created", pcapOut)
	}
	_ = os.Remove(pcapOut)
}
