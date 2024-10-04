package reservation

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"reservation-system/internal/domain/recruiter"
	"reservation-system/pkg/log"
	"reservation-system/pkg/store"
)

func (s *Service) ListRecruiters(ctx context.Context) (res []recruiter.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("ListRecruiters")

	data, err := s.recruiterRepository.List(ctx)
	if err != nil {
		logger.Error("failed to select", zap.Error(err))
		return
	}
	res = recruiter.ParseFromEntities(data)

	return
}

func (s *Service) AddRecruiter(ctx context.Context, req recruiter.Request) (res recruiter.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("AddRecruiter")

	data := recruiter.Entity{
		FullName: &req.FullName,
		Email:    &req.Email,
		Phone:    &req.Phone,
	}

	data.ID, err = s.recruiterRepository.Add(ctx, data)
	if err != nil {
		logger.Error("failed to create", zap.Error(err))
		return
	}
	res = recruiter.ParseFromEntity(data)

	return
}

func (s *Service) GetRecruiter(ctx context.Context, id string) (res recruiter.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("GetRecruiter").With(zap.String("id", id))

	data, err := s.recruiterRepository.Get(ctx, id)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		logger.Error("failed to get by id", zap.Error(err))
		return
	}
	res = recruiter.ParseFromEntity(data)

	return
}

func (s *Service) UpdateRecruiter(ctx context.Context, id string, req recruiter.Request) (err error) {
	logger := log.LoggerFromContext(ctx).Named("UpdateRecruiter").With(zap.String("id", id))

	data := recruiter.Entity{
		FullName: &req.FullName,
		Email:    &req.Email,
		Phone:    &req.Phone,
	}

	err = s.recruiterRepository.Update(ctx, id, data)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		logger.Error("failed to update by id", zap.Error(err))
		return
	}

	return
}

func (s *Service) DeleteRecruiter(ctx context.Context, id string) (err error) {
	logger := log.LoggerFromContext(ctx).Named("DeleteRecruiter").With(zap.String("id", id))

	err = s.recruiterRepository.Delete(ctx, id)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		logger.Error("failed to delete by id", zap.Error(err))
		return
	}

	return
}
