package server

import (
	"bytes"
	"context"
	"dd-log-proxy/datadog"
	"dd-log-proxy/logentry"
	"encoding/json"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	log "github.com/jlentink/yaglogger"
)

func Start() {
	ctx, cancel := context.WithCancel(context.Background())

	var waitGroup sync.WaitGroup

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	channel := make(chan logentry.LogEntry)

	log.Info("Starting UDP server on %s:%s", os.Getenv("HOST"), os.Getenv("PORT"))

	udpServer, err := net.ListenPacket("udp", os.Getenv("HOST")+":"+os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Could not start UDP server on '%s': %v", os.Getenv("HOST")+":"+os.Getenv("PORT"), err)
	}

	go waitForUDPMessage(channel, udpServer)

	// Batch procesing is handled in another thread than the UDP server,
	// which makes it possible to handle the log entries and udp messages at the same time
	go handleLogEntries(&waitGroup, ctx, channel)

	<-sigs
	log.Debug("Stopping udp server...")
	if err := udpServer.Close(); err != nil {
		log.Fatalf("Could not close UDP server correcty, error: %s", err.Error())
	}
	cancel()

	log.Debug("Waiting for async functions to shutdown....")
	waitGroup.Wait()
	log.Info("ByeBye!")
}

func waitForUDPMessage(channel chan logentry.LogEntry, udpServer net.PacketConn) {
	for {
		buf := make([]byte, 65023)
		_, addr, err := udpServer.ReadFrom(buf)
		if err != nil {
			continue
		}
		go handleUDPMessage(addr, buf, channel)
	}
}

func handleUDPMessage(addr net.Addr, buf []byte, channel chan logentry.LogEntry) {
	var logEntry logentry.LogEntry
	b := bytes.Trim(buf, "\x00")
	if err := json.Unmarshal(b, &logEntry); err != nil {
		log.Info("Could not read log entry (error: %s), original data: %s", err.Error(), b)
	} else {
		log.Debug("Received LogEntry from %s; Service: %s, Message: %s", addr.String(), logEntry.Service, logEntry.Message)
		channel <- logEntry
	}
}

func handleLogEntries(waitGroup *sync.WaitGroup, serverContext context.Context, channel chan logentry.LogEntry) {
	waitGroup.Add(1)
	defer waitGroup.Done()

	maxItemsInBatch, maxItemsInBatchError := strconv.Atoi(os.Getenv("BATCH_SIZE"))
	maxWaitTimeInSeconds, maxWaitTimeInSecondsError := strconv.Atoi(os.Getenv("BATCH_WAIT_IN_SECONDS"))
	maxWaitTime := time.Duration(maxWaitTimeInSeconds) * time.Second

	if maxItemsInBatchError != nil || maxWaitTimeInSecondsError != nil {
		log.Fatal("Unable to retrieve batch config from environment variables!")
	}

	for {
		logEntryBatch := createBatch(serverContext, channel, maxItemsInBatch, maxWaitTime)
		if len(logEntryBatch) > 0 {
			err := datadog.SendToDatadog(logEntryBatch)
			if err != nil {
				// Try again to mitigate incidental network issues
				err = datadog.SendToDatadog(logEntryBatch)
				if err != nil {
					log.Error("Could not send batch of size %d to Datadog after 2 attempts: %v", len(logEntryBatch), err)
				}
			}
		} else {
			log.Debug("Nothing to send!")
		}

		select {
		case <-serverContext.Done():
			return
		default:
		}
	}
}
