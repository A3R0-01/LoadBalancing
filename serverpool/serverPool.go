package serverpool

import (
	"sync"

	"github.com/A3R0-01/loadbalance/backend"
)

type ServerPool interface {
	GetBackends() []backend.Backend
	GetNextValidPeer() backend.Backend
	AddBackend(backend.Backend)
	GetServerPoolSize() int
}

type roundRobin struct {
	backends []backend.Backend
	mux      sync.RWMutex
	current  int
}

func (r *roundRobin) Rotate() backend.Backend {
	r.mux.Lock()
	r.current = (r.current + 1) % r.GetServerPoolSize()
	r.mux.Unlock()
	return r.backends[r.current]
}

func (r *roundRobin) GetNextValidPeer() backend.Backend {
	for i := 0; i < r.GetServerPoolSize(); i++ {
		nextPeer := r.Rotate()
		if nextPeer.IsAlive() {
			return nextPeer
		}
	}
	return nil
}

type leastNumberOfConnections struct {
	backends []backend.Backend
	mux      sync.RWMutex
}

func (l *leastNumberOfConnections) GetNextValidPeer() backend.Backend {
	var leastConnectedPeer backend.Backend
	for _, b := range l.backends {
		if b.IsAlive() {
			leastConnectedPeer = b
			break
		}
	}
	for _, b := range l.backends {
		if !b.IsAlive() {
			continue
		}
		if leastConnectedPeer.GetActiveConnections() > b.GetActiveConnections() {
			leastConnectedPeer = b
		}
	}

	return leastConnectedPeer
}
