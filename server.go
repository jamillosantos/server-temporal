package srvtemporal

import (
	"context"

	"github.com/pkg/errors"
	"go.temporal.io/sdk/worker"
)

var (
	ErrClientRequired = errors.New("client required: use WithClient to set the client")
)

// Worker is a services.Server that is able to initialize and manage the temporal Worker together with the
// github.com/jamillosantos/application.
type Worker struct {
	name string
	w    worker.Worker
}

// NewWorker implements
func NewWorker(registerer Registerer, cfg PlatformConfig, options ...Option) (Worker, error) {
	opts := defaultOpts()
	for _, opt := range options {
		opt(&opts)
	}
	if opts.client == nil {
		return Worker{}, ErrClientRequired
	}
	w := worker.New(opts.client, opts.taskQueue, worker.Options{
		MaxConcurrentActivityExecutionSize:      cfg.MaxConcurrentActivityExecutionSize,
		WorkerActivitiesPerSecond:               cfg.WorkerActivitiesPerSecond,
		MaxConcurrentLocalActivityExecutionSize: cfg.MaxConcurrentLocalActivityExecutionSize,
		WorkerLocalActivitiesPerSecond:          cfg.WorkerLocalActivitiesPerSecond,
		TaskQueueActivitiesPerSecond:            cfg.TaskQueueActivitiesPerSecond,
		MaxConcurrentActivityTaskPollers:        cfg.MaxConcurrentActivityTaskPollers,
		MaxConcurrentWorkflowTaskExecutionSize:  cfg.MaxConcurrentWorkflowTaskExecutionSize,
		MaxConcurrentWorkflowTaskPollers:        cfg.MaxConcurrentWorkflowTaskPollers,
		EnableLoggingInReplay:                   cfg.EnableLoggingInReplay,
		StickyScheduleToStartTimeout:            cfg.StickyScheduleToStartTimeout,
		BackgroundActivityContext:               opts.backgroundAcitivityContext,
		WorkflowPanicPolicy:                     worker.WorkflowPanicPolicy(cfg.WorkflowPanicPolicy),
		WorkerStopTimeout:                       cfg.WorkerStopTimeout,
		EnableSessionWorker:                     cfg.EnableSessionWorker,
		MaxConcurrentSessionExecutionSize:       cfg.MaxConcurrentSessionExecutionSize,
		DisableWorkflowWorker:                   cfg.DisableWorkflowWorker,
		LocalActivityWorkerOnly:                 cfg.LocalActivityWorkerOnly,
		Identity:                                cfg.Identity,
		DeadlockDetectionTimeout:                cfg.DeadlockDetectionTimeout,
		MaxHeartbeatThrottleInterval:            cfg.MaxHeartbeatThrottleInterval,
		DefaultHeartbeatThrottleInterval:        cfg.DefaultHeartbeatThrottleInterval,
		Interceptors:                            opts.interceptors,
		OnFatalError:                            opts.onFatalError,
		DisableEagerActivities:                  cfg.DisableEagerActivities,
		MaxConcurrentEagerActivityExecutionSize: cfg.MaxConcurrentEagerActivityExecutionSize,
		DisableRegistrationAliasing:             cfg.DisableRegistrationAliasing,
		BuildID:                                 cfg.BuildID,
		UseBuildIDForVersioning:                 cfg.UseBuildIDForVersioning,
	})
	registerer.Register(w)
	return Worker{
		name: opts.name,
		w:    w,
	}, nil
}

func (w *Worker) Name() string {
	return w.name
}

func (w *Worker) Listen(_ context.Context) error {
	return w.w.Start()
}

func (w *Worker) Close(_ context.Context) error {
	w.w.Stop()
	return nil
}
