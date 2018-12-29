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
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/m90/papertrailforwarder"
)

func main() {
	handler, err := papertrailforwarder.New(
		papertrailforwarder.WithPapertrailHost(<HOST_VALUE>),
		papertrailforwarder.WithPapertrailPort(<PORT_VALUE>),
	)
	if err != nil {
		panic(err)
	}
	lambda.Start(handler)
}
```

Refer to the packages [godoc](http://godoc.org/github.com/m90/papertrailforwarder)
for the entire documentation.

### License
MIT © [Frederik Ring](http://www.frederikring.com)
