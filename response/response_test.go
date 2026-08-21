package response

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponseOjectCreation(t *testing.T) {
	writer, err := FileWriter("example-writer.txt")

	if err != nil {
		t.Fatal(err)
	}

	response, err := NewResponse(writer, HttpVersion(), BuildStatusLine(StatusCode(OK)))

	require.NoError(t, err)
	assert.Equal(t, []byte("HTTP/1.1 200 OK\r\n"), response.statusLine)

}
