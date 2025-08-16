package softap

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type HttpTransmitter struct {
	hostname string
}

func NewHttpTransmitter(hostname string) (*HttpTransmitter, error) {
	if hostname == "" {
		return nil, errors.New("hostname is empty")
	}

	return &HttpTransmitter{hostname: hostname}, nil
}

func (httpTransmitter HttpTransmitter) Send(path string, data []byte) ([]byte, error) {
	if path == "" {
		return nil, errors.New("path is required")
	}

	if data == nil {
		return nil, errors.New("data is required")
	}

	resp, err := http.Post(
		httpTransmitter.hostname+path,
		"application/x-www-form-urlencoded",
		bytes.NewBuffer(data))

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, errors.New(fmt.Sprintf("Status code: %s - %s", resp.StatusCode, string(body)))
	}

	return body, nil
}
