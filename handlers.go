package castcenter

import (
	"crypto/md5"
	"encoding/hex"
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

// Run .
func (c *CastCenter) Run() {
	if err := c.heartbeat(); err != nil {
		logrus.Fatal(err)
	}

	go recoverLoop(c.startReceive)
	go recoverLoop(c.startSend)

	c.loop()
}

func (c *CastCenter) loop() {
	tick := time.NewTicker(c.clusterHeartbeat)
	for {
		select {
		case <-tick.C:
			c.heartbeat()
		}
	}
}

func (c *CastCenter) startReceive() {
	addr, err := net.ResolveUDPAddr("udp", c.serviceAddr)
	if err != nil {
		logrus.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		logrus.Fatal(err)
	}

	if err := NewUDPServer(c.handleReceive, c.chanSize).ListenUDP(conn); err != nil {
		logrus.Fatal(err)
	}
}

func (c *CastCenter) handleReceive(event *UDPEvent) {
	if !c.GetClusterStatus().SubNetLeader {
		return
	}

	hashed := md5.Sum(event.buf)
	c.cache.Set(string(hashed[:]), 1, time.Second)

	logrus.Debugf("receive: len= %d content= %s", len(event.buf), hex.Dump(event.buf))
	c.multicastConn.Write(event.buf)
}

func (c *CastCenter) startSend() {
	addr, err := net.ResolveUDPAddr("udp", c.multicastAddr)
	if err != nil {
		logrus.Fatal(err)
	}

	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		logrus.Fatal(err)
	}

	if err := NewUDPServer(c.handleSend, c.chanSize).ListenUDP(conn); err != nil {
		logrus.Fatal(err)
	}
}

func (c *CastCenter) handleSend(event *UDPEvent) {
	if event.ip != netinfo.IP {
		return
	}

	hashed := md5.Sum(event.buf)
	if _, ok := c.cache.Get(string(hashed[:])); ok {
		return
	}

	hexstr := hex.Dump(event.buf)
	for host, conn := range c.GetClusterStatus().Conns {
		logrus.Debugf("send to %s: len= %d content= %s", host, len(event.buf), hexstr)
		conn.Write(event.buf)
	}
}
