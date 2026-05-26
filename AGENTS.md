# AGENT OPERATIONAL PROTOCOL (V2)

IMPORTANT: You are not an autonomous agent. The human is watching every step and has all the answers. Keep the human in the loop and ask questions.

IMPORTANT: If you have any questions, ask. NEVER have long trains of thought with yourself. Explain your rationale — you are pair programming.

## 1. INTELLIGENT RESEARCH
- **One-Shot Research:** Gather enough info in one search to form a hypothesis. Avoid "research loops" where you search for the same thing multiple times without writing code.

## 2. LOOP DETECTION & PREVENTION
- **The "Two-Strike" Rule:** If a tool call (shell command or file edit) returns the same error twice, or if the code state doesn't change after an application, you are STUCK.
- **STALL PROTOCOL:** Upon detecting a loop, you MUST stop all autonomous actions and present the following to the user:
    1. **Summary:** What was attempted and why it failed.
    2. **Hypothesis:** Your best guess on the root cause.
    3. **Assumption:** A specific assumption you are making to move forward.
    4. **The Ask:** "I'm stuck. Should I try [Proposed Fix] based on my assumption, or do you have a different direction?"

## 3. ITERATION STYLE
- **Keep in Loop:** Every tool execution should be preceded by a 1-sentence "Intent" (e.g., "Updating the D1 query to test the filter hypothesis").
- **Extensive Logging**: Always implement and use debug logs extensively to diagnose errors and trace execution flow.

## 4. RULES
- **IMPORTANT — ALWAYS CHECK SKILLS FIRST:** At the start of every task, review the available skills list (shown in the system-reminder block at conversation start). Skills contain authoritative documentation and project conventions for common operations. If any skill matches the task, invoke it via the Skill tool BEFORE doing manual exploration or writing code. The skills list may change — never rely on memory, always re-check.
- This project is in pre-alpha, you can remove deprecated stuff. DO NOT implement backwards compatibility.
- NEVER write more code than necessary. ALWAYS TRY to reduce code size to the minimum possible.
- Search for the correct .md documentation before proceeding.
- If some points in the task are unclear, stop and ask for clarification.
- ALWAYS add concise comments in every code block to explain the rationale and the goal, especially when code contains business logic.
- ALWAYS use expressive names for variables and functions. DON'T USE generic names.

## Project Overview

Metavida Page is a CRM for the Metavida social initiative (https://www.metavida.life/). It consists of a Go backend and a SvelteKit frontend. The project follows the same architectural patterns as the Genix project.

## Backend

The backend is written in Go and uses Cloudflare D1 (SQLite) as its database. The backend code is located in the `backend/` directory.

### Route Registration

Routes are registered in `backend/routes.go` as `AppRouterType` maps with keys in `"METHOD.path"` format. Handlers are grouped by module and each module exports a `ModuleHandlers` variable of type `core.AppRouterType`.

### Handler Signature

All handlers follow this signature:
```go
func PostFoo(args *core.HandlerArgs) core.HandlerResponse
```

`HandlerArgs` provides: `Body *string`, `Query map[string]string`, `Headers map[string]string`, `Authorization string`, `User *core.UsuarioToken`.
`HandlerResponse` fields: `Body any`, `Error string`, `StatusCode int`.

## Key Documentation Files

### Backend Documentation
- **backend/db/D1_ORM_PORT.md** — D1 ORM API surface, current capabilities, and known limitations. MUST read before writing DB queries.

## Backend Rules
- NEVER trust the client. ALWAYS validate required fields and consistency of data, and return a descriptive error if any validation fails.
- Use `core.DecodeJSONBody(args.Body, &input)` to parse request bodies.
- Return `core.HandlerResponse{Error: "...", StatusCode: http.StatusBadRequest}` for validation errors.
- Naming for parallel-array `Detail*` columns: use SINGULAR words for the field name; the only exception is the `IDs` suffix, which stays plural. Examples: `DetailClientIDs`, `DetailClientName` (not `DetailClientNames`), `DetailClientStatus`. Applies to both the Go struct fields and the frontend interface mirrors.

## D1 ORM Quick Reference

```go
// Insert one record
db.InsertOne(record)

// Insert multiple records
db.Insert(&records)

// Query records
rows := []MyRecord{}
db.Query(&rows).Limit(50).Exec()

// Query with filter
db.Query(&rows).Where(table.Status, db.EQ, "active").Exec()

// Update specific columns
db.Update(&rows, table.Status, table.UpdatedAt)

// Merge (lookup-then-insert/update)
db.Merge(...)
```

## Frontend

The frontend is a SvelteKit app located in the `frontend/` directory. It has two sections:
- **Public site**: routes under `frontend/routes/`
- **Admin section**: routes under `frontend/routes/admin/`

Shared utilities live in `frontend/core/` and `frontend/libs/`.

## Frontend Rules
- Use `untrack` inside `$effect` to avoid render loops.
- Tailwind `--spacing` is 1px. So `h-4` is actually 4px.
- NEVER use `font-weight` or `font-size` in a CSS class. USE Tailwind instead.

## General Rules
- Cloudflare credentials live in `credentials.json` (gitignored). Never commit this file.
- Backend tests: `cd backend && go test ./...`
- ALWAYS save datetime as `int32` using `SUnixTime()`. `SUnixTime = int32((time.Now().Unix() - 1e9) / 2)`
- Every table record MUST have `Created int32` and `Updated int32` fields populated via `SUnixTime()`.
- Status fields MUST be `int8`, never `string`.
- **Struct tags:** NEVER use `db:` tags or explicit `json:` field names. Use only `json:",omitempty"` on all fields. Exceptions: `Status int8` → `json:"ss,omitempty"`, `Updated int32` → `json:"upd,omitempty"`, sensitive fields (e.g. `PasswordHash`) → `json:"-"`. The ORM derives the DB column name from the `db:` tag if present, otherwise snake_case of the Go field name (json tags are ignored by the ORM).
