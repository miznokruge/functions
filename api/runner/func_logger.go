package runner

import (
	"bufio"
	"io"

	"context"
	"github.com/iron-io/runner/common"
)

type FuncLogger interface {
	Writer(context.Context) io.Writer
}

// FuncLogger reads STDERR output from a container and outputs it in a parseable structured log format, see: https://github.com/iron-io/functions/issues/76
type DefaultFuncLogger struct {
}

func NewFuncLogger() FuncLogger {
	return &DefaultFuncLogger{}
}

func (l *DefaultFuncLogger) Writer(ctx context.Context) io.Writer {
	r, w := io.Pipe()

	log := common.Logger(ctx)
	// log = log.WithFields(logrus.Fields{"user_log": true, "app_name": appName, "path": path, "function": function, "call_id": requestID})

	go func(reader io.Reader) {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			log.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.WithError(err).Println("There was an error with the scanner in attached container")
		}
	}(r)

	return w
}