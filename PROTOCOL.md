# PROTOCOL.md

# LIBERTY-FIN 실행 규격

**Version:** 0.2.0  
**Status:** Design Specification

---

## 0. 범위

이 문서는 LIBERTY-FIN 6단계 연산의 실행 규격을 정의한다.

- `SIGNAL` — 첩보 수집
- `FUSE` — 정보 융합
- `MAP` — 상황도
- `ATTACK MAP` — 공격 지도
- `RECONSTRUCT` — 구조 복원
- `EXPAND` — 확장

각 단계는 이전 단계의 상태를 입력으로 받고, 다음 단계의 상태를 출력으로 낸다.

이 문서는 자동 실행 엔진의 사양서가 아니다. **판단 절차의 규격서다.**

## 1. 공통 원칙

**원칙 1. 상태는 버리지 않는다.** 모든 관측, 추론, 반례, 무효화는 상태로 남는다. 실패한 경로도 다음 탐색의 입력이 된다.

**원칙 2. 출처와 관측시점을 구분한다.** 사실과 추론을 구분한다. 관측시점이 다른 신호는 다른 신호로 취급한다.

**원칙 3. 자기 가설을 먼저 공격한다.** 검증되지 않은 가설은 확정 상태로 승격되지 않는다. 반례를 견딘 엣지만 RECONSTRUCT에 들어간다.

**원칙 4. 비밀은 필요 없다.** 모든 SIGNAL은 공개·합법적 관측에서 나온다. 시스템 무단침입, 계정 탈취, 사적 데이터 접근은 범위 밖이다.

**원칙 5. 다음 탐색은 불확실성이 최대인 지점에서 시작한다.** 확실한 곳은 다시 보지 않는다. 모순이 남은 곳, 근거가 약한 곳, 관측되지 않은 곳으로 간다. 단, 실제 다음 목표 선정에는 관측 가능성과 정보 획득 비용을 반영한다.

## 2. SIGNAL — 첩보 수집

### 2.1 정의

잡음이 섞인 원시 관측값을 수집하고, 출처와 관측시점을 기록한다.

### 2.2 입력

- 목표 (`GOAL / TARGET`)
- 관측 가능 범위
- 기존 상태 `D_t`

### 2.3 처리

1. 공개·합법적 출처에서 신호를 수집한다.
2. 각 신호에 출처, 관측시점, 수집방법을 태깅한다.
3. 사실(`fact`)과 추론(`inference`)을 구분한다.
4. 중복 신호를 식별한다.
5. 신호를 `signal_log`에 기록한다.

### 2.4 출력

```json
{
  "signal_id": "...",
  "source": "...",
  "observed_at": "...",
  "collected_at": "...",
  "type": "fact | inference",
  "content": "...",
  "confidence": null,
  "provenance": "..."
}
```

### 2.5 검증

- 출처가 확인되었는가
- 관측시점이 기록되었는가
- 사실과 추론이 구분되었는가
- 중복이 식별되었는가

## 3. FUSE — 정보 융합

### 3.1 정의

서로 다른 출처의 신호를 교차해 조건부 불확실성을 낮춘다.

### 3.2 입력

- `signal_log`
- 기존 상태 `D_t`

### 3.3 처리

1. 출처 독립성을 평가한다.
2. 시간축을 정렬한다.
3. 동일 행위자·동일 사건을 식별한다.
4. 서로 다른 신호가 같은 방향을 가리키는지 확인한다.
5. 정보 획득에 기여하는 조합을 `fusion_map`에 기록한다.

### 3.4 출력

```json
{
  "fusion_id": "...",
  "signals": ["signal_id_1", "signal_id_2"],
  "relation": "supports | contradicts | orthogonal",
  "mutual_information": null,
  "conditional_entropy_before": null,
  "conditional_entropy_after": null
}
```

### 3.5 검증

- 출처 독립성이 평가되었는가
- 중복 신호가 식별되었는가
- 불확실성이 감소했는가, 또는 남은 공백이 특정되었는가
- 모순 신호가 숨겨지지 않았는가

### 3.6 주의

독립 관측이라고 해서 상호정보량이 자동으로 증가하지 않는다. 중복 신호는 독립적인 정보 획득으로 계산하지 않는다. 수치 정보량은 정해진 확률모형이 있을 때만 기입한다.

## 4. MAP — 상황도

### 4.1 정의

노드와 엣지로 세계 상태를 표현하고, 근거가 허용하는 범위에서 사후분포를 구성한다.

### 4.2 입력

- `fusion_map`
- 기존 상태 `G_t = (V_t, E_t, D_t, U_t)`

### 4.3 처리

1. 노드를 식별한다. (기업, 조직, 시스템, 사건, 장소)
2. 엣지를 식별한다. (소유, 계약, 운영, 접속, 의존, 인과)
3. 각 엣지에 근거와 유효기간을 부여한다.
4. 엣지 상태를 분류한다.
   - `OBSERVED` — 직접 관측
   - `INFERRED` — 추론
   - `CONTESTED` — 모순
   - `REJECTED` — 반례로 무효화
   - `UNKNOWN` — 미해소
5. `state_map`을 생성한다.

### 4.4 출력

```json
{
  "state_id": "...",
  "nodes": [],
  "edges": [
    {
      "from": "...",
      "to": "...",
      "type": "...",
      "status": "OBSERVED",
      "evidence": [],
      "valid_from": "...",
      "valid_to": null
    }
  ],
  "uncertainty": []
}
```

### 4.5 검증

- 모든 확정 엣지에 근거가 있는가
- 유효기간이 기록되었는가
- 모순 엣지가 `CONTESTED`로 표시되었는가
- 미해소 항목이 `U_t`에 남았는가

## 5. ATTACK MAP — 공격 지도

### 5.1 정의

상황도의 취약한 엣지에 반례를 투입해 예측 오차를 검증한다.

### 5.2 공격 대상

- 가설
- 모델
- 인과관계
- 증거 연결

시스템 무단침입, 계정 탈취, 사적 데이터 접근은 범위 밖이다.

### 5.3 입력

- `state_map`
- 반례 후보

### 5.4 처리

1. 가장 취약한 엣지를 식별한다.
   - 근거가 약한 엣지
   - 출처가 단일한 엣지
   - 시간적 모순이 있는 엣지
2. 각 엣지에 반례를 투입한다.
3. 예측 오차를 검증한다.
4. 확률모형이 있으면 반례의 자기정보량을 계산한다.
5. 취약한 가정을 드러낸다.
6. `attack_map`을 생성한다.

### 5.5 출력

```json
{
  "attack_id": "...",
  "target_edge": "...",
  "counterexample": "...",
  "prediction_error": null,
  "self_information": null,
  "verdict": "survived | contested | rejected"
}
```

### 5.6 검증

- 반례가 실제로 엣지를 공격했는가
- 예측 오차가 측정되었거나 검증 공백이 기록되었는가
- 판정이 근거를 갖는가
- 살아남은 엣지가 실제로 반례를 견뎠는가

## 6. RECONSTRUCT — 구조 복원

### 6.1 정의

반례를 견딘 엣지를 재가중하고, 살아남은 구조를 최소기술길이 원리에 따라 압축한다.

### 6.2 입력

- `state_map`
- `attack_map`

### 6.3 처리

1. `REJECTED` 엣지를 활성 상황도에서 무효화한다.
2. `CONTESTED` 엣지를 재검증한다.
3. `OBSERVED`·`INFERRED` 엣지를 재가중한다.
4. 중복 엣지를 병합한다.
5. 불필요한 가설을 압축한다.
6. **증거가 지탱하는 구조만 남긴다.**
7. `reconstructed`를 생성한다.

### 6.4 출력

```json
{
  "reconstructed_id": "...",
  "surviving_edges": [],
  "removed_edges": [],
  "compressed_hypotheses": [],
  "remaining_uncertainty": []
}
```

### 6.5 검증

- `REJECTED` 엣지가 활성 상황도에서 무효화되었는가
- 살아남은 엣지가 증거에 의해 지탱되는가
- 불필요한 가설이 압축되었는가
- 남은 불확실성이 `U_t`에 반영되었는가

### 6.6 주의

압축은 삭제가 아니다. 제거된 엣지도 감사 로그에 남는다.

## 7. EXPAND — 확장

### 7.1 정의

남은 불확실성과 관측 가능성을 계산해 기대 정보 획득량이 최대인 지점에서 다음 SIGNAL을 탐색한다.

### 7.2 입력

- `reconstructed`
- 관측 비용
- 접근 가능성

### 7.3 처리

1. 남은 불확실성을 식별한다.
2. 각 지점의 관측 가능성을 평가한다.
3. 관측 비용을 산정한다.
4. 기대 정보 획득량을 계산하거나 판별력으로 평가한다.
5. 다음 탐색 좌표를 결정한다.
6. `collection_plan`을 생성한다.

### 7.4 출력

```json
{
  "plan_id": "...",
  "next_targets": [
    {
      "target": "...",
      "expected_information_gain": null,
      "cost": null,
      "priority": null
    }
  ]
}
```

### 7.5 설계 목표식

```text
next_target = argmax_a [ I(X; Y_a | D) - λ × cost(a) ]
```

### 7.6 주의

위 식은 설계 목표이며, 실제 최적화나 엔트로피 감소량이 측정된다는 뜻은 아니다.

## 8. 재귀

EXPAND의 출력은 다음 SIGNAL의 입력이 된다.

```text
SIGNAL → FUSE → MAP → ATTACK MAP → RECONSTRUCT → EXPAND → SIGNAL → ...
```

각 사이클은 이전 사이클에 남은 불확실성을 다음 사이클의 입력으로 재투입한다. 이것이 재귀적 베이즈 갱신의 설계 원리다.

### 운용 커널

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

위 함수는 개념적 판단 루프다. 실제 실행 엔진이나 수집 자동화 구현의 완료를 뜻하지 않는다.

## 9. 감사 로그

모든 단계는 다음을 기록한다.

- `PROVENANCE` — 출처
- `UNCERTAINTY` — 불확실성
- `AUDIT LOG` — 감사 로그

제거된 엣지, 무효화된 가설, 실패한 경로도 감사 로그에 남는다.

**삭제는 없다. 무효화만 있다.**

## 10. 금지 사항

LIBERTY-FIN은 다음을 하지 않는다.

- 시스템 무단침입
- 계정 탈취
- 사적 데이터 접근
- 비공개 정보의 무단 수집
- 자동 실행 엔진의 무단 배포

공격 대상은 가설이지, 시스템이 아니다.

## 11. LOCK

```text
SIGNAL → FUSE → MAP → ATTACK MAP → RECONSTRUCT → EXPAND → SIGNAL
```

재귀적 베이즈 갱신.

**비밀 없이 세계를 벗긴다.**

공격 대상은 가설이지, 시스템이 아니다.

다음 SIGNAL 나와도 된다.
