# Customer Support Conversation Dashboard

A small full-stack app: a Go REST API serving customer support conversations,
and an Angular dashboard to browse, filter, and update them.

- **Backend:** Go (standard library only)
- **Frontend:** Angular 20 (standalone components + signals)
- **Storage:** in-memory store with predefined seed data

---

## Project structure

```
.
├── backend/
│   ├── main.go                 # entry point: wires store + router, starts server
│   ├── models.go               # Conversation, Status, Priority + validation
│   ├── store.go                # in-memory store, seed data, filtering + search
│   ├── store_test.go           # filtering / search tests
│   ├── handlers.go             # the 3 endpoints + JSON/CORS helpers
│   ├── handlers_test.go        # HTTP-level tests
│   └── summary/
│       ├── summary.go          # Part 3 (no-AI task)
│       └── summary_test.go     # Part 3 tests
└── frontend/
    └── src/app/
        ├── models/conversation.ts
        ├── core/conversation.service.ts        # HTTP + state (signals)
        ├── core/conversation.service.spec.ts   # service test
        ├── features/conversation-list/
        └── features/conversation-detail/
```

---

## How to run

### Backend

Requires Go 1.22+ (developed on 1.27).

```bash
cd backend
go run .
```

The API starts on `http://localhost:8080`. Set `PORT` to change it.

### Frontend

Requires Node 20.19+ and npm.

```bash
cd frontend
npm install
npm start
```

The app starts on `http://localhost:4200` and expects the API on
`http://localhost:8080` (configured in `src/app/api.ts`).

### Tests

Backend:

```bash
cd backend
go test ./...
```

Frontend (headless Chrome):

```bash
cd frontend
npm test -- --watch=false --browsers=ChromeHeadless
```

---

## API

| Method | Path | Notes |
|--------|------|-------|
| `GET` | `/api/conversations` | supports `?status=`, `?priority=`, `?search=` |
| `GET` | `/api/conversations/:id` | `404` if not found |
| `PATCH` | `/api/conversations/:id` | body: `{ "status"?, "priority"? }` |

- Search is case-insensitive and matches `customerName`, `customerEmail`, `subject`.
- Filters combine with AND.
- `PATCH` accepts either field or both; an omitted field is left unchanged; an unknown enum value returns `400`.
- List results are sorted newest-first.

---

## Architecture notes & decisions

**Backend**

- **No web framework.** Go's standard library has matched routes on
  method + path pattern since 1.22 (`GET /api/conversations/{id}`), so a framework (gin/chi) would add a dependency for no real benefit here.
- **Two layers, not three.** `handlers.go` (HTTP) depends on a `Store`
  *interface*; `store.go` is the in-memory implementation. There's almost no business logic beyond validation, so a separate service layer would be ceremony. The interface still keeps handlers testable and would let a database implementation slot in later.
- **`Status` / `Priority` as named string types + constants.** Go has no `enum`. This is the idiomatic substitute: serializes cleanly to/from JSON, usable in `switch` and as map keys. The compiler can't guarantee a value is legal, so `.Valid()` is checked at the API boundary.
- `sync.RWMutex` guards the map because Go serves requests concurrently.

**Frontend**

- **State lives in one service** (`ConversationService`) as signals. The app is small enough that NgRx would be overkill.
- **No router.** List and detail are shown side by side; selecting a row sets a `selected` signal. Simpler than route params for this scope.
- **Plain CSS** (no component library), since visual design isn't being evaluated. Kept light: a header bar, card surfaces, colored status/priority badges, and loading/error/empty states.
- Search input is debounced (300 ms) so typing doesn't fire a request per keystroke.

---

## Time spent

| | Approx. |
|---|---|
| Total | ~2.5h |
| Of which, using AI | ~2h |
| No-AI section (Part 3) | ~20 min |

Most of the time overall went into understanding Go syntax itself (I work in Node/TypeScript day to day) — scoping rules, slice/composite literal syntax, and how table-driven tests are structured — rather than the actual logic or architecture decisions, which came together quickly once translated from patterns I already knew.

---

## What I'd improve with more time

- **Frontend component tests** for the list and detail components (rendering the loading/error/empty states, click-to-select), not just the service.
- **Optimistic updates** on `PATCH` with rollback on failure — currently the UI waits for the response.
- **Backend:** structured logging and a graceful shutdown on `SIGINT`.
- **Validation:** reject unknown JSON fields on `PATCH` rather than ignoring them.
- Wire the API base URL through Angular's environment files instead of a constant.
- A short e2e test covering the full "select → change status → list updates" flow.
- **Visual polish:** kept intentionally minimal since design wasn't the focus of the assignment — with more time I'd add responsive breakpoints and accessibility touches (focus states, color contrast) rather than a redesign.
