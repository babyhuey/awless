# DESIGN: Graph storage — should `awless` keep RDF?

**Status:** evaluation only. Not scheduled, not committed to.  
**Created:** 2026-08-06  
**Decision:** keep `triplestore`. This document is the record; it was evaluated once and not scheduled.

This document exists because the question "is there a better-maintained RDF library we could use instead of `triplestore`?" has a more interesting answer than expected: **no, and `awless` may not need RDF at all.**

---

## What `triplestore` actually provides

`awless` depends on `github.com/wallix/triplestore` (now forked to `bootswithdefer/triplestore`), ~2,087 LOC of stdlib-only Go, for five things:

| Capability | Where used |
|---|---|
| In-memory indexed graph — `Source` (mutable) → `RDFGraph` (immutable snapshot, six index maps) | `graph/` throughout |
| **On-disk cache format** — binary encoding of synced graphs in `~/.awless` | `graph/graph.go:366` (`NewAutoDecoder`), `graph/graph.go:375` (`NewDatasetDecoder`) |
| NTriples encode/decode | `awless show --format rdf`, `web/web.go:79` |
| DOT export | graph visualisation (`NewDotGraphEncoder`) |
| struct → triples, tree traversal | `TriplesFromStruct` (5 uses), `NewTree` (3 uses) |

## The finding: usage is far narrower than RDF implies

Every graph query in `awless`, counted across the whole tree:

| Method | Call sites |
|---|---|
| `WithSubjPred` | 22 |
| `WithPredObj` | 4 |
| `WithSubject` | 1 |
| **everything else** (`WithObject`, `WithSubjObj`, `WithPredicate`) | **0** |

There is **no SPARQL, no inference, and no reasoning** anywhere in `awless`. The `rdfs:`/`owl:` strings in `cloud/rdf/gen_rdf.go` are string constants used as type labels — they are not consumed by any reasoner.

`WithSubjPred(subject, predicate)` is "give me property P of entity S". That is a two-level map lookup. `awless` is using an RDF triple store as an indexed property bag plus a serialisation format.

## Why swapping RDF libraries is not an option

Checked default-branch commit activity (not `pushed_at`, which counts dependabot pushes to side branches):

| Library | Last real commit | Commits since 2025-08 | Verdict |
|---|---|---|---|
| `cayleygraph/cayley` | 2024-07-06 | **0** | Abandoned despite 15k stars. Issue #963: "Seeking for a new maintainer", self-described "zombie state" |
| `knakk/rdf` | 2026-03-17 | **1** | Barely alive |
| `deiu/rdf2go` | 2024-12-12 | **0** | Dormant |
| `piprate/json-gold` | 2026-02-23 | 13 | Active, but a **JSON-LD processor** — wrong tool |
| `wallix/triplestore` | 2019-02-19 | 0 | Dead (hence the fork) |

The Go RDF ecosystem is thin and mostly dormant; Cayley was its flagship. **Migrating to another RDF library would trade a dead dependency for a dying one.** That option is closed.

## The two real options

### Option 1 — Keep the fork (status quo)

Maintain `bootswithdefer/triplestore`. It is 2,087 LOC, stdlib-only, zero transitive dependencies, and the outstanding issues in its backlog are all small (delete an intrusive `init()`, wrap errors, add CI, modernise fuzzers).

**Pros:** near-zero cost; you already own it; no migration risk.  
**Cons:** you permanently own an RDF library to serve three map lookups.

### Option 2 — Replace RDF with typed Go structures

Model resources as typed structs with explicit indexes — e.g. `map[subject]map[predicate]value` for the `WithSubjPred` case, plus a secondary index for `WithPredObj`.

**Potentially retires:** the 2,087-LOC dependency, and much of the 3,196 LOC across `graph/` + `cloud/rdf/`.

**Pros:** removes a dependency; deletes the generated `cloud/rdf` layer; type safety instead of `interface{}` objects; almost certainly faster and lower-allocation than building six index maps per snapshot.  
**Cons:** the blockers below.

## Blockers for Option 2

These are the reasons this is a design doc and not a scheduled phase.

1. **On-disk cache migration.** `triplestore`'s binary format *is* the format of every synced graph in every user's `~/.awless` (`graph/graph.go:366`). Changing the storage model needs either a converter, or a detect-and-resync-on-first-run path. `(*triple).key()` is effectively a wire contract.
2. **Generated code.** `cloud/rdf/gen_rdf.go` and `cloud/properties/gen_properties.go` are generated from `gen/aws/properties_definitions.go`. Replacing the model means rewriting generator templates, not just output.
3. **Public output formats.** `awless show --format rdf` and the DOT export are user-facing. Dropping RDF internally means either re-implementing an NTriples serialiser for output, or removing those formats — the latter being a breaking change.
4. **Effort exceeds anything else attempted here.** Larger than the SDK v2 mocking rework, which was the biggest single task in the modernization.

## Recommendation

**Keep the fork (Option 1) for now.** Treat Option 2 as an independent initiative to be justified on its own merits — performance, type safety, or maintenance burden — rather than as dependency hygiene. The fork already resolves the immediate problem, which was having an unmaintained upstream.

**Revisit if any of these become true:**

- `triplestore` maintenance starts consuming real time.
- A graph-layer performance problem traces back to snapshot index construction.
- `awless` needs a query pattern the current three cannot express (that would be an argument to move *toward* a real graph engine, not away from one).
- The on-disk cache format has to change for an unrelated reason — that would make the migration cost incremental rather than additive, and is the cheapest moment to do this.

## If Option 2 is ever pursued

Suggested sequencing, so the risky part comes last:

1. Define the replacement model and indexes behind the **existing** `cloud.GraphAPI` interface, leaving call sites untouched.
2. Run both implementations in parallel behind a flag; assert identical query results across a real synced account.
3. Write the cache converter, plus a version marker in the on-disk format so old caches are detected rather than misparsed.
4. Keep NTriples/DOT as *output-only* serialisers so `--format rdf` survives.
5. Delete `cloud/rdf` and the `triplestore` dependency only once 1–4 are proven.
