package logs

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"os"
)

type Log struct {
	instance  *logrus.Logger
	traceName string
}

func NewLog() *Log {

	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logrus.SetLevel(logrus.TraceLevel)

	return &Log{
		instance:  logrus.StandardLogger(),
		traceName: "trace_id",
	}
}

func (l *Log) Trace(ctx context.Context, format string, args ...interface{}) {
	logrus.Tracef(format, args...)
}

func (l *Log) Debug(ctx context.Context, format string, args ...interface{}) {
	l.getTraceId(ctx)
	logrus.WithField(l.traceName, "123").Debugf(format, args...)
}

func (l *Log) Info(ctx context.Context, format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func (l *Log) Warning(ctx context.Context, format string, args ...interface{}) {
	logrus.Warningf(format, args...)
}

func (l *Log) Error(ctx context.Context, format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func (l *Log) Fatal(ctx context.Context, format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}

func (l *Log) Panic(ctx context.Context, format string, args ...interface{}) {
	logrus.Panicf(format, args...)
}

func (l *Log) getTraceId(ctx context.Context) string {
	traceId := ctx.Value("trace_id")
	if traceId != nil {
		if traceStr, ok := traceId.(string); ok {
			return traceStr
		} else {
			panic(fmt.Sprintf("trace id is not string , is value %v", traceId))
		}
	}

	return uuid.New().String()
}
