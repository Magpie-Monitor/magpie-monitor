package remoteWrite

import (
	"bytes"
	"log"
	"net/http"
	"time"
)

type RemoteWriter struct {
	urls  []string
	cache string // TODO - cache has to work for each url separately
}

func NewRemoteWriter(urls []string) RemoteWriter {
	return RemoteWriter{urls: urls}
}

func (w *RemoteWriter) Write(content string) {
	// TODO - reiterate on retry + cache
	for _, url := range w.urls {
		//content = w.cache + content

		err, code := w.sendRequest(url, content)
		retries := 0
		for code >= 400 {
			log.Println("Error sending request: ", err, ", status code: ", code)
			log.Println("Retrying request...")

			err, code = w.sendRequest(url, content)
			if code >= 400 {
				if retries > 5 {
					//w.cache += content
					break
				}
				retries++
				time.Sleep(10 * time.Second)
			}
		}
	}
}

//func (w *RemoteWriter) cacheContent(content string) {
//	w.cache = w.cache + content
//}
//
//func (w *RemoteWriter) clearCache() {
//	w.cache = ""
//}

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
