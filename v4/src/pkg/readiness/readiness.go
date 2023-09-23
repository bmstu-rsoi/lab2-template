package readiness

import (
	"sync"
)

type Probe struct {
	marks map[string]bool
	mx    sync.RWMutex
}

func New() *Probe {
	return &Probe{mx: sync.RWMutex{}, marks: map[string]bool{}}
}

func (p *Probe) Mark(key string, value bool) {
	p.mx.Lock()
	defer p.mx.Unlock()

	p.marks[key] = value
}

func (p *Probe) Ready() bool {
	p.mx.RLock()
	defer p.mx.RUnlock()

	for _, isReady := range p.marks {
		if !isReady {
			return false
		}
	}

	return true
}
