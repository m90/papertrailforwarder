package papertrailforwarder

import (
	"context"
	"errors"
	"fmt"
	"log/syslog"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
)

// Option is a function that transforms the default configuration
type Option func(configuration) configuration

// MessageTransform is a function that is used for applying conditional transforms
// to log messages. The first argument passed is the raw log message, the second
// one is the entire event. It can return false as its second return value to
// signal that the event should be skipped and not be forwarded to Papertrail.
type MessageTransform func(string, events.CloudwatchLogsLogEvent) (string, bool)

type configuration struct {
	papertrailHost   string
	papertrailPort   int
	messageTransform MessageTransform
}

func (c configuration) validate() error {
	if c.papertrailHost == "" {
		return errors.New("missing host value for Papertrail log target")
	}
	if c.papertrailPort == 0 {
		return errors.New("missing port value for Papertrail log target")
	}
	return nil
}

// WithPapertrailHost sets the host of the Papertrail log target
func WithPapertrailHost(host string) Option {
	return func(config configuration) configuration {
		config.papertrailHost = host
		return config
	}
}

// WithPapertrailPort sets the port value for the Papertrail log target
func WithPapertrailPort(port int) Option {
	return func(config configuration) configuration {
		config.papertrailPort = port
		return config
	}
}

// WithMessageTransform defines a transform that will be applied to the log message
func WithMessageTransform(transform MessageTransform) Option {
	return func(config configuration) configuration {
		config.messageTransform = transform
		return config
	}
}

// New creates a new AWS Lambda handler that forwards the passed CloudWatch
// Logs Event to Papertrail using the given configuration options. By default
// it will use the values set in the environment variables `PAPERTRAIL_HOST`
// and `PAPERTRAIL_PORT`
func New(options ...Option) (func(context.Context, events.CloudwatchLogsEvent) error, error) {
	portValue, _ := strconv.Atoi(os.Getenv("PAPERTRAIL_PORT"))
	config := configuration{
		papertrailHost: os.Getenv("PAPERTRAIL_HOST"),
		papertrailPort: portValue,
		messageTransform: func(message string, event events.CloudwatchLogsLogEvent) (string, bool) {
			return message, true
		},
	}

	for _, option := range options {
		config = option(config)
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return func(ctx context.Context, log events.CloudwatchLogsEvent) error {
		data, dataErr := log.AWSLogs.Parse()
		if dataErr != nil {
			return dataErr
		}

		logger, loggerErr := syslog.Dial(
			"udp", fmt.Sprintf("%s:%d", config.papertrailHost, config.papertrailPort), syslog.LOG_EMERG, data.LogGroup,
		)
		if loggerErr != nil {
			return loggerErr
		}

		for _, event := range data.LogEvents {
			if line, ok := config.messageTransform(event.Message, event); ok {
				logger.Info(line)
			}
		}

		return nil
	}, nil
}
