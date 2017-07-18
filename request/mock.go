package request

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
)

type AnythingClass struct{}

var Anything AnythingClass

type MockCondition struct {
	Method   interface{}
	Url      interface{}
	Checkers [](func(*Opts) bool)
}

type MockResponse struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Pass       bool
}

type Mock struct {
	Condition MockCondition
	Response  MockResponse
}

type MockConfigs struct {
	Enabled bool
	Mocks   []*Mock
}

var mc MockConfigs
var mapLock sync.Mutex

func EnableMock() {
	mc.Enabled = true
}

func DisableMock() {
	mc.Enabled = false
}

func AddMock() *Mock {
	mapLock.Lock()
	m := &Mock{
		Condition: MockCondition{
			Method:   Anything,
			Url:      Anything,
			Checkers: make([](func(*Opts) bool), 0),
		},
	}
	mc.Mocks = append(mc.Mocks, m)
	mapLock.Unlock()
	return m
}

func (m *Mock) OnUrl(url string) *Mock {
	mapLock.Lock()
	m.Condition.Url = url
	mapLock.Unlock()
	return m
}

func (m *Mock) OnMethod(method string) *Mock {
	mapLock.Lock()
	m.Condition.Method = method
	mapLock.Unlock()
	return m
}

func (m *Mock) OnChecker(c func(*Opts) bool) *Mock {
	mapLock.Lock()
	m.Condition.Checkers = append(m.Condition.Checkers, c)
	mapLock.Unlock()
	return m
}

func (m *Mock) AndReturn(statusCode int, header http.Header, body []byte) {
	mapLock.Lock()
	m.Response.StatusCode = statusCode
	m.Response.Header = header
	m.Response.Body = body
	mapLock.Unlock()
}

func (m *Mock) AndPass() {
	mapLock.Lock()
	m.Response.Pass = true
	mapLock.Unlock()
}

func ClearMock() {
	mapLock.Lock()
	mc.Mocks = make([]*Mock, 0)
	mapLock.Unlock()
}

func checkMock(url string, method string, opts *Opts) (bool, *MockResponse) {
	if !mc.Enabled {
		return false, nil
	}
	for _, mock := range mc.Mocks {
		if !checkMockCondition(url, method, opts, mock.Condition) {
			continue
		}
		if mock.Response.Pass {
			return false, nil
		}
		return true, &mock.Response
	}
	pm := fmt.Sprintf("Detect unmocked request!\nUrl=%s\nMethod=%s\nJson=%s", url, method, string(opts.JSON.([]byte)))
	panic(errors.New(pm))
}

func checkMockCondition(url string, method string, opts *Opts, condition MockCondition) bool {
	if condition.Url != Anything && condition.Url.(string) != url {
		return false
	}
	if condition.Method != Anything && condition.Method.(string) != method {
		return false
	}
	for _, c := range condition.Checkers {
		if !c(opts) {
			return false
		}
	}
	return true
}
