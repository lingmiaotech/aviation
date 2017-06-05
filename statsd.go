package tonic

import (
	"fmt"
	"gopkg.in/alexcesaro/statsd.v2"
	"os"
	"time"
)

type StatsdClass struct {
	AppName string
	Enabled bool
	Client  *statsd.Client
}

type Timer struct {
	start  time.Time
	client *StatsdClass
}

var Statsd StatsdClass

func (s *StatsdClass) Increment(bucket string) {
	b := s.getBucket(bucket)
	if !s.Enabled {
		fmt.Println(fmt.Sprintf("[STATS] key=%s count=1", b))
		return
	}
	s.Client.Increment(b)
}

// Timing takes bucket name and delta in milliseconds
func (s *StatsdClass) Timing(bucket string, delta int) {
	b := s.getBucket(bucket)
	if !s.Enabled {
		fmt.Println(fmt.Sprintf("[STATS] key=%s time_delta=%d(ms)", b, delta))
		return
	}
	s.Client.Timing(b, delta)
}

func (s *StatsdClass) Count(bucket string, n int) {
	b := s.getBucket(bucket)
	if !s.Enabled {
		fmt.Println(fmt.Sprintf("[STATS] key=%v count=%d", b, n))
		return
	}
	s.Client.Count(b, n)
}

func (s *StatsdClass) Gauge(bucket string, n int) {
	b := s.getBucket(bucket)
	if !s.Enabled {
		fmt.Println(fmt.Sprintf("[STATS] key=%v gauge=%d", b, n))
		return
	}
	s.Client.Gauge(b, n)
}

func (s *StatsdClass) getBucket(bucket string) string {
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "UNKNOWNHOST"
	}
	return fmt.Sprintf("%v.%v.%v", hostName, s.AppName, bucket)
}

func (s *StatsdClass) NewTimer() Timer {
	return Timer{start: time.Now(), client: s}
}

// Send sends the time elapsed since the creation of the Timing.
func (t Timer) Send(bucket string) {
	t.client.Timing(bucket, int(t.Duration()/time.Millisecond))
}

// Duration returns the time elapsed since the creation of the Timing.
func (t Timer) Duration() time.Duration {
	return time.Now().Sub(t.start)
}

func InitStatsd() error {

	Statsd.AppName = Configs.GetString("app_name")
	Statsd.Enabled = Configs.GetBool("statsd.enabled")

	if !Statsd.Enabled {
		return nil
	}

	host := Configs.GetString("statsd.host")
	port := Configs.GetString("statsd.port")
	address := fmt.Sprintf("%v:%v", host, port)
	c, err := statsd.New(statsd.Address(address))
	if err != nil || c == nil {
		return err
	}

	Statsd.Client = c
	return nil

}
