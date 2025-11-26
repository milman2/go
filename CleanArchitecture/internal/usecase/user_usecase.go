package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/milman2/go-api/clean-architecture/internal/domain"
)

// UserUseCase - 사용자 관련 유스케이스 (애플리케이션 비즈니스 규칙)
type UserUseCase struct {
	userRepo UserRepository
}

// NewUserUseCase - UserUseCase 생성자
func NewUserUseCase(userRepo UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

// CreateUser - 사용자 생성 유스케이스
func (uc *UserUseCase) CreateUser(ctx context.Context, email, name string) (*domain.User, error) {
	// 1. 도메인 엔티티 생성 (비즈니스 규칙 적용)
	user, err := domain.NewUser(email, name)
	if err != nil {
		return nil, err
	}

	// 2. 중복 체크 (애플리케이션 비즈니스 규칙)
	existingUser, _ := uc.userRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		return nil, domain.ErrUserExists
	}

	// 3. ID 생성
	user.ID = uuid.New().String()

	// 4. 저장
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser - 사용자 조회
func (uc *UserUseCase) GetUser(ctx context.Context, id string) (*domain.User, error) {
	if id == "" {
		return nil, domain.ErrInvalidUserID
	}

	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetAllUsers - 모든 사용자 조회
func (uc *UserUseCase) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	users, err := uc.userRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateUser - 사용자 정보 수정
func (uc *UserUseCase) UpdateUser(ctx context.Context, id, name string) (*domain.User, error) {
	// 1. 사용자 조회
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 2. 도메인 로직으로 업데이트
	if err := user.UpdateName(name); err != nil {
		return nil, err
	}

	// 3. 저장
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser - 사용자 삭제
func (uc *UserUseCase) DeleteUser(ctx context.Context, id string) error {
	if id == "" {
		return domain.ErrInvalidUserID
	}

	// 1. 존재 확인
	_, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 2. 삭제
	return uc.userRepo.Delete(ctx, id)
}
