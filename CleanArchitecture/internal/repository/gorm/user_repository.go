package gorm

import (
	"context"
	"errors"

	"github.com/milman2/go-api/clean-architecture/internal/domain"
	"gorm.io/gorm"
)

// UserModel - GORM 모델 (DB 테이블 매핑)
// Clean Architecture: 어댑터 레이어의 DB 모델
type UserModel struct {
	ID        string `gorm:"primaryKey;size:36"`
	Email     string `gorm:"uniqueIndex;not null;size:255"`
	Name      string `gorm:"not null;size:100"`
	CreatedAt int64  `gorm:"autoCreateTime"`
	UpdatedAt int64  `gorm:"autoUpdateTime"`
}

// TableName - 테이블 이름 지정
func (UserModel) TableName() string {
	return "users"
}

// UserRepository - GORM 기반 리포지토리 구현
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository - UserRepository 생성자
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// toModel - Domain Entity → DB Model 변환
func toModel(user *domain.User) *UserModel {
	return &UserModel{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAt: user.UpdatedAt.Unix(),
	}
}

// toDomain - DB Model → Domain Entity 변환
func (m *UserModel) toDomain() *domain.User {
	user := &domain.User{
		ID:    m.ID,
		Email: m.Email,
		Name:  m.Name,
	}
	user.CreatedAt = user.CreatedAt.Add(0) // time 초기화
	user.UpdatedAt = user.UpdatedAt.Add(0)
	return user
}

// Create - 사용자 생성
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	model := toModel(user)
	result := r.db.WithContext(ctx).Create(model)
	if result.Error != nil {
		// GORM 에러를 도메인 에러로 변환
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return domain.ErrUserExists
		}
		return result.Error
	}
	return nil
}

// GetByID - ID로 사용자 조회
func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var model UserModel
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&model)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, result.Error
	}
	
	return model.toDomain(), nil
}

// GetByEmail - 이메일로 사용자 조회
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var model UserModel
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&model)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, result.Error
	}
	
	return model.toDomain(), nil
}

// GetAll - 모든 사용자 조회
func (r *UserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	var models []UserModel
	result := r.db.WithContext(ctx).Order("created_at DESC").Find(&models)
	
	if result.Error != nil {
		return nil, result.Error
	}
	
	users := make([]*domain.User, len(models))
	for i, model := range models {
		users[i] = model.toDomain()
	}
	
	return users, nil
}

// Update - 사용자 정보 수정
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	model := toModel(user)
	result := r.db.WithContext(ctx).Model(&UserModel{}).
		Where("id = ?", user.ID).
		Updates(model)
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	
	return nil
}

// Delete - 사용자 삭제
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&UserModel{}, "id = ?", id)
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	
	return nil
}

// ===== GORM 고급 기능 예제 =====

// GetWithPagination - 페이지네이션
func (r *UserRepository) GetWithPagination(ctx context.Context, page, pageSize int) ([]*domain.User, int64, error) {
	var models []UserModel
	var total int64
	
	// 전체 개수 조회
	if err := r.db.WithContext(ctx).Model(&UserModel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 페이지네이션 조회
	offset := (page - 1) * pageSize
	result := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&models)
	
	if result.Error != nil {
		return nil, 0, result.Error
	}
	
	users := make([]*domain.User, len(models))
	for i, model := range models {
		users[i] = model.toDomain()
	}
	
	return users, total, nil
}

// Search - 이름으로 검색 (LIKE)
func (r *UserRepository) Search(ctx context.Context, keyword string) ([]*domain.User, error) {
	var models []UserModel
	
	result := r.db.WithContext(ctx).
		Where("name LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Order("created_at DESC").
		Find(&models)
	
	if result.Error != nil {
		return nil, result.Error
	}
	
	users := make([]*domain.User, len(models))
	for i, model := range models {
		users[i] = model.toDomain()
	}
	
	return users, nil
}

// Transaction - 트랜잭션 예제
func (r *UserRepository) CreateBatch(ctx context.Context, users []*domain.User) error {
	// GORM 트랜잭션 사용
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, user := range users {
			model := toModel(user)
			if err := tx.Create(model).Error; err != nil {
				return err // 롤백
			}
		}
		return nil // 커밋
	})
}

