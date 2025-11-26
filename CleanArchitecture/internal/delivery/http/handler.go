package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/milman2/go-api/clean-architecture/internal/domain"
	"github.com/milman2/go-api/clean-architecture/internal/usecase"
)

// UserHandler - HTTP 핸들러 (프레젠테이션 레이어)
type UserHandler struct {
	userUseCase *usecase.UserUseCase
}

// NewUserHandler - UserHandler 생성자
func NewUserHandler(userUseCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// CreateUserRequest - 사용자 생성 요청 DTO
type CreateUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// UpdateUserRequest - 사용자 수정 요청 DTO
type UpdateUserRequest struct {
	Name string `json:"name"`
}

// UserResponse - 사용자 응답 DTO
type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ErrorResponse - 에러 응답 DTO
type ErrorResponse struct {
	Error string `json:"error"`
}

// respondJSON - JSON 응답 헬퍼
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError - 에러 응답 헬퍼
func respondError(w http.ResponseWriter, status int, err error) {
	respondJSON(w, status, ErrorResponse{Error: err.Error()})
}

// toUserResponse - 도메인 엔티티를 DTO로 변환
func toUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// CreateUser - 사용자 생성 핸들러
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}

	user, err := h.userUseCase.CreateUser(r.Context(), req.Email, req.Name)
	if err != nil {
		switch err {
		case domain.ErrInvalidEmail, domain.ErrInvalidName:
			respondError(w, http.StatusBadRequest, err)
		case domain.ErrUserExists:
			respondError(w, http.StatusConflict, err)
		default:
			respondError(w, http.StatusInternalServerError, err)
		}
		return
	}

	respondJSON(w, http.StatusCreated, toUserResponse(user))
}

// GetUser - 사용자 조회 핸들러
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := h.userUseCase.GetUser(r.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			respondError(w, http.StatusNotFound, err)
		case domain.ErrInvalidUserID:
			respondError(w, http.StatusBadRequest, err)
		default:
			respondError(w, http.StatusInternalServerError, err)
		}
		return
	}

	respondJSON(w, http.StatusOK, toUserResponse(user))
}

// GetAllUsers - 모든 사용자 조회 핸들러
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userUseCase.GetAllUsers(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = toUserResponse(user)
	}

	respondJSON(w, http.StatusOK, responses)
}

// UpdateUser - 사용자 수정 핸들러
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}

	user, err := h.userUseCase.UpdateUser(r.Context(), id, req.Name)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			respondError(w, http.StatusNotFound, err)
		case domain.ErrInvalidName:
			respondError(w, http.StatusBadRequest, err)
		default:
			respondError(w, http.StatusInternalServerError, err)
		}
		return
	}

	respondJSON(w, http.StatusOK, toUserResponse(user))
}

// DeleteUser - 사용자 삭제 핸들러
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.userUseCase.DeleteUser(r.Context(), id); err != nil {
		switch err {
		case domain.ErrUserNotFound:
			respondError(w, http.StatusNotFound, err)
		case domain.ErrInvalidUserID:
			respondError(w, http.StatusBadRequest, err)
		default:
			respondError(w, http.StatusInternalServerError, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

