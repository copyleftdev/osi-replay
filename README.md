# osi-replay

[![Build Status](https://img.shields.io/github/actions/workflow/status/copyleftdev/osi-replay/go.yml?style=flat-square)](https://github.com/copyleftdev/osi-replay/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/copyleftdev/osi-replay.svg)](https://pkg.go.dev/github.com/copyleftdev/osi-replay)
[![License](https://img.shields.io/github/license/copyleftdev/osi-replay?style=flat-square)](LICENSE)

**osi-replay** is a modular **Go** project for capturing, replaying, transforming, and rewriting network packets at OSI layers 2–4, built on top of [gopacket](https://github.com/google/gopacket). It’s designed for security professionals, network engineers, and QA teams who need to:

- **Capture** real-world traffic in `.pcap` format  
- **Replay** traffic at will on a chosen interface  
- **Sanitize/transform** data (e.g., remove or drop sensitive IPs)  
- **Rewrite** IP and MAC addresses for anonymization or environment mirroring

Visit our website at **[osi-replay.vercel.app](https://osi-replay.vercel.app)** for detailed documentation, tutorials, and release updates.

---

## Table of Contents

1. [Overview](#overview)  
2. [Features](#features)  
3. [Installation & Setup](#installation--setup)  
4. [Usage](#usage)  
   - [Capture](#capture)  
   - [Replay](#replay)  
   - [Transform](#transform)  
   - [Rewriter](#rewriter)  
5. [Architecture](#architecture)  
6. [Advanced Topics](#advanced-topics)  
7. [Contributing](#contributing)  
8. [License](#license)  

---

## Overview

**osi-replay** aims to streamline packet-level testing and debugging workflows. Whether you’re capturing production traffic for local replication or sanitizing sensitive details before distribution, **osi-replay** offers a suite of command-line tools to get the job done quickly and effectively.  

**Repository**: [github.com/copyleftdev/osi-replay](https://github.com/copyleftdev/osi-replay)

---

## Features

- **Live Capture**: Grab packets from any network interface in promiscuous mode.
- **Replay**: Inject captured traffic onto a local interface for debugging, load testing, or simulation.
- **Transform**: Customize or sanitize `.pcap` files (e.g., remove/obfuscate certain IPs).
- **Rewrite**: Update MAC and IP addresses to anonymize data or migrate traffic captures between environments.
- **Modular**: Each operation is split into a separate subcommand, making it easier to integrate into automation scripts.
- **Extensible**: Built on [gopacket](https://github.com/google/gopacket), enabling deeper protocol inspection or custom transformations.

---

## Installation & Setup

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/copyleftdev/osi-replay.git
   cd osi-replay
   ```

2. **Fetch Dependencies**:
   ```bash
   go mod tidy
   ```
   Make sure you have an appropriate version of `libpcap` (on Unix) or equivalent on your operating system.

3. **Build** (via Makefile or directly with `go build`):
   ```bash
   make build
   ```
   This will produce binaries under `./bin/`:  
   - `capture`  
   - `replay`  
   - `transform`  
   - `rewriter`

4. **(Optional) Testing**:
   ```bash
   make test
   ```

For more details and configuration tips, check **[osi-replay.vercel.app](https://osi-replay.vercel.app)**.

---

## Usage

Below are quick usage examples for each subcommand. Adjust flags based on your environment.

### Capture

Capture packets on a specific interface and store them in a `.pcap`:

```bash
./bin/capture -i eth0 -o capture.pcap
```
- **`-i eth0`**: Interface to capture from  
- **`-o capture.pcap`**: Output file for captured packets  

Press **Ctrl+C** to stop the capture process.

---

### Replay

Read a `.pcap` file and replay packets onto an interface:

```bash
./bin/replay -i eth0 -f capture.pcap
```
- **`-i eth0`**: Interface to replay onto  
- **`-f capture.pcap`**: PCAP file to replay  

Ensure your user has the necessary network privileges (e.g., `sudo` or `CAP_NET_RAW`).

---

### Transform

Sanitize or filter sensitive data by reading from a `.pcap`, applying logic in `pkg/sanitizer`, and writing a new file:

```bash
./bin/transform -in capture.pcap -out sanitized_capture.pcap
```
- **`-in capture.pcap`**: Source PCAP  
- **`-out sanitized_capture.pcap`**: Where to store the result  

By default, it drops packets from certain blocked IPs (see `pkg/sanitizer/sanitizer.go`). Customize as needed!

---

### Rewriter

Rewrite IP or MAC addresses for anonymization or environment adaptation:

```bash
./bin/rewriter -in capture.pcap -out rewritten_capture.pcap
```
- **`-in capture.pcap`**: Original capture  
- **`-out rewritten_capture.pcap`**: Output with updated addresses  

Check `pkg/rewriter/rewriter.go` for how to adjust mappings.

---

## Architecture

```
osi-replay/
├── cmd/
│   ├── capture/      # capture tool
│   ├── replay/       # replay tool
│   ├── transform/    # sanitize/transform tool
│   └── rewriter/     # rewriting IP/MACs
├── pkg/
│   ├── capture/      # logic for capturing
│   ├── replay/       # logic for replaying
│   ├── transform/    # transform logic
│   ├── rewriter/     # rewriting logic
│   ├── sanitizer/    # filtering & sanitizing packets
│   └── common/       # shared config, logger, utilities
├── go.mod
├── go.sum
└── README.md
```

- **`cmd/`**: Command-line entry points, minimal main.go files.  
- **`pkg/`**: Core libraries with reusable functionality.  

Each subcommand uses a `common.CaptureConfig` struct for uniform configuration (e.g., interface name, promiscuous mode, etc.).

---

## Advanced Topics

1. **Rate Control**: Implement custom replay pacing in `pkg/replay` to simulate real-world timing.  
2. **Concurrent Pipelines**: For large `.pcap` files, consider pipelining read, transform, and write steps with channels and goroutines.  
3. **Deep Packet Inspection**: Extend `gopacket` layering to decode advanced protocols (HTTP, DNS, TLS) for deeper transformations.  
4. **Testing & Integration**: Place sample `.pcap` fixtures in a `testdata/` folder, then run them through the pipeline to ensure consistency.  
5. **Future Enhancements**: We welcome PRs adding subcommands or expanded OSI layer support (e.g., rewriting L7 data).

---

## Contributing

We welcome community contributions via pull requests on [GitHub](https://github.com/copyleftdev/osi-replay). Please open an issue first to discuss proposed changes or improvements.

**Guidelines**:

- Write clear commit messages and PR descriptions.  
- Add or update tests for new logic.  
- Follow standard Go formatting (`go fmt`).  

For major feature requests or support inquiries, open a GitHub issue. We appreciate your feedback and collaboration!

---

## License

This project is licensed under the **MIT License**. See the [LICENSE](LICENSE) file for details.

---

Happy packet capturing, replaying, and remixing!  
For documentation, tutorials, and updates, visit **[osi-replay.vercel.app](https://osi-replay.vercel.app)**.