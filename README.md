# Protocol Benchmark

A Go-based benchmarking tool that compares different communication protocols and serialization formats for message passing.

[![Go Report Card](https://goreportcard.com/badge/github.com/jpotts18/go-protobench)](https://goreportcard.com/report/github.com/jpotts18/go-protobench)

## Overview

This project implements and benchmarks various protocols to understand their performance characteristics when sending messages of configurable size. It measures throughput, reliability, and error rates across different communication methods.

## Installation

```bash
git clone https://github.com/jpotts18/go-protobench.git
cd go-protobench
go mod download
```

## Protocols Implemented

- **JSON over HTTP**: Traditional REST-style communication using Go's standard library
- **gRPC**: Google's RPC framework using Protocol Buffers
- **UDP with Acknowledgment**: Custom UDP implementation with basic reliability via acks and chunking
- **BSON**: Binary JSON format over TCP with length-prefixed framing
- **XML over HTTP**: Traditional XML-based communication

## Sample Results (1000 messages, 50KB each)

```bash
 ✗ go run cmd/benchmark/main.go --kb 50

Running benchmarks (1000 messages, 50KB each):

JSON 100% [===============] (1000/1000)
gRPC 100% [===============] (1000/1000)
UDP-ACK 100% [===============] (1000/1000)
BSON 100% [===============] (1000/1000)
XML 100% [===============] (1000/1000)

Results:
Protocol             Time        Msgs/sec     Errors    Missing
-----------------------------------------------------------------
JSON               2.669s          374.64          0          0
gRPC                512ms         1951.58          0          0
UDP-ACK            2.125s          470.53          0          0
BSON                366ms         2735.51          0          0
XML                2.811s          355.69          0          0
```

## Key Findings

1. **HTTP-based Protocols (JSON, XML)**

   - Perfect reliability (0 errors, 0 missing)
   - Lower throughput (~350-375 msgs/sec)
   - XML slightly slower than JSON due to more verbose format

2. **gRPC**

   - Excellent balance of speed and reliability
   - ~2000 msgs/sec with no errors
   - Benefits from Protocol Buffers' efficient serialization

3. **UDP with Acknowledgment**

   - Reliable delivery through chunking and retries
   - Moderate throughput (~470 msgs/sec)
   - Handles large messages by breaking them into smaller chunks

4. **BSON over TCP**

   - Highest throughput (~2700 msgs/sec)
   - Reliable delivery through TCP
   - Efficient binary serialization

## Usage

Run basic benchmark with default settings (1000 messages, 10KB each):

```bash
go run cmd/benchmark/main.go
```

Run with custom message count and size:

```bash
go run cmd/benchmark/main.go -n 500 -kb 50
```

Run with profiling:

```bash
go run cmd/benchmark/main.go --profile -n 100 -kb 200
```

All options:

- `-n`: Number of messages to send (default: 1000)
- `-kb`: Size of each message in kilobytes (default: 10)
- `--profile`: Enable CPU and memory profiling

## Future Work

- Optimize UDP chunking and acknowledgment strategy
- Add latency and jitter measurements
- Test with varying payload sizes
- Add raw TCP implementation
- Test under different network conditions and loads
- Add support for bidirectional streaming
