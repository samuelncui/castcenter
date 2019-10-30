package castcenter

import (
	"fmt"
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

// ClusterStatus .
type ClusterStatus struct {
	SubNetLeader bool
	Conns        map[string]*net.UDPConn
}

// GetClusterStatus .
func (c *CastCenter) GetClusterStatus() *ClusterStatus {
	cs := c.clusterStatus.Load()
	if cs == nil {
		return nil
	}

	return cs.(*ClusterStatus)
}

// SetClusterStatus .
func (c *CastCenter) SetClusterStatus(cs *ClusterStatus) {
	logrus.Infof("castcenter fresh cluster status: %+v", cs)
	c.clusterStatus.Store(cs)
}

func (c *CastCenter) heartbeatKey() string {
	return fmt.Sprintf("%s:heartbeat", c.clusterName)
}

func (c *CastCenter) subnetLeaderKey() string {
	return fmt.Sprintf("%s:subnet:%s:leader", c.clusterName, netinfo.SubNet)
}

func (c *CastCenter) heartbeat() error {
	now := time.Now().Unix()
	timeout := int64(c.clusterTimeout / time.Second)
	oldStatus := c.GetClusterStatus()

	_, err := c.redis.ZAdd(c.heartbeatKey(), c.serviceAddr, float64(now))
	if err != nil {
		logrus.WithError(err).Errorf("castcenter: heartbeat fail")
		return err
	}

	subnetLeader, err := c.checkSubNetLeader()
	if err != nil {
		logrus.WithError(err).Errorf("castcenter: fetch hosts list fail")
		return err
	}

	hosts, err := c.redis.ZRangeByScore(c.heartbeatKey(), float64(now-timeout), float64(now+timeout))
	if err != nil {
		logrus.WithError(err).Errorf("castcenter: fetch hosts list fail")
		return err
	}

	var oldConns map[string]*net.UDPConn
	if oldStatus != nil && oldStatus.Conns != nil {
		oldConns = oldStatus.Conns
	}

	conns := make(map[string]*net.UDPConn, len(hosts))
	for _, host := range hosts {
		addr, err := net.ResolveUDPAddr("udp", host)
		if err != nil {
			logrus.WithError(err).Errorf("castcenter: append conn fail when resolve udp addr, host= %s", host)
			return err
		}

		if netinfo.SubNet.Contains(addr.IP) {
			continue
		}

		conn, ok := oldConns[host]
		if ok {
			conns[host] = conn
			continue
		}

		conn, err = net.DialUDP("udp", nil, addr)
		if err != nil {
			logrus.WithError(err).Errorf("castcenter: append conn fail when dial udp conn, host= %s", host)
			return err
		}

		conns[host] = conn
	}

	c.SetClusterStatus(&ClusterStatus{
		SubNetLeader: subnetLeader,
		Conns:        conns,
	})

	return nil
}

func (c *CastCenter) checkSubNetLeader() (bool, error) {
	key := c.subnetLeaderKey()

	ok, err := c.redis.SetNX(key, c.serviceAddr, c.clusterTimeout)
	if err != nil {
		return false, err
	}

	if ok {
		logrus.Infof("castcenter: as subnet leader")
		return true, nil
	}

	oldStatus := c.GetClusterStatus()
	if oldStatus != nil && !oldStatus.SubNetLeader {
		return false, nil
	}

	host, err := c.redis.Get(key)
	if err != nil {
		return false, err
	}

	if host != c.serviceAddr {
		return false, nil
	}

	c.redis.Expire(key, c.clusterTimeout)
	logrus.Infof("castcenter: as subnet leader, refresh")
	return true, nil
}
