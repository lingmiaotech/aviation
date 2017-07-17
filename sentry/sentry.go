package sentry

import (
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/lingmiaotech/tonic/configs"
	"github.com/lingmiaotech/tonic/logging"
	"strings"
)

type InstanceClass struct {
	Enabled bool
	Dsn     string
	Client  *raven.Client
}

var Instance InstanceClass

type Extra struct {
	Data interface{} `json:"data"`
}

func (i Extra) Class() string { return "extra" }

// InitSentry : Initialize sentry DSN while sentry config is enabled
func InitSentry() (err error) {

	Instance.Enabled = configs.GetBool("sentry.enabled")
	Instance.Dsn = configs.GetString("sentry.dsn")

	if !Instance.Enabled {
		return nil
	}

	Instance.Client, err = raven.New(Instance.Dsn)
	if err != nil {
		return
	}

	return
}

// CaptureError : Capture Error and deliver an error to the Sentry server
func CaptureError(err error, params map[string]interface{}) {
	if !Instance.Enabled {
		logging.GetDefaultLogger().Infof("[SENTRY] error=%s , params=%v\n", err, printParams(params))
		return
	}
	Instance.Client.CaptureError(err, nil, Extra{params})
}

// CaptureMessage : Capture message and additional parametric, the deliver a string message to the Sentry server
func CaptureMessage(msg string, params map[string]interface{}) {
	if !Instance.Enabled {
		logging.GetDefaultLogger().Infof("[SENTRY] error=%s, params=%v\n", msg, printParams(params))
		return
	}
	Instance.Client.CaptureMessage(msg, nil, Extra{params})
}

func printParams(params map[string]interface{}) string {
	results := []string{}
	for key, value := range params {
		results = append(results, fmt.Sprintf("%s=%v", key, value))
	}
	return strings.Join(results, ", ")
}
