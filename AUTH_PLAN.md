# Auth Implementation Plan

## Current State

### What exists
- **Backend**: `POST /p-user-login` endpoint is fully implemented (`security/login.go`). Returns `UserToken` (HMAC-SHA256 signed CBOR), `TokenExpTime`, `UserInfo` (base64 JSON), `UserID`, `user` object.
- **Frontend**: `Env.getToken()` reads `localStorage['metavida-token']`. `http.svelte.ts` already adds `Authorization: Bearer` to all requests and calls `Env.clearAccesos?.()` on 401. `LoginForm.svelte` component exists in `ui-components/form/` but has no submission logic and no backing route.
- **No `/login` route** exists.
- **No auth guard** on `/admin` — all admin pages are publicly accessible.

### What's missing
1. `/login` route (page + login logic)
2. Auth state helpers (`checkIsLogin`, session storage/clear)
3. Auth guard in `frontend/routes/admin/+layout.svelte`

---

## Architecture (based on Genix pattern)

```
User visits /admin/*
    ↓
admin/+layout.svelte checks checkIsLogin()
    ↓ (not logged in)
navigate("/login")
    ↓
Login page calls POST /p-user-login
    ↓ (success)
Store token + expiry + userInfo in localStorage
    ↓
navigate("/admin")
    ↓
On 401 or logout: Env.clearAccesos() → clear localStorage → navigate("/login")
```

---

## Step-by-step implementation

### Step 1 — Auth helpers in `frontend/core/env.ts`

Add to `Env`:

- **`checkIsLogin()`** — returns `2` if logged in (token present and not expired), `3` if not. Mirrors Genix's convention so shared UI components stay compatible.
- **`setSession(token, expTime, userInfo)`** — stores `metavida-token`, `metavida-token-exp`, `metavida-user` in `localStorage`.
- **`clearAccesos`** — implement (currently `null`): clear the three keys above, then `window.location.href = '/login'`.
- **`navigate(path)`** — wrapper for `window.location.href` (used by login/logout).

```ts
// localStorage keys
const TOKEN_KEY = 'metavida-token'
const TOKEN_EXP_KEY = 'metavida-token-exp'
const USER_KEY = 'metavida-user'

export const checkIsLogin = (): number => {
  if (!browser) return 0
  const token = localStorage.getItem(TOKEN_KEY)
  const exp = parseInt(localStorage.getItem(TOKEN_EXP_KEY) || '0')
  if (!token || !exp || Math.floor(Date.now() / 1000) > exp) return 3
  return 2
}
```

Update `Env.getToken` to use the same `TOKEN_KEY` constant.

---

### Step 2 — Create `frontend/routes/login/+page.svelte`

- Reuses `LoginForm.svelte` from `ui-components/form/LoginForm.svelte`.
- On mount: if `checkIsLogin() === 2`, redirect to `/admin`.
- On submit: call `POST /p-user-login` with `{ email, password }`.
- On success: call `Env.setSession(res.UserToken, res.TokenExpTime, res.UserInfo)` then navigate to `/admin`.
- On error: show `Notify.failure(...)`.

The LoginForm component needs a callback prop added (`onSubmit`) so the page can wire it to the API call, or the page can own the form state directly (simpler — do this).

**Modified `LoginForm.svelte`** — add `onlogin` callback prop:
```svelte
let { onlogin } = $props<{ onlogin: (e: {email:string, password:string}) => void }>()
```

Or keep `LoginForm` purely presentational and own the form state in the page — **preferred** since it avoids changing a shared component.

---

### Step 3 — Auth guard in `frontend/routes/admin/+layout.svelte`

Add `onMount` check:

```svelte
<script lang="ts">
  import { onMount } from 'svelte';
  import { browser } from '$app/environment';
  import { checkIsLogin } from '$core/env';

  onMount(() => {
    if (checkIsLogin() !== 2) {
      window.location.href = '/login';
    }
  });
</script>
```

This runs client-side on every admin page load. Because the static adapter serves `index.html` for all routes, the redirect happens before any admin content is rendered (same pattern as Genix).

---

### Step 4 — Backend: verify login endpoint is complete

The `POST /p-user-login` endpoint in `security/login.go` is already registered as a public route (prefix `p-`). Verify:
- `routes.go` has the handler registered.
- Response includes all fields the frontend needs: `UserToken`, `TokenExpTime`, `UserInfo`.
- Bootstrap admin (`EnsureBootstrapAdmin()`) creates the initial user from `credentials.json`.

No backend changes needed unless the endpoint is missing from `routes.go`.

---

### Step 5 — Logout button

Add a logout button to `frontend/routes/admin/+layout.svelte` sidebar that calls `Env.clearAccesos?.()`.

---

## File change summary

| File | Action |
|------|--------|
| `frontend/core/env.ts` | Add `checkIsLogin()`, `setSession()`, implement `clearAccesos`, add `navigate()` |
| `frontend/routes/login/+page.svelte` | **Create** — login page with form + API call |
| `frontend/routes/admin/+layout.svelte` | Add `onMount` auth guard + logout button |
| `backend/routes.go` | Verify `p-user-login` is registered (no change expected) |

---

## Out of scope (not in this plan)
- Token auto-refresh (Genix's 4-min interval approach)
- Fine-grained route access control (access catalog / bit-packed permissions)
- Role-based visibility inside admin
