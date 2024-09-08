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
	dir         string
	application string
	node        string
	content     string
}

type RemoteWriter struct {
	urls  []string
	cache map[string]string
}

func NewRemoteWriter(urls []string) RemoteWriter {
	return RemoteWriter{urls: urls, cache: make(map[string]string)}
}

// TODO - decouple writer from node.IncrementalFetch, make it "struct agnostic"
func (w *RemoteWriter) Write(content node.IncrementalFetch) {
	for _, url := range w.urls {
		content.Content = w.getCachedContent(url) + content.Content
		jsonContent, err := json.Marshal(content)
		if err != nil {
			log.Println("Error converting content to JSON: ", err)
		}

		err, code := w.sendRequest(url, string(jsonContent))
		retries := 0
		for err != nil {
			log.Println("Error sending request: ", err, ", status code: ", code)
			log.Println("Retrying request...")

			err, code = w.sendRequest(url, string(jsonContent))
			if err != nil {
				if retries > 5 {
					w.cacheContent(url, content.Content)
					break
				}
				retries++
				time.Sleep(2 * time.Second)
			}
		}

		w.clearCache(url)
	}
}

// TODO - think about cache in redis
func (w *RemoteWriter) cacheContent(url, content string) {
	val, ok := w.cache[url]
	if ok {
		w.cache[url] = val + content
	} else {
		w.cache[url] = content
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

	if err != nil {
		log.Println("Error sending http request: ", err)
		return err, 0
	}
	defer res.Body.Close()

	return nil, res.StatusCode
}
