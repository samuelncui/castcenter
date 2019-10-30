package main

import (
	"flag"

	"github.com/abc950309/castcenter"
	"github.com/abc950309/castcenter/tools/rediswrap"
	"github.com/sirupsen/logrus"
)

var (
	h = flag.Bool("h", false, "show help")

	redisURL = flag.String("redis", "redis://127.0.0.1:6379", "Redis service url")

	servicePort          = flag.Int("port", 0, "Service port")
	multicastAddr        = flag.String("group", "", "Multicast group")
	chanSize             = flag.Int("chsize", 0, "Chan size for message")
	cacheTimeout         = flag.Duration("cache-timeout", 0, "Cache timeout for deduplication")
	cacheCleanupInterval = flag.Duration("cache-cleanup-interval", 0, "Cache cleanup interval for deduplication")

	clusterName      = flag.String("cluster", "", "Cluster name, as redis key prefix")
	clusterHeartbeat = flag.Duration("heartbeat", 0, "Cluster heartbeat interval")
	clusterTimeout   = flag.Duration("timeout", 0, "Cluster timeout")
)

func main() {
	flag.Parse()

	if *h {
		flag.Usage()
		return
	}

	opts := []castcenter.Option{}

	r, err := rediswrap.NewWithURL(*redisURL)
	if err != nil {
		logrus.Fatal(err)
	}
	opts = append(opts, castcenter.SetRedis(r))

	if *servicePort != 0 {
		opts = append(opts, castcenter.SetServicePort(*servicePort))
	}
	if *multicastAddr != "" {
		opts = append(opts, castcenter.SetMulticastAddr(*multicastAddr))
	}
	if *chanSize != 0 {
		opts = append(opts, castcenter.SetChanSize(*chanSize))
	}
	if *cacheTimeout != 0 {
		opts = append(opts, castcenter.SetCacheTimeout(*cacheTimeout))
	}
	if *cacheCleanupInterval != 0 {
		opts = append(opts, castcenter.SetCacheCleanupInterval(*cacheCleanupInterval))
	}
	if *clusterName != "" {
		opts = append(opts, castcenter.SetClusterName(*clusterName))
	}
	if *clusterHeartbeat != 0 {
		opts = append(opts, castcenter.SetClusterHeartbeat(*clusterHeartbeat))
	}
	if *clusterTimeout != 0 {
		opts = append(opts, castcenter.SetClusterTimeout(*clusterTimeout))
	}

	castcenter.New(opts...).Run()
}
