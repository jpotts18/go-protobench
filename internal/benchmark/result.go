package benchmark

import "time"

type Result struct {
	Protocol          string
	TotalTime         time.Duration
	MessagesPerSecond float64
	Errors            int
	Missing           int
}
