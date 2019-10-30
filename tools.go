package castcenter

import (
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	pool = &sync.Pool{New: func() interface{} {
		return make([]byte, datagramReadBufferSize)
	}}
)

// Redis .
type Redis interface {
	Get(key string) (string, error)
	SetNX(key string, value interface{}, expiration time.Duration) (bool, error)
	Expire(key string, expiration time.Duration) (bool, error)
	ZAdd(key, member string, score float64) (int64, error)
	ZRangeByScore(key string, min, max float64) ([]string, error)
}

// GetBytes .
func GetBytes() []byte {
	return pool.Get().([]byte)
}

// PutBytes .
func PutBytes(buf []byte) {
	pool.Put(buf[:cap(buf)])
}

func recoverLoop(f func()) {
	defer func() {
		go recoverLoop(f)

		e := recover()
		if e == nil {
			return
		}

		err, ok := e.(error)
		if !ok {
			err = fmt.Errorf("%v", e)
		}

		logrus.WithError(err).Errorf("udp loop fail, recovered")
	}()

	f()
}
