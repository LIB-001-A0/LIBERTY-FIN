# STATE-SPACE.md

# LIBERTY-FIN 상태공간 규격

**Version:** 0.2.0  
**Status:** Design Specification  
**System:** Recursive Information Weapon System

---

## 0. 범위

이 문서는 LIBERTY-FIN이 추적하는 세계 상태와 그 갱신 규칙을 정의한다. `PROTOCOL.md`가 연산 순서를 정의한다면, 이 문서는 여섯 연산의 입출력에 걸쳐 보존되는 상태의 구조를 규정한다.

고정 코어: `SIGNAL → FUSE → MAP → ATTACK MAP → RECONSTRUCT → EXPAND → SIGNAL`

공격 대상은 **상황도의 가설·엣지·증거 연결**이다. 모든 관측은 공개되었거나 접근이 허가된 자료를 사용한다.

## 1. 상태공간

시점 `t`의 정보 상태는 다음과 같다.

```text
G_t = (V_t, E_t, D_t, U_t)
```

- `V_t` — **노드:** 인물 아닌 역할·조직·기업·제도·장비·자산·사건·시스템을 포함한 개체 집합.
- `E_t` — **엣지:** 노드 사이의 유형화된 소유·계약·운영·접속·자금·권한·인과 관계.
- `D_t` — **증거:** 원시 신호, 출처, 반례, 수집 이력, 재해석 이력, 관측 시점의 집합.
- `U_t` — **불확실성:** 아직 미확정인 노드·엣지·시간·권한·경쟁 가설 및 관측 공백.

실행 상태는 아래 항목을 별도로 보존한다.

```text
R_t = {
  goal,
  adapter,
  G_t,
  signal_log,
  fusion_map,
  attack_map,
  reconstruction_diff,
  collection_queue,
  audit_log,
  cycle_id
}
```

`G_t`는 정보 상태이고 `R_t`는 운용 상태다. `GOAL`과 어댑터는 탐색 범위·허용 관측원·증거 기준을 결정하며, 세계의 관측 사실을 임의로 바꾸지 않는다.

## 2. 노드(Node)

```json
{
  "node_id": "N-001",
  "class": "ACTOR | ASSET | SYSTEM | POLICY | EVENT | LOCATION",
  "label": "...",
  "aliases": [],
  "identity_evidence": ["S-001"],
  "observed_at": "...",
  "valid_from": null,
  "valid_to": null
}
```

동명이인, 조직개편 전후의 법인, 제품명과 법인명, 과거·현재 기관은 개체 동일성 검사를 거쳐 병합한다. 같은 이름은 동일한 노드라는 증거가 아니다.

## 3. 엣지(Edge)

엣지는 **출발 노드 + 관계 유형 + 도착 노드 + 시점 + 근거**로 정의한다.

```json
{
  "edge_id": "E-001",
  "from": "N-001",
  "to": "N-002",
  "relation": "OWNS",
  "qualifiers": {"share": "4%"},
  "status": "OBSERVED",
  "process_state": "SUPPORTED",
  "evidence_ids": ["S-001"],
  "counterexample_ids": [],
  "valid_from": "...",
  "valid_to": null,
  "observed_at": "...",
  "confidence_basis": "직접 공시",
  "assumptions": [],
  "revision": 1
}
```

### 3.1 관계형

| 관계 | 뜻 | 차단할 거짓 엣지 |
|---|---|---|
| `OWNS` | 지분·자산의 법적 소유 | 소수지분 → 운영권 |
| `FUNDS` | 투자·대출·자금 지원 | 지원 → 지배 |
| `CONTRACTS_WITH` | 계약 체결 | 협의·보도 → 계약 확정 |
| `OPERATES` | 실제 운영 권한 | 설비 투자 → 운영 |
| `CONNECTS_TO` | 기술·조직의 접속 | 연동 계획 → 실제 연결 |
| `CAN_ACCESS` | 명시적 접근권한 | 접속 가능 → 실제 접근 |
| `CONTROLS` | 승인·차단·통제권 | 접근권한 → 전체 통제권 |
| `DEPENDS_ON` | 작동의 의존 관계 | 의존 → 독점 |
| `CLAIMS` | 특정 주체의 주장 | 주장 → 검증 사실 |
| `MAY_ENABLE` | 조건부 가능성 | 가능 → 실현 |
| `CAUSED` | 입증된 인과 기제 | 동시 발생 → 인과 |
| `CONTRADICTS` | 서로 충돌하는 관측 | 부분 모순 → 전체 부정 |

`OWNS`, `OPERATES`, `CONTROLS`는 독립적인 관계다. `CLAIMS(A, P)`가 관측됐다는 사실과 명제 `P`가 실제 참이라는 판정도 분리한다.

### 3.2 증거 상태(status)

- `OBSERVED` — 직접 확인한 관측 또는 명시적 기록.
- `INFERRED` — 관측에서 도출했으나 추가 가정이 필요한 연결.
- `CONTESTED` — 반례·출처 충돌로 현재 판정이 다투어지는 연결.
- `REJECTED` — 정한 증거 기준에서 배제된 연결.
- `UNKNOWN` — 판단을 위한 관측이 부족한 연결.

**증거 상태**와 **검증 절차 상태**는 다른 축이다. `process_state`는 `PROPOSED → SUPPORTED → TESTED → CONFIRMED` 또는 `CONTESTED → REVISED / REJECTED`로 추적할 수 있다. `CONFIRMED`는 절대적 진실이 아니라 현재 증거 기준의 충족을 뜻한다.

## 4. 시간축

모든 증거와 관계는 최소 두 시간을 구분한다.

- `observed_at` — 신호가 확인된 시각.
- `valid_from / valid_to` — 해당 관계가 실제 성립한 기간.
- `collected_at` — 관측을 수집·기록한 시각.

과거 도입량, 현재 운용량, 미래 예정량은 시간축이 다르다. 날짜가 불명확하면 `null` 또는 `UNKNOWN`으로 남긴다. 서로 다른 시점의 숫자를 현재 총량으로 합산하지 않는다.

## 5. 증거 계보(Provenance)

모든 엣지는 최소 하나 이상의 `evidence_id` 또는 `UNKNOWN` 근거 공백을 가진다. 직접 관측과 기사·발표의 재인용은 구분한다. 두 매체가 같은 발표문을 재인용했다면 출처 계보는 하나다.

```text
SOURCE → SIGNAL → FUSION → EDGE → HYPOTHESIS
                 ↘ COUNTEREXAMPLE ↗
```

관측 내용·파일의 원본, 출처, 시각, 해석 이력을 보존한다. 증거가 반례로 바뀌어도 최초 관측을 지우지 않는다.

## 6. 가설과 불확실성

가설 `X_i`는 여러 엣지와 가정을 포함한 **세계 설명 후보**다. 경쟁 가설을 병렬로 유지하고, 판별에 필요한 다음 관측을 기록한다.

```json
{
  "hypothesis_id": "H-001",
  "proposition": "N-001이 N-002를 운영한다",
  "supporting_edges": ["E-001"],
  "assumptions": ["지분 보유가 운영권을 수반한다"],
  "counterexamples": [],
  "status": "INFERRED",
  "unresolved_questions": ["운영계약은 누구와 체결됐는가?"]
}
```

`U_t`에는 **관측 불확실성·구조 불확실성·시간 불확실성·인과 불확실성·권한 불확실성**을 기록한다. 충분한 확률모형과 자료가 있을 때만 `P(X|D_t)`를 수치화한다. 근거 없는 확률과 정보량 숫자는 생성하지 않는다.

## 7. ATTACK MAP — 거짓 엣지 분쇄

반례 공격은 기술 시스템이 아니라 현재 상황도의 설명 구조를 겨냥한다.

**공격 좌표**

1. **시간:** 과거 기록을 현재 상태로 둔갑시켰는가.
2. **개체:** 법인·사업부·파트너·제품을 하나로 합쳤는가.
3. **출처:** 여러 보도가 한 원문을 재인용했는가.
4. **권한:** 지분·접근·운영·차단의 권한자를 혼동했는가.
5. **실행:** 기술적으로 가능한 기능을 실제 수행으로 단정했는가.
6. **인과:** 상관·선후관계를 인과로 승격했는가.
7. **부재:** 발견되지 않은 기록을 사건 부재로 해석했는가.
8. **대안:** 더 적은 가정으로 관측을 설명할 경쟁 가설이 있는가.

```json
{
  "attack_id": "A-001",
  "target_edge_id": "E-002",
  "counterexample_evidence_ids": ["S-002"],
  "attack_class": "AUTHORITY",
  "prediction_error": null,
  "self_information": null,
  "verdict": "contested",
  "reason": "소유와 운영은 별도 권한",
  "affected_hypotheses": ["H-001"]
}
```

`prediction_error`와 `self_information`은 실제 확률·관측 모형으로 계산될 때만 숫자를 넣는다. 값이 없으면 `null`이다. 반례 하나가 엣지의 증거 수준을 무너뜨렸더라도 사건 자체를 자동으로 부정하지 않는다.

## 8. RECONSTRUCT — 구조 복원

`ATTACK MAP`의 판정을 근거로 상태를 갱신한다.

1. 폐기 대상 엣지를 **활성 상황도**에서 비활성화한다.
2. 무효화 엣지에 의존하는 가설의 하류 관계를 다시 검증한다.
3. 반례를 견딘 엣지의 증거 상태와 가중치를 갱신한다.
4. 중복·모순 표현을 병합하되 출처 계보를 유지한다.
5. 불필요한 가정을 압축하고 경쟁 가설을 남긴다.
6. `reconstruction_diff`에 전후 상태와 변경 근거를 기록한다.

```json
{
  "from_state": "G-004",
  "to_state": "G-005",
  "surviving_edges": ["E-001"],
  "invalidated_edges": ["E-002"],
  "reweighted_edges": [],
  "remaining_uncertainty": ["실제 운영 주체"],
  "reason_evidence_ids": ["S-002"],
  "audit_ref": "AUD-005"
}
```

**삭제는 없다. 무효화만 있다.** 제거된 엣지는 활성 지도에서 제외하되 `audit_log`에 계속 남는다. 최소기술길이(MDL)는 설명의 불필요한 복잡성을 줄이는 기준이며 사실 여부를 단독 판정하는 규칙이 아니다.

## 9. EXPAND — 다음 SIGNAL의 좌표

`EXPAND`는 `U_t`에 남은 질문에 대해 관측 가능성·증거 판별력·비용을 평가한다.

```text
next_target = argmax_a [ I(X; Y_a | D_t) - λ × cost(a) ]
```

이 식은 설계 목표다. 기대 정보 획득량의 측정이 불가능하면 **관측 가능성, 경쟁 가설의 판별력, 확인 비용**을 명시한 근거 중심 순위를 사용한다. 가장 불확실한 지점이 가장 좋은 다음 목표라는 보장은 없다.

`collection_queue`의 각 항목은 질문, 기대되는 증거 유형, 허용된 출처, 관측 비용, 성공·실패 시 상황도 변화를 포함한다.

## 10. 운용 커널과 상태 불변식

`PROTOCOL.md`의 `operate(goal, adapter, evidence_store)`는 `R_t`를 갱신하는 제어 루프다.

```pseudo
function operate(goal, adapter, evidence_store):
    map = initial_model(goal)
    audit_log = []
    collection_queue = goal.collection_targets

    while not stop_rule(goal, map):
        signals = SIGNAL(goal, adapter, evidence_store)
        evidence_store.append(signals)
        fused = FUSE(signals, evidence_store)
        map = MAP(fused, map)
        attacks = ATTACK_MAP(map, adapter)
        evidence_store.append(attacks.counterexamples)
        map = RECONSTRUCT(map, attacks)
        collection_queue = EXPAND(map, adapter)
        audit_log.append(record_audit(signals, fused, attacks, map, collection_queue))
        goal.collection_targets = collection_queue

    return map, audit_log, collection_queue
```

위 표기는 **판단 절차를 위한 의사코드**다. `evidence_store.append`와 `record_audit`의 구체적 저장 인터페이스는 구현 단계에서 정의한다.

**상태 불변식**

- 모든 신규 엣지는 증거·시각·출처 계보와 연결된다.
- 반례는 기존 증거를 소거하지 않고 상태를 갱신한다.
- 무효화된 엣지는 활성 지도에서 제외하고 감사 기록에 남긴다.
- 새 사이클은 이전 사이클의 `G_t`, `D_t`, `U_t`를 입력으로 받는다.
- 새 `GOAL`로 전환할 때는 독립된 `cycle_id` 또는 작전 이력을 남긴다.
- 승인되지 않은 접근이나 개인정보 취득을 탐색 대상으로 생성하지 않는다.

## 11. 예시 — 소수지분과 운영권

가상 자료에서 기업 A가 기업 B 지분 4%를 보유한다는 공시 `S1`을 확인했다.

```text
SIGNAL:        S1 = A가 B 지분 4% 보유
FUSE:          다른 보도의 원출처가 S1인지 확인
MAP:           A --OWNS(4%)--> B [OBSERVED]
               A --OPERATES--> B [INFERRED]
ATTACK MAP:    지분 4%만으로 운영권이 성립하는가?
RECONSTRUCT:   OWNS 유지 / OPERATES 비활성화 / 운영권 UNKNOWN
EXPAND:        이사회 구성·주주간 계약·운영 위탁계약 조사
```

반례가 제거하는 것은 **근거 없는 `OPERATES` 엣지**다. `OWNS(4%)`라는 관측 사실은 보존한다. 다음 SIGNAL은 운영권의 실제 보유자를 판별한다.

## 12. LOCK

```text
SIGNAL → FUSE → MAP → ATTACK MAP → RECONSTRUCT → EXPAND → SIGNAL
```

상태를 버리지 않는다. 출처와 시간을 구분한다. 자신의 상황도를 먼저 공격한다. 살아남은 구조를 복원하고 다음 관측을 생성한다.

**비밀 없이 세계를 벗긴다.**
