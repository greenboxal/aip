package jobs

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"github.com/hashicorp/go-multierror"
	"github.com/jbenet/goprocess"
	goprocessctx "github.com/jbenet/goprocess/context"
	"go.uber.org/zap"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-forddb/pkg/tracing"
	"github.com/greenboxal/aip/aip-sdk/pkg/utils"
)

type Supervisor struct {
	logger *zap.SugaredLogger

	db forddb.Database

	registry utils.BindingRegistry[JobHandlerBinding]
	handlers map[string]JobHandlerBinding

	m           sync.Mutex
	runningJobs map[JobID]*runningJob
}

func NewSupervisor(
	logger *zap.SugaredLogger,
	db forddb.Database,
	registry utils.BindingRegistry[JobHandlerBinding],
) *Supervisor {
	handlers := map[string]JobHandlerBinding{}

	for _, binding := range registry.Bindings() {
		handlers[binding.ID().String()] = binding
	}

	return &Supervisor{
		logger: logger.Named("job-supervisor"),

		db:       db,
		registry: registry,

		runningJobs: map[JobID]*runningJob{},
		handlers:    handlers,
	}
}

func (sup *Supervisor) CheckJobState(job *Job) (JobHandle, error) {
	shouldBeRunning := false

	switch job.Status.State {
	case JobStateScheduled:
		shouldBeRunning = true
	case JobStatePending:
		shouldBeRunning = true
	}

	state := sup.getRunningJob(job.ID, shouldBeRunning)

	if state == nil {
		return nil, nil
	}

	if err := state.reconcile(job); err != nil {
		return nil, err
	}

	return state, nil
}

func (sup *Supervisor) getRunningJob(id JobID, create bool) *runningJob {
	if existing, ok := sup.runningJobs[id]; ok {
		return existing
	}

	sup.m.Lock()
	defer sup.m.Unlock()

	if existing, ok := sup.runningJobs[id]; ok {
		return existing
	}

	if !create {
		return nil
	}

	job := &runningJob{
		sup: sup,
		id:  id,

		doneCh: make(chan struct{}),
	}

	sup.runningJobs[id] = job

	return job
}

type runningJob struct {
	ctx  context.Context
	proc goprocess.Process
	sup  *Supervisor

	m         sync.Mutex
	isRunning atomic.Bool

	id      JobID
	job     Job
	handler JobHandler

	err    error
	result any

	doneCh chan struct{}
}

func (r *runningJob) Await() (any, error) {
	_, _ = <-r.doneCh

	return r.result, r.err
}

func (r *runningJob) Payload() any {
	return r.job.Spec.Payload
}

func (r *runningJob) SetError(err error) {
	if !r.isRunning.Load() {
		panic(errors.New("job is not running"))
	}

	r.err = multierror.Append(err)
}

func (r *runningJob) SetResult(result any) {
	if !r.isRunning.Load() {
		panic(errors.New("job is not running"))
	}

	r.result = result
}

func (r *runningJob) Result() any {
	return r.result
}

func (r *runningJob) Error() error {
	return r.err
}

func (r *runningJob) Job() Job {
	return r.job
}

func (r *runningJob) Context() context.Context {
	return r.ctx
}

func (r *runningJob) Process() goprocess.Process {
	return r.proc
}

func (r *runningJob) shouldBeRunning() bool {
	return r.job.Status.State == JobStatePending || r.job.Status.State == JobStateScheduled || r.job.Status.State == JobStateRunning
}

func (r *runningJob) run(proc goprocess.Process) {
	if !r.isRunning.CompareAndSwap(false, true) {
		r.m.Unlock()
		return
	}

	defer func() {
		r.proc = nil
		r.ctx = nil
		r.isRunning.Store(false)

		if r.doneCh != nil {
			close(r.doneCh)
		}
	}()

	r.err = nil
	r.result = nil

	ctx := goprocessctx.OnClosingContext(proc)

	if r.job.Status.State == JobStatePending {
		r.updateState(ctx, JobStateScheduled)
		return
	} else if r.job.Status.State == JobStateScheduled {
		r.proc = proc.Go(func(proc goprocess.Process) {
			r.ctx = goprocessctx.OnClosingContext(proc)

			r.updateState(r.ctx, JobStateRunning)

			spanCtx, span := tracing.StartTrace(r.ctx, "Job: "+r.job.Spec.Handler)
			defer span.End()

			r.ctx = spanCtx

			r.err = r.handler.Run(r)
		})

		if err := r.proc.Err(); err != nil {
			r.err = multierror.Append(r.err, err)
		}

		r.setCompleted(ctx, r.err)
	}
}

func (r *runningJob) reconcile(job *Job) error {
	if job.ID != r.id {
		panic("job ID mismatch")
	}

	r.m.Lock()
	defer r.m.Unlock()

	if job.GetResourceVersion() >= r.job.GetResourceVersion() {
		r.job = *job
	}

	if r.handler == nil {
		binding := r.sup.handlers[r.job.Spec.Handler]

		if binding == nil {
			return errors.New("no handler found")
		}

		r.handler = binding.Handler()
	}

	if r.shouldBeRunning() && !r.isRunning.Load() {
		goprocess.Go(r.run)
	}

	return nil
}

func (r *runningJob) setCompleted(ctx context.Context, err error) {
	r.err = err

	if r.err == nil {
		if r.result != nil {
			r.job.Status.Result = r.result
			r.updateState(ctx, JobStateCompleted)
		} else {
			r.updateState(ctx, JobStateCompleted)
		}
	}

	if r.err != nil {
		r.job.Status.LastError = r.err.Error()
		r.updateState(ctx, JobStateFailed)
	}

	if r.doneCh != nil {
		close(r.doneCh)
		r.doneCh = nil
	}
}

func (r *runningJob) updateState(ctx context.Context, state JobState) {
	r.job.Status.State = state

	updated, err := forddb.Put(ctx, r.sup.db, &r.job)

	if err != nil {
		panic(err)
	}

	r.job = *updated
}
