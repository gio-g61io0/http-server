package request

import (
	"fmt"
	"personal-http-server/parser/constants"
	"strconv"
	"strings"
)

type RequestParserState string

const (
	ParseInit        RequestParserState = "init"
	ParseRequestLine RequestParserState = "request_line"
	ParseHeaders     RequestParserState = "headers"
	ParseBody        RequestParserState = "body"
	ParseDone        RequestParserState = "done"
)

type Request struct {
	requestLine RequestLine
	headers     Header
	body        Body
	state       RequestParserState
}
type Body struct {
	content       []byte
	contentLength int
}
type RequestLine struct {
	method       string
	httpVersion  string
	methodTarget string
}
type Header struct {
	fields map[string]string
}

func NewBody() Body {
	return Body{
		content:       []byte{},
		contentLength: 0,
	}

}

func (body *Body) parse(buffer []byte, contentLength int) (int, error) {
	body.content = append(body.content, buffer...)
	body.contentLength += contentLength

	return len(buffer), nil

}

func (r *Request) parse(buffer []byte) (int, error) {

	bytes_read := 0
	for {
		switch r.state {
		case ParseInit:
			r.state = ParseHeaders
		case ParseRequestLine:
			n_read, err := r.parseRequestLine(buffer[bytes_read:])
			if err != nil {
				return 0, err
			}
			if n_read <= 0 {
				return bytes_read, nil
			}
			bytes_read += n_read
			r.state = ParseHeaders

		case ParseHeaders:
			n_read, done, err := r.headers.parse(buffer[bytes_read:])

			if err != nil {
				return bytes_read, err
			}

			bytes_read += n_read

			if n_read == 0 {
				return bytes_read, nil
			}

			if done {
				r.state = ParseBody
				contentLength, err := r.headers.ContentLength()

				//No content length. I have to check the spec if err if invalid content length value
				if err != nil {
					r.state = ParseDone
					return bytes_read, nil
				}

				//mutate body's content length field
				r.body.contentLength = contentLength
				return bytes_read, nil
			}

		case ParseBody:

			//If no content length header then no body bytes should be parsed
			if !r.headers.Exist("Content-Length") {
				r.state = ParseDone
				return bytes_read, nil
			}
			contentLength, err := strconv.Atoi(r.headers.Get("Content-Length"))
			if err != nil {
				return bytes_read, fmt.Errorf("Invalid Content Length field ")
			}

			//mutate request's body field content length
			r.body.contentLength = contentLength

			n_read, err := r.body.parse(buffer[bytes_read:], contentLength)

			if err != nil {
				return bytes_read, err
			}

		}

	}

}

func (r *Request) parseRequestLine(buffer []byte) (int, error) {
	completeRequestLine := strings.Split(string(buffer), "\r\n")
	if len(completeRequestLine) < 2 {
		return 0, nil
	}
	parts := strings.Split(completeRequestLine[0], " ")

	if len(parts) < 3 {
		return 0, fmt.Errorf("Invalid request line")
	}

	r.requestLine = RequestLine{
		method:       parts[0],
		methodTarget: parts[1],
		httpVersion:  parts[2],
	}

	return len(buffer) + len(constants.SEPARATOR), nil
}

func (h *Header) ContentLength() (int, error) {
	contentLength, err := strconv.Atoi(h.fields[strings.ToLower("Content-Length")])
	if err != nil {
		return 0, err
	}
	return contentLength, nil

}

func (h *Header) parseHeader(buf []byte) (string, string, error) {
	split := strings.Split(string(buf), constants.KEY_VALUE_SPLIT)

	if len(split) < 2 {
		return "", "", fmt.Errorf("Invalid header line")
	}
	if strings.Contains(split[0], " ") {
		return "", "", fmt.Errorf("Invalid Key! whitespace not allowed")
	}
	cleanedHeaderValue := strings.TrimSpace(split[1])
	lowerCasedHeaderKey := strings.ToLower(split[0])
	h.fields[lowerCasedHeaderKey] = cleanedHeaderValue

	return lowerCasedHeaderKey, cleanedHeaderValue, nil

}

func (h *Header) Get(key string) string {
	return h.fields[strings.ToLower(key)]
}

func (h *Header) Exist(key string) bool {
	_, exist := h.fields[strings.ToLower(key)]
	return exist
}

func (h *Header) parse(buf []byte) (int, bool, error) {
	separatorIdx := strings.Index(string(buf), constants.SEPARATOR)
	if separatorIdx == -1 {
		return 0, false, nil
	}

	//Done parsing the header
	if separatorIdx == 0 {
		return len(constants.SEPARATOR), true, nil
	}

	_, _, err := h.parseHeader(buf[:separatorIdx])
	if err != nil {
		return 0, false, err

	}

	return separatorIdx + len(constants.SEPARATOR), false, nil
}
