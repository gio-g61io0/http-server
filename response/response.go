package response

import (
	"fmt"
	"io"
	"personal-http-server/parser/constants"
	"strings"
)

type SupportedMethod string
type SupportedSupportedStatusCode int

const (
	GET  SupportedMethod = "GET"
	POST SupportedMethod = "POST"
)

const (
	OK             SupportedSupportedStatusCode = 200
	BAD_REQUEST    SupportedSupportedStatusCode = 400
	INTERNAL_ERROR SupportedSupportedStatusCode = 500
)

type Response struct {
	statusLine []byte
	headers    map[string][]byte
	body       []byte
	writer     io.Writer
}
type ResponseStatusLine struct {
	version    string
	method     SupportedMethod
	reason     string
	statusCode SupportedSupportedStatusCode
}

type RequestLineOption func(*ResponseStatusLine) error
type Option func(*Response) error

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

func (res *Response) HeaderSet(key string, value string) {
	res.headers[strings.ToLower(key)] = []byte(value)
}

func (res *Response) ForEach(cb func(k, v string)) {
	for key, value := range res.headers {
		cb(key, string(value))
	}

}

func (res *Response) WriteHeader() string {
	headerString := ""
	res.ForEach(func(key, value string) {
		headerString += fmt.Sprintf("%s:%s\r\n", key, value)
	})

	return fmt.Sprintf("%s\r\n", headerString)

}

func Method(method string) RequestLineOption {
	return func(reqStat *ResponseStatusLine) error {
		methodUpper := strings.ToUpper(method)
		if methodUpper != "GET" && methodUpper != "POST" {
			return fmt.Errorf("Unsupported method %s", method)
		}

		reqStat.method = SupportedMethod(methodUpper)
		return nil

	}
}

func StatusCode(code SupportedSupportedStatusCode) RequestLineOption {
	return func(reqStat *ResponseStatusLine) error {
		reqStat.statusCode = code
		return nil
	}
}

func Reason(resStat *ResponseStatusLine) error {

	switch resStat.statusCode {
	case OK:
		resStat.reason = "OK"
	case BAD_REQUEST:
		resStat.reason = "Bad Request"
	case INTERNAL_ERROR:
		resStat.reason = "Internal Server Error"
	}
	return nil
}

func Version(reqStat *ResponseStatusLine) error {
	reqStat.version = "HTTP/1.1"
	return nil
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
		Version(&statusLine)
		Reason(&statusLine)

		statusLineString := statusLine.BuildString()
		r.statusLine = []byte(statusLineString)
		return nil
	}
}

// Body should be buffer reader of the body
func BuildBody(bufferedBody io.Reader) Option {
	buf := make([]byte, 2)

	return func(r *Response) error {
		for {
			n, err := bufferedBody.Read(buf)
			if n <= 0 || err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}

			r.body = append(r.body, buf[:n]...)
		}
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
	return fmt.Sprintf("%s %d %s\r\n", statusLine.version, statusLine.statusCode, statusLine.reason)
}
