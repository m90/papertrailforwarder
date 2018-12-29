package papertrailforwarder

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestConfiguration(t *testing.T) {
	tests := []struct {
		name        string
		config      configuration
		expectError bool
	}{
		{
			"ok",
			configuration{
				papertrailHost: "logs1.papertrail.com",
				papertrailPort: 1234,
			},
			false,
		},
		{
			"empty struct",
			configuration{},
			true,
		},
		{
			"no host",
			configuration{
				papertrailPort: 1234,
			},
			true,
		},
		{
			"no port",
			configuration{
				papertrailHost: "logs1.papertrail.com",
			},
			true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.config.validate(); test.expectError != (err != nil) {
				t.Errorf("Unexpected error %v", err)
			}
		})
	}
}

func TestOptions(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		config := configuration{
			papertrailHost: "logs1.papertrail.com",
			papertrailPort: 9999,
			messageTransform: func(string, events.CloudwatchLogsLogEvent) (string, bool) {
				return "before", true
			},
		}
		config = WithPapertrailPort(1111)(config)
		config = WithPapertrailHost("logs9.papertrail.com")(config)
		config = WithMessageTransform(func(string, events.CloudwatchLogsLogEvent) (string, bool) {
			return "after", true
		})(config)

		if config.papertrailHost != "logs9.papertrail.com" {
			t.Errorf("Unexpected host value %s", config.papertrailHost)
		}
		if config.papertrailPort != 1111 {
			t.Errorf("Unexpected port value %d", config.papertrailPort)
		}
	})
}

func TestNew(t *testing.T) {
	t.Run("no configuration", func(t *testing.T) {
		if _, err := New(); err == nil {
			t.Error("Expected error, got nil")
		}
	})
	t.Run("missing host", func(t *testing.T) {
		if _, err := New(WithPapertrailPort(1111)); err == nil {
			t.Error("Expected error, got nil")
		}
	})
	t.Run("default", func(t *testing.T) {
		if _, err := New(
			WithPapertrailPort(1234),
			WithPapertrailHost("logs7.papertrail.com"),
			WithMessageTransform(func(string, events.CloudwatchLogsLogEvent) (string, bool) {
				return "ok", true
			}),
		); err != nil {
			t.Errorf("Unexpected error %v", err)
		}
	})
}
