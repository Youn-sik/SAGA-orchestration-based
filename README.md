# SAGA Orchestration-Based Service

분산 트랜잭션 관리를 위한 오케스트레이션 기반 SAGA 패턴 구현 서비스입니다.

## 개요

이 서비스는 마이크로서비스 환경에서 분산 트랜잭션을 관리하기 위한 SAGA 패턴의 오케스트레이션 방식을 구현합니다. 여러 서비스에 걸친 트랜잭션을 순차적으로 실행하고, 실패 시 자동으로 보상 트랜잭션을 역순으로 실행하여 데이터 일관성을 보장합니다.

## 주요 기능

- **순차적 트랜잭션 실행**: 정의된 순서대로 여러 서비스의 API를 호출
- **자동 보상 처리**: 트랜잭션 실패 시 자동으로 보상 트랜잭션 실행
- **비동기 처리**: Context Deadline 방지를 위한 비동기 실행
- **트랜잭션 추적**: 고유한 Global Transaction ID를 통한 트랜잭션 추적
- **에러 처리**: 상세한 스택 트레이스를 포함한 에러 로깅

## 기술 스택

- **언어**: Go 1.23.0
- **웹 프레임워크**: Gin-Gonic v1.10.0
- **기타 라이브러리**:
  - google/uuid: 고유 ID 생성
  - go-playground/validator: 구조체 검증

## 프로젝트 구조

```
.
├── main.go              # 애플리케이션 진입점
├── router/
│   └── http.go          # HTTP 라우터 및 미들웨어 설정
├── saga/
│   └── saga.go          # SAGA 패턴 핵심 구현
├── workflow/
│   └── workflow.go      # 워크플로우 오케스트레이션 로직
├── utils/
│   └── utils.go         # 유틸리티 함수 (응답 처리, ID 생성)
└── logger/
    └── logger.go        # 로깅 유틸리티
```

## API 명세

### SAGA 트랜잭션 제출

**엔드포인트**: `POST /api/rtu/submit`

**요청 본문**:
```json
{
  "requests": [
    {
      "http://service1/api/order": {
        "orderId": "12345",
        "amount": 10000
      }
    },
    {
      "http://service2/api/payment": {
        "orderId": "12345",
        "paymentMethod": "card"
      }
    }
  ],
  "requests_compensation": [
    {
      "http://service1/api/order/cancel": {
        "orderId": "12345"
      }
    },
    {
      "http://service2/api/payment/refund": {
        "orderId": "12345"
      }
    }
  ]
}
```

**응답**:
```json
{
  "status": "ok",
  "data": {
    "gid": "unique-transaction-id"
  }
}
```

## 동작 방식

1. **트랜잭션 시작**: 클라이언트가 실행할 요청과 보상 요청을 정의하여 제출
2. **순차 실행**: `requests` 배열의 각 요청을 순서대로 실행
3. **실패 감지**: 어느 단계에서든 실패 시 즉시 보상 프로세스 시작
4. **보상 실행**: `requests_compensation` 배열의 요청을 역순으로 실행
5. **완료**: 모든 보상이 실행된 후 종료 (개별 보상 실패 시에도 계속 진행)

## 설치 및 실행

### 요구사항

- Go 1.23.0 이상

### 설치

```bash
# 저장소 클론
git clone https://github.com/your-org/SAGA-orchestration-based.git
cd SAGA-orchestration-based

# 의존성 설치
go mod download
```

### 실행

```bash
# 서버 실행 (기본 포트: 8998)
go run main.go
```

### 빌드

```bash
# 실행 파일 빌드
go build -o saga-service
```

## 환경 설정

현재 하드코딩된 설정:
- **포트**: 8998
- **HTTP 타임아웃**: 30초
- **실행 모드**: Release (Gin)
