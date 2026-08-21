package response

import (
	"bufio"
	"io"
	"os"
)

func FileWriter(filename string) (io.Writer, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	return writer, nil
}
