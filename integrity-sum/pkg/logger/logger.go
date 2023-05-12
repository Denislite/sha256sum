package logger

import (
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(config *LoggerConfig) *logrus.Logger {
	l := logrus.New()
	l.Level = logrus.Level(config.Level)
	l.SetReportCaller(true)

	logfile := &lumberjack.Logger{
		Filename:   "./logs/integritySum.log",
		MaxSize:    5, // MB
		MaxBackups: 10,
		MaxAge:     30,   // days
		Compress:   true, // disabled by default
	}

	l.SetOutput(io.MultiWriter(logfile, os.Stdout))
	l.Formatter = &formatter{"[integritySum]"}
	return l
}

// Formatter implements logrus.Formatter interface.
type formatter struct {
	prefix string
}

// Format building log message.
func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var sb bytes.Buffer

	sb.WriteString(strings.ToUpper(entry.Level.String()) + " " + entry.Time.Format(time.RFC3339) + " " + f.prefix + " " + entry.Message + " ")
	file, ok := entry.Data["file"].(string)
	if ok {
		sb.WriteString("file:" + file)
	}
	line, ok := entry.Data["line"].(int)
	if ok {
		sb.WriteString(":" + strconv.Itoa(line))
	}
	function, ok := entry.Data["function"].(string)
	if ok {
		sb.WriteString(" " + "func:" + function)
	}
	sb.WriteString("\n")

	return sb.Bytes(), nil
}
