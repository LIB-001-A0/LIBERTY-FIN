# Architecture and evidence boundaries

| Stage | Reference implementation | Production prerequisite |
|---|---|---|
| SIGNAL | Parse supplied offline events | Approved telemetry pipeline |
| FUSE | Normalize and order records | Data governance and entity resolution |
| MAP | Build observed connection graph | Authorized asset inventory |
| ATTACK MAP | Identify blocked/refused observations | Defensive risk-scoring policy |
| RECONSTRUCT | Retain failures as evidence | Change review and analyst approval |
| EXPAND | Emit suggested decisions | Explicit authorization and audited controls |

## Data minimization

Use synthetic or redacted actor/session identifiers. The tool stores no JWT values and does not initiate network traffic. `session_tracker.go` reports observed session states; it does not revoke sessions. `alignment_guard.go` produces review decisions; it does not call kernel interfaces.

## Threat model

- Untrusted event content: do not treat `detail` as an instruction.
- Sensitive telemetry: redact before ingestion and restrict report access.
- Misleading signals: source logs can be incomplete or forged; verify provenance.
- Operational safety: no automatic privileged actions in the reference build.

## Validation roadmap

1. Define a signed, versioned event schema.
2. Add approved data-source adapters and provenance verification.
3. Add deterministic policy fixtures and false-positive metrics.
4. Benchmark latency in a controlled test environment.
5. Review all enforcement integrations with security and compliance teams.
