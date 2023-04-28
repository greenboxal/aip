package ford

import (
	"context"

	"go.uber.org/zap"

	"github.com/greenboxal/aip/pkg/collective"
	"github.com/greenboxal/aip/pkg/ford/forddb"
)

type TaskReconciler struct {
	logger *zap.SugaredLogger

	db    forddb.Database
	cache map[collective.TaskID]*collective.Task
}

func NewTaskReconciler(
	logger *zap.SugaredLogger,
	db forddb.Database,
) *TaskReconciler {
	tr := &TaskReconciler{
		logger: logger.Named("task-reconciler"),

		db: db,

		cache: map[collective.TaskID]*collective.Task{},
	}

	watchCh := make(chan collective.TaskID, 128)

	db.AddListener(
		forddb.TypedListenerFunc[collective.TaskID, *collective.Task](
			func(id collective.TaskID, previous, current *collective.Task) {
				watchCh <- id
			},
		),
	)

	go func() {
		for id := range watchCh {
			previous := tr.cache[id]

			current, err := forddb.Get[*collective.Task](db, id)

			if err != nil {
				logger.Error(err)
				continue
			}

			if previous == nil || current == nil || current.GetVersion() > previous.GetVersion() {
				encoded, err := forddb.Encode(current)

				if err != nil {
					logger.Error(err)
					continue
				}

				decoded, err := forddb.Decode(encoded)

				if err != nil {
					logger.Error(err)
					continue
				}

				_, err = tr.Reconcile(context.Background(), previous, decoded.(*collective.Task))

				if err != nil {
					logger.Error(err)
				}
			}

			tr.cache[id] = current
		}
	}()

	return tr
}

func (tr *TaskReconciler) Reconcile(ctx context.Context, previous, current *collective.Task) (*collective.Task, error) {
	if current == nil && previous != nil {
		tr.logger.Info("task deleted", "task_id", current.ID)

		return nil, nil
	}

	tr.logger.Info("entering reconciliation loop", "task_id", current.ID)

	pipeline, err := forddb.Get[*collective.Pipeline](tr.db, current.Spec.PipelineID)

	if err != nil {
		return nil, err
	}

	if current.Status.State == "" {
		current.Status.State = collective.TaskStateCreated
	}

	if previous != nil || previous.Status.State != current.Status.State {
		switch current.Status.State {
		case collective.TaskStateCreated:
			current.Status.State = collective.TaskStatePending

		case collective.TaskStatePending:
			fallthrough
		case collective.TaskStateInProgress:
			for _, stage := range pipeline.Spec.Stages {
				err = tr.ReconcileStage(ctx, current, pipeline, stage)

				if err != nil {
					return nil, err
				}
			}

			mainTaskStatus := current.Status.GetTaskStatus(current.Spec.OutputStageID)

			if mainTaskStatus != nil && mainTaskStatus.State == collective.TaskPhaseStateCompleted {
				current.Status.State = collective.TaskStateCompleted
			}
		}
	}

	return forddb.CreateOrUpdate(tr.db, current)
}

func (tr *TaskReconciler) ReconcileStage(
	ctx context.Context,
	task *collective.Task,
	pipeline *collective.Pipeline,
	stage collective.Stage,
) error {
	status := task.Status.GetOrCreateTaskStatus(stage.ID)

	switch status.State {
	case collective.TaskPhaseStateCreated:
		for _, dep := range stage.DependsOn {
			depStatus := task.Status.GetTaskStatus(dep)

			if depStatus == nil {
				return nil
			}

			if depStatus.State != collective.TaskPhaseStateCompleted {
				return nil
			}
		}

		status.State = collective.TaskPhaseStatePending

	case collective.TaskPhaseStatePending:
		agent := &collective.Agent{
			Spec:   collective.AgentSpec{},
			Status: collective.AgentStatus{},
		}

		agent, err := forddb.CreateOrUpdate(tr.db, agent)

		if err != nil {
			return err
		}

		status.State = collective.TaskPhaseStateScheduled

	case collective.TaskPhaseStateScheduled:
	}

	return nil
}
