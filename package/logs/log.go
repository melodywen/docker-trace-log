package logs

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strings"
	"time"
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
	logrus.WithField(l.traceName, l.getTraceId(ctx)).Debugf(format, args...)
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
	traceId := ctx.Value(l.traceName)
	if traceStr, ok := traceId.(string); ok {
		return traceStr
	} else {
		panic(fmt.Sprintf("trace id is not string , is value %v", traceId))
	}
}

func (l *Log) SetTraceId(ctx context.Context) context.Context {
	id := strings.Replace(uuid.New().String(), "-", "", -1)

	return context.WithValue(ctx, l.traceName, id)
}

// GetContextOfLog 获取包含链路id 的 ctx
//  @return context.Context
func GetContextOfLog() context.Context {
	instance := NewLog()
	return instance.SetTraceId(context.Background())
}

// EnterExitFunc 打印函数进出日志:
//     使用方法
// 	defer EnterExitFunc()()
func (l *Log) EnterExitFunc(ctx context.Context) func() {
	funcName, file, line := getCallerInfo(true)
	start := time.Now()

	l.Debug(ctx, "enter %s func (%s:%d)", funcName, file, line)
	return func() {
		_, file, line = getCallerInfo(false)
		l.Debug(ctx, "exit %s (%s) func (%s:%d)", funcName, time.Since(start), file, line)
	}
}

func getCallerInfo(needFuncName bool) (string, string, int) {
	pc, file, line, _ := runtime.Caller(2)

	temp := strings.Split(file, "/")
	file = temp[len(temp)-1]

	var funcName string
	if needFuncName {
		temp = strings.Split(runtime.FuncForPC(pc).Name(), ".")
		funcName = temp[len(temp)-1]
	}

	return funcName, file, line
}
