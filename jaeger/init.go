package jaeger

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"time"

	"github.com/lingmiaotech/tonic/configs"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/transport"
)

func InitJaeger() {

	agentURI := os.Getenv("JAEGER_AGENT")
	if agentURI == "" {
		agentURI = fmt.Sprintf("http://jaeger-agent.%s:6831", configs.GetString("app_name"))
	}

	var backendHostPort string
	env, _ := configs.Get("env").(string)
	if env == "development" {
		backendHostPort = "127.0.0.1:6831"
	} else {
		backendHostPort = agentURI
	}

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

	sampleType := os.Getenv("JAEGER_SAMPLE_TYPE")
	if sampleType == "" {
		sampleType = "probabilistic"
	}
	sampleParamString := os.Getenv("JAEGER_SAMPLE_PARAM")
	if sampleParamString == "" {
		sampleParamString = "0.5"
	}
	sampleParam, _ := strconv.ParseFloat(sampleParamString, 64)

	cfg := &config.Configuration{
		ServiceName: configs.GetString("app_name"),
		Sampler: &config.SamplerConfig{
			Type:  sampleType,
			Param: sampleParam,
		},
	}
	tracer, _, err := cfg.NewTracer(
		config.Reporter(jaeger.NewRemoteReporter(
			sender,
			jaeger.ReporterOptions.BufferFlushInterval(1*time.Second),
		)))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	opentracing.InitGlobalTracer(tracer)
	return
}
