#!/usr/bin/env python3
"""Offline fixture validator. Standard library only; schema JSON is normative."""
import json, re, sys
from datetime import datetime
from pathlib import Path

STAGES = {"OBSERVE", "FUSE", "MAP", "ATTACK_MAP", "RECONSTRUCT", "EXPAND"}
STATUSES = {"OBSERVED", "INFERRED", "CONTESTED", "REJECTED", "UNKNOWN"}
EVENTS = {"API_RESPONSE", "SESSION_STATE", "AUTHZ_DECISION", "CONNECTION_REFUSED", "HTTP_403", "HTTP_404", "RESOURCE_DELETED", "GRAPH_UPDATED", "HYPOTHESIS_REVISED", "CONTROL_VERIFIED", "OTHER"}
SOURCES = {"SYNTHETIC_LOG", "REDACTED_LOG", "PUBLIC_DOCUMENT", "LAB_API_RESPONSE", "ANALYST_NOTE"}
REF = re.compile(r"^[A-Za-z0-9._:-]{1,120}$")

def check_time(s):
    if not isinstance(s, str): return False
    try: return datetime.fromisoformat(s.replace("Z", "+00:00")).tzinfo is not None
    except ValueError: return False

def validate(bundle):
    errors=[]
    if not isinstance(bundle,dict): return ["bundle must be object"]
    if bundle.get("schema_version")!="1.0.0": errors.append("unsupported schema_version")
    scope=bundle.get("scope")
    if not isinstance(scope,dict): return errors+["scope must be object"]
    sid=scope.get("scope_id")
    if not isinstance(sid,str) or not REF.fullmatch(sid): errors.append("invalid scope_id")
    if scope.get("environment") not in {"SYNTHETIC","ISOLATED_LAB","AUTHORIZED_TESTBED"}: errors.append("invalid environment")
    if not scope.get("approved_by"): errors.append("missing approval")
    if not check_time(scope.get("approved_at")): errors.append("invalid approved_at")
    observations=bundle.get("observations")
    if not isinstance(observations,list) or not (1<=len(observations)<=10000): return errors+["observations must contain 1..10000 items"]
    seen=set()
    for i,obs in enumerate(observations):
        p=f"observations[{i}]"
        if not isinstance(obs,dict): errors.append(p+" must be object");continue
        oid=obs.get("observation_id")
        if not isinstance(oid,str) or not REF.fullmatch(oid): errors.append(p+" invalid observation_id")
        elif oid in seen: errors.append(p+" duplicate observation_id")
        else: seen.add(oid)
        if not check_time(obs.get("observed_at")): errors.append(p+" invalid observed_at")
        if obs.get("stage") not in STAGES: errors.append(p+" invalid stage")
        if obs.get("event_type") not in EVENTS: errors.append(p+" invalid event_type")
        if obs.get("evidence_status") not in STATUSES: errors.append(p+" invalid evidence_status")
        if obs.get("authorized_scope_id")!=sid: errors.append(p+" out-of-scope observation")
        if not isinstance(obs.get("subject_ref"),str) or not REF.fullmatch(obs["subject_ref"]): errors.append(p+" invalid subject_ref")
        source=obs.get("source")
        if not isinstance(source,dict) or source.get("source_type") not in SOURCES or not isinstance(source.get("source_ref"),str): errors.append(p+" invalid source")
        if not isinstance(obs.get("summary"),str) or not obs["summary"].strip(): errors.append(p+" missing summary")
        refs=obs.get("evidence_refs")
        if not isinstance(refs,list) or any(not isinstance(r,str) for r in refs): errors.append(p+" invalid evidence_refs")
        elif obs.get("evidence_status")=="INFERRED" and not refs: errors.append(p+" inferred item requires evidence_refs")
        c=obs.get("confidence")
        if c is not None and (not isinstance(c,(int,float)) or isinstance(c,bool) or not 0<=c<=1): errors.append(p+" invalid confidence")
    for i,obs in enumerate(observations):
        if isinstance(obs,dict):
            for ref in obs.get("evidence_refs",[]):
                if ref not in seen: errors.append(f"observations[{i}] unknown evidence reference {ref}")
    return errors

def main():
    if len(sys.argv)!=2: print("usage: python3 tools/validate_observations.py FILE",file=sys.stderr);return 2
    data=json.loads(Path(sys.argv[1]).read_text(encoding="utf-8"))
    errors=validate(data)
    if errors:
        for e in errors: print("ERROR:",e,file=sys.stderr)
        return 1
    print(f"VALID: {len(data['observations'])} observations in scope {data['scope']['scope_id']}")
    return 0
if __name__=="__main__": sys.exit(main())
