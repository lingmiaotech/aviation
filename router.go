package tonic

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/dyliu/tonic/statsd"
)

func (s *Server) InitRoutes() error {

	app, ok := (s.App).(*gin.Engine)
	if !ok {
		return errors.New("invalid_app_engine")
	}
	app.Use(RequestHandler)
	return nil
}

//RequestHandler handles each request and make some records
func RequestHandler(c *gin.Context) {

	timer := statsd.NewTimer()

	c.Next()

	statsd.Increment(getCountBucket(c))

	timer.Send(getTimingBucket(c))

}

func getCountBucket(c *gin.Context) string {
	return fmt.Sprintf(
		"views.%s.%s.status_code.%d",
		strings.Trim(strings.Join(strings.Split(c.Request.RequestURI, "/"), "."), "."),
		c.Request.Method,
		c.Writer.Status(),
	)
}

func getTimingBucket(c *gin.Context) string {
	return fmt.Sprintf(
		"views.%s.%s",
		strings.Trim(strings.Join(strings.Split(c.Request.RequestURI, "/"), "."), "."),
		c.Request.Method,
	)
}
