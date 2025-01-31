# Protocol Benchmark

A Go-based benchmarking tool that compares different communication protocols and serialization formats for message passing.

[![Go Report Card](https://goreportcard.com/badge/github.com/jpotts18/go-protobench)](https://goreportcard.com/report/github.com/jpotts18/go-protobench)

## Overview

This project implements and benchmarks various protocols to understand their performance characteristics when sending ~100KB messages. It measures throughput, reliability, and error rates across different communication methods.

## Installation

```bash
git clone https://github.com/jpotts18/go-protobench.git
cd go-protobench
go mod download
```

## Protocols Implemented

- **JSON over HTTP**: Traditional REST-style communication using Go's standard library
- **gRPC**: Google's RPC framework using Protocol Buffers
- **UDP with Acknowledgment**: Custom UDP implementation with basic reliability via acks
- **UDP without Acknowledgment**: Raw UDP for maximum throughput
- **BSON**: Binary JSON format over UDP
- **XML over HTTP**: Traditional XML-based communication

## Sample Results (1000 messages, 100KB each)

## Key Findings

1. **HTTP-based Protocols (JSON, XML)**

   - Perfect reliability (0 errors, 0 missing)
   - Lowest throughput (~200-250 msgs/sec)
   - XML slightly slower than JSON due to more verbose format

2. **gRPC**

   - Excellent balance of speed and reliability
   - ~1200 msgs/sec with no errors
   - Benefits from Protocol Buffers' efficient serialization

3. **UDP Protocols**
   - Highest raw throughput (~1600-2000 msgs/sec)
   - Complete message loss at 100KB payload size
   - Even with acknowledgments, unable to handle large messages reliably
   - BSON serialization fastest but limited by UDP transport

## Usage

Run basic benchmark

`go run cmd/benchmark/main.go`

Run with CPU and memory profiling

`go run cmd/benchmark/main.go --profile`

Specify number of messages

`go run cmd/benchmark/main.go -n 500`

## Conclusions

The results demonstrate clear trade-offs between reliability and performance:

- HTTP-based protocols provide guaranteed delivery but with significant overhead
- gRPC offers an excellent balance of performance and reliability
- UDP protocols show that for large messages (100KB), reliable delivery requires additional protocol support
- Raw UDP, while fastest, is unsuitable for large messages without additional reliability mechanisms

## Future Work

- Implement retries and backoff for UDP protocols
- Add latency and jitter measurements
- Test with varying payload sizes
- Add raw TCP implementation
- Test under different network conditions and loads
- Add support for bidirectional streaming

## Project Structure

```
.
├── cmd/
│   └── benchmark/          # Main benchmark executable
├── internal/
│   ├── benchmark/         # Benchmark runner and results
│   ├── model/            # Shared data models
│   └── protocols/        # Protocol implementations
│       ├── bson/
│       ├── grpc/
│       ├── json/
│       ├── udpack/       # UDP with acknowledgments
│       ├── udpnoack/     # UDP without acknowledgments
│       └── xml/
├── profiles/             # Generated profiling data (gitignored)
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```
