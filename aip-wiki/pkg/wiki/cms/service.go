package cms

import (
	"context"
)

type Service struct {
	pm *PageManager
}

func NewService(pm *PageManager) *Service {
	return &Service{pm: pm}
}

type EmptyRequest struct {
	Empty string `json:"empty"`
}

func (svc *Service) Empty(ctx context.Context, req *EmptyRequest) (*EmptyRequest, error) {
	return &EmptyRequest{}, nil
}
