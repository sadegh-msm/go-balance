package balancers

import (
	"sync"
)

type RoundRobin struct {
	sync.RWMutex
	num   uint64
	hosts []string
}

// Add a server to available servers list
func (r *RoundRobin) Add(host string) {
	r.Lock()
	defer r.Unlock()

	for _, h := range r.hosts {
		if h == host {
			return
		}
	}
	r.hosts = append(r.hosts, host)
}

// Remove a server to available servers list
func (r *RoundRobin) Remove(host string) {
	r.Lock()
	defer r.Unlock()
	for i, h := range r.hosts {
		if h == host {
			r.hosts = append(r.hosts[:i], r.hosts[i+1:]...)
			return
		}
	}
}

// Balance the requests equally between hosts
func (r *RoundRobin) Balance(host string) (string, error) {
	r.RLock()
	defer r.RUnlock()
	if len(r.hosts) == 0 {
		return "", NoHostError
	}
	h := r.hosts[r.num%uint64(len(r.hosts))]
	r.num++
	return h, nil
}

func (r *RoundRobin) Inc(host string) {
	// no need to implement
}

func (r *RoundRobin) Done(host string) {
	// no need to implement
}

func initBalancer() {
	Factories[RRBalancer] = NewRoundRobin
}

func NewRoundRobin(hosts []string) Balancer {
	return &RoundRobin{
		num:   0,
		hosts: hosts,
	}
}
