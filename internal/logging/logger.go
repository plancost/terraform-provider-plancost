package logging

import (
	"github.com/rs/zerolog"
)

var (
	// Logger is the global logger for the logging package. Callers should use this Logger
	// to ensure that logging functionality is consistent across the Infracost codebase.
	//
	// It is advised to create child Loggers as needed and pass them into packages with
	// relevant log metadata. This can be done by using the With method:
	//
	//		childLogger := logging.Logger.With().Str("additional", "field")
	//		foo := MyStruct{Logger: childLogger}
	//		foo.DoSomething()
	//
	// Child loggers will inherit the parent metadata fields, unless the child logger sets metadata
	// field information with the same key. In this case child fields will overwrite the parent field.
	Logger = zerolog.New(TfLogAdapter{}).With().Timestamp().Logger()
)
