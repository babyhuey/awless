# PLAN.md

Dependency-ordered execution plan across the three repositories in this project.

This plan lives in `awless` because `awless` is the consumer at the root of the dependency graph — sequencing is driven by what it needs from its dependencies.

**Referenced backlogs:**

| Repo | Backlog | Local path |
|---|---|---|
| `awless` | [`ISSUES.md`](./ISSUES.md) | `~/awless` |
| `triplestore` | `ISSUES.md` | `~/triplestore` |

Items are cited as `repo#ID` — e.g. `triplestore#D1`, `awless#B5`.

**`awless-scheduler` is retired.** The fork and its local clone were deleted on 2026-08-06, so it has no backlog file. Everything still needed is captured in Phase 4 below, and all of it lives in `awless`. Note that `awless` continues to build today because its `go.mod` requires **upstream** `github.com/wallix/awless-scheduler v0.0.6`, which is unaffected by the fork's deletion — that requirement is what Phase 4 removes.

---

## Dependency graph

```
                    ┌─────────────────┐
                    │   triplestore   │  stdlib only, no deps
                    │  (fork, live)   │
                    └────────┬────────┘
                             │ 14 files in graph/, web/, aws/services/
                             ▼
                    ┌─────────────────┐
                    │     awless      │  the CLI — builds today
                    │      (CLI)      │
                    └────────┬────────┘
                             │ still requires upstream wallix/awless-scheduler
                             │ v0.0.6 for client/ + model/ only  ── Phase 4 removes
                             ▼
              ┌──────────────────────────────┐
              │  wallix/awless-scheduler     │  upstream, unmaintained
              │  (fork DELETED 2026-08-06)   │  vestigial dependency
              └──────────────────────────────┘
```

Three facts drive the whole ordering:

1. **`triplestore` has no reverse dependencies.** Nothing depends on it except `awless`. It can be fixed in complete isolation and is therefore the safest starting point.
2. **`awless` builds today** (verified: `go build ./...` succeeds). Nothing here is on fire; this is improvement work, not recovery.
3. **The remaining scheduler dependency is inert.** `awless` imports only `client/` and `model/` from the *upstream* module, and neither imports `awless` back — so the daemon's long-standing breakage never affected `awless`, and deleting the fork did not either. Phase 4 removes the requirement outright.

---

## Critical path

Only one hard chain exists. Everything else is parallelizable:

```
triplestore#D1 (go.mod)  →  triplestore#I1 (CI)  →  triplestore#{B1,D3,D4}  →  triplestore#I5 (tag v1.0.0)
                                                                                        │
                                                                                        ▼
                                                              awless: adopt forked triplestore
```

`triplestore#D1` is the single highest-leverage item in the entire project: without a `go.mod` that repo cannot be built, tested, linted, or fuzzed, and the fork cannot be adopted cleanly.

---

## Phase 0 — Decisions — **ALL RESOLVED 2026-08-06**

Every judgment call that shapes later phases is settled. Recorded here as the decision log.

### 0.1 — Module paths — **DECIDED: rename both to `bootswithdefer`**

| Repo | Decision |
|---|---|
| `awless` | `module github.com/wallix/awless` → **`github.com/bootswithdefer/awless`** |
| `triplestore` | fresh `go.mod` declaring **`github.com/bootswithdefer/triplestore`** |

Rejected: keeping `wallix/...` with a `replace` directive. `replace` is ignored when the module is consumed as a library, and it would leave one deliberate piece of `wallix` coupling after removing it everywhere else.

**Why `awless` needed renaming at all:** `README.md:33` tells users `go install github.com/wallix/awless@latest`, which installs **upstream's 2018 code**, not this fork. Renaming is the only option that makes `go install` correct — pointing the README at `bootswithdefer` without renaming the module simply fails, because Go resolves modules by declared path.

**Blast radius:** 519 occurrences of `github.com/wallix/awless/...` across 188 files, plus 14 `github.com/wallix/triplestore` sites. Mechanical via `gofmt -r`. Both rewrites land in **one sweep**, as a single commit, **before Phase 3** — it touches nearly every file and would conflict with anything else in flight.

**Generated-code caveat:** 7 of the affected files are `gen_*.go`, and 2 of the triplestore sites are in `aws/services/gen_services.go`. The generator templates under `gen/aws/generators/` must be updated and the output regenerated — never hand-edited.

**Ordering constraint:** `triplestore#D1` (`go mod init`) must land **before** the `awless` import rewrite, or there is no resolvable module to point at.

### 0.2 — Retire or restore the scheduler daemon — **DECIDED: RETIRE (done)**

**Decision made 2026-08-06: retire.** The fork and local clone were deleted the same day, so this is settled rather than pending. Rationale:

- The daemon had been unbuildable since ~2018: it imported `awless/aws/driver` and `awless/template/driver`, both deleted during the SDK v2 migration, and called `DefaultTemplateEnv()` / `NewDriver()`, neither of which exists in `awless` any more. Upstream's README said as much — it opened with `DEPRECATED UNTIL SCHEDULER IS UP TO DATE WITH LATEST AWLESS VERSION`.
- Restoring it would have been the largest single task across all repos — a port to the `aws/spec` + `template/env` API, not a dependency bump.
- Its security model was poor: an unauthenticated `POST /tasks` endpoint executing templates against the host's AWS credentials. This was the highest-severity network exposure in the project, and deletion is a stronger fix than hardening.
- The scheduling flags in `awless` are non-functional today anyway, since they require a daemon that cannot be built.
- Deferred/scheduled infrastructure changes are adequately served by EventBridge Scheduler, cron + `awless`, or CI schedules.

**Remaining consequence:** `awless` still carries the dependency and the dead flags. Phase 4 removes them, which also drops a direct dependency and resolves part of `awless#D6`.

### 0.3 — golangci-lint scope — **DECIDED: migrate now, expand later**

`awless`'s lint job is **broken right now**: `.golangci.yml` is v1 format and CI installs golangci-lint `@latest`, which resolves to v2.12+. Verified locally — it fails with `unsupported version of the configuration`.

- **Phase 1A:** v2 config migration **only**, keeping the existing 10 linters. Pin to **`v2.12.2`** exactly (matches local install so CI and local agree). `@latest` is what caused this break.
- **Phase 5 (`awless#I9`):** add `errcheck`, `errorlint`, `noctx`, `gocritic` incrementally.

Rejected: doing both at once. `errorlint` alone flags all 127 `%s`-instead-of-`%w` sites, so the commit that unbreaks CI could not merge until 127 unrelated fixes landed — leaving CI red *longer*.

**Interaction:** v2 moves `gofmt`/`goimports` from `linters` to a `formatters` section. Their `local-prefixes` setting must become `github.com/bootswithdefer/awless` once 0.1 lands.

### 0.4 — Commit workflow — **DECIDED: direct to `master`**

Commit directly to `master` in both repos, in logical units. No feature branches, no PRs.

- **Phases 1–2:** one commit per issue (unrelated correctness fixes — keeps `git bisect` useful)
- **Phase 5:** one commit per sweep (`%s`→`%w` across 127 sites is one logical change)
- Conventional Commits format throughout

This overrides the standing "do not commit to main or master" convention, by explicit instruction.

### 0.5 — Version bump for the flag removal — **DECIDED: `v1.1.0`**

Phase 4 removes `--run-in`, `--revert-in`, the hidden `scheduler` command, and the `scheduler.url` config key. Treated as a minor bump, not major: semver governs working behavior, and none of these can succeed today without a buildable daemon.

Rejected: `v2.0.0` (nothing functional is lost) and a deprecation cycle (writing new code whose only job is to reject flags that already fail).

CHANGELOG entry required regardless — the flags currently *appear* in `--help`, so their disappearance needs to be discoverable.

**Release grouping:** land this in the same release as the 0.1 module rename, so users absorb one disruption (`go install` path change + flag removal) instead of two.

### 0.6 — Graph storage / de-RDF evaluation — **DECIDED: design doc**

Captured in [`DESIGN-graph-storage.md`](./DESIGN-graph-storage.md) rather than scheduled as a phase. `awless` uses only three of `triplestore`'s query patterns and no RDF semantics, but replacing it requires an on-disk cache migration for every user's `~/.awless`. Documented, not committed to.

### 0.7 — Untracked tooling output — **DECIDED: gitignore**

Add `.codegraph/` to `.gitignore` in both repos, folded into the Phase 1A housekeeping commit alongside `awless#B1` (which removes and ignores the two committed binaries).

---

## Phase 1 — Unblock CI and tooling — **DONE**

Highest value per unit effort. Nothing here depends on anything else, so all three tracks run in parallel.

### 1A — `awless`: repair CI

| Item | Why first |
|---|---|
| `awless#B5` | **Lint job is broken.** v1→v2 config migration; full replacement config is in the issue. Also pin the golangci-lint version — `@latest` caused this. |
| `awless#I17` | `main.go` discards `Execute()`'s error, so `awless` exits 0 on failure. Two-line fix, breaks shell scripting today. |
| `awless#B1` | Remove two committed binaries (11 MB) and `.gitignore` them. |
| `awless#D3` | Delete stale `.travis.yml`. |
| Phase 0.7 | Add `.codegraph/` to `.gitignore` (fold into the `awless#B1` housekeeping commit). |
| `awless#I18` | Pin the Go toolchain; switch CI to `go-version-file: go.mod`. |

### 1B — `triplestore`: make the repo buildable

| Item | Note |
|---|---|
| `triplestore#D1` | **`go mod init github.com/bootswithdefer/triplestore`** (path fixed by 0.1). Unblocks everything else in that repo, and must land before the `awless` import rewrite. Add `.codegraph/` to `.gitignore` here too. |
| `triplestore#I1` | Add CI **with `-race` from day one** — this library has an `atomic.Value` snapshot cache and `awless` previously hit a data race fixed by bumping this dependency. |
| `triplestore#D6`, `#D7` | Delete `.travis.yml`; point README badges at the fork. |

### 1C — `awless-scheduler`: no work (repo deleted)

The fork and clone were deleted, so there is nothing to build, test, or lint. All removal work lives in Phase 4 and touches only `awless`.

**Exit criteria:** all three repos have green CI. `triplestore` is buildable for the first time.

---

## Phase 2 — Security and correctness — **DONE**

Real defects. Ordered by user impact, and grouped where fixes share code.

### 2A — Secret handling

`awless#B6` — `awless` writes full command lines (including `password=`) into its BoltDB template log via `cmd.String()`, and stores `Fillers` verbatim. `awless log` replays them.

**Fix:** mark sensitive params in the command struct tags, redact in `cmd.String()` and `Fillers` before marshalling.

*(The sibling issue in the scheduler — template content written `0644` — is void now that the scheduler is retired and deleted. Previously these had to be fixed together so a secret was not merely relocated; that no longer applies.)*

### 2B — Network exposure

| Item | Note |
|---|---|
| `awless#B9` | Web UI binds `0.0.0.0:8080`, unauthenticated, no timeouts. Default to loopback; add explicit server timeouts. |

*(The scheduler's unauthenticated write endpoint — `POST /tasks`, which executed templates against the host's AWS credentials — was the highest-severity network exposure in the project. Retiring and deleting the repo eliminated it outright, which is a stronger fix than hardening.)*

### 2C — Crashes and leaks

| Item | Note |
|---|---|
| `awless#B7` | `--sort` panics on mixed-type columns. Fall back to string compare; never panic in display code. A sort panic was already fixed once (commit `71665081`), so this class is proven live. |
| `awless#B8` + `awless#I13` | Goroutine leak on first error in 5 fan-out sites, plus unbounded concurrency causing API throttling on large accounts. **Same fix** — migrate both to `errgroup` with `SetLimit`. Do as one change. |
| `awless#I12` | Restore `-race` in CI (non-blocking at first). Deliberately removed in commit `6830d3d0`; its absence is why B8 went unnoticed. Land **after** B8 so it starts closer to green. |
| `triplestore#B1` | Library seeds global `math/rand` in `init()` — intrusive side effect on every importer, and `rand.Seed` is deprecated. Delete the `init()`. |

*(The scheduler's dropped-error handlers and its `close(eventc)` send-on-closed-channel race are void — retired.)*

### 2D — Low-risk cleanups

`awless#B11` (generated source written `0666`), `triplestore#B2`/`#B3` (panics in library code).

**Exit criteria:** no known crash, leak, or plaintext-secret path.

---

## Phase 3 — Dependency modernization — **DONE** (3C partially: I6 needs a repo setting)

Depends on Phase 1 (CI must be green to trust these changes) and Phase 0.1 (module path decision).

### 3A — Adopt the forks

Sequenced so each fork is worth depending on before `awless` switches to it:

1. **`triplestore`**: finish `#D3` (error wrapping) and `#B1`, then `#I5` — tag `v1.0.0`. First tagged release in the library's history.
2. **`awless`**: the **single module-rename sweep** from 0.1 — rewrite 519 `wallix/awless` self-imports *and* 14 `wallix/triplestore` imports in one commit, update `gen/aws/generators/` templates, regenerate, then set `go.mod` to `github.com/bootswithdefer/awless` and require the tagged `bootswithdefer/triplestore`. Also update `.golangci.yml` `local-prefixes` and `README.md:33`.

The scheduler's 3 import sites are handled by deletion in Phase 4, not by repointing.

Verify: `go mod graph | grep -c wallix` trends toward zero.

### 3B — Replace archived dependencies (`awless#D6`)

Independent of 3A, parallelizable. Ordered by effort-to-value:

| Order | Change | Files | Note |
|---|---|---|---|
| 1 | `boltdb/bolt` → `go.etcd.io/bbolt` | 2 | Archived Mar 2019. API-compatible fork; mostly an import path change. |
| 2 | `yaml.v2` → `go.yaml.in/yaml/v3` | 1 | `go-yaml` archived Apr 2025 — **`yaml.v3` is also unmaintained**, so v2→v3 is not enough. Target module is already in the tree transitively. |
| 3 | `oklog/ulid` → `/v2` | 2 | Straightforward. |
| 4 | `mitchellh/ioprogress` | 1 | 2018 orphan, ~100 LOC. Inline rather than depend. |
| 5 | `src-d/go-git.v4` → `go-git/go-git/v5` | 1 | **Biggest payoff** (`awless#D6a`): collapses 26 stale transitives including archived `pkg/errors`, and removes a request for 2019-era `x/crypto`. API changed v4→v5. Target v5; v6 is alpha. |

### 3C — Supply chain

| Item | Note |
|---|---|
| `awless#I10` | Add `govulncheck` to CI. Do it **after** 3B so the initial run is not dominated by findings from already-known-archived deps. |
| `awless#B4` | Pin CI actions to commit SHAs. |
| `awless#I6` | Dependabot auto-merge for patch updates — 30+ AWS SDK modules generate real noise. |

**Exit criteria:** no archived direct dependencies; `govulncheck` green in CI.

---

## Phase 4 — Retire the scheduler — **DONE**

Decision recorded in Phase 0.2: **retire**. This phase is a deletion task, not a port. The scheduler's own backlog is gone along with the repo; nothing below depends on it.

Sequenced here rather than earlier because it is not blocking anything — but it is self-contained and could be pulled forward at any time. It has no dependency on Phases 1–3.

### 4.1 — Remove scheduler support from `awless`

Exact scope, verified against the current tree. **Note:** there is no `--schedule` flag in this fork; scheduling intent is derived from `--run-in` / `--revert-in` via `isSchedulingMode()`. Earlier drafts of this plan said otherwise — that came from the upstream scheduler README describing upstream `awless`, not this fork.

| File | Action |
|---|---|
| `commands/scheduler.go` | **Delete entire file.** Holds the hidden `scheduler` command, `--list-tasks` / `--list-failures` flags, and `printTasks`. Sole importer of `awless-scheduler/model`. |
| `commands/run.go:40` | Remove `"github.com/wallix/awless-scheduler/client"` import. |
| `commands/run.go:55-56` | Remove `scheduleRunInFlag`, `scheduleRevertInFlag` vars. |
| `commands/run.go:66-67` | Remove the `run-in` / `revert-in` flag registrations on `runCmd`. |
| `commands/run.go:80-81` | Remove the same two `PersistentFlags` registrations applied to every generated driver command. |
| `commands/run.go:585-607` | Remove `scheduleTemplate()`. |
| `commands/run.go:630-638` | Remove `isSchedulingMode()`. |
| `commands/runner.go:85` | Remove the `if isSchedulingMode() { return false, scheduleTemplate(...) }` block, leaving `return true, nil`. |
| `config/config.go:28` | Remove the `schedulerURL` const. |
| `config/config.go:53` | Remove the `schedulerURL` entry from `configDefinitions`. |
| `config/getters.go:41-47` | Remove `GetSchedulerURL()`. |
| `config/config_extra_test.go:114` | Remove the `"scheduler.url"` expectation from the defaults assertion. |
| `config/config_extra_test.go:233-248` | Remove `TestGetSchedulerURL`. |
| `go.mod` | Remove the `github.com/wallix/awless-scheduler v0.0.6` require. |

Then:

```sh
go mod tidy
go build ./... && go test ./...
```

**Watch for:** removing the `client` import from `run.go` must not orphan `config` or `logger` — both are used elsewhere in that file, so they stay. `go build` will catch it either way.

**Test handling:** the two `config_extra_test.go` edits delete coverage for functionality that no longer exists, which is correct. They are not test removals to force a pass — verify `go test ./config/...` passes on its own merits afterward.

**Version:** `v1.1.0` per decision 0.5, released together with the 0.1 module rename.

**User-visible change:** `awless run --run-in 2h` and `--revert-in 4h` stop existing, as does the hidden `awless scheduler` command and the `scheduler.url` config key. Since all of them require a daemon that cannot be built, none of them work today — so this removes dead flags rather than working features. Worth a CHANGELOG entry regardless, because the flags currently *appear* in `--help`.

### 4.2 — Fork retirement (done)

The `bootswithdefer/awless-scheduler` fork and its local clone were deleted on 2026-08-06. No further action.

For anyone looking for the capability later: `awless` templates can be scheduled with EventBridge Scheduler, cron + `awless`, or CI schedules. The timed-revert behavior has no direct substitute, but `awless revert` on a logged template execution covers the manual case.

---

## Phase 5 — Broad modernization

Large, mechanical, low-risk. Deliberately last: these touch many files, so doing them earlier would conflict with every other change.

| Item | Note |
|---|---|
| `awless#I8` (part 1) | `%s` → `%w` in 127 places. **Semantic fix, not style** — error chains are silently broken today while `errors.Is`/`errors.As` are used in 61 places. Also `triplestore#D3`. |
| `awless#I9` | Expand linters incrementally: `errcheck` → `errorlint` → `noctx` → `gocritic` (the deferred half of decision 0.3). `errcheck` would have caught `awless#I17`. Sequence **after** the `%w` sweep above, since `errorlint` flags all 127 of those sites. |
| `awless#I8` (part 2) | `interface{}` → `any` (527 sites). Mechanical. Also `triplestore#I2`. |
| `awless#D2` | `strings.Title` (deprecated Go 1.18) in 12 sites. |
| `awless#D4` | Thread `context.Context` through commands and fetchers (101 `context.Background()` sites). Prerequisite for `awless#I4` (signal handling) and `#I16` (removing 9 `os.Exit` calls). |
| `awless#I4` + `#I16` | Signal handling and single exit path. Depend on D4. |
| `awless#I11` | GoReleaser — replaces 200-line `release.go`, adds arm64, checksums, changelog. Resolves `#I1`, `#D7`. |
| `awless#I15` | CI check for generated-code drift. Enforces the "never hand-edit `gen_*.go`" rule. |
| `awless#I14`, `triplestore#D4` | Migrate fuzzers to native Go fuzzing. `triplestore`'s parser handles untrusted input from disk and network — the better target of the two. |
| `awless#D1`, `#D8` | Rework SDK v2 mocking; unblocks 30 dead acceptance-test factories. Largest item here. |
| `awless#I7` | Test coverage (currently ~22.4%). Easier after D4/I16, since `os.Exit` is a major reason `commands/` is untestable. |
| `awless#D5`, `#I3` | CRUD commands for the 8 list-only services. Feature work — schedule independently. |

---

## Suggested order at a glance

```
Phase 0   decisions ......... module path, scheduler fate (DECIDED: retire), CI-red acknowledgement
Phase 1   unblock ........... awless#B5 CI ‖ triplestore#D1 go.mod        (scheduler: no work)
Phase 2   correctness ....... secrets (awless#B6) → exposure → crashes/leaks
Phase 3   dependencies ...... adopt triplestore fork → replace archived deps → govulncheck
Phase 4   retire scheduler .. delete from awless (14 edits, 1 file deleted) → archive fork
Phase 5   modernization ..... %w → linters → any → context → release tooling
```

**If only one thing gets done:** `awless#B5`, because CI lint is broken right now and nothing else can be validated until it is fixed.

**If only one week is available:** all of Phase 1, plus `awless#B6` and `awless#B7` from Phase 2. That gets green CI on all three repos, stops writing secrets to disk, and removes a live CLI crash.

---

## Cross-cutting notes

- **Never hand-edit `gen_*.go`** in `awless`. Any change touching generated output must go through `gen/aws/` definitions or `gen/aws/generators/` templates, then regenerate (`cd gen/aws/generators && go run *.go`). This affects Phase 3A (2 of 14 import sites) and Phase 5 (`strings.Title` appears in generator templates).
- **`triplestore`'s `(*triple).key()` format is a wire contract.** Changing it invalidates every binary-encoded graph already cached in users' `~/.awless`. Treat as a breaking change requiring a migration path.
- **`triplestore` has no upstream backlog.** Upstream `master` is identical to the commit `awless` pins (`compare` reports `ahead_by: 0`), so the fork starts exactly at what is in use. There is nothing to catch up on — only the absence of a maintainer.
- **Test before and after each phase.** `awless` builds green today; that is the baseline to protect.

---

## Status — 2026-08-06

Phases 1–4 complete. Phase 5 outstanding.

| Phase | State |
|---|---|
| 0 Decisions | all resolved |
| 1 Unblock CI | done — B5, I17, B1, D3, I18; triplestore D1, I1, D6, D7 |
| 2 Correctness | done — B6, B9, B7, B8/I13 (partial), I12, B11; triplestore B1 |
| 3 Dependencies | done — triplestore tagged v0.1.0 and adopted, module renamed, D6/D6a, I10, B4 |
| 4 Retire scheduler | done — v1.1.0, last wallix dependency gone |
| 5 Modernization | not started |

**Deviations from this plan, all recorded in the relevant commits:**

- `triplestore` was tagged **v0.1.0**, not `v1.0.0`. Its `B2`, `B3` and `D5` are pending
  breaking API changes, so 1.0.0 would promise stability that does not hold yet.
- `boltdb` → `bbolt` was pulled forward from Phase 3B into Phase 2, because boltdb
  crashes under the race detector (`checkptr`) and therefore blocked `I12`.
- `B8`/`I13` are **partially** done. The issue claimed 5 fan-out sites; `manual_fetchers.go`
  had never been audited and holds 5 more.
- Generated files were rewritten in place during the module rename rather than regenerated,
  because the generators produce non-compiling output. Escalated as `I15`.
- Phase 4 removed **two** `isSchedulingMode()` call sites; this plan listed one.

**Still needs input:** `I6` (Dependabot auto-merge) requires `allow_auto_merge` on the
repository, currently `false`. Everything else in Phase 5 is unblocked.
