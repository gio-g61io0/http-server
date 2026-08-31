package response

import (
	// "bufio"
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestResponseOjectCreationOk(t *testing.T) {
	writer, err := FileWriter("example-writer.txt")

	if err != nil {
		t.Fatal(err)
	}

	response, err := NewResponse(writer, HttpVersion(), BuildStatusLine(StatusCode(OK)))

	require.NoError(t, err)
	assert.Equal(t, []byte("HTTP/1.1 200 OK\r\n"), response.statusLine)

}
func TesTestResponseOjectCreationBad(t *testing.T) {

	writer, err := FileWriter("example-writer.txt")

	if err != nil {
		t.Fatal(err)
	}

	response, err := NewResponse(writer, HttpVersion(), BuildStatusLine(StatusCode(BAD_REQUEST)))

	require.NoError(t, err)
	assert.Equal(t, []byte("HTTP/1.1 400 Bad Request\r\n"), response.statusLine)
}
func TesTestResponseOjectCreationServerErr(t *testing.T) {

	writer, err := FileWriter("example-writer.txt")

	if err != nil {
		t.Fatal(err)
	}

	response, err := NewResponse(writer, HttpVersion(), BuildStatusLine(StatusCode(INTERNAL_ERROR)))

	require.NoError(t, err)
	assert.Equal(t, []byte("HTTP/1.1 500 Internal Server Error\r\n"), response.statusLine)
}
func TestHeader(t *testing.T) {

	writer, err := FileWriter("example-writer.txt")

	if err != nil {
		t.Fatal(err)
	}

	response, err := NewResponse(writer, HttpVersion(), BuildStatusLine(StatusCode(INTERNAL_ERROR)))

	response.HeaderSet("Content-Type", "application/json")

	headerString := response.WriteHeader()

	assert.Equal(t, "content-type:application/json\r\n\r\n", headerString)
}
func TestBody(t *testing.T) {

	writer, err := FileWriter("example-writer.txt")
	data := bytes.NewReader([]byte("Hello world. Gwapo si Gio"))
	// bufferedReader := bufio.NewReader(data)

	if err != nil {
		t.Fatal(err)
	}
	response, err := NewResponse(writer, HttpVersion(), BuildStatusLine(StatusCode(INTERNAL_ERROR)), BuildBody(data))

	response.HeaderSet("Content-length", strconv.Itoa(len("Hello world. Gwapo si Gio")))

	require.NoError(t, err)
	assert.Equal(t, "Hello world. Gwapo si Gio", string(response.body))

}
