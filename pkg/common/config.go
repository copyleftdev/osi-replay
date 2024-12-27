package common

import "time"

// CaptureConfig holds global capture/replay config options.
type CaptureConfig struct {
	InterfaceName string
	Promiscuous   bool
	SnapLen       int32
	Timeout       time.Duration
	PcapFile      string
}
