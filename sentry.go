package tonic

import (
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/lingmiaotech/tonic/configs"
	"strings"
)

type SentryClass struct {
	Enabled bool
	Dsn     string
	Client  *raven.Client
}

var Sentry SentryClass

type sentryExtra struct {
	Data interface{} `json:"data"`
}

func (i sentryExtra) Class() string { return "extra" }

// InitSentry : Initialize sentry DSN while sentry config is enabled
func InitSentry() (err error) {

	Sentry.Enabled = configs.GetBool("sentry.enabled")
	Sentry.Dsn = configs.GetString("sentry.dsn")

	if !Sentry.Enabled {
		return nil
	}

	Sentry.Client, err = raven.New(Sentry.Dsn)
	if err != nil {
		return
	}

	return
}

// CaptureError : Capture Error and deliver an error to the Sentry server
func (s *SentryClass) CaptureError(err error, params map[string]interface{}) {
	if !Sentry.Enabled {
		Logging.GetDefaultLogger().Infof("[SENTRY] error=%s , params=%v\n", err, printParams(params))
		return
	}
	Sentry.Client.CaptureError(err, nil, sentryExtra{params})
}

// CaptureMessage : Capture message and additional parametric, the deliver a string message to the Sentry server
func (s *SentryClass) CaptureMessage(msg string, params map[string]interface{}) {
	if !Sentry.Enabled {
		Logging.GetDefaultLogger().Infof("[SENTRY] error=%s, params=%v\n", msg, printParams(params))
		return
	}
	Sentry.Client.CaptureMessage(msg, nil, sentryExtra{params})
}

func printParams(params map[string]interface{}) string {
	results := []string{}
	for key, value := range params {
		results = append(results, fmt.Sprintf("%s=%v", key, value))
	}
	return strings.Join(results, ", ")
}
