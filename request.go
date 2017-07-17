package tonic

import (
	"bytes"
	"encoding/json"
	"github.com/levigross/grequests"
	"net/http"
)

type RequestOptions grequests.RequestOptions

type Response struct {
	StatusCode  int
	Header      http.Header
	Body        []byte
	RawResponse *http.Response
}

func Get(url string, ro *RequestOptions) (*Response, error) {
	ro.fixJsonRequestEscapeIssue()
	response, err := grequests.Get(url, (*grequests.RequestOptions)(ro))
	return &Response{
		StatusCode:  response.StatusCode,
		Header:      response.Header,
		Body:        response.Bytes(),
		RawResponse: response.RawResponse,
	}, err
}

func Put(url string, ro *RequestOptions) (*Response, error) {
	ro.fixJsonRequestEscapeIssue()
	response, err := grequests.Put(url, (*grequests.RequestOptions)(ro))
	return &Response{
		StatusCode:  response.StatusCode,
		Header:      response.Header,
		Body:        response.Bytes(),
		RawResponse: response.RawResponse,
	}, err
}

func Patch(url string, ro *RequestOptions) (*Response, error) {
	ro.fixJsonRequestEscapeIssue()
	response, err := grequests.Patch(url, (*grequests.RequestOptions)(ro))
	return &Response{
		StatusCode:  response.StatusCode,
		Header:      response.Header,
		Body:        response.Bytes(),
		RawResponse: response.RawResponse,
	}, err
}

func Delete(url string, ro *RequestOptions) (*Response, error) {
	ro.fixJsonRequestEscapeIssue()
	response, err := grequests.Delete(url, (*grequests.RequestOptions)(ro))
	return &Response{
		StatusCode:  response.StatusCode,
		Header:      response.Header,
		Body:        response.Bytes(),
		RawResponse: response.RawResponse,
	}, err
}

func Post(url string, ro *RequestOptions) (*Response, error) {
	ro.fixJsonRequestEscapeIssue()
	response, err := grequests.Post(url, (*grequests.RequestOptions)(ro))
	return &Response{
		StatusCode:  response.StatusCode,
		Header:      response.Header,
		Body:        response.Bytes(),
		RawResponse: response.RawResponse,
	}, err
}

func Head(url string, ro *RequestOptions) (*Response, error) {
	ro.fixJsonRequestEscapeIssue()
	response, err := grequests.Head(url, (*grequests.RequestOptions)(ro))
	return &Response{
		StatusCode:  response.StatusCode,
		Header:      response.Header,
		Body:        response.Bytes(),
		RawResponse: response.RawResponse,
	}, err
}

func Options(url string, ro *RequestOptions) (*Response, error) {
	ro.fixJsonRequestEscapeIssue()
	response, err := grequests.Options(url, (*grequests.RequestOptions)(ro))
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

func (ro *RequestOptions) fixJsonRequestEscapeIssue() error {

	switch ro.JSON.(type) {

	case string:
	case []byte:
		return nil

	default:
		buffer := bytes.Buffer{}
		encoder := json.NewEncoder(buffer)
		encoder.SetEscapeHTML(false)
		err := encoder.Encode(ro.JSON)
		if err != nil {
			return err
		}
		ro.JSON = buffer.Bytes()
	}
	return nil

}
