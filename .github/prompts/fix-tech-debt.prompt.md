---
mode: agent
description: Fix a tech debt item from ai/TECH_DEBT.md by ID (e.g. TD-001).
tools: [codebase, editFiles, runCommands]
---

You are fixing a known tech debt item from `ai/TECH_DEBT.md`.

## Inputs required
Ask the user for:
- **Tech debt ID** (e.g. `TD-001`, `TD-003`)

---

## Process

1. Read `ai/TECH_DEBT.md` and locate the section for the given ID.
2. Read all files mentioned in the description to understand the current broken state.
3. Apply the minimal fix — no refactoring beyond what is described.
4. If the item is in **Wave 1** (XS effort), fix it in a single focused edit.
5. If the item is **Critical** severity, confirm the fix with the user before writing.
6. After editing, run:
   ```bash
   make verify
   ```
7. If all checks pass, mark the item as resolved in `ai/TECH_DEBT.md` by moving it to the **Resolved** section with a one-sentence summary of what was changed.

---

## Known item index (for reference)

| ID | Summary | Severity | Effort |
|---|---|---|---|
| TD-001 | File token repo param order swap | Critical | XS |
| TD-002 | Static token silently ignored | Critical | XS |
| TD-003 | `lastUsage` race under RLock | Critical | S |
| TD-005 | Typo `buildDetebaseClient` | High | XS |
| TD-007 | Variable shadowing in `WaitFor` | High | XS |
| TD-009 | Caller headers override `Content-Type` | High | XS |
| TD-010 | 2 000+ lines duplicated response parsing | High | L |
| TD-012 | Expired token injected after failed refresh | Medium | S |
| TD-014 | `ParseResponseBody` panics on nil response | Medium | XS |
| TD-015 | `DefaultWaitFor` timeout too short | Medium | XS |
| TD-016 | No structured logging | Medium | L |
| TD-017 | `WARN` writes to stdout | Medium | XS |
| TD-020 | Test coverage gaps | Low | XL |
| TD-021 | Create responses validate metadata ID/URI/Name | Low | M |
