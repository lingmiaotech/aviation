package tonic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lingmiaotech/tonic/configs"
)

type Server struct {
	App  *gin.Engine
	Port int
}

func New() (server *Server, err error) {

	err = configs.InitConfigs()
	if err != nil {
		return
	}

	gin.SetMode(GetServerMode())

	err = InitKafka()
	if err != nil {
		return
	}

	err = InitLogging()
	if err != nil {
		return
	}

	err = InitStatsd()
	if err != nil {
		return
	}

	err = InitSentry()
	if err != nil {
		return
	}

	err = InitRedis()
	if err != nil {
		return
	}

	err = InitDatabase()
	if err != nil {
		return
	}

	server = &Server{
		App:  gin.Default(),
		Port: 8080,
	}

	err = server.InitRoutes()
	if err != nil {
		return
	}

	return
}

func (s *Server) SetPort(p int) {
	s.Port = p
}

func (s *Server) Start() (err error) {
	err = s.App.Run(fmt.Sprintf(":%d", s.Port))
	return
}

func GetServerMode() string {
	env, ok := configs.Get("env").(string)
	if !ok {
		Logging.GetDefaultLogger().Warn("Cannot find env setting in config, will use development mode.")
	}
	switch env {
	case "production":
		return gin.ReleaseMode
	case "test":
		return gin.TestMode
	}
	return gin.DebugMode
}
