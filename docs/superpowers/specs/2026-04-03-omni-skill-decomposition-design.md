# Omni Engine Migration Skill Decomposition

**Date:** 2026-04-03
**Scope:** Rewrite the monolithic `omni-engine-migration` skill into 1 orchestrator + 3 sub-skills

## Problem

The current `omni-engine-migration` skill is ~300 lines. It loads entirely into context every time,
including document templates and implementation details for all phases. Most conversations only need
one phase. The skill also re-documents the brainstorming/planning/implementing cycle that superpowers
sub-skills already handle.

## Design

### Skill Architecture

```
~/.claude/skills/
  omni-engine-migration/SKILL.md        # Orchestrator (~60 lines)
  omni-engine-analyzing/SKILL.md        # Phase 1-2 (~120 lines)
  omni-engine-planning/SKILL.md         # Phase 3 (~80 lines)
  omni-engine-implementing/SKILL.md     # Phase 4 (~40 lines)
```

### Data Flow

```
analyzing --> analysis.md --> planning --> dag.md --> implementing --> brainstorming
```

Each stage reads the previous stage's document and produces its own output.

### Document Paths

All documents are engine-centric:

```
docs/migration/<engine>/analysis.md    # Output of analyzing
docs/migration/<engine>/dag.md         # Output of planning
```

Implementation specs/plans go through the normal superpowers flow
(`docs/superpowers/specs/` and `docs/superpowers/plans/`).

Existing docs in `docs/superpowers/` for cosmosdb/mongo stay as-is.

---

## Skill 1: Orchestrator (`omni-engine-migration`)

### Responsibility

Phase detection and routing. No document templates. No implementation details.

### Behavior

On load:
1. Check `docs/migration/<engine>/` for analysis.md, dag.md
2. Check engine directories under repo root for existing packages
3. Check `docs/superpowers/specs/` and `docs/superpowers/plans/` for existing work
4. Infer current phase
5. Present 2-3 options with recommendation
6. Route to the appropriate sub-skill

### Phase Detection Rules

| Condition | Inferred Phase | Route To |
|-----------|---------------|----------|
| No analysis.md, no parser in engine dir | Analyzing | omni-engine-analyzing |
| analysis.md exists, no dag.md | Planning | omni-engine-planning |
| dag.md exists with uncompleted nodes | Implementing | omni-engine-implementing |
| Engine dir has parser but no analysis features | Planning (for next feature) | omni-engine-planning |
| All DAG nodes done | Migration (bytebase imports) | Manual guidance |

### Content

- Repository table (omni, parser, bytebase)
- Engine package structure reference
- Quick Start logic (phase detection)
- Linear tracking reference (BYT-9000)
- Routes to sub-skills — does NOT duplicate their content

---

## Skill 2: Analyzing (`omni-engine-analyzing`)

### Responsibility

Phase 1-2: Analyze the legacy parser and bytebase consumption for an engine.

### Input

Engine name (from user or orchestrator).

### Output

`docs/migration/<engine>/analysis.md` with these sections:

1. **Grammar Coverage** — statement types supported (DDL, DML, DCL, etc.)
2. **Parse API Surface** — public functions, input/output types, error handling
3. **AST Types** — node hierarchy, key fields consumed by bytebase
4. **Gaps and Limitations** — unsupported syntax, workarounds
5. **Bytebase Import Sites** — file paths, functions, which parser APIs called
6. **Feature Dependency Map** — bytebase feature -> parser APIs -> data extracted
7. **Priority Ranking** — features by migration priority
8. **Minimum Viable Surface** — subset needed before bytebase can drop legacy import

### Behavior

- Use Explore agents to search `bytebase/parser` and `bytebase/bytebase`
- If repos not cloned locally, use `gh` CLI or web search
- Use worktree `investigate/<engine>` or Agent with `isolation: "worktree"`

---

## Skill 3: Planning (`omni-engine-planning`)

### Responsibility

Phase 3: Build a feature DAG from the analysis document.

### Input

`docs/migration/<engine>/analysis.md`

### Output

`docs/migration/<engine>/dag.md` with these sections:

1. **Nodes table** — node name, package path, dependencies, parallelization
2. **Execution order** — grouped by dependency level
3. **Common dependency patterns** reference (parser -> analysis, catalog; catalog + analysis -> completion, etc.)

### Behavior

- Read the analysis doc
- Identify features needed based on bytebase consumption
- Map dependencies between features
- Identify which nodes can run in parallel
- Output the DAG document

---

## Skill 4: Implementing (`omni-engine-implementing`)

### Responsibility

Phase 4: Pick a DAG node and delegate to `superpowers:brainstorming`.

### Input

`docs/migration/<engine>/dag.md`

### Output

Delegates to brainstorming, which produces specs and plans via the superpowers chain.

### Behavior

1. Read the DAG document
2. Identify next actionable node(s) — nodes whose dependencies are all complete
3. Present options to user (which node to work on next)
4. Invoke `superpowers:brainstorming` for the chosen node
5. Reference existing engine packages (pg, mysql, mongo) as patterns

This skill is intentionally thin. The brainstorming skill handles the full
spec -> plan -> implement cycle.
