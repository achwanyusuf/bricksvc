package schedulerengine

import (
	"time"
)

type SchedulerEngineInterface interface {
	Close()
	Schedule(name string, duration time.Duration, fn HandlerFunc)
}

type schedulerEngine struct {
	signal []chan struct{}
}

func New() SchedulerEngineInterface {
	return &schedulerEngine{
		signal: make([]chan struct{}, 0),
	}
}

type HandlerFunc func()

func (s *schedulerEngine) Schedule(name string, duration time.Duration, fn HandlerFunc) {
	signal := make(chan struct{}, 1)
	s.signal = append(s.signal, signal)
	go func() {
		ticker := time.NewTicker(duration)
		defer ticker.Stop()
		for {
			select {
			case <-signal:
				return
			case <-ticker.C:
				fn()
			}
		}
	}()
}

func (s *schedulerEngine) Close() {
	for _, s := range s.signal {
		s <- struct{}{}
	}
}
