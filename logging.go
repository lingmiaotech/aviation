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
	Hook      logrus.Hook
	Formatter logrus.Formatter
}

type LoggingClass struct {
	AppName    string
	Loggers    map[string]*logrus.Logger
	Handler    map[string]*LogHandler
	Formatters map[string]logrus.Formatter
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

	Statsd.AppName = Configs.GetString("app_name")

	formatters, ok := Configs.Get("logging.formatters").([]map[string]interface{})
	if !ok {
		return errors.New("tonic_error.log.invalid_config_format")
	}

	for _, formatter := range formatters {
		name, ok := formatter["name"].(string)
		if !ok {
			return errors.New("tonic_error.log.invalid_config_format")
		}

		format, ok := formatter["format"].(string)
		if !ok {
			return errors.New("tonic_error.log.invalid_config_format")
		}

		color, ok := formatter["color"].(bool)
		if !ok {
			color = false
		}

		f, err := getFormatter(format, color)
		if err != nil {
			return err
		}

		Logging.Formatters[name] = f
	}

	handlers, ok := Configs.Get("logging.handlers").([]map[string]interface{})
	if !ok {
		return errors.New("tonic_error.log.invalid_config_format")
	}

	for _, handler := range handlers {
		name, ok := handler["name"].(string)
		if !ok {
			return errors.New("tonic_error.log.invalid_config_format")
		}

		formatter, ok := handler["formatter"].(string)
		if !ok {
			return errors.New("tonic_error.log.invalid_config_format")
		}

		h, err := getHandler(name, formatter)
		if err != nil {
			return err
		}

		Logging.Handler[name] = h
	}

	loggers, ok := Configs.Get("logging.loggers").([]map[string]interface{})
	if !ok {
		return errors.New("tonic_error.log.invalid_config_format")
	}

	for _, logger := range loggers {
		name := logger["name"].(string)
		if !ok {
			return errors.New("tonic_error.log.invalid_config_format")
		}

		handlers, ok := logger["handlers"].([]string)
		if !ok {
			return errors.New("tonic_error.log.invalid_config_format")
		}

		level, ok := logger["level"].(string)
		if !ok {
			return errors.New("tonic_error.log.invalid_config_format")
		}

		l, err := getLogger(name, level, handlers)
		if err != nil {
			return err
		}

		Logging.Loggers[name] = l
	}

	_, ok = Logging.Loggers["default"]
	if !ok {
		return errors.New("tonic_error.log.missing_default_logger")
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
	return nil, errors.New("tonic_error.log.unsupported_formatter")
}

func getHandler(name string, formatter string) (*LogHandler, error) {
	f, ok := Logging.Formatters[name]
	if !ok {
		return nil, errors.New("tonic_error.log.invalid_formatter")
	}
	switch name {
	case "console":
		return &LogHandler{Name: name, Formatter: f}, nil
	case "kafka":
		topic := fmt.Sprintf("%s.%s", Logging.AppName, name)
		kafkaHook, err := NewKafkaHook(topic, logrus.AllLevels, f)
		if err != nil {
			return nil, err
		}
		return &LogHandler{Name: name, Hook: kafkaHook, Formatter: f}, nil
	}
	return nil, errors.New("tonic_error.log.unsupported_handler")
}

func getLogger(name string, level string, handlers []string) (*logrus.Logger, error) {
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

		loggerHandler, ok := Logging.Handler[handler]
		if !ok {
			return nil, fmt.Errorf("tonic_error.log.invalid_handler.%s", handler)
		}
		loggerHandlers = append(loggerHandlers, loggerHandler)

	}

	for _, loggerHandler := range loggerHandlers {

		if loggerHandler.Name == "console" {
			logger.Out = os.Stdout
			logger.Formatter = loggerHandler.Formatter
			continue
		}

		logger.Hooks.Add(loggerHandler.Hook)
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
