package sentry

import (
	"fmt"
	"strings"

	"github.com/getsentry/raven-go"
	"github.com/dyliu/tonic/configs"
	"github.com/dyliu/tonic/logging"
)

type Sender interface {
	CaptureError(err error, params map[string]interface{})
	CaptureMessage(msg string, params map[string]interface{})
}

type DefaultSender struct {
	Enabled bool
	Dsn     string
	Client  *raven.Client
}

var S Sender
var DS *DefaultSender

type Extra struct {
	Data interface{} `json:"data"`
}

func (i Extra) Class() string { return "extra" }

// InitSentry : Initialize sentry DSN while sentry config is enabled
func InitSentry() (err error) {

	DS = new(DefaultSender)
	S = DS

	DS.Enabled = configs.GetBool("sentry.enabled")
	DS.Dsn = configs.GetString("sentry.dsn")

	if !DS.Enabled {
		return nil
	}

	DS.Client, err = raven.New(DS.Dsn)
	if err != nil {
		return
	}

	return
}

func (s *DefaultSender) CaptureError(err error, params map[string]interface{}) {
	if !s.Enabled {
		logging.GetDefaultLogger().Errorf("[SENTRY] error=%s , params=%v\n", err, printParams(params))
		return
	}
	s.Client.CaptureError(err, nil, Extra{params})
}

func (s *DefaultSender) CaptureMessage(msg string, params map[string]interface{}) {
	if !s.Enabled {
		logging.GetDefaultLogger().Errorf("[SENTRY] error=%s, params=%v\n", msg, printParams(params))
		return
	}
	s.Client.CaptureMessage(msg, nil, Extra{params})
}

func printParams(params map[string]interface{}) string {
	results := []string{}
	for key, value := range params {
		results = append(results, fmt.Sprintf("%s=%v", key, value))
	}
	return strings.Join(results, ", ")
}

func CaptureError(err error, params map[string]interface{}) {
	S.CaptureError(err, params)
}

func CaptureMessage(msg string, params map[string]interface{}) {
	S.CaptureMessage(msg, params)
}
