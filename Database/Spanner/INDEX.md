# 📚 Spanner + yo 문서 인덱스

## 🎯 어떤 문서를 읽어야 할까요?

### 처음 시작하는 경우
1. **README.md** - 전체 개요, yo란 무엇인가?
2. **QUICK_START.md** - 30초 빠른 시작
3. **SETUP_CHECKLIST.md** - 설치 체크리스트

### 이미 설치했고 사용법을 알고 싶은 경우
1. **USAGE.md** - 상세 사용법
2. **YO_GUIDE.md** - yo 완전 가이드

### Docker/Spanner 설정이 궁금한 경우
1. **DOCKER_GUIDE.md** - Docker Spanner 가이드

### 빠르게 명령어만 보고 싶은 경우
```bash
make help
```

## 📖 문서 목록

### 핵심 문서

| 문서 | 내용 | 대상 |
|------|------|------|
| **README.md** | 프로젝트 전체 개요, yo 소개 | 모든 사용자 |
| **QUICK_START.md** | 30초 빠른 시작 가이드 | 처음 시작하는 사용자 |
| **SETUP_CHECKLIST.md** | 설치 단계별 체크리스트 | 처음 설치하는 사용자 |
| **USAGE.md** | 상세 사용법, 워크플로우 | 실제 개발하는 사용자 |
| **YO_GUIDE.md** | yo 완전 가이드, 고급 기능 | yo를 깊이 사용하려는 사용자 |
| **DOCKER_GUIDE.md** | Docker Spanner 설정 가이드 | Docker 관련 이슈가 있는 사용자 |

### 참고 문서

| 문서 | 내용 |
|------|------|
| **Makefile** | 모든 자동화 명령어 (주석 포함) |
| **test.sh** | API 테스트 스크립트 |
| **migrations/*.sql** | 마이그레이션 파일 예제 |

## 🎓 학습 경로

### 초보자 (처음 시작)

```
1. README.md 읽기 (5분)
   ↓
2. SETUP_CHECKLIST.md 따라하기 (10분)
   ↓
3. QUICK_START.md 실행 (3분)
   ↓
4. test.sh 실행해보기 (1분)
```

**총 소요 시간: 약 20분**

### 중급자 (개발 시작)

```
1. USAGE.md 읽기 (10분)
   ↓
2. 새 테이블 추가해보기 (15분)
   ↓
3. YO_GUIDE.md 훑어보기 (10분)
   ↓
4. 복잡한 쿼리 작성해보기 (20분)
```

**총 소요 시간: 약 1시간**

### 고급자 (깊이 활용)

```
1. YO_GUIDE.md 정독 (30분)
   ↓
2. 커스텀 타입 사용 (20분)
   ↓
3. 커스텀 템플릿 작성 (30분)
   ↓
4. Clean Architecture 적용 (1시간)
```

**총 소요 시간: 약 2시간**

## 🔍 상황별 문서 찾기

### "설치가 안 돼요"
→ **SETUP_CHECKLIST.md** 확인

### "Docker가 이상해요"
→ **DOCKER_GUIDE.md** 확인

### "yo를 어떻게 써요?"
→ **YO_GUIDE.md** 확인

### "마이그레이션은 어떻게 해요?"
→ **USAGE.md** > "개발 워크플로우" 섹션

### "생성된 코드를 어떻게 써요?"
→ **USAGE.md** > "생성된 코드 사용법" 섹션

### "빠르게 시작하고 싶어요"
→ **QUICK_START.md** 실행

### "명령어만 보고 싶어요"
→ `make help`

## 📝 문서 요약

### README.md
```
- yo란 무엇인가?
- 프로젝트 구조
- 주요 특징
- Hammer vs Wrench
- yo 생성 코드 예제
- 학습 자료
```

### QUICK_START.md
```
- 30초 빠른 시작
- 3단계 초기화
- 주요 명령어
- 팁
```

### SETUP_CHECKLIST.md
```
- 시작 전 확인사항
- 도구 설치
- Spanner 설정
- Instance/Database 생성
- 마이그레이션
- yo 실행
- 최종 체크리스트
```

### USAGE.md
```
- 전체 워크플로우
- 단계별 상세 설명
- 개발 워크플로우
- 생성된 코드 사용법
- Makefile 명령어
- 문제 해결
```

### YO_GUIDE.md
```
- yo 핵심 개념
- 기본 사용법
- 생성되는 코드 구조
- 실제 사용 예제
- 고급 기능
- yo vs ORM
- 모범 사례
```

### DOCKER_GUIDE.md
```
- 현재 Spanner 사용
- 새 Spanner 띄우기
- 상태 확인
- gcloud 설정
- Spanner CLI
- 문제 해결
- Emulator vs 실제 Spanner
```

## 🎯 핵심 명령어 (자주 사용)

```bash
# 전체 초기화
make init

# 서버 실행
make run

# 테스트
make test

# 마이그레이션 + 코드 재생성
make reset

# 도움말
make help
```

## 🔗 외부 링크

- [yo 공식 문서](https://pkg.go.dev/go.mercari.io/yo)
- [yo GitHub](https://github.com/mercari/yo)
- [Cloud Spanner 문서](https://cloud.google.com/spanner/docs)
- [Hammer GitHub](https://github.com/daichirata/hammer)
- [Wrench GitHub](https://github.com/cloudspannerecosystem/wrench)

## 💡 문서 읽기 팁

### 시간이 없다면
1. **QUICK_START.md** (3분)
2. `make init` 실행
3. `make test` 실행
4. 나중에 필요할 때 다른 문서 참고

### 제대로 배우고 싶다면
1. **README.md** 정독
2. **SETUP_CHECKLIST.md** 따라하기
3. **USAGE.md** 읽으며 실습
4. **YO_GUIDE.md** 심화 학습

### 문제가 생겼다면
1. 해당 섹션의 "문제 해결" 확인
2. **SETUP_CHECKLIST.md** 다시 체크
3. **DOCKER_GUIDE.md** 확인
4. `make info` 실행

## 🎉 시작하기

**가장 빠른 시작:**
```bash
cd /home/milman2/go-api/go/Database/Spanner
make init
```

**그리고:**
- 서버가 실행되면 → **USAGE.md** 읽기
- 문제가 있으면 → **SETUP_CHECKLIST.md** 확인
- yo에 대해 궁금하면 → **YO_GUIDE.md** 읽기

Happy Coding! 🚀

