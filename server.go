package tonic

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type Server struct {
	app    *gin.Engine
	Port   int
	Routes []Route
}

func New() *Server {
	return &Server{
		app:    gin.Default(),
		Port:   8080,
		Routes: []Route{},
	}
}

func (s *Server) SetPort(p int) {
	s.Port = p
}

func (s *Server) Start() (err error) {

	err = InitConfigs()
	if err != nil {
		return
	}

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

	err = InitRedis()
	if err != nil {
		return
	}

	err = InitDatabase()
	if err != nil {
		return
	}

	err = s.InitRoutes()
	if err != nil {
		return
	}

	err = s.app.Run(fmt.Sprintf(":%d", s.Port))
	return
}
