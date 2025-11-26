package memory

import (
	"context"
	"sync"

	"github.com/milman2/go-api/clean-architecture/internal/domain"
)

// UserRepository - 메모리 기반 리포지토리 구현 (어댑터)
// 인터페이스를 구현하여 의존성 역전 원칙 적용
type UserRepository struct {
	mu    sync.RWMutex
	users map[string]*domain.User
}

// NewUserRepository - UserRepository 생성자
func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]*domain.User),
	}
}

// Create - 사용자 생성
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; exists {
		return domain.ErrUserExists
	}

	// 복사본 저장 (불변성 보장)
	userCopy := *user
	r.users[user.ID] = &userCopy

	return nil
}

// GetByID - ID로 사용자 조회
func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	// 복사본 반환 (불변성 보장)
	userCopy := *user
	return &userCopy, nil
}

// GetByEmail - 이메일로 사용자 조회
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			userCopy := *user
			return &userCopy, nil
		}
	}

	return nil, domain.ErrUserNotFound
}

// GetAll - 모든 사용자 조회
func (r *UserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*domain.User, 0, len(r.users))
	for _, user := range r.users {
		userCopy := *user
		users = append(users, &userCopy)
	}

	return users, nil
}

// Update - 사용자 정보 수정
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return domain.ErrUserNotFound
	}

	userCopy := *user
	r.users[user.ID] = &userCopy

	return nil
}

// Delete - 사용자 삭제
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return domain.ErrUserNotFound
	}

	delete(r.users, id)
	return nil
}

