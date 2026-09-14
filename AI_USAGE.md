# AI Usage

I used Claude (via Claude Code) throughout Parts 1 and 2. Part 3 was done without any AI assistance. Below are the decisions that mattered.

---

## 1. Backend: framework vs. standard library

**Goal:** Decide how to build the REST API without over-engineering it.

**Prompt:** "Should I use gin/chi/echo for a 3-endpoint API, or is the standard library enough? What are the tradeoffs?"

**Result:** Claude recommended the standard library — since Go 1.22, `http.ServeMux` matches on method + path pattern (`GET /api/conversations/{id}`) and exposes path params via `r.PathValue`, which covers everything this API needs.

**Decision:** Accepted. The backend has zero dependencies. If routing needs grew
(middleware groups, param constraints) a router would start to pay off, but not
at this size.

---

## 2. Project structure: pushed back on the layering

**Goal:** Lay out the backend packages.

**Prompt:** "Sketch a project structure" — followed by me questioning it:
"isn't a separate `tests/` folder better? do we need `cmd/` and `internal/`?"

**Result:** The first proposal had a `cmd/server` entry point, an `internal/`
tree split into `conversation`, `summary`, and `httpx` packages, and a
controller/service/repository split.

**Decision:** Changed it. I flattened the backend to a single `package main`
with one file per concern (`models`, `store`, `handlers`) and a two-layer
split (HTTP handlers + a `Store` interface). For a ~3-hour assignment the
multi-package layout was ceremony that would read as cargo-culting rather than
judgment. Claude also confirmed that Go *requires* test files to sit next to
the code they test (`store_test.go`), so the "separate tests folder" idea I had
from Node doesn't apply here.

---

## 3. Frontend: disagreed with the AI's dropdown binding

**Goal:** Bind the detail panel's Status/Priority `<select>`s to the currently selected conversation.

**Prompt:** "Build the detail component with dropdowns to change status and priority."

**Result:** The generated template used `[value]="c.status"` on the `<select>` elements, with `<option>`s produced by a `@for` loop.

**Decision:** Rejected after testing. When I opened the app, a conversation with priority `HIGH` showed `LOW` in the dropdown. `[value]` on a `<select>` isn't a reliable way to set the selected option in Angular when the options are rendered dynamically — the binding can apply before the options exist. I switched to
`[ngModel]="c.priority"` + `(ngModelChange)="..."` (importing `FormsModule`), which is the correct pattern and fixed it. This is the one place the AI's output looked correct but was subtly wrong, and only manual testing caught it.

---

## 4. Frontend: checking whether one service spec was enough test coverage

**Goal:** Decide whether `conversation.service.spec.ts` alone was sufficient,
or whether the list/detail components also needed their own tests.

**Prompt:** "Is it enough to have only one test file, for the service only?"

**Result:** Claude suggested leaving it as-is and documenting the missing
component tests under "What I'd improve with more time" in the README — as
scoped-out future work rather than something to rush in under time pressure.

**Decision:** Left as-is for this submission — the service is tested, the
components aren't, and that gap is documented rather than silently missing.

---

## Other AI use

- Explaining Go syntax and idioms (I work in Node/TypeScript, not Go), which I then reviewed line by line rather than taking on trust.
- Generating the seed data and the boilerplate for tests, which I checked and adjusted (e.g. making search assertions check field-by-field rather than hard-coding expected IDs).
- Summarizing the existing test scenarios in `handlers_test.go` and
  `store_test.go` back to me, as a sanity check on coverage rather than to write new tests.
