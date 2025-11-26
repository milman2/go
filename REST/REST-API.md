# Gin
# Echo
# Fiber
# Beego
# Chi
# Gorilla Toolkit
# Buffalo
# FastHTTP

# 선택 기준
## 성능 최우선
- Gin, FastHTTP
## 개발 생산성
- Echo, Fiber, Buffalo
## 엔터프라이즈/풀스택
- Beego
## 마이크로서비스
- Chi, Gorilla

# Middelware
## Logging
- Gin.Logger()
- middleware.Logger()
## Error Handling
## Authentication & Authorization
- JWT
- OAuth2
- API Key
- middleware.JWTWithConfig()
## CORS (Cross-Origin Resource Sharing)
- middleware.CORS()
- cors.Default()
## Request Validation
- go-playground/validator
## Security
- HTTPS
- CSRF
- XSS
- secure
## Rate Limiting
- ulule/limiter
## Compression
- middleware.Gzip()
## Session Management
- Gorilla Sessions