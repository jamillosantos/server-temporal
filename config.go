package srvtemporal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/interceptor"
	"go.temporal.io/sdk/worker"
)

type WorkflowPanicPolicy worker.WorkflowPanicPolicy

var (
	workflowPanicPolicyBlock = []byte("block")
	workflowPanicPolicyFail  = []byte("fail")

	ErrUnknownWorkflowPanicPolicy = errors.New("unknown workflow panic policy")
)

func (w *WorkflowPanicPolicy) UnmarshalText(text []byte) error {
	switch {
	case bytes.Equal(text, workflowPanicPolicyBlock):
		*w = WorkflowPanicPolicy(worker.BlockWorkflow)
	case bytes.Equal(text, workflowPanicPolicyFail):
		*w = WorkflowPanicPolicy(worker.FailWorkflow)
	default:
		return fmt.Errorf("%w: %s", ErrUnknownWorkflowPanicPolicy, text)
	}
	return nil
}

type PlatformConfig struct {
	MaxConcurrentActivityExecutionSize      int                 `config:"max_concurrent_activity_execution_size"`
	WorkerActivitiesPerSecond               float64             `config:"worker_activities_per_second"`
	MaxConcurrentLocalActivityExecutionSize int                 `config:"max_concurrent_local_activity_execution_size"`
	WorkerLocalActivitiesPerSecond          float64             `config:"worker_local_activities_per_second"`
	TaskQueueActivitiesPerSecond            float64             `config:"task_queue_activities_per_second"`
	MaxConcurrentActivityTaskPollers        int                 `config:"max_concurrent_activity_task_pollers"`
	MaxConcurrentWorkflowTaskExecutionSize  int                 `config:"max_concurrent_workflow_task_execution_size"`
	MaxConcurrentWorkflowTaskPollers        int                 `config:"max_concurrent_workflow_task_pollers"`
	EnableLoggingInReplay                   bool                `config:"enable_logging_in_replay"`
	StickyScheduleToStartTimeout            time.Duration       `config:"sticky_schedule_to_start_timeout"`
	WorkflowPanicPolicy                     WorkflowPanicPolicy `config:"workflow_panic_policy"`
	WorkerStopTimeout                       time.Duration       `config:"worker_stop_timeout"`
	EnableSessionWorker                     bool                `config:"enable_session_worker"`
	MaxConcurrentSessionExecutionSize       int                 `config:"max_concurrent_session_execution_size"`
	DisableWorkflowWorker                   bool                `config:"disable_workflow_worker"`
	LocalActivityWorkerOnly                 bool                `config:"local_activity_worker_only"`
	Identity                                string              `config:"identity"`
	DeadlockDetectionTimeout                time.Duration       `config:"deadlock_detection_timeout"`
	MaxHeartbeatThrottleInterval            time.Duration       `config:"max_heartbeat_throttle_interval"`
	DefaultHeartbeatThrottleInterval        time.Duration       `config:"default_heartbeat_throttle_interval"`
	DisableEagerActivities                  bool                `config:"disable_eager_activities"`
	MaxConcurrentEagerActivityExecutionSize int                 `config:"max_concurrent_eager_activity_execution_size"`
	DisableRegistrationAliasing             bool                `config:"disable_registration_aliasing"`
	BuildID                                 string              `config:"build_id"`
	UseBuildIDForVersioning                 bool                `config:"use_build_id_for_versioning"`
}

type config struct {
	name                       string
	taskQueue                  string
	backgroundAcitivityContext context.Context
	interceptors               []interceptor.WorkerInterceptor
	onFatalError               func(error)
	client                     client.Client
}

func defaultOpts() config {
	return config{
		name: "Temporal Worker Server",
	}
}

type Option func(*config)

func WithName(name string) Option {
	return func(c *config) {
		c.name = name
	}
}

func WithTaskQueue(taskQueue string) Option {
	return func(c *config) {
		c.taskQueue = taskQueue
	}
}

func WithBackgroundActivityContext(ctx context.Context) Option {
	return func(c *config) {
		c.backgroundAcitivityContext = ctx
	}
}

func WithInterceptors(interceptors ...interceptor.WorkerInterceptor) Option {
	return func(c *config) {
		c.interceptors = interceptors
	}
}

func WithOnFatalError(onFatalError func(error)) Option {
	return func(c *config) {
		c.onFatalError = onFatalError
	}
}

func WithClient(client client.Client) Option {
	return func(c *config) {
		c.client = client
	}
}
