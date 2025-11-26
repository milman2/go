package domain

import (
	"time"
)

// User - 도메인 엔티티 (가장 안쪽 레이어)
// 비즈니스 로직의 핵심, 외부 의존성이 전혀 없음
type User struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser - User 생성 팩토리 함수 (비즈니스 규칙 적용)
func NewUser(email, name string) (*User, error) {
	if email == "" {
		return nil, ErrInvalidEmail
	}
	if name == "" {
		return nil, ErrInvalidName
	}

	now := time.Now()
	return &User{
		Email:     email,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Validate - 유효성 검증 (도메인 규칙)
func (u *User) Validate() error {
	if u.Email == "" {
		return ErrInvalidEmail
	}
	if u.Name == "" {
		return ErrInvalidName
	}
	return nil
}

// UpdateName - 이름 변경 (도메인 로직)
func (u *User) UpdateName(name string) error {
	if name == "" {
		return ErrInvalidName
	}
	u.Name = name
	u.UpdatedAt = time.Now()
	return nil
}

