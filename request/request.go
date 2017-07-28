package request

import (
	"bytes"
	"encoding/json"
	"github.com/levigross/grequests"
	"net/http"
)

type Opts grequests.RequestOptions

type Response struct {
	StatusCode  int
	Header      http.Header
	Body        []byte
	RawResponse *http.Response
}

type Sender interface {
	Get(url string, opts *Opts) (*Response, error)
	Put(url string, opts *Opts) (*Response, error)
	Patch(url string, opts *Opts) (*Response, error)
	Delete(url string, opts *Opts) (*Response, error)
	Post(url string, opts *Opts) (*Response, error)
	Head(url string, opts *Opts) (*Response, error)
	Options(url string, opts *Opts) (*Response, error)
}

type DefaultSender struct{}

var DS *DefaultSender = new(DefaultSender)
var S Sender = DS

func (s *DefaultSender) Get(url string, opts *Opts) (*Response, error) {
	opts.fixJsonRequestEscapeIssue()
	response, err := grequests.Get(url, (*grequests.RequestOptions)(opts))
	return &Response{
		StatusCode:  response.StatusCode,
		Header:      response.Header,
		Body:        response.Bytes(),
		RawResponse: response.RawResponse,
	}, err
}

func (s *DefaultSender) Put(url string, opts *Opts) (*Response, error) {
	opts.fixJsonRequestEscapeIssue()
	response, err := grequests.Put(url, (*grequests.RequestOptions)(opts))
	return &Response{
		StatusCode:  response.StatusCode,
		Header:      response.Header,
		Body:        response.Bytes(),
		RawResponse: response.RawResponse,
	}, err
}

func (s *DefaultSender) Patch(url string, opts *Opts) (*Response, error) {
	opts.fixJsonRequestEscapeIssue()
	response, err := grequests.Patch(url, (*grequests.RequestOptions)(opts))
	return &Response{
		StatusCode:  response.StatusCode,
		Header:      response.Header,
		Body:        response.Bytes(),
		RawResponse: response.RawResponse,
	}, err
}

func (s *DefaultSender) Delete(url string, opts *Opts) (*Response, error) {
	opts.fixJsonRequestEscapeIssue()
	response, err := grequests.Delete(url, (*grequests.RequestOptions)(opts))
	return &Response{
		StatusCode:  response.StatusCode,
		Header:      response.Header,
		Body:        response.Bytes(),
		RawResponse: response.RawResponse,
	}, err
}

func (s *DefaultSender) Post(url string, opts *Opts) (*Response, error) {
	opts.fixJsonRequestEscapeIssue()
	response, err := grequests.Post(url, (*grequests.RequestOptions)(opts))
	return &Response{
		StatusCode:  response.StatusCode,
		Header:      response.Header,
		Body:        response.Bytes(),
		RawResponse: response.RawResponse,
	}, err
}

func (s *DefaultSender) Head(url string, opts *Opts) (*Response, error) {
	opts.fixJsonRequestEscapeIssue()
	response, err := grequests.Head(url, (*grequests.RequestOptions)(opts))
	return &Response{
		StatusCode:  response.StatusCode,
		Header:      response.Header,
		Body:        response.Bytes(),
		RawResponse: response.RawResponse,
	}, err
}

func (s *DefaultSender) Options(url string, opts *Opts) (*Response, error) {
	opts.fixJsonRequestEscapeIssue()
	response, err := grequests.Options(url, (*grequests.RequestOptions)(opts))
	return &Response{
		StatusCode:  response.StatusCode,
		Header:      response.Header,
		Body:        response.Bytes(),
		RawResponse: response.RawResponse,
	}, err
}

func (r *Response) JSON(userStruct interface{}) error {
	err := json.Unmarshal(r.Body, &userStruct)
	if err != nil {
		return err
	}
	return nil

}

func (r *Response) String() string {
	return string(r.Body)
}

func (opts *Opts) fixJsonRequestEscapeIssue() error {
	if opts.JSON == nil {
		return nil
	}
	switch v := opts.JSON.(type) {
	case string:
		opts.JSON = []byte(v)
		return nil
	case []byte:
		return nil
	default:
		buffer := bytes.Buffer{}
		encoder := json.NewEncoder(&buffer)
		encoder.SetEscapeHTML(false)
		err := encoder.Encode(v)
		if err != nil {
			return err
		}
		opts.JSON = buffer.Bytes()
		return nil
	}
}

func Get(url string, opts *Opts) (*Response, error) {
	return S.Get(url, opts)
}

func Put(url string, opts *Opts) (*Response, error) {
	return S.Put(url, opts)
}

func Patch(url string, opts *Opts) (*Response, error) {
	return S.Patch(url, opts)
}

func Delete(url string, opts *Opts) (*Response, error) {
	return S.Delete(url, opts)
}

func Post(url string, opts *Opts) (*Response, error) {
	return S.Post(url, opts)
}

func Head(url string, opts *Opts) (*Response, error) {
	return S.Head(url, opts)
}

func Options(url string, opts *Opts) (*Response, error) {
	return S.Options(url, opts)
}
