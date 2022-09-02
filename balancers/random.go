package balancers

import (
	"math/rand"
	"sync"
	"time"
)

// Random will randomly select a http server from the server
type Random struct {
	sync.RWMutex
	hosts []string
	rnd   *rand.Rand
}

func init() {
	Factories[RandomBalancer] = NewRandom
}

// NewRandom create new Random balancer
func NewRandom(hosts []string) Balancer {
	return &Random{
		hosts: hosts,
		rnd:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Add new host to the balancer
func (r *Random) Add(host string) {
	r.Lock()
	defer r.Unlock()
	for _, h := range r.hosts {
		if h == host {
			return
		}
	}
	r.hosts = append(r.hosts, host)
}

// Remove new host from the balancer
func (r *Random) Remove(host string) {
	r.Lock()
	defer r.Unlock()
	for i, h := range r.hosts {
		if h == host {
			r.hosts = append(r.hosts[:i], r.hosts[i+1:]...)
		}
	}
}

// Balance selects a suitable host according
func (r *Random) Balance(host string) (string, error) {
	r.RLock()
	defer r.RUnlock()
	if len(r.hosts) == 0 {
		return "", NoHostError
	}
	return r.hosts[r.rnd.Intn(len(r.hosts))], nil
}

func (r *Random) Inc(host string) {
	// no need to implement
}

func (r *Random) Done(host string) {
	// no need to implement
}
