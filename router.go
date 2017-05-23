package aviation

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

type Route struct {
	pattern string
	method  string
	handler func(c *gin.Context)
}

func (s *Server) AddRoute(pattern string, method string, handler func(c *gin.Context)) {
	s.Routes = append(s.Routes, Route{pattern: pattern, method: method, handler: handler})
}

func (s *Server) InitRoutes() error {

	s.app.Use(RequestHandler)

	for _, route := range s.Routes {
		switch route.method {
		case "GET":
			s.app.GET(route.pattern, route.handler)
		case "POST":
			s.app.POST(route.pattern, route.handler)
		case "PUT":
			s.app.PUT(route.pattern, route.handler)
		case "OPTIONS":
			s.app.OPTIONS(route.pattern, route.handler)
		case "DELETE":
			s.app.DELETE(route.pattern, route.handler)
		default:
			lowerMethod := strings.ToLower(route.method)
			return fmt.Errorf("aviation_error.routes.invalid_method.%s", lowerMethod)

		}
	}
	return nil
}

//RequestHandler handles each request and make some records
func RequestHandler(c *gin.Context) {

	timer := Statsd.NewTimer()

	c.Next()

	Statsd.Increment(getCountBucket(c))

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
