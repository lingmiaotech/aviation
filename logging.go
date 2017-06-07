package tonic

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

type LogHandler struct {
	Name      string
	Formatter logrus.Formatter
	GetHook   func(loggerName string) (logrus.Hook, error)
}

type LoggingClass struct {
	AppName    string
	Loggers    map[string]*logrus.Logger
	Handler    map[string]*LogHandler
	Formatters map[string]logrus.Formatter
}

func (logging LoggingClass) GetDefaultLogger() *logrus.Logger {
	defaultLogger, _ := logging.Loggers["default"]
	return defaultLogger
}

func (logging LoggingClass) GetLogger(name string) *logrus.Logger {
	logger, ok := logging.Loggers[name]
	if ok {
		return logger
	}
	defaultLogger, _ := logging.Loggers["default"]
	return defaultLogger
}

var Logging LoggingClass

func InitLogging() (err error) {

	Logging.Loggers = make(map[string]*logrus.Logger)
	Logging.Handler = make(map[string]*LogHandler)
	Logging.Formatters = make(map[string]logrus.Formatter)

	Statsd.AppName = Configs.GetString("app_name")

	formatters, ok := Configs.Get("logging.formatters").([]interface{})
	if !ok {
		return errors.New("tonic_error.logging.invalid_config_format.formatters")
	}

	for _, formatter := range formatters {
		formatterMap, ok := formatter.(map[interface{}]interface{})
		if !ok {
			return errors.New("tonic_error.logging.invalid_config_format.formatters.type")
		}

		name, ok := formatterMap["name"].(string)
		if !ok {
			return errors.New("tonic_error.logging.invalid_config_format.formatters.name")
		}

		format, ok := formatterMap["format"].(string)
		if !ok {
			return errors.New("tonic_error.logging.invalid_config_format.formatters.format")
		}

		color, ok := formatterMap["color"].(bool)
		if !ok {
			color = false
		}

		f, err := getFormatter(format, color)
		if err != nil {
			return err
		}

		Logging.Formatters[name] = f
	}

	handlers, ok := Configs.Get("logging.handlers").([]interface{})
	if !ok {
		return errors.New("tonic_error.logging.invalid_config_format.handlers")
	}

	for _, handler := range handlers {
		handlerMap, ok := handler.(map[interface{}]interface{})
		if !ok {
			return errors.New("tonic_error.logging.invalid_config_format.handlers.type")
		}

		name, ok := handlerMap["name"].(string)
		if !ok {
			return errors.New("tonic_error.logging.invalid_config_format.handlers.name")
		}

		handle, ok := handlerMap["handle"].(string)
		if !ok {
			return errors.New("tonic_error.logging.invalid_config_format.handlers.handle")
		}

		formatter, ok := handlerMap["formatter"].(string)
		if !ok {
			return errors.New("tonic_error.logging.invalid_config_format.handlers.formatter")
		}

		h, err := getHandler(handle, formatter)
		if err != nil {
			return err
		}

		Logging.Handler[name] = h
	}

	loggers, ok := Configs.Get("logging.loggers").([]interface{})
	if !ok {
		return errors.New("tonic_error.logging.invalid_config_format.loggers")
	}

	for _, logger := range loggers {
		loggerMap, ok := logger.(map[interface{}]interface{})
		if !ok {
			return errors.New("tonic_error.logging.invalid_config_format.loggers.type")
		}

		name := loggerMap["name"].(string)
		if !ok {
			return errors.New("tonic_error.logging.invalid_config_format.loggers.name")
		}

		handlers, ok := loggerMap["handlers"].([]interface{})
		if !ok {
			return errors.New("tonic_error.logging.invalid_config_format.loggers.handlers")
		}

		level, ok := loggerMap["level"].(string)
		if !ok {
			return errors.New("tonic_error.logging.invalid_config_format.loggers.level")
		}

		l, err := getLogger(name, level, handlers)
		if err != nil {
			return err
		}

		Logging.Loggers[name] = l
	}

	_, ok = Logging.Loggers["default"]
	if !ok {
		return errors.New("tonic_error.logging.missing_default_logger")
	}

	return nil

}

func getFormatter(format string, color bool) (logrus.Formatter, error) {
	switch format {
	case "text":
		return &logrus.TextFormatter{
			ForceColors:     color,
			TimestampFormat: "2006-01-02T15:04:05.999Z07:00",
		}, nil

	}
	return nil, errors.New("tonic_error.logging.unsupported_formatter")
}

func getHandler(name string, formatter string) (*LogHandler, error) {
	f, ok := Logging.Formatters[formatter]
	if !ok {
		return nil, errors.New("tonic_error.logging.invalid_formatter")
	}
	switch name {
	case "console":
		return &LogHandler{
			Name:      name,
			Formatter: f,
			GetHook: func(loggerName string) (logrus.Hook, error) {
				return nil, errors.New("tonic_error.logging.abuse_console_handler")
			},
		}, nil
	case "kafka":
		return &LogHandler{
			Name:      name,
			Formatter: f,
			GetHook: func(loggerName string) (logrus.Hook, error) {
				topic := fmt.Sprintf("%s.%s", Logging.AppName, loggerName)
				kafkaHook, err := NewKafkaHook(topic, logrus.AllLevels, f)
				if err != nil {
					return nil, err
				}
				return kafkaHook, nil
			},
		}, nil
	}
	return nil, errors.New("tonic_error.logging.unsupported_handler")
}

func getLogger(name string, level string, handlers []interface{}) (*logrus.Logger, error) {
	logger := logrus.New()

	/* Add levels */
	l, err := strLevelConv(level)
	if err != nil {
		fmt.Println(err.Error())
	}
	logger.Level = l

	/* Add handlers */
	logger.Out = ioutil.Discard

	loggerHandlers := []*LogHandler{}

	for _, handler := range handlers {

		handlerStr, ok := handler.(string)
		if !ok {
			return nil, fmt.Errorf("tonic_error.logging.invalid_handler_type.%s", handler)
		}

		loggerHandler, ok := Logging.Handler[handlerStr]
		if !ok {
			return nil, fmt.Errorf("tonic_error.logging.invalid_handler.%s", handlerStr)
		}
		loggerHandlers = append(loggerHandlers, loggerHandler)

	}

	for _, loggerHandler := range loggerHandlers {

		if loggerHandler.Name == "console" {
			logger.Out = os.Stdout
			logger.Formatter = loggerHandler.Formatter
			continue
		}

		hook, err := loggerHandler.GetHook(name)
		if err != nil {
			return nil, err
		}

		logger.Hooks.Add(hook)
	}

	return logger, nil
}

func strLevelConv(level string) (logrus.Level, error) {
	switch level {
	case "DEBUG":
		return logrus.DebugLevel, nil
	case "INFO":
		return logrus.InfoLevel, nil
	case "WARN":
		return logrus.WarnLevel, nil
	case "ERROR":
		return logrus.ErrorLevel, nil
	default:
		return logrus.DebugLevel, errors.New("invalid_log_level")
	}
}
