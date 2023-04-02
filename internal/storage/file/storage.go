package file

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/eugeniylennik/alertics/internal/metrics"
	"os"
)

type Writer struct {
	file *os.File
}

type Reader struct {
	file    *os.File
	scanner *bufio.Scanner
	decoder *json.Decoder
}

func NewWriter(fileName string) (*Writer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &Writer{
		file: file,
	}, nil
}

func (w *Writer) WriteMetrics(m []byte) error {
	if err := w.file.Truncate(0); err != nil {
		return err
	}
	if _, err := w.file.Seek(0, 0); err != nil {
		return err
	}
	if _, err := w.file.Write(m); err != nil {
		return err
	}
	return nil
}

func (w *Writer) Close() error {
	return w.file.Close()
}

func NewReader(fileName string) (*Reader, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &Reader{
		file:    file,
		scanner: bufio.NewScanner(file),
		decoder: json.NewDecoder(file),
	}, nil
}

func (r *Reader) ReadMetrics() ([]metrics.Data, error) {
	var m map[string]interface{}
	var data []metrics.Data

	if !r.scanner.Scan() {
		return nil, r.scanner.Err()
	}

	mBz := r.scanner.Bytes()

	err := json.Unmarshal(mBz, &m)
	if err != nil {
		fmt.Println("Error parsing metrics data:", err)
		return nil, err
	}
	for name, value := range m["Gauge"].(map[string]interface{}) {
		data = append(data, metrics.Data{
			Name:  name,
			Type:  "gauge",
			Value: value.(float64),
		})
	}

	for name, value := range m["Counter"].(map[string]interface{}) {
		data = append(data, metrics.Data{
			Name:  name,
			Type:  "counter",
			Value: value.(float64),
		})
	}

	return data, err
}

func (r *Reader) Close() error {
	return r.file.Close()
}
