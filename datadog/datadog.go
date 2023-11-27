package datadog

import (
	"context"
	"dd-log-proxy/logentry"
	"encoding/json"
	"fmt"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	log "github.com/jlentink/yaglogger"
)

func SendToDatadog(batch []logentry.LogEntry) error {
	var datadogLogItems []datadogV2.HTTPLogItem

	for _, logEntry := range batch {
		datadogLogItems = append(datadogLogItems, mapLogEntryToDatadogLogItem(logEntry))
	}

	log.Debug("Sending batch with %d logItems to datadog...", len(datadogLogItems))

	ctx := datadog.NewDefaultContext(context.Background())
	configuration := datadog.NewConfiguration()
	apiClient := datadog.NewAPIClient(configuration)
	api := datadogV2.NewLogsApi(apiClient)
	resp, r, err := api.SubmitLog(ctx, datadogLogItems, *datadogV2.NewSubmitLogOptionalParameters().WithContentEncoding(datadogV2.CONTENTENCODING_GZIP))

	if err != nil {
		log.Warn("Error when calling `LogsApi.SubmitLog`: %v\n Full HTTP handleUDPMessage: %v\n", err, r)
		return fmt.Errorf("could not submit logs to Datadog: %w", err)
	} else {
		responseContent, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			log.Warn("Response from `LogsApi.SubmitLog`:\n%s\n", responseContent)
			return fmt.Errorf("invalid response from Datadog logs submit API: %w", err)
		}
	}

	return nil
}

func mapLogEntryToDatadogLogItem(logEntry logentry.LogEntry) datadogV2.HTTPLogItem {
	return datadogV2.HTTPLogItem{
		Ddsource: datadog.PtrString(logEntry.Ddsource),
		Ddtags:   datadog.PtrString(logEntry.Ddtags),
		Hostname: datadog.PtrString(logEntry.Hostname),
		Message:  logEntry.Message,
		Service:  datadog.PtrString(logEntry.Service),
		AdditionalProperties: map[string]string{
			"status":      logEntry.Level,
			"dd.trace_id": logEntry.TraceId,
			"dd.span_id":  logEntry.SpanId,
		},
	}
}
