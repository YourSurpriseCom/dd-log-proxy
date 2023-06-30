package server

import (
	"context"
	"dd-log-proxy/logentry"
	"os"
	"strconv"
	"time"

	log "github.com/jlentink/yaglogger"
)

func createBatch(serverContext context.Context, channel chan logentry.LogEntry) []logentry.LogEntry {
	batchStartTime := time.Now()
	itemsInBatch := 0
	maxItemsInBatch, maxItemsInBatchError := strconv.Atoi(os.Getenv("BATCH_SIZE"))
	maxWaitTimeInSeconds, maxWaitTimeInSecondsError := strconv.Atoi(os.Getenv("BATCH_WAIT_IN_SECONDS"))

	if maxItemsInBatchError != nil || maxWaitTimeInSecondsError != nil {
		log.Fatal("Unable to retrieve batch config from environment variables!")
	}

	var batch []logentry.LogEntry

	for {
		select {
		case logEntry := <-channel:
			batch = append(batch, logEntry)
			itemsInBatch++
			log.Debug("Added message to batch, total batch size %d", itemsInBatch)
		case <-serverContext.Done():
			log.Info("Shutting down, returning open batch with %d messages", itemsInBatch)
			return batch
		default:
			time.Sleep(100 * time.Millisecond)
		}

		if time.Now().After(batchStartTime.Add(time.Duration(maxWaitTimeInSeconds)*time.Second)) && itemsInBatch > 0 {
			log.Info("Max wait time reached, sending %d messages which are waiting in the batch", itemsInBatch)
			break
		}
		if itemsInBatch >= maxItemsInBatch {
			log.Info("Batch is full, sending %d messages to datadog.", maxItemsInBatch)
			break
		}
	}

	return batch
}
