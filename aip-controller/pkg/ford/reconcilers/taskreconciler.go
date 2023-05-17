package reconcilers

import (
	"context"

	"go.uber.org/zap"

	collective2 "github.com/greenboxal/aip/aip-controller/pkg/collective"
	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-forddb/pkg/reconciliation"
)

type TaskReconciler struct {
	*reconciliation.ReconcilerBase[collective2.TaskID, *collective2.Task]

	logger *zap.SugaredLogger
	db     forddb.Database
}

func NewTaskReconciler(
	logger *zap.SugaredLogger,
	db forddb.Database,
) *TaskReconciler {
	tr := &TaskReconciler{
		logger: logger.Named("task-reconciler"),
		db:     db,
	}

	tr.ReconcilerBase = reconciliation.NewReconciler(
		tr.logger,
		db,
		"task-reconciler",
		tr.Reconcile,
	)

	return tr
}

func (tr *TaskReconciler) Reconcile(
	ctx context.Context,
	id collective2.TaskID,
	previous *collective2.Task,
	current *collective2.Task,
) (*collective2.Task, error) {
	if current == nil && previous != nil {
		tr.logger.Info("task deleted", "task_id", current.ID)

		return nil, nil
	}

	tr.logger.Info("entering reconciliation loop", "task_id", current.ID)

	pipeline, err := forddb.Get[*collective2.Pipeline](ctx, tr.db, current.Spec.PipelineID)

	if err != nil {
		return nil, err
	}

	if current.Status.State == "" {
		current.Status.State = collective2.TaskStateCreated
	}

	if previous != nil || previous.Status.State != current.Status.State {
		switch current.Status.State {
		case collective2.TaskStateCreated:
			current.Status.State = collective2.TaskStatePending

		case collective2.TaskStatePending:
			fallthrough
		case collective2.TaskStateInProgress:
			for _, stage := range pipeline.Spec.Stages {
				err = tr.ReconcileStage(ctx, current, pipeline, stage)

				if err != nil {
					return nil, err
				}
			}

			mainTaskStatus := current.Status.GetTaskStatus(current.Spec.OutputStageID)

			if mainTaskStatus != nil && mainTaskStatus.State == collective2.TaskPhaseStateCompleted {
				current.Status.State = collective2.TaskStateCompleted
			}
		}
	}

	return forddb.Put(ctx, tr.db, current)
}

func (tr *TaskReconciler) ReconcileStage(
	ctx context.Context,
	task *collective2.Task,
	pipeline *collective2.Pipeline,
	stage collective2.Stage,
) error {
	status := task.Status.GetOrCreateTaskStatus(stage.ID)

	switch status.State {
	case collective2.TaskPhaseStateCreated:
		for _, dep := range stage.DependsOn {
			depStatus := task.Status.GetTaskStatus(dep)

			if depStatus == nil {
				return nil
			}

			if depStatus.State != collective2.TaskPhaseStateCompleted {
				return nil
			}
		}

		status.State = collective2.TaskPhaseStatePending

	case collective2.TaskPhaseStatePending:
		agent := &collective2.Agent{
			Spec:   collective2.AgentSpec{},
			Status: collective2.AgentStatus{},
		}

		agent, err := forddb.Put(ctx, tr.db, agent)

		if err != nil {
			return err
		}

		status.State = collective2.TaskPhaseStateScheduled

	case collective2.TaskPhaseStateScheduled:
	}

	return nil
}
