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

func Get(url string, opts *Opts) (*Response, error) {
	opts.fixJsonRequestEscapeIssue()

	shouldMock, mresponse := checkMock(url, "GET", opts)
	if shouldMock {
		return &Response{
			StatusCode:  mresponse.StatusCode,
			Header:      mresponse.Header,
			Body:        mresponse.Body,
			RawResponse: nil,
		}, nil
	}

	response, err := grequests.Get(url, (*grequests.RequestOptions)(opts))
	return &Response{
		StatusCode:  response.StatusCode,
		Header:      response.Header,
		Body:        response.Bytes(),
		RawResponse: response.RawResponse,
	}, err
}

func Put(url string, opts *Opts) (*Response, error) {
	opts.fixJsonRequestEscapeIssue()

	shouldMock, mresponse := checkMock(url, "PUT", opts)
	if shouldMock {
		return &Response{
			StatusCode:  mresponse.StatusCode,
			Header:      mresponse.Header,
			Body:        mresponse.Body,
			RawResponse: nil,
		}, nil
	}

	response, err := grequests.Put(url, (*grequests.RequestOptions)(opts))
	return &Response{
		StatusCode:  response.StatusCode,
		Header:      response.Header,
		Body:        response.Bytes(),
		RawResponse: response.RawResponse,
	}, err
}

func Patch(url string, opts *Opts) (*Response, error) {
	opts.fixJsonRequestEscapeIssue()

	shouldMock, mresponse := checkMock(url, "PATCH", opts)
	if shouldMock {
		return &Response{
			StatusCode:  mresponse.StatusCode,
			Header:      mresponse.Header,
			Body:        mresponse.Body,
			RawResponse: nil,
		}, nil
	}

	response, err := grequests.Patch(url, (*grequests.RequestOptions)(opts))
	return &Response{
		StatusCode:  response.StatusCode,
		Header:      response.Header,
		Body:        response.Bytes(),
		RawResponse: response.RawResponse,
	}, err
}

func Delete(url string, opts *Opts) (*Response, error) {
	opts.fixJsonRequestEscapeIssue()

	shouldMock, mresponse := checkMock(url, "DELETE", opts)
	if shouldMock {
		return &Response{
			StatusCode:  mresponse.StatusCode,
			Header:      mresponse.Header,
			Body:        mresponse.Body,
			RawResponse: nil,
		}, nil
	}

	response, err := grequests.Delete(url, (*grequests.RequestOptions)(opts))
	return &Response{
		StatusCode:  response.StatusCode,
		Header:      response.Header,
		Body:        response.Bytes(),
		RawResponse: response.RawResponse,
	}, err
}

func Post(url string, opts *Opts) (*Response, error) {
	opts.fixJsonRequestEscapeIssue()

	shouldMock, mresponse := checkMock(url, "POST", opts)
	if shouldMock {
		return &Response{
			StatusCode:  mresponse.StatusCode,
			Header:      mresponse.Header,
			Body:        mresponse.Body,
			RawResponse: nil,
		}, nil
	}

	response, err := grequests.Post(url, (*grequests.RequestOptions)(opts))
	return &Response{
		StatusCode:  response.StatusCode,
		Header:      response.Header,
		Body:        response.Bytes(),
		RawResponse: response.RawResponse,
	}, err
}

func Head(url string, opts *Opts) (*Response, error) {
	opts.fixJsonRequestEscapeIssue()

	shouldMock, mresponse := checkMock(url, "HEAD", opts)
	if shouldMock {
		return &Response{
			StatusCode:  mresponse.StatusCode,
			Header:      mresponse.Header,
			Body:        mresponse.Body,
			RawResponse: nil,
		}, nil
	}

	response, err := grequests.Head(url, (*grequests.RequestOptions)(opts))
	return &Response{
		StatusCode:  response.StatusCode,
		Header:      response.Header,
		Body:        response.Bytes(),
		RawResponse: response.RawResponse,
	}, err
}

func Options(url string, opts *Opts) (*Response, error) {
	opts.fixJsonRequestEscapeIssue()

	shouldMock, mresponse := checkMock(url, "OPTIONS", opts)
	if shouldMock {
		return &Response{
			StatusCode:  mresponse.StatusCode,
			Header:      mresponse.Header,
			Body:        mresponse.Body,
			RawResponse: nil,
		}, nil
	}

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

	switch v := opts.JSON.(type) {

	case string:
		opts.JSON = []byte(v)

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
	}
	return nil

}
