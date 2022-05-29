package logger

import (
	"fmt"
	"io"
	stdLog "log"
	"strings"
)

//LogLevel constment for logger level
type LogLevel int

//go:generate stringer -type LogLevel -linecomment
const (
	//ERROR only show error
	ERROR LogLevel = iota
	//NOTICE normal level
	NOTICE
	//DEBUG show debug
	DEBUG
)

//Logger logger
type Logger interface {
	Error(string)
	Errorf(string, ...interface{})
	Notice(string)
	Noticef(string, ...interface{})
	Debug(string)
	Debugf(string, ...interface{})
}

//GetLogger get a Logger
func GetLogger(level LogLevel, ws ...io.Writer) Logger {
	l := stdLog.New(io.MultiWriter(ws...), "", stdLog.LstdFlags|stdLog.Lshortfile)
	return &logger{
		output: l,
		level:  level,
	}
}

type logger struct {
	output *stdLog.Logger
	level  LogLevel
}

func (l *logger) Error(s string) {
	_ = l.output.Output(2, strings.Join([]string{ERROR.String(), s}, " "))
}

func (l *logger) Errorf(t string, i ...interface{}) {
	_ = l.output.Output(2, strings.Join([]string{ERROR.String(), fmt.Sprintf(t, i...)}, " "))
}

func (l *logger) Notice(s string) {
	if l.level < NOTICE {
		return
	}
	_ = l.output.Output(2, strings.Join([]string{NOTICE.String(), s}, " "))
}

func (l *logger) Noticef(t string, i ...interface{}) {
	if l.level < NOTICE {
		return
	}
	_ = l.output.Output(2, strings.Join([]string{NOTICE.String(), fmt.Sprintf(t, i...)}, " "))
}

func (l *logger) Debug(s string) {
	if l.level < DEBUG {
		return
	}
	_ = l.output.Output(2, strings.Join([]string{DEBUG.String(), s}, " "))
}

func (l *logger) Debugf(t string, i ...interface{}) {
	if l.level < DEBUG {
		return
	}
	_ = l.output.Output(2, strings.Join([]string{DEBUG.String(), fmt.Sprintf(t, i...)}, " "))
}
