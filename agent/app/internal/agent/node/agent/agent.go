package agent

import (
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/agent/node/data"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/config"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/database"
)

type IncrementalReader struct {
	nodeName                      string
	files                         []string
	scrapeIntervalSeconds         int
	metadataScrapeIntervalSeconds int
	results                       chan<- data.Chunk
	metadata                      chan<- data.NodeState
	redis                         database.Redis
}

func NewReader(cfg *config.Config, logsChan chan<- data.Chunk, metadataChan chan<- data.NodeState) IncrementalReader {
	return IncrementalReader{
		nodeName:                      cfg.Global.NodeName,
		files:                         cfg.WatchedFiles,
		scrapeIntervalSeconds:         cfg.Global.LogScrapeIntervalSeconds,
		metadataScrapeIntervalSeconds: cfg.Global.MetadataScrapeIntervalSeconds,
		results:                       logsChan,
		metadata:                      metadataChan,
		redis:                         database.NewRedis(cfg.Redis.Url, cfg.Redis.Password, cfg.Redis.Database),
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

func (r *IncrementalReader) watchFile(dir string, cooldownSeconds int, results chan<- data.Chunk) {
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
				Kind:          "Node",
				Name:          r.nodeName,
				CollectedAtMs: time.Now().UnixMilli(),
				Filename:      dir,
				Content:       string(buf),
			}
		}

		time.Sleep(time.Duration(cooldownSeconds * 1000))
	}
}

func (r *IncrementalReader) gatherNodeMetadata() {
	for {
		state := data.NewNodeState(r.nodeName, r.files)
		state.SetTimestamp()
		r.metadata <- state

		time.Sleep(time.Duration(r.metadataScrapeIntervalSeconds) * time.Second)
	}
}