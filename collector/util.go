package collector

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

const (
	Namespace          string = "ecs"
	ContainerSubsystem string = "container"
	TaskSubsystem      string = "task"
)

func GetHTTPBytes(client http.Client, url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%v: %v", resp.Status, resp.StatusCode)
	}

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
