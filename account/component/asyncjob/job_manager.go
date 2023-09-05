package asyncjob

import (
	"context"
	"log"
	"sync"

	"github.com/tronglv92/accounts/common"
)

type group struct {
	isConcurrrent bool
	jobs          []Job
	wg            *sync.WaitGroup
}

func NewGroup(isConcurrrent bool, jobs ...Job) *group {
	g := &group{
		isConcurrrent: isConcurrrent,
		jobs:          jobs,
		wg:            new(sync.WaitGroup),
	}
	return g
}

func (g *group) Run(ctx context.Context) error {
	g.wg.Add(len(g.jobs))

	errChan := make(chan error, len(g.jobs))

	for i, _ := range g.jobs {
		if g.isConcurrrent {
			// Do this instead
			go func(aj Job) {
				defer common.AppRecover()
				errChan <- g.runJob(ctx, aj)
				g.wg.Done()
			}(g.jobs[i])

			continue
		}

		job := g.jobs[i]

		// err := g.runJob(ctx, job)
		// if err != nil {
		// 	return err
		// }
		// errChan <- err
		errChan <- g.runJob(ctx, job)

		g.wg.Done()
	}
	var err error

	for i := 1; i <= len(g.jobs); i++ {
		if v := <-errChan; v != nil {
			err = v
		}
	}

	g.wg.Wait()
	return err
}

// Retry if needed
func (g *group) runJob(ctx context.Context, j Job) error {
	if err := j.Execute(ctx); err != nil {
		for {
			log.Println(err)
			if j.State() == StateRetryFailed {
				return err
			}

			if j.Retry(ctx) == nil {
				return nil
			}
		}
	}
	return nil
}
