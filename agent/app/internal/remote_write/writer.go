package remote_write

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

type RemoteWriter struct {
	urls []string
}

func (w *RemoteWriter) write(content string) {
	// TODO - retry + cache
	for _, url := range w.urls {
		err, code := w.sendRequest(url, content)
		fmt.Println(err, code)
	}
}

func (w *RemoteWriter) sendRequest(url string, content string) (error, int) {
	r, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(content)))
	if err != nil {
		log.Println("Error creating http request: ", err)
		return err, 0
	}

	r.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		log.Println("Error sending http request: ", err)
		return err, 0
	}

	defer res.Body.Close()
	return nil, res.StatusCode
}
