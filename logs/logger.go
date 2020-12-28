package logs

import (
	"fmt"
	"github.com/huija/summer/utils"
	"io"
	"log"
	"os"
)

const (
	Console = "console"
	File    = "file"

	// file
	defaultPath = "./conf/test.log"
)

const (
	DebugLevel int8 = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

// Log config
type Log struct {
	// non-zero value
	Type string `json:",omitempty"`

	// zero val is valid
	Level int8
	File  *FileLog `yaml:",omitempty"`
}

// FileLog config
type FileLog struct {
	// non-zero value
	Path       string `json:",omitempty"`
	Maxsize    int    `json:",omitempty"`
	MaxBackups int    `json:",omitempty"`
	MaxAge     int    `json:",omitempty"`

	// zero val is valid
	Compress  bool
	LocalZone bool //backup time zone
}

// SugaredLogger global sugar logger
// provide based log print func
var SugaredLogger Logger

// Writer global root writer
var Writer io.Writer

// default logger & writer
func init() {
	Writer = os.Stdout
	SugaredLogger = &innerLog{log.New(Writer, "", log.LstdFlags|log.Lshortfile)}
}

// Defaults
func Defaults(ls []*Log) ([]*Log, error) {
	if ls == nil || len(ls) == 0 {
		return []*Log{defaultsConsole()}, nil
	}
	var err error
	for _, l := range ls {
		switch l.Type {
		case File:
			if l.File != nil {
				err = utils.MergeStructByMarshal(l.File, defaultsFile())
			} else {
				l.File = defaultsFile()
			}
			if err != nil {
				return ls, err
			}
		default:
			err = utils.MergeStructByMarshal(l, defaultsConsole())
			if err != nil {
				return ls, err
			}
		}
	}
	return ls, nil
}

func defaultsConsole() *Log {
	return &Log{
		Level: DebugLevel,
		Type:  Console,
	}
}

func defaultsFile() *FileLog {
	return &FileLog{
		Path:       defaultPath,
		Maxsize:    1024, // 1024 mb
		MaxBackups: 7,    // seven files
		MaxAge:     30,   // one month
		Compress:   false,
		LocalZone:  false,
	}
}

// Logger sugared logger interface
type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}

type innerLog struct {
	*log.Logger
}

func (i *innerLog) Debug(v ...interface{}) {
	i.Output(2, fmt.Sprintln(v...))
}

func (i *innerLog) Debugf(format string, v ...interface{}) {
	i.Output(2, fmt.Sprintf(format, v...))
}

func (i *innerLog) Info(v ...interface{}) {
	i.Output(2, fmt.Sprintln(v...))
}

func (i *innerLog) Infof(format string, v ...interface{}) {
	i.Output(2, fmt.Sprintf(format, v...))
}

func (i *innerLog) Warn(v ...interface{}) {
	i.Output(2, fmt.Sprintln(v...))
}

func (i *innerLog) Warnf(format string, v ...interface{}) {
	i.Output(2, fmt.Sprintf(format, v...))
}

func (i *innerLog) Error(v ...interface{}) {
	i.Output(2, fmt.Sprintln(v...))
}

func (i *innerLog) Errorf(format string, v ...interface{}) {
	i.Output(2, fmt.Sprintf(format, v...))
}
