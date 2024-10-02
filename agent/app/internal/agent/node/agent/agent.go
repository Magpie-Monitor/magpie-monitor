package agent

import (
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/agent/node/data"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/database"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type IncrementalReader struct {
	nodeName                          string
	files                             []string
	scrapeIntervalSeconds             int
	metadataCollectionIntervalSeconds int
	results                           chan data.Chunk
	metadata                          chan data.NodeState
	redis                             database.Redis
}

func NewReader(nodeName string, files []string, scrapeIntervalSeconds int, metadataCollectionIntervalSeconds int,
	results chan data.Chunk, metadata chan data.NodeState,
	redisUrl, redisPassword string, redisDb int) IncrementalReader {
	return IncrementalReader{
		nodeName:                          nodeName,
		files:                             files,
		scrapeIntervalSeconds:             scrapeIntervalSeconds,
		metadataCollectionIntervalSeconds: metadataCollectionIntervalSeconds,
		results:                           results,
		metadata:                          metadata,
		redis:                             database.NewRedis(redisUrl, redisPassword, redisDb),
	}
}

func (r *IncrementalReader) Start() {
	go r.watchFiles()
	r.gatherNodeMetadata()
}

func (r *IncrementalReader) watchFiles() {
	for _, file := range r.files {
		go r.watchFile(file, r.scrapeIntervalSeconds, r.results)
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

	// Seek to previous ending spot.
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

func (r *IncrementalReader) watchFile(dir string, cooldownSeconds int, results chan data.Chunk) {
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

			// Move reader cursor for byteDiff bytes (being the size increase since last read) from the end of file.
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

			results <- data.Chunk{
				Kind:      "Node",
				Name:      r.nodeName,
				Timestamp: time.Now().UnixNano(),
				Namespace: dir,
				Content:   string(buf),
			}
		}

		time.Sleep(time.Duration(cooldownSeconds * 1000))
	}
}

func (r *IncrementalReader) gatherNodeMetadata() {
	for {
		state := data.NodeState{NodeName: r.nodeName, WatchedFiles: r.files}
		state.SetTimestamp()
		r.metadata <- state

		time.Sleep(time.Duration(r.metadataCollectionIntervalSeconds) * time.Second)
	}
}
