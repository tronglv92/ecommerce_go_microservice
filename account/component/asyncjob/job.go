package asyncjob

import (
	"context"
	"log"
	"time"
)

// Job Requirement:
// 1. Job can do something (handler)
// 2. Job can retry
//  2.1 Config retry times and duration
// 3. Should be stateful
// 4. We should have job manager to manage jobs (*)

type Job interface {
	Execute(ctx context.Context) error
	Retry(ctx context.Context) error
	State() JobState
	SetMaxRetry(count int)
	SetSleepTime(time time.Duration)
	SetMultiplier(multiplier int)
}

const (
	defaultMaxTimeout    = time.Second * 10
	defaultMaxRetryCount = 3
)

var (
	defaultSleepTime = 100 * time.Millisecond
	defaultMultiple  = 2
)

type JobHandler func(ctx context.Context) error

type JobState int

const (
	StateInit JobState = iota
	StateRunning
	StateFailed
	StateTimeout
	StateCompleted
	StateRetryFailed
)

type jobConfig struct {
	Name       string
	MaxTimeout time.Duration
	MaxRetry   int
	SleepTime  time.Duration
	Multiplier int
}

func (js JobState) String() string {
	return []string{"Init", "Running", "Failed", "Timeout", "Completed", "RetryFailed"}[js]
}

type job struct {
	config     jobConfig
	handler    JobHandler
	state      JobState
	retryIndex int
	stopChan   chan bool
}

func NewJob(handler JobHandler, options ...OptionHdl) *job {
	j := job{
		config: jobConfig{
			MaxTimeout: defaultMaxTimeout,
			SleepTime:  defaultSleepTime,
			Multiplier: defaultMultiple,
			MaxRetry:   defaultMaxRetryCount,
		},
		handler:    handler,
		retryIndex: -1,
		state:      StateInit,
		stopChan:   make(chan bool),
	}

	for i := range options {
		options[i](&j.config)
	}

	return &j
}
func (j *job) Execute(ctx context.Context) error {
	log.Printf("execute %s\n", j.config.Name)
	j.state = StateRunning
	var err error
	err = j.handler(ctx)

	if err != nil {
		j.state = StateFailed
		return err
	}
	j.state = StateCompleted

	return nil

	// ch := make(chan error)
	// ctxJob, doneFunc := context.WithCancel(ctx)

	// go func() {
	// 	j.state = StateRunning
	// 	var err error

	// 	err = j.handler(ctxJob)

	// 	if err != nil {
	// 		j.state = StateFailed
	// 		ch <- err
	// 		return
	// 	}

	// 	j.state = StateCompleted
	// 	ch <- err
	// }()

	// // for {
	// // 	select {
	// // 		case <-j.stopChan:
	// // 			break
	// // 		default:
	// // 			fmt.Println("Hello world")

	// // 	}
	// // }
	// select {
	// case err := <-ch:
	// 	doneFunc()
	// 	return err
	// case <-j.stopChan:
	// 	doneFunc()
	// 	return nil
	// }
	//
	//return <-ch

}
func (j *job) Retry(ctx context.Context) error {
	j.retryIndex += 1
	var sleepTime time.Duration
	if j.retryIndex == 0 {
		sleepTime = j.config.SleepTime

	} else {
		mul := j.retryIndex * j.config.Multiplier
		sleepTime = j.config.SleepTime * time.Duration(mul)

	}

	log.Printf("sleepTime %s\n", sleepTime)
	time.Sleep(sleepTime)
	err := j.Execute(ctx)

	if err == nil {
		j.state = StateCompleted
		return nil
	}

	if j.retryIndex == j.config.MaxRetry-1 {
		j.state = StateRetryFailed
		return err
	}
	j.state = StateFailed
	return err
}

func (j *job) State() JobState { return j.state }
func (j *job) RetryIndex() int { return j.retryIndex }

type OptionHdl func(*jobConfig)

func WithName(name string) OptionHdl {
	return func(cf *jobConfig) {
		cf.Name = name
	}
}

func (j *job) SetMaxRetry(count int) {
	j.config.MaxRetry = count
}
func (j *job) SetSleepTime(time time.Duration) {
	j.config.SleepTime = time
}
func (j *job) SetMultiplier(multiplier int) {
	j.config.Multiplier = multiplier
}
