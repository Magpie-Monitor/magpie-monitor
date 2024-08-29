package node

import (
	"io"
	"log"
	"logather/internal/transformer"
	"os"
	"time"
)

type IncrementalFetch struct {
	Dir     string
	Content string
}

type IncrementalReader struct {
	files           []string
	transformers    []transformer.Transformer
	remoteWriteUrls []string
	results         chan IncrementalFetch
	progress        map[string]int64
}

func NewReader(files []string, transformers []transformer.Transformer, remoteWriteUrls []string,
	results chan IncrementalFetch) IncrementalReader {
	return IncrementalReader{files: files, transformers: transformers, remoteWriteUrls: remoteWriteUrls,
		results: results, progress: make(map[string]int64)}
}

func (r *IncrementalReader) WatchFiles() {
	for _, file := range r.files {
		go r.watchFile(file, 1, r.results)
	}
}

func (r *IncrementalReader) prepareFile(dir string) (*os.File, int64) {
	f, err := os.Open(dir)
	if err != nil {
		log.Println("Error opening file = ", dir)
		panic(err)
	}

	fi, err := f.Stat()
	if err != nil {
		log.Println("Error reading stat for file: ", dir)
		panic(err)
	}

	// seek to previous ending spot
	var currentSize int64
	// TODO - redis here
	previousSize, ok := r.progress[dir]
	if ok {
		_, err = f.Seek(previousSize, io.SeekStart)
		if err != nil {
			panic(err)
		}
		currentSize = previousSize
	} else {
		f.Seek(0, io.SeekEnd)
		currentSize = fi.Size()
	}

	return f, currentSize
}

func (r *IncrementalReader) watchFile(dir string, cooldown int, results chan IncrementalFetch) {
	f, currentSize := r.prepareFile(dir)
	defer f.Close()

	for {
		fi, err := f.Stat()
		if err != nil {
			log.Println("Error reading stat for file: ", dir)
			panic(err)
		}

		size := fi.Size()
		byteDiff := size - currentSize
		if byteDiff > 0 {
			buf := make([]byte, byteDiff+1)

			_, err := f.Seek(-byteDiff, io.SeekCurrent)
			if err != nil {
				log.Println("Error seeking diff for file: ", dir)
				panic(err)
			}

			_, err = f.Read(buf)
			if err != nil {
				log.Println("Error reading buffer for file: ", dir)
				panic(err)
			}

			log.Println("READ = ", string(buf))
			log.Println("BYTE DIFF = ", byteDiff)

			currentSize = size
			r.progress[dir] = size

			results <- IncrementalFetch{Dir: dir, Content: r.transform(string(buf))}
		}

		time.Sleep(time.Duration(cooldown * 1000))
	}

}

func (r *IncrementalReader) transform(content string) string {
	for _, t := range r.transformers {
		content = t.Transform(content)
	}
	return content
}

//func RemoteWrite(url string) {
//
//}

//func DatabaseWrite() {
//
//}
