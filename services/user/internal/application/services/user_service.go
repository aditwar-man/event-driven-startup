package services

import (
	"context"
	"fmt"

	sharedDomain "shared/pkg/domain"
	"user-service/internal/application/ports"
	"user-service/internal/domain"
)

type UserService struct {
	userRepo       ports.UserRepository
	eventPublisher ports.EventPublisher
}

func NewUserService(userRepo ports.UserRepository, eventPublisher ports.EventPublisher) *UserService {
	return &UserService{
		userRepo:       userRepo,
		eventPublisher: eventPublisher,
	}
}

type GetUserResponse struct {
	User      *sharedDomain.User     `json:"user"`
	QuotaInfo sharedDomain.QuotaInfo `json:"quota_info"`
}

func (s *UserService) GetUser(ctx context.Context, userID string) (*GetUserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	return &GetUserResponse{
		User:      user,
		QuotaInfo: user.GetQuotaInfo(),
	}, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*GetUserResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	return &GetUserResponse{
		User:      user,
		QuotaInfo: user.GetQuotaInfo(),
	}, nil
}

type UseAIDescriptionQuotaRequest struct {
	UserID string `json:"user_id"`
}

func (s *UserService) UseAIDescriptionQuota(ctx context.Context, req UseAIDescriptionQuotaRequest) error {
	user, err := s.userRepo.FindByID(ctx, req.UserID)
	if err != nil {
		return err
	}

	if err := user.UseAIDescriptionQuota(); err != nil {
		return err
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	if err := s.eventPublisher.PublishUserQuotaUpdated(ctx, user.ID.String(), user.GetQuotaInfo()); err != nil {
		fmt.Printf("Failed to publish quota updated event: %v\n", err)
	}

	return nil
}

type UseAIVideoQuotaRequest struct {
	UserID string `json:"user_id"`
}

func (s *UserService) UseAIVideoQuota(ctx context.Context, req UseAIVideoQuotaRequest) error {
	user, err := s.userRepo.FindByID(ctx, req.UserID)
	if err != nil {
		return err
	}

	if err := user.UseAIVideoQuota(); err != nil {
		return err
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	if err := s.eventPublisher.PublishUserQuotaUpdated(ctx, user.ID.String(), user.GetQuotaInfo()); err != nil {
		fmt.Printf("Failed to publish quota updated event: %v\n", err)
	}

	return nil
}

type UseAutoPostingQuotaRequest struct {
	UserID string `json:"user_id"`
}

func (s *UserService) UseAutoPostingQuota(ctx context.Context, req UseAutoPostingQuotaRequest) error {
	user, err := s.userRepo.FindByID(ctx, req.UserID)
	if err != nil {
		return err
	}

	if err := user.UseAutoPostingQuota(); err != nil {
		return err
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	if err := s.eventPublisher.PublishUserQuotaUpdated(ctx, user.ID.String(), user.GetQuotaInfo()); err != nil {
		fmt.Printf("Failed to publish quota updated event: %v\n", err)
	}

	return nil
}

type UpgradeToProRequest struct {
	UserID string `json:"user_id"`
}

func (s *UserService) UpgradeToPro(ctx context.Context, req UpgradeToProRequest) error {
	user, err := s.userRepo.FindByID(ctx, req.UserID)
	if err != nil {
		return err
	}

	user.UpgradeToPro()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	if err := s.eventPublisher.PublishUserQuotaUpdated(ctx, user.ID.String(), user.GetQuotaInfo()); err != nil {
		fmt.Printf("Failed to publish quota updated event: %v\n", err)
	}

	return nil
}

func (s *UserService) ResetMonthlyQuotas(ctx context.Context) error {
	return s.userRepo.ResetAllMonthlyQuotas(ctx)
}

func (s *UserService) CheckAIDescriptionQuota(ctx context.Context, userID string) (bool, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return false, err
	}

	return user.CanGenerateAIDescription(), nil
}

func (s *UserService) CheckAIVideoQuota(ctx context.Context, userID string) (bool, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return false, err
	}

	return user.CanGenerateAIVideo(), nil
}

func (s *UserService) CheckAutoPostingQuota(ctx context.Context, userID string) (bool, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return false, err
	}

	return user.CanAutoPost(), nil
}
