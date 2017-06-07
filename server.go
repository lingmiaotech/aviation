package tonic

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type Server struct {
	App    *gin.Engine
	Port   int
	Routes []Route
}

func New() (server *Server, err error) {

	err = InitConfigs()
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
		app:    gin.Default(),
		Port:   8080,
		Routes: []Route{},
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
	err = s.app.Run(fmt.Sprintf(":%d", s.Port))
	return
}

func GetServerMode() string {
	env, ok := Configs.Get("env").(string)
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
