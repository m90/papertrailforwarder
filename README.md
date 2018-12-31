# papertrailforwarder

[![Build Status](https://travis-ci.org/m90/papertrailforwarder.svg?branch=master)](https://travis-ci.org/m90/papertrailforwarder)
[![godoc](https://godoc.org/github.com/m90/papertrailforwarder?status.svg)](http://godoc.org/github.com/m90/papertrailforwarder)

> Forward CloudWatch logs to Papertrail

## Installation:

Install the library:
```sh
go get github.com/m90/papertrailforwarder
```

## Usage:

Package `papertrailforwarder` allows for easy creation of AWS Lambda handlers that
forward CloudWatch Log Events to Papertrail.

This is an example of an entire Lambda function:

```go
package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/m90/papertrailforwarder"
)

func main() {
	handler, err := papertrailforwarder.New(
		// the host of the papertrail log target
		papertrailforwarder.WithPapertrailHost(<HOST_VALUE>),
		// the port of the papertrail log target as an integer value
		papertrailforwarder.WithPapertrailPort(<PORT_VALUE>),
		// a function that returns a message and whether it should be logged
		// it is passed the raw log message as well as the full event
		papertrailforwarder.WithMessageTransform(func(message string, event events.CloudwatchLogsLogEvent) (string, bool) {
			// prefix all log messages with the event id, forward all messages
			return fmt.Sprintf("[%s] %s", event.ID, message), true
		}),
	)
	if err != nil {
		panic(err)
	}
	lambda.Start(handler)
}
```

When invoked without any options, `New()` will return a handler that uses the
values set in the environment variables `PAPERTRAIL_HOST` and `PAPERTRAIL_PORT`
and the unmodified log message will be forwarded.

Refer to the packages [godoc](http://godoc.org/github.com/m90/papertrailforwarder)
for the entire documentation.

### License
MIT © [Frederik Ring](http://www.frederikring.com)
