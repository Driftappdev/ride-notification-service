package pkg

import gologger "github.com/driftappdev/libpackage/gologger"

type Logger = gologger.Logger
type LogLevel = gologger.Level
type Field = gologger.Field

const (
	LevelDebug = gologger.LevelDebug
	LevelInfo  = gologger.LevelInfo
	LevelWarn  = gologger.LevelWarn
	LevelError = gologger.LevelError
)

func NewJSONLogger(level gologger.Level) *gologger.Logger {
	return gologger.New(gologger.Options{Level: level, Formatter: &gologger.JSONFormatter{}, AddCaller: true})
}

func LogField(key string, value interface{}) gologger.Field { return gologger.F(key, value) }
