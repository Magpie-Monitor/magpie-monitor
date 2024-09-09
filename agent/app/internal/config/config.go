package config

import (
	"flag"
	"log"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "Array flags"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type Config struct {
	Mode            string
	RedisUrl        string
	WatchedFiles    []string
	RemoteWriteUrls []string
}

func GetConfig() Config {
	mode := flag.String("scrape", "nodes", "Mode in which log collector runs, either \"nodes\" to scrape nodes or \"pods\" to scrape pods.")
	redisUrl := flag.String("redisUrl", "", "Redis URL in cluster DNS format, that is: service.namespace.svc.cluster.local:port")

	var watchedFiles arrayFlags
	flag.Var(&watchedFiles, "file", "Log files that are watched for log collector running in \"nodes\" mode.")

	var remoteWriteUrls arrayFlags
	flag.Var(&remoteWriteUrls, "remoteWriteUrl", "URL to which logs are pushed using remote write protocol.")

	flag.Parse()

	log.Println("Agent configured to run in mode: ", *mode)
	log.Println("Redis url: ", *redisUrl)

	return Config{Mode: *mode, RedisUrl: *redisUrl, WatchedFiles: watchedFiles, RemoteWriteUrls: []string{}}
}
