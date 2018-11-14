package collector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func getHTTPBytes(client http.Client, url string) ([]byte, error) {
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

func readTaskMetadata(client http.Client, endpoint string, output interface{}) error {
	bs, err := getHTTPBytes(client, endpoint)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bs, output); err != nil {
		return err
	}

	return nil
}

func readTaskStats(client http.Client, endpoint string, output interface{}) error {
	bs, err := getHTTPBytes(client, endpoint)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bs, output); err != nil {
		return err
	}

	return nil
}

func fPrettyPrint(w io.Writer, obj interface{}) error {
	bs, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return err
	}
	_, err = w.Write(bs)
	return err
}
