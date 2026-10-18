# AGENTS.md — khaiii-api

Kakao khaiii 형태소 분석기를 HTTP API로 제공하는 서버. AI 에이전트·기여자용 규칙입니다.

## 프로젝트 한 줄

khaiii C 라이브러리를 Go로 감싼 REST API (`/v1`) + Go 클라이언트 패키지.

## 읽기 순서

1. [`README.md`](README.md)
2. [`main.go`](main.go) — 서버 엔트리, 인증 미들웨어
3. [`pkg/`](pkg/) — 외부에서 import하는 클라이언트
4. [`internal/c_khaiii/`](internal/c_khaiii/) — CGO 바인딩
5. [`Dockerfile`](Dockerfile) / [`Makefile`](Makefile)

## 레이아웃

| 경로 | 책임 |
|------|------|
| `main.go` | Gin 서버, Bearer 인증, Swagger |
| `pkg/client` | HTTP 클라이언트 |
| `internal/c_khaiii` | khaiii C API (CGO, 헤더 `khaiii/khaiii_api.h` 필요) |
| `docs/` | swag 생성 OpenAPI |
| `examples/` | 클라이언트 예제 |

## 빌드 주의

- 로컬 `go build`는 **시스템 khaiii 헤더/라이브러리**가 없으면 CGO 단계에서 실패합니다. 이미지 빌드(`make build_image`)가 기본 경로입니다.
- API 문서 갱신: `swag init` (swag CLI 필요).
- Go **1.27**.

## 인증·환경

| 변수 | 용도 |
|------|------|
| `KHAIII_APIKEY` | Bearer 토큰 (서버·클라이언트 예제 공통). 비어 있으면 인증 비활성 |
| `ROOT_PATH` | Gin root path (리버스 프록시용) |
| `KHAIII_PORT` | compose 호스트 포트 (기본 `8082`) |

로컬 실행:

```sh
cp .env.example .env   # KHAIII_APIKEY 설정
docker compose up --build
```

## 코딩 규칙

- CGO 바인딩(`internal/c_khaiii`) 변경은 최소로; API 표면은 `pkg/`와 HTTP 핸들러에서 확장.
- API 키는 env/`docker-compose`로만; `/secret/token` 파일 경로를 되살리지 말 것.
- `.env`를 커밋하지 말 것.
- 공개 엔드포인트/스키마 변경 시 `swag init` + README 예제 출력을 맞출 것.

## 하지 말 것

- 호스트에 khaiii 없이 CGO 오류를 “코드 버그”로 단정하고 대규모 리팩터
- `.env` 실비밀 커밋
- 사용자 요청 없이 커밋·푸시

## 커밋

사용자 요청 시에만. 메시지 스타일: `feat|fix|chore: …`
