package tonic

import (
	"fmt"
	"github.com/CrowdSurge/banner.git"
	"github.com/gin-gonic/gin"
	"github.com/lingmiaotech/tonic/configs"
	"github.com/lingmiaotech/tonic/database"
	"github.com/lingmiaotech/tonic/kafka"
	"github.com/lingmiaotech/tonic/logging"
	"github.com/lingmiaotech/tonic/redis"
	"github.com/lingmiaotech/tonic/sentry"
	"github.com/lingmiaotech/tonic/statsd"
)

type Server struct {
	App  *gin.Engine
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
		App:  gin.Default(),
		Port: 8080,
	}

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
	banner.Print("CHEERS!")

	err := s.App.Run(fmt.Sprintf(":%d", s.Port))
	if err != nil {
		return err
	}

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
