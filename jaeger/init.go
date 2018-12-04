package jaeger

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"time"

	"github.com/dyliu/tonic/configs"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/transport"
)

func Initialize() {

	var sampleStrategy string
	if sampleStrategy = os.Getenv("SAMPLE_STRATEGY"); sampleStrategy == ""{
		sampleStrategy = "probabilistic,0.2"
	}
	args := strings.Split(sampleStrategy, ",")
	if len(args) != 2 {
		panic(fmt.Sprintf("ERROR: invalid SAMPLE_STRATEGY format %s", sampleStrategy))
	}
	samplerType := args[0]
	param, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		panic(fmt.Sprintf("ERROR: invalid SAMPLE_STRATEGY param %s, param must be string of float64", sampleStrategy))
	}

	cfg := &config.Configuration{
		ServiceName: configs.GetString("app_name"),
		Sampler: &config.SamplerConfig{
			Type:  samplerType,
			Param: param,
		},
	}

	var agentURI string
	if agentURI = os.Getenv("JAEGER_AGENT"); agentURI == "" {
		agentURI = fmt.Sprintf("http://jaeger-agent.%s:6831", configs.GetString("app_name"))
	}

	var backendHostPort string
	if env, _ := configs.Get("env").(string); env == "production" {
		backendHostPort = agentURI
	}else {
		backendHostPort = "127.0.0.1:6831"
	}

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
