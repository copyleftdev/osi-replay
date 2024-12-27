package common_test

import (
	"os"
	"path/filepath"
	"testing"

	"osi-replay/pkg/common"
)

func TestEnsureFileWritable_NewFile(t *testing.T) {
	testFile := filepath.Join(os.TempDir(), "ensure_writable_test.file")
	os.Remove(testFile) // ensure a clean state

	err := common.EnsureFileWritable(testFile)
	if err != nil {
		t.Fatalf("Expected no error creating test file, got: %v", err)
	}
	// Clean up
	os.Remove(testFile)
}

func TestEnsureFileWritable_ExistingFile(t *testing.T) {
	testFile := filepath.Join(os.TempDir(), "ensure_writable_test.file")
	// Pre-create the file
	f, _ := os.Create(testFile)
	f.Close()

	err := common.EnsureFileWritable(testFile)
	if err != nil {
		t.Errorf("Expected no error for existing file, got: %v", err)
	}
	// Clean up
	os.Remove(testFile)
}
