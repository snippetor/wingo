package wingo

import (
	"fmt"
	"time"
)

// New creates and returns test.go new request logger middleware.
// Do not confuse it with the framework's Logger.
// This is for the http requests.
//
// Receives an optional configuation.
func Latency(ctx *Context) {
	//all except latency to string
	var latency time.Duration
	var startTime, endTime time.Time
	startTime = time.Now()

	ctx.Next()

	//no time.Since in order to format it well after
	endTime = time.Now()
	latency = endTime.Sub(startTime)

	// no new line, the framework's logger is responsible how to render each log.
	line := fmt.Sprintf(">>> Latency: %4v", latency)
	ctx.LogD(line)
}
