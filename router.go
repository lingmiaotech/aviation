package tonic

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/lingmiaotech/tonic/prom"
)

func (s *Server) InitRoutes() error {

	app, ok := (s.App).(*gin.Engine)
	if !ok {
		return errors.New("invalid_app_engine")
	}
	app.GET("/metrics", MetricsHandler)
	app.Use(RequestHandler)
	return nil
}

func MetricsHandler(c *gin.Context) {
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}

//RequestHandler handles each request and make some records
func RequestHandler(c *gin.Context) {

	timer := prom.NewTimer()

	c.Next()

	prom.Increment(getCountBucket(c))

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
