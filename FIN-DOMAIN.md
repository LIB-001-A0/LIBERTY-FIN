# FIN-DOMAIN.md

# LIBERTY-FIN 목표 대상별 특화 적용 모듈

**Version:** 0.2.0  
**Status:** Design Specification  
**System:** Recursive Information Weapon System

---

## 0. 범위

`FIN-DOMAIN`은 LIBERTY-FIN의 **목표별 탐색 어댑터 규격**이다. `FIN`은 최초 적용 영역인 금융에서 출발한 명칭이며, 목표 대상을 금융으로 제한하지 않는다.

고정 코어:

```text
SIGNAL → FUSE → MAP → ATTACK MAP → RECONSTRUCT → EXPAND → SIGNAL
```

어댑터는 **무엇을 관측하고, 어떤 관계를 검증하고, 어떤 반례로 상황도를 공격할지**를 결정한다. 코어의 순서, `G_t = (V_t, E_t, D_t, U_t)`의 정의, 출처·불확실성·감사 기록은 바꾸지 않는다.

**공격 대상은 가설·모델·엣지·증거 연결이다.** 관측 범위는 공개자료와 접근 권한이 있는 자료다.

## 1. 임무 설정

어댑터는 목표를 하나의 **구조적 판별 질문**으로 정한다.

```text
GOAL = 판별할 세계 상태
TARGET = 관측 대상과 범위
BOUNDARY = 시간·지역·조직·시스템의 경계
CLAIM = 공격할 최초의 설명 또는 가설
SUCCESS = 반증 가능하고 근거가 남는 판정 기준
```

예: `기업 A가 발전소에 투자했다 → 지역 전력망을 통제한다`는 주장에서 **투자·운영·통제권의 소유자가 실제로 같은가?**

`CLAIM`은 검증 전 `INFERRED` 또는 `UNKNOWN`으로 관리한다. 가설을 사실로 승격시키는 기준은 어댑터가 정한 증거 규칙을 따라야 한다.

## 2. 어댑터 계약

모든 목표 어댑터는 아래 필드를 제공한다.

```yaml
adapter_id: "INDUSTRY"
version: "0.2.0"
mission: "목표 구조를 판별하는 질문"
target_scope:
  entities: ["관측 대상"]
  time_range: "관측 시점과 유효기간"
  geography: "적용 지역"
node_types: ["ACTOR", "ASSET", "SYSTEM", "EVENT"]
edge_types: ["OWNS", "CONTRACTS_WITH", "OPERATES", "DEPENDS_ON"]
observables:
  - "실제로 얻을 수 있는 관측값"
sources:
  allowed: ["공시", "공식 문서", "허가된 기록"]
  forbidden: ["무단 접근", "비공개 개인정보"]
evidence_standard:
  observed: "직접 관측으로 인정할 조건"
  inferred: "추론으로 남길 조건"
  rejected: "반례로 무효화할 조건"
attack_tests:
  - "거짓 연결을 공격하는 반례 질문"
expansion_targets:
  - "불확실성을 줄일 다음 관측"
outputs: ["state_map", "attack_map", "reconstruction_diff", "collection_queue"]
stop_conditions:
  - "판별 질문이 충족되거나 관측이 고갈된 상태"
```

계약은 **판단 규격**이며 실제 API 호출이나 자동 실행을 보장하지 않는다. 관측할 수 없는 항목은 값 `0`이나 `false`로 채우지 않고 `UNKNOWN`으로 유지한다.

## 3. 공통 노드 및 엣지

| 공통 노드 | 의미 |
|---|---|
| `ACTOR` | 기업·기관·운영자·공급업체·감독기구 |
| `ASSET` | 설비·시설·금융·데이터·무형자산 |
| `SYSTEM` | 소프트웨어·클라우드·단말·물리 인프라 |
| `POLICY` | 법·규제·표준·허가·조달 조건 |
| `INTERFACE` | 조직·기술·자본의 접속 경계 |
| `CAPITAL` | 투자·대출·보증·준비금·수익 귀속 |
| `EVENT` | 발표·계약·가동·변경·사고·중단 |
| `PUBLIC` | 공시·IR·입찰·인허가·기술·공식 발표 |
| `COUNTER` | 경쟁 가설·반례·통제장치·대안 구조 |

엣지 유형은 `STATE-SPACE.md`의 관계 규격을 따른다. 목표별 별칭은 추가할 수 있지만 다음 관계는 반드시 구분한다.

```text
OWNS              ≠ OPERATES
FUNDS             ≠ CONTROLS
ANNOUNCES         ≠ IMPLEMENTS
MAY_CONNECT       ≠ CONNECTS_TO
DEMO              ≠ PRODUCTION
VULNERABLE_TO     ≠ COMPROMISED_BY
CORRELATES_WITH   ≠ CAUSES
```

동일한 사실이라도 **계약일·발표일·관측일·실제 유효기간**을 다른 필드로 남긴다. 출처가 서로를 인용했다면 독립 증거의 수를 늘리지 않는다.

## 4. 목표별 작전 어댑터

### 4.1 `FINANCE` — 금융·디지털자산

**임무:** 발행·보관·정산·유통·수익·통제권의 실제 보유자를 분리해 금융회로 복원.

- **노드:** `ISSUER`, `RESERVE`, `CUSTODY`, `BANK`, `CLEARING`, `SETTLEMENT`, `EXCHANGE`, `WALLET`, `DISTRIBUTION`, `REGULATOR`
- **엣지:** 발행권·준비금 운용·수탁 계약·결제 승인·지분·수수료 귀속
- **관측:** 사업보고서·공시·약관·규제 문서·정산 규격·공식 제휴 발표
- **반례 공격:** 소수지분→경영권, 파트너 명단→계약 확정, 기술 실증→상용화, 수탁→준비금 운용 수익
- **확장:** 실제 정산은행, 수익 귀속 주체, 환매 책임, 서비스 실거래 내역
- **전용 출력:** `capital_flow_map`, `settlement_map`, `control_map`

### 4.2 `INDUSTRY` — 산업·공급망

**임무:** 설계·제조·장비·부품·전력·물류가 실제 생산으로 연결되는 구조 복원.

- **노드:** `DESIGN`, `FAB`, `SUPPLIER`, `EQUIPMENT`, `POWER`, `LOGISTICS`, `CUSTOMER`, `CONTRACT`
- **엣지:** 공급계약·생산용량·기술 의존·대체 가능성·납품 경로
- **관측:** 생산 공시·공장 인허가·장비 발주·전력 인입·고객 인증 자료
- **반례 공격:** 건설 계획→가동, 설계 채택→대량 발주, 명목 생산량→실출하
- **확장:** 생산 병목, 고객 확인, 상업 가동일, 투입 에너지
- **전용 출력:** `supply_chain_map`, `capacity_map`, `power_flow_map`

### 4.3 `TECH` — AI·컴퓨팅·기술

**임무:** 모델·반도체·메모리·연산·데이터센터·전력의 실제 실행 경로 복원.

- **노드:** `MODEL`, `COMPUTE`, `MEMORY`, `STORAGE`, `NETWORK`, `POWER`, `API`, `OPERATOR`
- **엣지:** 학습·추론 사용, 장비 보유·임차, 시스템 통합, 실측 성능, 제어 권한
- **관측:** 기술 문서·재현 가능한 실험·장비 사양·전력·API 규격
- **반례 공격:** 발표된 연산량→가동 GPU, 데모→운영 성능, 모델 가능성→제품 실행 권한
- **확장:** 실측 처리량, 메모리 병목, 연산비, 권한 분리와 통제권
- **전용 출력:** `compute_map`, `dependency_map`, `execution_map`

### 4.4 `POLICY_INFRA` — 정책·국가 인프라

**임무:** 법적 결정권·예산·실증·상용 서비스·공공/민간 연결 구조 복원.

- **노드:** `AUTHORITY`, `REGULATOR`, `BUDGET`, `PROJECT`, `VENDOR`, `STANDARD`, `OPERATOR`, `USER`
- **엣지:** 인가·예산 집행·개발 계약·표준 접속권·시범운영·실제 이용
- **관측:** 법령·예산·국회 제출자료·조달·운영 통계·서비스 약관
- **반례 공격:** 예산→성과, 지갑 개설→결제 이용, 실증 거래→상용 거래, 표준 승인→시장 채택
- **확장:** 반복 이용률, 거래당 비용, 타 사업자 접속조건, 경쟁 기술과의 비교
- **전용 출력:** `policy_map`, `budget_flow_map`, `adoption_map`

### 4.5 `INCIDENT` — 사건·보안·장애

**임무:** 관측된 사건과 가능한 침해 경로, 실제 권한 이동, 대응 조치를 분리해 사건 상황도 복원.

- **노드:** `EVENT`, `ASSET`, `ACCOUNT`, `PERMISSION`, `SESSION`, `CONTROL`, `SOURCE`, `RESPONSE`
- **엣지:** 관측된 영향·허가된 접근·권한 보유·통제조치·공식 확인
- **관측:** 공개 사고 공지·감사 결과·패치 공고·접근이 허가된 사고 보고서
- **반례 공격:** 취약점 존재→실제 침해, 제조국→데이터 유출, 공격 가능성→실행 사실
- **확장:** 공식 확인된 영향 범위, 책임 운영주체, 차단·복구 검증 결과
- **전용 출력:** `incident_map`, `permission_map`, `control_map`

`ATTACK MAP`은 사고 설명을 반증하는 **정보 검증** 단계다. 타인의 시스템을 공격하거나 침투 경로를 실행하는 단계로 변환하지 않는다.

### 4.6 `CORPORATE` — 기업 지배구조·자본

**임무:** 지분·의결권·이사회·계열사·계약·수익 귀속이 만드는 실질적 결정권 복원.

- **노드:** `SHAREHOLDER`, `BOARD`, `SUBSIDIARY`, `FUND`, `CONTRACT`, `RIGHT`, `CASHFLOW`
- **엣지:** 소유·의결권·이사 지명권·보증·내부거래·현금흐름
- **관측:** 공시·주주총회 안건·사업보고서·합작 계약·투자 조건
- **반례 공격:** 경제적 지분→의결권, 대규모 투자→경영 지배, 계열사 이름→연결 재무실체
- **확장:** 실질 의결권, 이사회 구성, 계약상 거부권, 현금흐름 귀속
- **전용 출력:** `ownership_map`, `governance_map`, `capital_flow_map`

## 5. 도메인 간 융합

하나의 목표에 복수 어댑터를 적용할 수 있다. 이때 **같은 노드를 공유하되 엣지 유형과 증거 기준을 합쳐 버리지 않는다.**

```text
GOAL: 디지털 금융망의 실제 통제 구조

FINANCE      → 발행·수탁·청산·수익
TECH         → 월렛·인증·블록체인 접속
CORPORATE    → 지분·계약·운영권
POLICY_INFRA → 인허가·표준 접근·규제

공통 상태공간 G_t에 각각의 관측·불확실성을 기록한다.
```

동일한 발표문을 여러 어댑터가 사용해도 하나의 원출처로 계산한다. 분야를 넘는 관계는 `cross_domain_edge`로 표시하고 양쪽 증거 기준을 충족해야 `OBSERVED`로 승격한다.

## 6. 여섯 연산과의 접속

| 연산 | 어댑터가 제공할 것 | 코어가 유지할 것 |
|---|---|---|
| `SIGNAL` | 합법적 출처·검색어·관측 형식 | 출처·시점·계보 기록 |
| `FUSE` | 동일 대상 판별 규칙 | 독립성·중복·시간축 검증 |
| `MAP` | 노드·엣지 온톨로지 | 상태공간·근거 상태 |
| `ATTACK MAP` | 도메인별 대표 거짓 엣지 | 반례·예측 오차·판정 기록 |
| `RECONSTRUCT` | 유효한 설명의 도메인 조건 | 무효화·재가중·복원 차분 |
| `EXPAND` | 관측 가능한 후속 목표·비용 | 기대 정보 획득량에 따른 탐색 |

`EXPAND`의 우선순위는 다음 설계 목표에 따른다.

```text
next_target = argmax_a [ I(X; Y_a | D) - λ × cost(a) ]
```

실제 정보량이 추정되지 않는 경우 **관측 가능성·판별력·비용의 근거를 서술적으로 기록**한다. 수치가 없는 상태에서 측정된 값처럼 `information_gain`을 채우지 않는다.

## 7. 공통 산출물과 품질 판정

공통 출력:

```text
state_map
attack_map
reconstruction_diff
collection_queue
audit_log
```

목표별 전용 지도를 추가할 수 있다. 단, 모든 결론은 엣지와 증거 ID로 역추적할 수 있어야 한다.

**최종 판정 기준:**

1. 핵심 주장에 근거가 연결돼 있는가.
2. 관측 시점이 달라져 발생한 거짓 연결을 제거했는가.
3. 가장 강한 반례를 투입했는가.
4. 무효화된 엣지와 가설을 감사 로그에 남겼는가.
5. 미해결 불확실성을 `UNKNOWN`으로 보존했는가.
6. 다음 SIGNAL의 수집 대상과 판별 조건을 제시했는가.

## 8. 어댑터 확장 규칙

새 목표가 등장하면 **임무 질문 → 노드 → 엣지 → 관측원 → 증거 기준 → 반례 → 확장 좌표**의 일곱 항목만 추가한다. 원래의 여섯 연산과 상태공간을 다시 정의하지 않는다.

어댑터 변경 시에는 버전·변경 이유·기존 상태에 대한 영향·감사 로그를 기록한다. 새로운 정의가 기존 상태를 뒤집으면 과거 결론도 삭제하지 않고 무효화 이력으로 연결한다.

---

## LOCK

**CORE = FIXED. TARGET = VARIABLE. EVIDENCE = TRACEABLE.**

```text
SIGNAL → FUSE → MAP → ATTACK MAP → RECONSTRUCT → EXPAND → SIGNAL
```

**목표가 바뀌어도 병기의 작동 원리는 바뀌지 않는다.**
