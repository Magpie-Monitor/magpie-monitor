package node

import (
	"io"
	"log"
	"logather/internal/agent/pods"
	"logather/internal/database"
	"logather/internal/transformer"
	"os"
	"strconv"
	"time"
)

type IncrementalReader struct {
	files        []string
	transformers []transformer.Transformer
	results      chan pods.Chunk
	redis        database.Redis
}

func NewReader(files []string, transformers []transformer.Transformer, results chan pods.Chunk,
	redisUrl string) IncrementalReader {
	return IncrementalReader{files: files, transformers: transformers, results: results,
		redis: database.NewRedis(redisUrl, "", 0)} // TODO - reiterate on Redis password
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
	previousSize := r.redis.Get(dir)
	if previousSize != "" {
		size, _ := strconv.Atoi(previousSize)
		_, err = f.Seek(int64(size), io.SeekStart)
		if err != nil {
			panic(err)
		}
		currentSize = int64(size)
	} else {
		_, err = f.Seek(0, io.SeekEnd)
		if err != nil {
			log.Println("Error seeking end of file: ", dir)
		}
		currentSize = fi.Size()
	}

	return f, currentSize
}

func (r *IncrementalReader) watchFile(dir string, cooldown int, results chan pods.Chunk) {
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

			_, err = f.Seek(-byteDiff, io.SeekEnd)
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
			err = r.redis.Set(dir, strconv.FormatInt(currentSize, 10), -1)
			if err != nil {
				log.Println("Error persisting read progress for: ", dir)
			}

			// TODO - fetch real node name
			results <- pods.Chunk{Kind: "Node", Name: "mock-node-name", Namespace: dir, Content: r.transform(string(buf))}
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
