package server

import (
	"context"
	"dd-log-proxy/logentry"
	"os"
	"sync"
	"testing"
)

func Test_createBatch(t *testing.T) {
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

	os.Setenv("BATCH_SIZE", "1")
	os.Setenv("BATCH_WAIT_IN_SECONDS", "1")

	go testBatch(t, ctx, channel)

	channel <- logEntry

}
func Test_createBatchNotFull(t *testing.T) {
	var waitGroup sync.WaitGroup

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

	os.Setenv("BATCH_SIZE", "10")
	os.Setenv("BATCH_WAIT_IN_SECONDS", "1")

	go testBatchNotFull(&waitGroup, t, ctx, channel)

	channel <- logEntry

	waitGroup.Wait()

}
func Test_createBatchNotFullCancel(t *testing.T) {
	var waitGroup sync.WaitGroup

	channel := make(chan logentry.LogEntry)
	ctx, cancel := context.WithCancel(context.Background())

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

	os.Setenv("BATCH_SIZE", "10")
	os.Setenv("BATCH_WAIT_IN_SECONDS", "10")

	go testBatchNotFull(&waitGroup, t, ctx, channel)

	channel <- logEntry

	//Cancel test and wait on batch to be returned
	cancel()
	waitGroup.Wait()

}

func testBatch(t *testing.T, serverContext context.Context, channel chan logentry.LogEntry) {

	batch := createBatch(serverContext, channel)
	got := len(batch)
	want := 1

	if got != want {
		t.Errorf("Output %q not equal to expected %q", got, want)
	}
}

func testBatchNotFull(waitGroup *sync.WaitGroup, t *testing.T, serverContext context.Context, channel chan logentry.LogEntry) {
	waitGroup.Add(1)
	defer waitGroup.Done()

	batch := createBatch(serverContext, channel)
	got := len(batch)
	want := 1

	if got != want {
		t.Errorf("Output %q not equal to expected %q", got, want)
	}
}
