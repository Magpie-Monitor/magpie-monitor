package remote_write

import (
	"bytes"
	"log"
	"net/http"
)

type MetadataWriter struct {
	url string
}

func NewMetadataWriter(url string) RemoteWriter {
	return &MetadataWriter{url: url}
}

func (m *MetadataWriter) Write(content string) {
	// TODO - add buffering
	code, err := m.sendRequest(m.url, content)
	if err != nil || code >= 400 {
		log.Println("Error sending data to metadata service.")
	}
}

func (m *MetadataWriter) sendRequest(url string, content string) (int, error) {
	r, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(content)))
	if err != nil {
		log.Println("Error creating http request: ", err)
		return 0, err
	}

	r.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(r)

	if err != nil {
		log.Println("Error sending http request: ", err)
		return 0, err
	}
	defer res.Body.Close()

	return res.StatusCode, nil
}
