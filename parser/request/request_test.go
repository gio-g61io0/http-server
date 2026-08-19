package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderParser(t *testing.T) {
	headerStringByte := []byte("Throttle-Key:Testing\r\n")

	header := Header{
		fields: make(map[string]string),
	}

	nParse, err := header.parse(headerStringByte)
	require.NoError(t, err)
	assert.Equal(t, nParse, len(headerStringByte))
}
func TestHeaderParserError(t *testing.T) {
	headerStringByte := []byte("Throttle-Key :Testing\r\n")

	header := Header{
		fields: make(map[string]string),
	}

	nParse, err := header.parse(headerStringByte)
	require.Error(t, err)
	assert.Equal(t, nParse, 0) //means nothing is parsed
}
