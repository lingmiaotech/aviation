package jaeger

import (
	"fmt"
	"strings"
	"time"

	"github.com/lingmiaotech/tonic/configs"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/transport"
)

func InitJaeger() {

	backendHostPort := fmt.Sprintf("%s:%s", configs.GetString("backendHostPort.host"), configs.GetString("backendHostPort.port"))

	var err error
	var sender jaeger.Transport
	if strings.HasPrefix(backendHostPort, "http://") {
		sender = transport.NewHTTPTransport(
			backendHostPort,
			transport.HTTPBatchSize(1),
		)
	} else {
		sender, err = jaeger.NewUDPTransport(backendHostPort, 0)
		if err != nil {
			panic(fmt.Sprintf("ERROR: cannot initialize UDP sender: %v\n", err))
		}
	}

	var sampler config.SamplerConfig
	env := getServerMode()
	switch env {
	case "production":
		sampler.Type = "const"
		sampler.Param = 1
	default:
		sampler.Type = "const"
		sampler.Param = 1
	}

	cfg := &config.Configuration{
		ServiceName: configs.GetString("app_name"),
		Sampler: &sampler,
	}
	tracer, _, err := cfg.NewTracer(
		config.Reporter(jaeger.NewRemoteReporter(
			sender,
			jaeger.ReporterOptions.BufferFlushInterval(1*time.Second),
		)))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	opentracing.SetGlobalTracer(tracer)
	return
}

func getServerMode() string {
	env, ok := configs.Get("env").(string)
	if !ok {
		logging.GetDefaultLogger().Warn("Cannot find env setting in config, will use development mode.")
	}
	return env
}
