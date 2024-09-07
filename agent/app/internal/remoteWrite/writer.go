package remoteWrite

import (
	"bytes"
	"encoding/json"
	"log"
	"logather/internal/agent/node"
	"net/http"
	"time"
)

type Content struct {
	dir     string
	content string
}

type RemoteWriter struct {
	urls  []string
	cache map[string]string
}

func NewRemoteWriter(urls []string) RemoteWriter {
	return RemoteWriter{urls: urls}
}

func (w *RemoteWriter) Write(content node.IncrementalFetch) {
	for _, url := range w.urls {
		content.Content = w.getCachedContent(url) + content.Content
		jsonContent, err := json.Marshal(content)
		if err != nil {
			log.Println("Error converting content to JSON: ", err)
		}

		err, code := w.sendRequest(url, string(jsonContent))
		retries := 0
		for code >= 400 {
			log.Println("Error sending request: ", err, ", status code: ", code)
			log.Println("Retrying request...")

			err, code = w.sendRequest(url, string(jsonContent))
			if code >= 400 {
				if retries > 5 {
					w.cacheContent(url, content.Content)
					break
				}
				retries++
				time.Sleep(10 * time.Second)
			}
		}

		w.clearCache(url)
	}
}

func (w *RemoteWriter) cacheContent(url, content string) {
	val, ok := w.cache[url]
	if ok {
		w.cache[url] = val + content
	} else {
		w.cache[url] = val
	}
}

func (w *RemoteWriter) getCachedContent(url string) string {
	return w.cache[url]
}

func (w *RemoteWriter) clearCache(url string) {
	w.cache[url] = ""
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
	defer res.Body.Close()

	if err != nil {
		log.Println("Error sending http request: ", err)
		return err, 0
	}

	return nil, res.StatusCode
}
