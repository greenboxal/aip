package ford

import (
	"context"

	"github.com/greenboxal/aip/pkg/collective"
	"github.com/greenboxal/aip/pkg/ford/forddb"
)

type TaskReconciler struct {
	db forddb.Database
}

func NewTaskReconciler(
	db forddb.Database,
) *TaskReconciler {
	return &TaskReconciler{
		db: db,
	}
}

func (tr *TaskReconciler) Reconcile(ctx context.Context, previous, current *collective.Task) (*collective.Task, error) {
	if current == nil && previous != nil {
		return nil, nil
	}

	pipeline, err := forddb.Get[*collective.Pipeline](tr.db, current.Spec.PipelineID)

	if err != nil {
		return nil, err
	}

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
		status.State = collective.TaskPhaseStateScheduled

	}

	return nil
}
