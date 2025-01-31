package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"

	"protobench/internal/benchmark"
	"protobench/internal/model"
	"protobench/internal/protocols/bson"
	"protobench/internal/protocols/grpc"
	"protobench/internal/protocols/json"
	udp "protobench/internal/protocols/udpack"
	"protobench/internal/protocols/xml"

	"github.com/schollz/progressbar/v3"
)

func runProtocolBenchmark(name string, protocol model.Protocol, messageCount int, messageSize int, shouldProfile bool) benchmark.Result {
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

	runner := benchmark.NewRunner(messageCount, messageSize)
	runner.AddProtocol(name, protocol)

	// Create progress bar
	bar := progressbar.NewOptions(messageCount,
		progressbar.OptionSetDescription(name),
		progressbar.OptionEnableColorCodes(false),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	// Run benchmark with progress updates
	results := runner.RunBenchmarkWithProgress(func(sent, errors int) {
		bar.Set(sent)
		if errors > 0 {
			bar.Describe(fmt.Sprintf("%s (Errors: %d)", name, errors))
		}
	})

	bar.Finish()
	fmt.Println() // Add newline after progress bar

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
	messageSize := flag.Int("kb", 10, "Size of each message in kilobytes")
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
		{"BSON", "8084", func(p string) model.Protocol { return bson.NewClient(p) }},
		{"XML", "8085", func(p string) model.Protocol { return xml.NewClient(p) }},
	}

	var results []benchmark.Result

	fmt.Printf("\nRunning benchmarks (%d messages, %dKB each):\n\n", *messageCount, *messageSize)

	for _, c := range clients {
		client := c.new(c.port)
		if err := client.StartServer(); err != nil {
			log.Fatalf("Failed to start %s server: %v", c.name, err)
		}

		result := runProtocolBenchmark(c.name, client, *messageCount, *messageSize, *shouldProfile)
		results = append(results, result)

		client.StopServer()
	}

	// Print final results table
	fmt.Println("\nResults:")
	fmt.Printf("%-12s %12s %15s %10s %10s\n", "Protocol", "Time", "Msgs/sec", "Errors", "Missing")
	fmt.Println(strings.Repeat("-", 65))

	for _, result := range results {
		fmt.Printf("%-12s %12s %15.2f %10d %10d\n",
			result.Protocol,
			result.TotalTime.Round(time.Millisecond),
			result.MessagesPerSecond,
			result.Errors,
			result.Missing,
		)
	}
}
