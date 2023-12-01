package server

import (
	"context"
	"dd-log-proxy/logentry"
	"sync"
	"testing"
	"time"
)

func Test_createBatch(t *testing.T) {
	var waitGroup sync.WaitGroup
	defer waitGroup.Wait()

	channel := make(chan logentry.LogEntry)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logEntry := logentry.LogEntry{
		Service:  "test",
		Message:  "test",
		Hostname: "unittest",
		Ddsource: "unittest",
		Ddtags:   "",
		Level:    "info",
		TraceId:  "trace1",
		SpanId:   "span1",
	}

	go testBatch(&waitGroup, t, ctx, channel, 1, 1*time.Second, 1)

	channel <- logEntry
}

func Test_createBatchNotFull(t *testing.T) {
	var waitGroup sync.WaitGroup
	defer waitGroup.Wait()

	channel := make(chan logentry.LogEntry)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logEntry := logentry.LogEntry{
		Service:  "test",
		Message:  "test",
		Hostname: "unittest",
		Ddsource: "unittest",
		Ddtags:   "",
		Level:    "info",
		TraceId:  "trace1",
		SpanId:   "span1",
	}

	go testBatch(&waitGroup, t, ctx, channel, 10, 1*time.Second, 1)

	channel <- logEntry
}

func Test_createBatchNotFullCancel(t *testing.T) {
	var waitGroup sync.WaitGroup
	defer waitGroup.Wait()

	channel := make(chan logentry.LogEntry)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logEntry := logentry.LogEntry{
		Service:  "test",
		Message:  "test",
		Hostname: "unittest",
		Ddsource: "unittest",
		Ddtags:   "",
		Level:    "info",
		TraceId:  "trace1",
		SpanId:   "span1",
	}

	go testBatch(&waitGroup, t, ctx, channel, 10, 1*time.Second, 1)

	channel <- logEntry
}

func Test_batchSentWhenFullOrTimeout(t *testing.T) {
	var waitGroup sync.WaitGroup
	defer waitGroup.Wait()

	channel := make(chan logentry.LogEntry, 50)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Fill the channel's buffer
	for i := 0; i < 25; i++ {
		channel <- logentry.LogEntry{
			Service:  "test",
			Message:  "test",
			Hostname: "unittest",
			Ddsource: "unittest",
			Ddtags:   "",
			Level:    "info",
			TraceId:  "trace1",
			SpanId:   "span1",
		}
	}

	// Two batches should be full
	testBatch(&waitGroup, t, ctx, channel, 10, 1*time.Second, 10)
	testBatch(&waitGroup, t, ctx, channel, 10, 1*time.Second, 10)

	// One should timeout and send partial entries
	testBatch(&waitGroup, t, ctx, channel, 10, 1*time.Second, 5)
}

func testBatch(waitGroup *sync.WaitGroup, t *testing.T, serverContext context.Context, channel chan logentry.LogEntry, maxItemsInBatch int, maxWaitTime time.Duration, expectedBatchSize int) {
	waitGroup.Add(1)
	defer waitGroup.Done()

	batch := createBatch(serverContext, channel, maxItemsInBatch, maxWaitTime)
	got := len(batch)
	want := expectedBatchSize

	if got != want {
		t.Errorf("Output %q not equal to expected %q", got, want)
	}
}
