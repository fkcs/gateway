package http

import (
	"github.com/fkcs/gateway/internal/context"
	"github.com/gin-gonic/gin"
	"sync"
)

var once sync.Once

type Job interface {
	Init(ctx *context.Ctx) error
	Route(router *gin.RouterGroup) error
}

type ScheduleJob struct {
	mu        *sync.RWMutex
	ctx       *context.Ctx
	router    *gin.RouterGroup
	Observers []Job
}

func NewScheduleJob(ctx *context.Ctx, router *gin.RouterGroup) *ScheduleJob {
	return &ScheduleJob{
		mu:     new(sync.RWMutex),
		ctx:    ctx,
		router: router,
	}
}

func (s *ScheduleJob) RegisterRoute(observer Job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Observers = append(s.Observers, observer)
}

func (s *ScheduleJob) Observe() error {
	once.Do(func() {

	})
	for _, observer := range s.Observers {
		if err := observer.Init(s.ctx); err != nil {
			return err
		}
		if err := observer.Route(s.router); err != nil {
			return err
		}
	}
	return nil
}
