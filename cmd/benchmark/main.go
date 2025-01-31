package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"time"

	"protobench/internal/benchmark"
	"protobench/internal/model"
	"protobench/internal/protocols/bson"
	"protobench/internal/protocols/grpc"
	"protobench/internal/protocols/json"
	udp "protobench/internal/protocols/udpack"
	udpfast "protobench/internal/protocols/udpnoack"
	"protobench/internal/protocols/xml"
)

func runProtocolBenchmark(name string, protocol model.Protocol, messageCount int, shouldProfile bool) benchmark.Result {
	if shouldProfile {
		// Create profile directory
		if err := os.MkdirAll("profiles", 0755); err != nil {
			log.Fatal(err)
		}

		// CPU Profile
		cpuFile, err := os.Create(filepath.Join("profiles", fmt.Sprintf("%s_cpu.prof", name)))
		if err != nil {
			log.Fatal(err)
		}
		defer cpuFile.Close()
		
		if err := pprof.StartCPUProfile(cpuFile); err != nil {
			log.Fatal(err)
		}
		defer pprof.StopCPUProfile()
	}

	// Run benchmark for this protocol
	runner := benchmark.NewRunner(messageCount)
	runner.AddProtocol(name, protocol)
	results := runner.RunBenchmark()

	if shouldProfile {
		// Memory Profile
		runtime.GC()
		memFile, err := os.Create(filepath.Join("profiles", fmt.Sprintf("%s_mem.prof", name)))
		if err != nil {
			log.Fatal(err)
		}
		defer memFile.Close()
		
		if err := pprof.WriteHeapProfile(memFile); err != nil {
			log.Fatal(err)
		}
	}

	return results[0]
}

func main() {
	shouldProfile := flag.Bool("profile", false, "Enable CPU and memory profiling")
	messageCount := flag.Int("n", 1000, "Number of messages to send")
	flag.Parse()

	// Setup protocols
	clients := []struct {
		name string
		port string
		new  func(string) model.Protocol
	}{
		{"JSON", "8080", func(p string) model.Protocol { return json.NewClient(p) }},
		{"gRPC", "8081", func(p string) model.Protocol { return grpc.NewClient(p) }},
		{"UDP-ACK", "8082", func(p string) model.Protocol { return udp.NewClient(p) }},
		{"UDP-NOACK", "8083", func(p string) model.Protocol { return udpfast.NewClient(p) }},
		{"BSON", "8084", func(p string) model.Protocol { return bson.NewClient(p) }},
		{"XML", "8085", func(p string) model.Protocol { return xml.NewClient(p) }},
	}

	var results []benchmark.Result

	// Run benchmark for each protocol
	for _, c := range clients {
		fmt.Printf("Starting %s protocol test...\n", c.name)
		client := c.new(c.port)
		if err := client.StartServer(); err != nil {
			log.Fatalf("Failed to start %s server: %v", c.name, err)
		}

		time.Sleep(time.Second) // Wait for server to start
		fmt.Printf("Running benchmark for %s...\n", c.name)
		result := runProtocolBenchmark(c.name, client, *messageCount, *shouldProfile)
		results = append(results, result)

		fmt.Printf("Stopping %s server...\n", c.name)
		if err := client.StopServer(); err != nil {
			log.Printf("Warning: failed to stop %s server: %v", c.name, err)
		}
		fmt.Printf("Finished %s protocol test\n", c.name)
	}

	// Print results
	fmt.Println("\nBenchmark Results:")
	for _, result := range results {
		fmt.Printf("\nProtocol: %s\n", result.Protocol)
		fmt.Printf("Total Time: %v\n", result.TotalTime)
		fmt.Printf("Messages/second: %.2f\n", result.MessagesPerSecond)
		fmt.Printf("Errors: %d\n", result.Errors)
		fmt.Printf("Missing: %d\n", result.Missing)
	}
} 
