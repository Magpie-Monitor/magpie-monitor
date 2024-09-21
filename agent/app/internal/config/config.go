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
	Mode                     string
	ClusterName              string
	ScrapeInterval           int
	RedisUrl                 string
	RedisPassword            string
	WatchedFiles             []string
	RemoteWriteUrls          []string
	RemoteWriteRetryInterval int
	RemoteWriteMaxRetries    int
	ExcludedNamespaces       []string
}

func GetConfig() Config {
	mode := flag.String("scrape", "pods", "Mode in which log collector runs, either \"nodes\" to scrape nodes or \"pods\" to scrape pods.")
	clusterName := flag.String("clusterFriendlyName", "unknown", "Friendly name of your cluster, visible in Magpie Cloud.")

	scrapeInterval := flag.Int("scrapeInterval", 10, "Interval between scraping logs from files in \"nodes\" mode or pods in \"pods\" mode.")
	remoteWriteRetryInterval := flag.Int("remoteWriteRetryInterval", 2, "Interval between retries in case of Remote Write error.")
	remoteWriteMaxRetries := flag.Int("remoteWriteMaxRetries", 5, "Maximal number of retries in case of Remote Write error.")

	redisUrl := flag.String("redisUrl", "", "Redis URL in cluster DNS format, that is: service.namespace.svc.cluster.local:port")
	redisPassword := flag.String("redisPassword", "", "Password to Redis instance pointed by --redisUrl flag.")

	var watchedFiles arrayFlags
	flag.Var(&watchedFiles, "file", "Log files that are watched for log collector running in \"nodes\" mode.")

	var remoteWriteUrls arrayFlags
	flag.Var(&remoteWriteUrls, "remoteWriteUrl", "URL to which logs are pushed using remote write protocol.")

	var excludedNamespaces arrayFlags
	flag.Var(&excludedNamespaces, "excludedNamespace", "Namespace which is excluded from log scraping when agent runs in \"pods\" mode.")

	flag.Parse()

	log.Println("Agent configured to run in mode: ", *mode)
	log.Println("Redis url: ", *redisUrl)

	return Config{
		Mode:                     *mode,
		ClusterName:              *clusterName,
		ScrapeInterval:           *scrapeInterval,
		RedisUrl:                 *redisUrl,
		RedisPassword:            *redisPassword,
		WatchedFiles:             watchedFiles,
		RemoteWriteUrls:          remoteWriteUrls,
		RemoteWriteRetryInterval: *remoteWriteRetryInterval,
		RemoteWriteMaxRetries:    *remoteWriteMaxRetries,
		ExcludedNamespaces:       excludedNamespaces,
	}
}
