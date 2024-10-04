package reservation

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"reservation-system/pkg/log"
	"reservation-system/pkg/store"

	"reservation-system/internal/domain/candidate"
)

func (s *Service) ListCandidates(ctx context.Context) (res []candidate.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("ListCandidates")

	data, err := s.candidateRepository.List(ctx)
	if err != nil {
		logger.Error("failed to select", zap.Error(err))
		return
	}
	res = candidate.ParseFromEntities(data)

	return
}

func (s *Service) AddCandidate(ctx context.Context, req candidate.Request) (res candidate.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("AddCandidate")

	data := candidate.Entity{
		FullName: &req.FullName,
		Email:    &req.Email,
		Phone:    &req.Phone,
	}

	data.ID, err = s.candidateRepository.Add(ctx, data)
	if err != nil {
		logger.Error("failed to create", zap.Error(err))
		return
	}
	res = candidate.ParseFromEntity(data)

	return
}

func (s *Service) GetCandidate(ctx context.Context, id string) (res candidate.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("GetCandidate").With(zap.String("id", id))

	data, err := s.candidateRepository.Get(ctx, id)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		logger.Error("failed to get by id", zap.Error(err))
		return
	}
	res = candidate.ParseFromEntity(data)

	return
}

func (s *Service) UpdateCandidate(ctx context.Context, id string, req candidate.Request) (err error) {
	logger := log.LoggerFromContext(ctx).Named("UpdateCandidate").With(zap.String("id", id))

	data := candidate.Entity{
		FullName: &req.FullName,
		Email:    &req.Email,
		Phone:    &req.Phone,
	}

	err = s.candidateRepository.Update(ctx, id, data)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		logger.Error("failed to update by id", zap.Error(err))
		return
	}

	return
}

func (s *Service) DeleteCandidate(ctx context.Context, id string) (err error) {
	logger := log.LoggerFromContext(ctx).Named("DeleteCandidate").With(zap.String("id", id))

	err = s.candidateRepository.Delete(ctx, id)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		logger.Error("failed to delete by id", zap.Error(err))
		return
	}

	return
}
