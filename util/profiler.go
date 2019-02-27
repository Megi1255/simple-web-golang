package util

import (
	"context"
	"errors"
	"simple-web-golang/config"
	"time"
)

type Profiler struct {
	perf        map[string]*Performance
	handleFlush FlushHandler
	Enable      bool
}

type Performance struct {
	StartTime   time.Time     `json:"start_time"`
	ElapsedTime time.Duration `json:"elapsed_time"`
	Ended       bool          `json:"-"`
}

type FlushHandler func(string, *Performance) error

func NewProfiler() *Profiler {
	return &Profiler{
		perf: make(map[string]*Performance),
	}
}

func (p *Profiler) Start(key string) {
	if !p.Enable {
		return
	}
	p.perf[key] = &Performance{
		StartTime: time.Now(),
		Ended:     false,
	}
}

func (p *Profiler) End(key string) {
	if !p.Enable {
		return
	}
	if perf, ok := p.perf[key]; ok {
		perf.ElapsedTime = time.Since(perf.StartTime)
		perf.Ended = true
	}
}

func (p *Profiler) SetFlushHandler(h FlushHandler) {
	p.handleFlush = h
}

func (p *Profiler) Flush() error {
	if !p.Enable {
		return nil
	}
	if p.handleFlush == nil {
		return errors.New("handler is nil")
	}

	for k, perf := range p.perf {
		if perf.Ended {
			if err := p.handleFlush(k, perf); err != nil {
				return err
			}
		}
	}
	p.perf = make(map[string]*Performance)
	return nil
}

func TsFromContext(c context.Context) time.Time {
	val := c.Value(config.KeyTimestamp)
	return val.(time.Time)
}

func ProfilerFromContext(c context.Context) *Profiler {
	val := c.Value(config.KeyProfiler)
	return val.(*Profiler)
}
