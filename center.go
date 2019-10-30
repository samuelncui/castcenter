package castcenter

import (
	"fmt"
	"net"
	"sync/atomic"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

const (
	defaultServicePort          = 52593
	defaultMulticastAddr        = "239.95.3.9:9539"
	defaultChanSize             = 2048
	defaultCacheTimeout         = time.Second
	defaultCacheCleanupInterval = time.Second * 10

	defaultClusterName      = "abc950309/castcenter"
	defaultClusterHeartbeat = time.Second * 10
	defaultClusterTimeout   = time.Second * 30
)

// CastCenter .
type CastCenter struct {
	redis                Redis
	servicePort          int
	multicastAddr        string
	chanSize             int
	cacheTimeout         time.Duration
	cacheCleanupInterval time.Duration

	clusterName      string
	clusterHeartbeat time.Duration
	clusterTimeout   time.Duration

	serviceAddr   string
	clusterStatus atomic.Value // *ClusterStatus

	multicastConn *net.UDPConn
	cache         *cache.Cache
}

// Option .
type Option func(*CastCenter) *CastCenter

// New .
func New(options ...Option) *CastCenter {
	c := &CastCenter{}

	for _, option := range options {
		c = option(c)
	}

	if c.redis == nil {
		logrus.Fatal(fmt.Errorf("castcenter new: expected option SetRedis"))
	}

	if c.servicePort == 0 {
		c.servicePort = defaultServicePort
	}

	if c.multicastAddr == "" {
		c.multicastAddr = defaultMulticastAddr
	}

	if c.clusterName == "" {
		c.clusterName = defaultClusterName
	}

	if c.chanSize == 0 {
		c.chanSize = defaultChanSize
	}

	if c.cacheTimeout == 0 {
		c.cacheTimeout = defaultCacheTimeout
	}

	if c.cacheCleanupInterval == 0 {
		c.cacheCleanupInterval = defaultCacheCleanupInterval
	}

	if c.clusterHeartbeat == 0 {
		c.clusterHeartbeat = defaultClusterHeartbeat
	}

	if c.clusterTimeout == 0 {
		c.clusterTimeout = defaultClusterTimeout
	}

	mcAddr, err := net.ResolveUDPAddr("udp", c.multicastAddr)
	if err != nil {
		logrus.Fatal(err)
	}

	mcConn, err := net.DialUDP("udp", nil, mcAddr)
	if err != nil {
		logrus.Fatal(err)
	}

	c.multicastConn = mcConn
	c.serviceAddr = fmt.Sprintf("%s:%d", netinfo.IP, c.servicePort)
	c.cache = cache.New(c.cacheTimeout, c.cacheCleanupInterval)

	return c
}

// SetRedis .
func SetRedis(redis Redis) Option {
	return func(c *CastCenter) *CastCenter {
		c.redis = redis
		return c
	}
}

// SetClusterName .
func SetClusterName(name string) Option {
	return func(c *CastCenter) *CastCenter {
		c.clusterName = name
		return c
	}
}

// SetServicePort .
func SetServicePort(port int) Option {
	return func(c *CastCenter) *CastCenter {
		c.servicePort = port
		return c
	}
}

// SetMulticastAddr .
func SetMulticastAddr(addr string) Option {
	return func(c *CastCenter) *CastCenter {
		c.multicastAddr = addr
		return c
	}
}

// SetChanSize .
func SetChanSize(size int) Option {
	return func(c *CastCenter) *CastCenter {
		c.chanSize = size
		return c
	}
}

// SetCacheTimeout .
func SetCacheTimeout(timeout time.Duration) Option {
	return func(c *CastCenter) *CastCenter {
		c.cacheTimeout = timeout
		return c
	}
}

// SetCacheCleanupInterval .
func SetCacheCleanupInterval(interval time.Duration) Option {
	return func(c *CastCenter) *CastCenter {
		c.cacheCleanupInterval = interval
		return c
	}
}

// SetClusterHeartbeat .
func SetClusterHeartbeat(interval time.Duration) Option {
	return func(c *CastCenter) *CastCenter {
		c.clusterHeartbeat = interval
		return c
	}
}

// SetClusterTimeout .
func SetClusterTimeout(timeout time.Duration) Option {
	return func(c *CastCenter) *CastCenter {
		c.clusterTimeout = timeout
		return c
	}
}
