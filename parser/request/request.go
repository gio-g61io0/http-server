package request

import (
	"fmt"
	"io"
	"log"
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
	requestLine *RequestLine
	headers     *Header
	body        *Body
	state       RequestParserState
	reader      io.Reader
}
type Body struct {
	content       []byte
	contentLength int
	readNBytes    int
}
type RequestLine struct {
	method       string
	httpVersion  string
	methodTarget string
}
type Header struct {
	fields map[string]string
}

func NewBody() *Body {
	return &Body{
		content:       []byte{},
		contentLength: 0,
		readNBytes:    0,
	}

}

func (body *Body) parse(buffer []byte) (int, error) {

	//Only consume bytes that are remaining
	if body.readNBytes+len(buffer) > body.contentLength {
		remaining := body.contentLength - body.readNBytes
		if remaining <= 0 {
			return 0, nil //Done. Nothing to consume
		}
		buffer = buffer[:remaining]
	}

	body.content = append(body.content, buffer...)
	body.readNBytes += len(buffer)

	return len(buffer), nil

}
func (r *Request) parse(buffer []byte) (int, error) {
	bytes_read := 0
	for {
		switch r.state {
		case ParseInit:
			r.state = ParseRequestLine
		case ParseRequestLine:
			n_read, err := r.parseRequestLine(buffer[bytes_read:])
			if err != nil {
				return 0, err
			}
			if n_read <= 0 {
				return bytes_read, nil
			}
			fmt.Printf("Read bytes from request line %d\n", n_read)
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

			n_read, err := r.body.parse(buffer[bytes_read:])

			if err != nil {
				return bytes_read, err
			}
			if n_read == 0 {
				r.state = ParseDone
			}
			bytes_read += n_read

		case ParseDone:
			return bytes_read, nil
		}

	}
}

func NewRequest(reader io.Reader) *Request {
	return &Request{
		requestLine: nil,
		headers:     &Header{fields: make(map[string]string)},
		body:        NewBody(),
		state:       ParseInit,
		reader:      reader,
	}
}

// printing only
func (r *Request) InspectRequest() {
	log.Printf("Request line %v", r.requestLine)
	log.Printf("Request line %v", r.headers)
}

func (r *Request) ReadRequest(buffer []byte) (int, error) {

	if r.reader == nil {
		return 0, fmt.Errorf("No socket reader set")
	}

	bytes_read := 0
	bufLength := 0
	for {

		if r.IsDone() {
			return bytes_read, nil
		}

		n, err := r.reader.Read(buffer[bufLength:])

		bufLength += n

		if err == io.EOF {

			//for cases where there is a content length specified for a GET request but there is no body.
			if r.headers.Exist("Content-Length") && r.requestLine.method == "GET" {
				r.state = ParseDone
				return bytes_read, nil
			}
			return 0, fmt.Errorf("End of File! already")
		}

		if err != nil {
			return 0, err
		}

		if n <= 0 {
			return bytes_read, err
		}

		n_parsed, err := r.parse(buffer[:bufLength])

		if err != nil {
			return 0, err
		}

		copy(buffer, buffer[n_parsed:bufLength])
		bufLength -= n_parsed

		fmt.Printf("N parsed %d\n", n_parsed)
		fmt.Printf("Buffer after 1 parse pass %s\n", string(buffer[n_parsed:]))

		bytes_read += n_parsed
	}

}

func (r *Request) IsDone() bool {
	return r.state == ParseDone
}

func (r *Request) parseRequestLine(buffer []byte) (int, error) {

	idx := strings.Index(string(buffer), constants.SEPARATOR)

	//Nothing to parse yet. Buffer is not complete
	if idx == -1 {
		return 0, nil
	}

	completeRequestLine := strings.Split(string(buffer), "\r\n")
	if len(completeRequestLine) < 2 {
		return 0, nil
	}
	parts := strings.Split(completeRequestLine[0], " ")

	if len(parts) < 3 {
		return 0, fmt.Errorf("Invalid request line")
	}

	r.requestLine = &RequestLine{
		method:       parts[0],
		methodTarget: parts[1],
		httpVersion:  parts[2],
	}

	return idx + len(constants.SEPARATOR), nil
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

	fmt.Printf("Data here %s\n", string(buf))
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
