package util

import (
	"errors"
	"fmt"
	"time"
)

type Profiler struct {
	perf        map[string]Performance
	handleFlush FlushHandler
}

type Performance struct {
	StartTime   time.Time     `json:"start_time"`
	ElapsedTime time.Duration `json:"elapsed_time"`
	Ended       bool          `json:"-"`
}

type FlushHandler func(string, Performance) error

func NewProfiler() *Profiler {
	return &Profiler{
		perf: make(map[string]Performance),
	}
}

func (p *Profiler) Start(key string) {
	p.perf[key] = Performance{
		StartTime: time.Now(),
		Ended:     false,
	}
}

func (p *Profiler) End(key string) {
	if _, ok := p.perf[key]; ok {
		p.perf[key].ElapsedTime = time.Since(perf.StartTime)
		perf.Ended = true
		fmt.Printf("%s %v", key, perf)
	}
}

func (p *Profiler) SetFlushHandler(h FlushHandler) {
	p.handleFlush = h
}

func (p *Profiler) Flush() error {
	if p.handleFlush == nil {
		return errors.New("handler is nil")
	}

	for k, perf := range p.perf {
		fmt.Printf("%s %v", k, perf)
		if perf.Ended {
			if err := p.handleFlush(k, perf); err != nil {
				return err
			}
		}
	}
	p.perf = make(map[string]Performance)
	return nil
}
