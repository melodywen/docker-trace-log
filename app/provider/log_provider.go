package provider

import (
	"github.com/melodywen/docker-trace-log/contracts"
	log "github.com/sirupsen/logrus"
	"os"
)

type LogProvider struct {
}

func NewLogProvider() *LogProvider {
	return &LogProvider{}
}

func (l LogProvider) StartServerBeforeEvent(app contracts.AppAttributeInterface) error {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.TraceLevel)

	app.SetLog(log.StandardLogger())

	log.Debugf("成功注册了logrus")

	return nil
}

func (l LogProvider) StartServerAfterEvent(app contracts.AppAttributeInterface) error {
	defer app.EnterExitFunc()()
	return nil
}
