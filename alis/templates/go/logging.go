package main

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"os"
	"strings"
)

// LogSeverity is used to map the logging levels consistent with Google Cloud Logging.
type LogSeverity string

const (
	// LogSeverity_DEBUG Debug or trace information.
	LogSeverity_DEBUG LogSeverity = "DEBUG"
	// LogSeverity_INFO Routine information, such as ongoing status or performance.
	LogSeverity_INFO LogSeverity = "INFO"
	// LogSeverity_NOTICE Normal but significant events, such as start up, shut down, or
	// a configuration change.
	LogSeverity_NOTICE LogSeverity = "NOTICE"
	// LogSeverity_WARNING Warning events might cause problems.
	LogSeverity_WARNING LogSeverity = "WARNING"
	// LogSeverity_ERROR Error events are likely to cause problems.
	LogSeverity_ERROR LogSeverity = "ERROR"
	// LogSeverity_CRITICAL Critical events cause more severe problems or outages.
	LogSeverity_CRITICAL LogSeverity = "CRITICAL"
	// LogSeverity_ALERT A person must take an action immediately.
	LogSeverity_ALERT LogSeverity = "ALERT"
	// LogSeverity_EMERGENCY One or more systems are unusable.
	LogSeverity_EMERGENCY LogSeverity = "EMERGENCY"
)

// Entry defines a log entry.
// If logs are provided in this format, Google Cloud Logging automatically
// parses the attributes into their LogEntry format as per
// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry which then automatically
// makes the logs available in Google Cloud Logging and Tracing.
type Entry struct {
	Message  string          `json:"message"`
	Severity LogSeverity     `json:"severity,omitempty"`
	Trace    string          `json:"logging.googleapis.com/trace,omitempty"`
	Ctx      context.Context `json:"-"`
	// To extend details sent to the logs, you may add the attributes here.
	//MyAttr1 string `json:"component,omitempty"`
}

// String renders an entry structure to the JSON format expected by Cloud Logging.
func (e Entry) String() string {

	// Defaults to INFO level.
	if e.Severity == "" {
		e.Severity = LogSeverity_INFO
	}

	// Attempt to extract the trace from the context.
	if e.Trace == "" && e.Ctx != nil {
		e.Trace = getTrace(e.Ctx)
	}

	// if Development is local then print out all logs
	if os.Getenv("ENV") == "LOCAL" {
		var prefix string
		switch e.Severity {
		case LogSeverity_DEBUG:
			prefix = colorize("DBG:      ", 90)
		case LogSeverity_INFO:
			prefix = colorize("INFO:     ", 32)
		case LogSeverity_NOTICE:
			prefix = colorize("NOTICE:   ", 34)
		case LogSeverity_WARNING:
			prefix = colorize("WARNING:  ", 33)
		case LogSeverity_ERROR:
			prefix = colorize("ERROR:    ", 31)
		case LogSeverity_ALERT:
			prefix = colorize("ALERT:    ", 91)
		case LogSeverity_CRITICAL:
			prefix = colorize("CRITICAL: ", 41)
		case LogSeverity_EMERGENCY:
			prefix = colorize("EMERGENCY:", 101)
		}
		return prefix + " " + e.Message
	} else {
		out, err := json.Marshal(e)
		if err != nil {
			log.Printf("json.Marshal: %v", err)
		}
		return string(out)
	}
}

// getTrace retrieves a trace header from the provided context.
// Returns an empty string if not found.
func getTrace(ctx context.Context) string {
	// Derive the traceID associated with the current request.
	var trace string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		traceHeaders := md.Get("x-cloud-trace-context")
		if len(traceHeaders) > 0 {
			traceParts := strings.Split(traceHeaders[0], "/")
			if len(traceParts) > 0 && len(traceParts[0]) > 0 {
				trace = fmt.Sprintf("projects/%s/traces/%s", os.Getenv("ALIS_OS_PROJECT"), traceParts[0])
			}
		}
	}
	return trace
}

// serverInterceptor is an example of a Server Interceptor which could be used to 'inject'
// for example logs and/or tracing details to incoming server requests.
// Add this method to your grpc server connection, for example
// grpcServer := grpc.NewServer(grpc.UnaryInterceptor(serverInterceptor))
//	pb.RegisterServiceServer(grpcServer, &myService{})
func serverInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	// Calls the handler
	h, err := handler(ctx, req)
	if err != nil {
		log.Println(&Entry{Message: fmt.Sprintf("%v", req), Severity: LogSeverity_DEBUG, Trace: getTrace(ctx)})
		log.Println(&Entry{
			Message:  err.Error(),
			Severity: LogSeverity_WARNING,
			Trace:    getTrace(ctx),
		})
	}
	return h, err
}

// colorize returns the string s wrapped in ANSI code c
// Codes available at https://en.wikipedia.org/wiki/ANSI_escape_code#Colors
func colorize(s interface{}, c int) string {
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}
