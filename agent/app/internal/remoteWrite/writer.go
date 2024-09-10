package remoteWrite

import (
	"bytes"
	"log"
	"net/http"
	"time"
)

type RemoteWriter struct {
	urls          []string
	retryInterval int
	maxRetries    int
	cache         map[string]string
}

func NewRemoteWriter(urls []string, retryInterval, maxRetries int) RemoteWriter {
	return RemoteWriter{
		urls:          urls,
		retryInterval: retryInterval,
		maxRetries:    maxRetries,
		cache:         make(map[string]string),
	}
}

func (w *RemoteWriter) Write(content string) {
	for _, url := range w.urls {
		content = w.getCachedContent(url) + content

		code, err := w.sendRequest(url, content)
		retries := 0
		for err != nil {
			log.Println("Error sending request: ", err, ", status code: ", code)
			log.Println("Retrying request...")

			code, err = w.sendRequest(url, content)
			if err != nil {
				if retries > w.maxRetries {
					w.cacheContent(url, content)
					break
				}
				retries++
				time.Sleep(time.Duration(w.retryInterval * 1000))
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

func (w *RemoteWriter) sendRequest(url string, content string) (int, error) {
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
