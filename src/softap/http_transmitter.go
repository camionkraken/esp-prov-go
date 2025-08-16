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

func NewHttpTransmitter(hostname string) *HttpTransmitter {
	return &HttpTransmitter{hostname: hostname}
}

func (httpTransmitter HttpTransmitter) Send(path string, data []byte) ([]byte, error) {
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
		return nil, errors.New(fmt.Sprint("Status code: %s - %s", resp.StatusCode, string(body)))
	}

	return body, nil
}
