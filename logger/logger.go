package logger

type Logger interface {
	Error(string)
	Errorf(string, ...interface{})
	Notice(string)
	Noticef(string, ...interface{})
	Debug(string)
	Debugf(string, ...interface{})
}
