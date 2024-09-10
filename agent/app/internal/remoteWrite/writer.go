package remoteWrite

import (
	"bytes"
	"log"
	"net/http"
	"time"
)

type RemoteWriter struct {
	urls  []string
	cache map[string]string
}

func NewRemoteWriter(urls []string) RemoteWriter {
	return RemoteWriter{
		urls:  urls,
		cache: make(map[string]string),
	}
}

func (w *RemoteWriter) Write(content string) {
	for _, url := range w.urls {
		content = w.getCachedContent(url) + content

		err, code := w.sendRequest(url, content)
		retries := 0
		for err != nil {
			log.Println("Error sending request: ", err, ", status code: ", code)
			log.Println("Retrying request...")

			err, code = w.sendRequest(url, content)
			if err != nil {
				if retries > 5 {
					w.cacheContent(url, content)
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
