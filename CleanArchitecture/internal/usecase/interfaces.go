package usecase

import (
	"context"

	"github.com/milman2/go-api/clean-architecture/internal/domain"
)

// UserRepository - 리포지토리 인터페이스 (포트)
// Use Case 레이어가 외부 레이어에 의존하지 않도록 인터페이스 정의
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetAll(ctx context.Context) ([]*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
}
