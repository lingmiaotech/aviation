package statsd

import (
	"fmt"
	"github.com/lingmiaotech/tonic/configs"
	"gopkg.in/alexcesaro/statsd.v2"
	"os"
	"time"
)

type InstanceClass struct {
	AppName string
	Enabled bool
	Client  *statsd.Client
}

type Timer struct {
	start time.Time
}

var Instance InstanceClass

func Increment(bucket string) {
	b := getBucket(bucket)
	if !Instance.Enabled {
		fmt.Println(fmt.Sprintf("[STATSD] key=%s count=1", b))
		return
	}
	Instance.Client.Increment(b)
}

// Timing takes bucket name and delta in milliseconds
func Timing(bucket string, delta int) {
	b := getBucket(bucket)
	if !Instance.Enabled {
		fmt.Println(fmt.Sprintf("[STATSD] key=%s time_delta=%d(ms)", b, delta))
		return
	}
	Instance.Client.Timing(b, delta)
}

func Count(bucket string, n int) {
	b := getBucket(bucket)
	if !Instance.Enabled {
		fmt.Println(fmt.Sprintf("[STATSD] key=%v count=%d", b, n))
		return
	}
	Instance.Client.Count(b, n)
}

func Gauge(bucket string, n int) {
	b := getBucket(bucket)
	if !Instance.Enabled {
		fmt.Println(fmt.Sprintf("[STATSD] key=%v gauge=%d", b, n))
		return
	}
	Instance.Client.Gauge(b, n)
}

func getBucket(bucket string) string {
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "UNKNOWNHOST"
	}
	return fmt.Sprintf("%v.%v.%v", hostName, Instance.AppName, bucket)
}

func NewTimer() Timer {
	return Timer{start: time.Now()}
}

// Send sends the time elapsed since the creation of the Timing.
func (t Timer) Send(bucket string) {
	Timing(bucket, int(t.Duration()/time.Millisecond))
}

// Duration returns the time elapsed since the creation of the Timing.
func (t Timer) Duration() time.Duration {
	return time.Now().Sub(t.start)
}

func InitStatsd() error {

	Instance.AppName = configs.GetString("app_name")
	Instance.Enabled = configs.GetBool("statsd.enabled")

	if !Instance.Enabled {
		return nil
	}

	host := configs.GetString("statsd.host")
	port := configs.GetString("statsd.port")
	address := fmt.Sprintf("%v:%v", host, port)
	c, err := statsd.New(statsd.Address(address))
	if err != nil || c == nil {
		return err
	}

	Instance.Client = c
	return nil

}
