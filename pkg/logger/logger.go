package logger

import (
	"os"

	"github.com/charmbracelet/log"
)

var instance log.Logger

func init() {
	instance = *log.NewWithOptions(os.Stdout, log.Options{
		ReportTimestamp: false,
	})
}

var Info = instance.Info
var Debug = instance.Debug
var Warn = instance.Warn
var Error = instance.Error
var Fatal = instance.Fatal
