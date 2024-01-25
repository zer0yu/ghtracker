package runner

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"io"
	"os"
	"path/filepath"
)

// OutputWriter outputs content to writers.
type OutputWriter struct {
	JSON bool
}

// NewOutputWriter creates a new OutputWriter
func NewOutputWriter(json bool) *OutputWriter {
	return &OutputWriter{JSON: json}
}

func (o *OutputWriter) createFile(filename string, appendToFile bool) (*os.File, error) {
	if filename == "" {
		return nil, errors.New("empty filename")
	}

	dir := filepath.Dir(filename)

	if dir != "" {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return nil, err
			}
		}
	}

	var file *os.File
	var err error
	if appendToFile {
		file, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	} else {
		file, err = os.Create(filename)
	}
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (o *OutputWriter) writeJSONResults(results interface{}, writer io.Writer) error {
	encoder := jsoniter.NewEncoder(writer)
	err := encoder.Encode(results)
	if err != nil {
		return err
	}
	return nil
}
