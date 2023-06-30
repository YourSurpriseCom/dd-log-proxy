package logentry

type LogEntry struct {
	Service  string
	Message  string
	Hostname string
	Ddsource string
	Ddtags   string
	Level    string
	TraceId  string
	SpanId   string
}
