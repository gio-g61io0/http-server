package request

import (
	"bufio"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func bufferedReader(byte_read int, filename string) *bufio.Reader {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReaderSize(file, byte_read)
	return reader

}
func stringBufferedReader(data string, bufferLength int) *bufio.Reader {
	return bufio.NewReaderSize(strings.NewReader(data), bufferLength)
}

func TestHeaderParser(t *testing.T) {
	headerStringByte := []byte("Throttle-Key:Testing\r\n")

	header := Header{
		fields: make(map[string]string),
	}

	nParse, done, err := header.parse(headerStringByte)
	assert.False(t, done)
	require.NoError(t, err)
	assert.Equal(t, nParse, len(headerStringByte))
}
func TestHeaderParserError(t *testing.T) {
	headerStringByte := []byte("Throttle-Key :Testing\r\n")

	header := Header{
		fields: make(map[string]string),
	}

	nParse, done, err := header.parse(headerStringByte)
	require.Error(t, err)
	assert.False(t, done)
	assert.Equal(t, nParse, 0) //means nothing is parsed
}
func TestBodyParser(t *testing.T) {

	contentToParse := []byte("This is a content here hehe")
	body := Body{
		content:       make([]byte, 0),
		contentLength: len(contentToParse),
		readNBytes:    0,
	}
	nParse, err := body.parse(contentToParse)

	require.NoError(t, err)
	assert.Equal(t, nParse, len(contentToParse))
}

func TestRequestLineParser(t *testing.T) {
	file, err := os.Open("example-http-request.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, 34)
	buffer := make([]byte, 34)

	request := NewRequest(reader)

	_, errParse := request.ReadRequest(buffer)

	assert.Equal(t, request.requestLine.method, "GET")
	assert.Equal(t, request.requestLine.methodTarget, "/api/v1/users")
	assert.Equal(t, request.requestLine.httpVersion, "HTTP/1.1")
	require.NoError(t, errParse)

}
func TestRequestParserWithBody(t *testing.T) {
	sampleHTTPRequest := "GET /api/v1/users HTTP/1.1\r\n" +
		"Host: example.com\r\n" +
		"Accept: application/json\r\n" +
		"Authorization: Bearer <token>\r\n" +
		"Content-Length: 9\r\n" +
		"User-Agent: CustomClient/1.0\r\n" +
		"Connection: keep-alive\r\n\r\n"
		// "Gio Gwapo"

	reader := stringBufferedReader(sampleHTTPRequest, 64)
	buffer := make([]byte, 64)

	request := NewRequest(reader)
	_, errParse := request.ReadRequest(buffer)

	assert.Equal(t, request.requestLine.method, "GET")
	assert.Equal(t, request.requestLine.methodTarget, "/api/v1/users")
	assert.Equal(t, request.requestLine.httpVersion, "HTTP/1.1")
	assert.Equal(t, "9", request.headers.Get("Content-Length"))
	require.NoError(t, errParse) //should error for now when there is a content length but no body

}
