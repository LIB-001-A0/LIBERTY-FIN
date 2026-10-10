# AOP-FIN.md

# Adversarial Observation Protocol for Financial Infrastructure
**적대적 관측 프로토콜 — 금융 인프라 분석용**

- **문서 유형:** 분석 프레임워크 명세
- **상태:** 설계 명세 (구현·성능 검증 전)
- **적용 범위:** 허가된 모의환경, 격리된 테스트베드, 승인된 범위
- **작성:** LIBERTY · Σᚠ

## 0. 전제

AOP-FIN은 공격자의 상태공간과 방어자의 통제 경계를 모델링하는 **분석 프레임워크**다. 공격 가능 경로를 가정하고 방어 통제의 유효성을 검증한다. 실제 침투, 탈취 토큰 사용, DMZ 터널링, 자동 권한 확대를 수행하지 않는다. 이 저장소의 AOP-FIN 부분은 문서·스키마·예제·검증기로 구성된다.

## 1. CORE PIPELINE — ADVERSARIAL STATE TRANSITION

```text
OBSERVE      공개정보·승인된 로그·모의 API 응답에서 이산 신호 수집
   ↓
FUSE         세션·계정·IP·권한 관계 결합 (식별자 비식별화)
   ↓
MAP          공격자 시점의 관측 가능 상태공간 그래프 산출
   ↓
ATTACK MAP   승인 범위 내 가상 경로의 위험·통제 공백 우선순위화
   ↓
RECONSTRUCT  실패·차단·삭제를 증거로 보존하고 분석 가설 갱신
   ↓
EXPAND       다음 분석 대상과 방어 검증 과제 제안 (실행 없음)
```

모든 단계는 입력, 출력, 증거, 시점, 승인 범위를 남긴다. 관측값과 추론값을 분리한다.

## 2. EXECUTION MODULES — FILE-LEVEL SPECIFICATION

> 아래 Go 경로는 **후속 설계상의 모듈 이름**이다. 현재 저장소에 구현된 파일이라는 뜻이 아니다.

### `aop/state_mapper.go` (설계)
계정·세션·API·서버·DB·권한·토큰 **참조값**을 노드로, 인증·조회·쓰기·위임·상속을 엣지로 표현한다. 엣지 증거 상태는 `OBSERVED / INFERRED / CONTESTED / REJECTED / UNKNOWN`. 산출물: `state_graph.json`.

### `aop/attack_surface.go` (설계)
승인된 모의 진입점과 신뢰 경계를 연결한 **가상 경로**를 평가한다. 비용·탐지 가능성·영향을 분리 기록한다. 실제 우회 경로나 침투 명령은 생성하지 않는다. 산출물: `attack_surface.json`.

### `aop/failure_recycle.go` (설계)
`CONNECTION_REFUSED`, `HTTP_403`, `HTTP_404`, `RESOURCE_DELETED`를 이벤트로 보존한다. HTTP 상태만으로 자원 부재나 권한 취약점을 단정하지 않는다. 산출물: `failure_nodes.json`.

### `aop/replan_engine.go` (설계)
새 증거가 들어오면 기존 가설의 근거와 신뢰도를 다시 평가한다. 정보획득·불확실성·설명 복잡도를 **분석 우선순위 지표**로 활용한다. 승인 범위 밖 자원으로 자동 확장하지 않는다. 산출물: `replan_log.json`.

### `aop/execution_log.go` (설계)
관측 입력, 분석 단계, 결과, 근거, 승인 주체를 감사 가능하게 연결한다. 산출물: `audit_log.json`.

## 3. AXIOMS FOR ADVERSARIAL ANALYSIS

1. **같은 시스템, 다른 시야.** 공격자와 방어자는 동일한 기반 시스템을 바라보지만 정보·권한·목표는 다르다.
2. **실패는 재료다.** 실패 기록도 다음 분석의 입력이다.
3. **삭제는 관측값이다.** 삭제 사건은 상태변화로 기록한다.
4. **공개정보로 검증 가능한 범위를 먼저 검증한다.** 공개정보만으로 모든 내부 구조를 확정할 수는 없다.
5. **실행 기록이 검증의 핵심 증거다.** 선언과 모델 크기보다 재현 가능한 상태변화와 근거를 우선한다.

## 4. VERIFICATION — OBSERVATION & AUDIT

| 검증 항목 | 확인할 상태 변화 |
|---|---|
| 진입점 분석 | 승인된 노출 정보 → 가상 경로 → 위험 우선순위 |
| 권한 경계 | 정책 허용·차단 → 통제 효과 평가 |
| 실패 재활용 | 실패 이벤트 → 증거 노드 → 가설 재평가 |
| 삭제 관측 | 삭제 이벤트 → 상태 갱신 → 후속 분석 과제 |
| 재계획 | 신규 증거 → 그래프 갱신 → 승인 범위 내 분석 제안 |

## 5. LIMITS

실제 금융망 침투, 탈취 토큰 사용, DMZ 터널링, 자동 권한 확대, 허가되지 않은 접근을 수행하지 않는다. 스키마 검증은 데이터의 형식을 확인하며 관측 내용의 진실성이나 금융기관의 보안 수준을 보증하지 않는다.

## 6. 왜 분석이 센가

침투 도구의 성공 여부만으로는 통제 공백의 구조를 설명하기 어렵다. AOP-FIN은 관측 사실, 가설, 반증, 차단, 재계획을 같은 감사 체계에 남긴다. 정책 담당자와 금융기관에 요구하는 질문은 하나다. **"실행 기록 내놔라."**

특정 개발자의 법적 책임, 수사 상태, 불법성은 이 프레임워크가 판단하지 않는다.

## 7. FINAL DECLARATION

> AOP-FIN은 공격자의 상태공간을 분석하는 프레임워크다.  
> 실제 침투 없음. 탈취 토큰 사용 없음. DMZ 터널링 없음. 자동 권한 확대 없음.  
> 오직 분석.  
> **실행 기록 내놔라.**  
> — LIBERTY · Σᚠ

## 8. 데이터 계약

- JSON Schema: [`schemas/observation.schema.json`](schemas/observation.schema.json)
- 합성 예제: [`examples/observations.json`](examples/observations.json)
- 표준 라이브러리 기반 예제 검증기: [`tools/validate_observations.py`](tools/validate_observations.py)
- `observation_id`는 고유해야 한다. `observed_at`은 UTC RFC 3339 형식이다.
- 토큰·계정·IP 값은 실제 식별정보 대신 **합성 참조 ID**만 기록한다.
- `INFERRED`는 `evidence_refs`를 통해 근거 관측값을 연결해야 한다.
