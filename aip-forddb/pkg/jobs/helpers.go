package jobs

import (
	"context"
)

func DefineHandler[TPayload, TResult any](id string) JobHandlerKey[TPayload, TResult] {
	return JobHandlerKey[TPayload, TResult]{JobHandlerID(id)}
}

func Await[TResult any](handle TypedJobHandle[TResult]) (def TResult, _ error) {
	result, err := handle.Await()

	if err != nil {
		return def, err
	}

	return result.(TResult), nil
}

func DispatchJob[TPayload, TResult any](
	ctx context.Context,
	manager *Manager,
	handlerId JobHandlerKey[TPayload, TResult],
	payload TPayload,
) (TypedJobHandle[TResult], error) {
	spec := JobSpec{
		Handler: handlerId.ID().String(),
		Payload: payload,
	}

	handle, err := manager.DispatchEphemeral(ctx, spec)

	if err != nil {
		return nil, err
	}

	return TypedJobHandle[TResult](handle), nil
}
