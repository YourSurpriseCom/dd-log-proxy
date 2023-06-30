package datadog

import (
	"dd-log-proxy/logentry"
	"testing"
)

func Test_mapLogEntryToDatadogLogItem(t *testing.T) {
	logEntry := logentry.LogEntry{
		Service:  "test",
		Message:  "this is a test",
		Ddsource: "unittest",
		Hostname: "ci",
		Ddtags:   "test=yes",
		Level:    "info",
		TraceId:  "id123",
		SpanId:   "span123",
	}

	result := mapLogEntryToDatadogLogItem(logEntry)
	switch {
	case *result.Service != logEntry.Service:
		t.Errorf("Service not equal %s %s", *result.Service, logEntry.Service)
	case result.Message != logEntry.Message:
		t.Errorf("Message not equal %s %s", result.Message, logEntry.Message)
	case *result.Ddsource != logEntry.Ddsource:
		t.Errorf("Ddsource not equal %s %s", *result.Ddsource, logEntry.Ddsource)
	case *result.Hostname != logEntry.Hostname:
		t.Errorf("Hostname not equal %s %s", *result.Hostname, logEntry.Hostname)
	case *result.Ddtags != logEntry.Ddtags:
		t.Errorf("Ddtags not equal %s %s", *result.Ddtags, logEntry.Ddtags)
	case result.AdditionalProperties["status"] != logEntry.Level:
		t.Errorf("Level not equal %s %s", result.AdditionalProperties["status"], logEntry.Level)
	case result.AdditionalProperties["dd.trace_id"] != logEntry.TraceId:
		t.Errorf("TraceId not equal %s %s", result.AdditionalProperties["dd.trace_id"], logEntry.TraceId)
	case result.AdditionalProperties["dd.span_id"] != logEntry.SpanId:
		t.Errorf("SpanId not equal %s %s", result.AdditionalProperties["dd.span_id"], logEntry.SpanId)

	}
}
