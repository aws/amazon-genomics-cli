// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package logging

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	gray   = 90
	red    = 31
	yellow = 33
	white  = 97
)

var Verbose bool

// ApplicationConsoleLogger generates a zerolog Logger that pretty prints text (not JSON) to STDERR.
//
// This can be used as the default logger as follows:
// 		log.Logger = logging.ApplicationConsoleLogger()
func ApplicationConsoleLogger() zerolog.Logger {
	logger := log.Output(zerolog.ConsoleWriter{
		Out:           os.Stderr,
		TimeFormat:    time.RFC3339,
		FormatLevel:   levelIconFormatter(),
		FormatMessage: coloredMessageFormatter(),
	})
	return logger
}

func colorize(s interface{}, c int) string {
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

// levelIconFormatter provides a custom Level formatter that maps logging levels to colored icons
func levelIconFormatter() zerolog.Formatter {
	return func(i interface{}) string {
		var formatted string
		if ll, ok := i.(string); ok {
			switch ll {
			case "trace":
				formatted = colorize(tracePrefix, gray)
			case "debug":
				formatted = colorize(debugPrefix, gray)
			case "info":
				formatted = colorize(infoPrefix, white)
			case "warn":
				formatted = colorize(warningPrefix, yellow)
			case "error":
				formatted = colorize(errorPrefix, red)
			case "fatal":
				formatted = colorize(fatalPrefix, red)
			case "panic":
				formatted = colorize(panicPrefix, red)
			default:
				formatted = colorize("??", yellow)
			}
		} else {
			if i == nil {
				formatted = colorize("??", yellow)
			} else {
				formatted = fmt.Sprintf("%s", i)
			}
		}
		return formatted
	}
}

func coloredMessageFormatter() zerolog.Formatter {
	return func(i interface{}) string {
		if i == nil {
			return ""
		}
		return colorize(i, gray)
	}
}
