package tonic

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CrowdSurge/banner"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/dyliu/tonic/configs"
	"github.com/dyliu/tonic/database"
	"github.com/dyliu/tonic/jaeger"
	"github.com/dyliu/tonic/kafka"
	"github.com/dyliu/tonic/logging"
	"github.com/dyliu/tonic/redis"
	"github.com/dyliu/tonic/sentry"
	"github.com/dyliu/tonic/statsd"
)

type Server struct {
	App  interface{}
	Port int
}

func New() (*Server, error) {

	var server *Server
	var err error

	err = configs.InitConfigs()
	if err != nil {
		return server, err
	}

	gin.SetMode(GetServerMode())

	err = kafka.InitKafka()
	if err != nil {
		return nil, err
	}

	err = logging.InitLogging()
	if err != nil {
		return nil, err
	}

	err = statsd.InitStatsd()
	if err != nil {
		return nil, err
	}

	err = sentry.InitSentry()
	if err != nil {
		return nil, err
	}

	err = redis.InitRedis()
	if err != nil {
		return nil, err
	}

	err = database.InitDatabase()
	if err != nil {
		return nil, err
	}

	server = &Server{
		App:  gin.New(),
		Port: 8080,
	}
	InitMiddlewares(server.App.(*gin.Engine))

	err = server.InitRoutes()
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (s *Server) SetPort(p int) {
	s.Port = p
}

func (s *Server) Start() error {
	banner.Print("cheers")

	//app, ok := (s.App).(*gin.Engine)
	//if !ok {
	//	return errors.New("invalid_app_engine")
	//}

	//err := app.Run(fmt.Sprintf(":%d", s.Port))
	//if err != nil {
	//	return err
	//}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.Port),
		Handler: (s.App).(*gin.Engine),
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 60 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM,syscall.SIGINT)
	<-quit
	fmt.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("Server Shutdown:", err)
	}
	fmt.Println("Server exiting")

	return nil
}

func GetServerMode() string {
	env, ok := configs.Get("env").(string)
	if !ok {
		logging.GetDefaultLogger().Warn("Cannot find env setting in config, will use development mode.")
	}
	switch env {
	case "production":
		return gin.ReleaseMode
	case "test":
		return gin.TestMode
	}
	return gin.DebugMode
}

//func InitJaeger() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		jaeger.Initialize()
//		c.Next()
//	}
//}

func InitJaegerSpan() gin.HandlerFunc {
	return func(c *gin.Context) {
		jaeger.Initialize()
		tracer := opentracing.GlobalTracer()
		var span opentracing.Span
		spanContext, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil {
			span = tracer.StartSpan("HTTP " + c.Request.Method)
		} else {
			span = tracer.StartSpan("HTTP "+c.Request.Method, ext.RPCServerOption(spanContext))
		}
		defer span.Finish()
		ext.HTTPMethod.Set(span, c.Request.Method)
		ext.HTTPUrl.Set(span, c.Request.URL.String())
		ext.Component.Set(span, "net/http")
		c.Request = c.Request.WithContext(opentracing.ContextWithSpan(c.Request.Context(), span))
		c.Next()
	}
}

func InitMiddlewares(app *gin.Engine) {
	env, _ := configs.Get("env").(string)

	switch env {
	case "test":
		app.Use(gin.LoggerWithWriter(ioutil.Discard), gin.Recovery())
	default:
		app.Use(gin.Logger(), gin.Recovery(),InitJaegerSpan())
	}
}
