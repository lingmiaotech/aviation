package tonic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lingmiaotech/tonic/statsd"
	"strings"
)

func (s *Server) InitRoutes() error {

	s.App.Use(RequestHandler)

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
