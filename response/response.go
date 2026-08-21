package response

import (
	"fmt"
	"io"
	"personal-http-server/parser/constants"
	"strings"
)

type SupportedMethod string

const (
	GET  SupportedMethod = "GET"
	POST SupportedMethod = "POST"
)

type Response struct {
	statusLine []byte
	headers    map[string][]byte
	body       []byte
	writer     io.Writer
}
type ResponseStatusLine struct {
	version string
	method  SupportedMethod
	reason  string
}

type RequestLineOption func(*ResponseStatusLine) error
type Option func(*Response) error

func Method(method string) Option {
	return func(res *Response) error {
		methodUpper := strings.ToUpper(method)

		if methodUpper != "GET" && methodUpper != "POST" {
			return fmt.Errorf("Unsupported method! Only GET and POST for now")
		}
		res.statusLine = append(res.statusLine, []byte(method)...)
		return nil
	}
}

func HttpVersion() Option {
	return func(res *Response) error {
		res.statusLine = append(res.statusLine, []byte(constants.HTTP_VERSION)...)
		return nil
	}
}

func Body(body []byte) Option {
	return func(res *Response) error {
		res.body = body
		return nil
	}

}

func (res *Response) HTTPVersion(version string) *Response {
	if version != "HTTP 1.1" {
		return nil
	}
	res.statusLine = append(res.statusLine, []byte(version)...)
	return res
}

func BuildStatusLine(reqLineOpts ...RequestLineOption) Option {
	return func(r *Response) error {
		statusLine := ResponseStatusLine{
			version: "",
			method:  "GET",
			reason:  "",
		}

		for _, reqOpt := range reqLineOpts {
			if err := reqOpt(&statusLine); err != nil {
				return fmt.Errorf("applying option: %w", err)
			}

		}
		statusLineString := statusLine.BuildString()
		r.statusLine = []byte(statusLineString)
		return nil
	}
}

func NewResponse(writer io.Writer, opts ...Option) (Response, error) {

	res := Response{
		statusLine: make([]byte, 0),
		headers:    make(map[string][]byte),
		body:       make([]byte, 0),
		writer:     writer,
	}

	for _, opt := range opts {
		if err := opt(&res); err != nil {
			return Response{}, fmt.Errorf("applying option: %w", err)
		}

	}
	return res, nil
}

func (statusLine *ResponseStatusLine) BuildString() string {
	return fmt.Sprintf("%s %s %s\r\n")

}
