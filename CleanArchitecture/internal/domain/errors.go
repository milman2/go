package domain

import "errors"

// 도메인 에러 정의
var (
	ErrUserNotFound  = errors.New("user not found")
	ErrUserExists    = errors.New("user already exists")
	ErrInvalidEmail  = errors.New("invalid email")
	ErrInvalidName   = errors.New("invalid name")
	ErrInvalidUserID = errors.New("invalid user id")
)

