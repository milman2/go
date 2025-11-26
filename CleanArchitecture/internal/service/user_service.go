package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/milman2/go-api/clean-architecture/internal/domain"
	"github.com/milman2/go-api/clean-architecture/internal/usecase"
)

// UserService - "Service" 용어를 선호한다면 이렇게도 사용 가능
// 실제로는 UserUseCase와 동일한 역할
type UserService struct {
	userRepo usecase.UserRepository
}

// NewUserService - UserService 생성자
func NewUserService(userRepo usecase.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser - 사용자 생성 (비즈니스 로직)
func (s *UserService) CreateUser(ctx context.Context, email, name string) (*domain.User, error) {
	// 1. 도메인 엔티티 생성
	user, err := domain.NewUser(email, name)
	if err != nil {
		return nil, err
	}

	// 2. 중복 체크
	existingUser, _ := s.userRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		return nil, domain.ErrUserExists
	}

	// 3. ID 생성
	user.ID = uuid.New().String()

	// 4. 저장
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser - 사용자 조회
func (s *UserService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	if id == "" {
		return nil, domain.ErrInvalidUserID
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetAllUsers - 모든 사용자 조회
func (s *UserService) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateUser - 사용자 정보 수정
func (s *UserService) UpdateUser(ctx context.Context, id, name string) (*domain.User, error) {
	// 1. 사용자 조회
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 2. 도메인 로직으로 업데이트
	if err := user.UpdateName(name); err != nil {
		return nil, err
	}

	// 3. 저장
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser - 사용자 삭제
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	if id == "" {
		return domain.ErrInvalidUserID
	}

	// 1. 존재 확인
	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 2. 삭제
	return s.userRepo.Delete(ctx, id)
}
