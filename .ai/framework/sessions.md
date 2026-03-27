# Sessions

Related: `overview.md`, `controllers.md`, `middleware.md`, `rendering.md`

## Purpose

The framework session system in `internal/session/` provides request-scoped state, flash messaging, CSRF storage, old-input helpers, and field-level error storage for application handlers and views.

## Core Rule

Always access session state through request context:

```go
sess := session.FromContext(r.Context())
```

Do not reach into stores directly from controllers.

## Current Session Flow

1. `SessionMiddleware` loads or creates a session.
2. The middleware injects `*session.Session` into request context.
3. Controllers and middleware read and update state through the `Session` API.
4. Writes persist through the session manager and refresh the cookie.

## Current Capabilities

- file-backed or Redis-backed storage
- session ID regeneration
- flash messages
- old form input storage
- field-level validation errors
- CSRF token storage

## Cookie Behavior

The session cookie is configured by the session manager and currently behaves as follows:

- `HttpOnly` is always `true`
- `Secure` is enabled in production
- `SameSite` is `Lax` in development and `Strict` in production
- `Path` is `/`
- expiry is refreshed from session length configuration

Treat the session cookie as the auth bearer and keep its configuration framework-owned.

## Session API That Controllers Should Use

### Basic values

- `sess.Set(key, value)`
- `sess.Get(key)`
- `sess.Remove(key)`
- `sess.ID()`

### Flash messages

- `sess.SetFlash(key, value)`
- `sess.GetFlash(key)` to read and consume
- `sess.GetAllFlash()` to read and consume all
- `sess.PeekFlash(key)` and `sess.PeekAllFlash()` to inspect without consuming

### Old input

- `sess.SetOld(key, value)` before redirecting after validation failure
- `sess.GetOld(key)` or `sess.GetAllOld()` when consuming old input
- `sess.PeekOld(key)` and `sess.PeekAllOld()` when rendering without consuming yet

### Field errors

- `sess.SetFieldErrors(field, errors)`
- `sess.AddFieldError(field, message)`
- `sess.GetFieldErrors(field)` or `sess.GetAllFieldErrors()` when consuming
- `sess.PeekFieldErrors(field)` when inspecting without clearing

### Lifecycle

- `sess.Regenerate()` to rotate the session ID while preserving data
- `sess.RegenerateNew()` to rotate and clear state except framework-preserved values

Use regeneration for login, logout, and privilege changes.

## Auth Integration

Authentication state is layered on top of session state through `internal/auth/`.

- `auth.GetAuthData(sess)` retrieves the application auth payload.
- `auth.IsAuthenticated(sess)` is the guard used by middleware.
- `auth.Logout(sess)` clears auth state through session lifecycle helpers.

Controllers should not duplicate auth/session bookkeeping that already exists in the auth package.

## Form Handling Rules

For redirect-based forms:

1. Parse and validate the request.
2. Persist old input with `SetOld` for the fields you need to repopulate.
3. Persist user-facing errors through flash or field error helpers.
4. Redirect back to the form.
5. Let the next request render from session state.

This is the current pattern used by the auth flows.

## Rendering Interaction

The render layer consumes session state to build request context for templ pages. Because `View.Render` calls `GetAllFlash()` and `GetAllFieldErrors()`, those values are consumed during full-page rendering. Design controllers with that one-request lifecycle in mind.

## Session Regeneration

Session regeneration is still relevant and should remain documented because the framework implements both preserved-data and fresh-session rotation.

### `sess.Regenerate()`

Use `sess.Regenerate()` when you need a new session ID but want to preserve existing session data.

Current implementation behavior:

- generates a new session ID
- copies existing session values into the new session
- preserves the CSRF token
- deletes the old session from storage
- reissues the session cookie

Use it for:

- successful login
- privilege elevation
- other sensitive identity or authorization state changes where the session should rotate without losing state

### `sess.RegenerateNew()`

Use `sess.RegenerateNew()` when you need a completely fresh session.

Current implementation behavior:

- generates a new session ID
- clears all previous session values
- creates a fresh CSRF token
- deletes the old session from storage
- reissues the session cookie

Use it for:

- logout
- security resets
- account compromise responses
- any workflow that should fully discard prior session state

### Regeneration Rule Of Thumb

- Preserve data with `Regenerate()` when the user is continuing a valid authenticated flow.
- Clear everything with `RegenerateNew()` when the user is ending or invalidating the session.

## Logout Rule

Keep the logout guidance explicit: logouts should clear session state rather than merely removing one auth key.

The current auth package already does this:

```go
func Logout(sess *session.Session) {
	sess.RegenerateNew()
}
```

Agents should treat that as the default logout pattern.

## CSRF Token Management

CSRF handling is still relevant and should stay documented because the session layer stores and rotates CSRF values.

Available session methods:

- `sess.GetCSRF()`
- `sess.SetCSRF(token)`
- `sess.RemoveCSRF()`

Operational guidance:

- let middleware manage CSRF creation for ordinary request flows
- do not manually remove CSRF tokens during normal controller logic
- expect `RegenerateNew()` to create a fresh CSRF token automatically
- use `RemoveCSRF()` only for intentional low-level session handling, not routine app code

## Session Middleware

Session middleware still needs to be documented clearly because it is what makes request-context session access possible.

Current behavior of `SessionMiddleware`:

- loads existing session data from the configured store
- creates a new session if one does not exist
- wraps the session in `*session.Session`
- injects it into request context

Without session middleware, controllers and middleware cannot safely rely on `session.FromContext(r.Context())`.

## Practical Session Rules

Use these as enforcement guidance for agents:

- Regenerate sessions on login, logout, and privilege changes.
- Prefer `auth.Logout(sess)` over hand-rolled logout behavior.
- Let middleware manage CSRF token creation unless there is a deliberate low-level reason not to.
- Do not bypass session middleware for routes that depend on auth, CSRF, flash, or old input.
- Keep session guidance focused on the framework contract so downstream apps implement the same flow consistently.

## Session Checklist

Use this when building or reviewing session-backed flows:

- Access session only through `session.FromContext(r.Context())`.
- Apply session middleware to any route group that depends on session, auth, CSRF, flash, or old input.
- Use flash values for redirect-driven user messages.
- Use old-input helpers to preserve submitted form values across redirects.
- Use field-error helpers for field-level validation feedback.
- Regenerate the session on login, logout, and privilege changes.
- Use `Regenerate()` when data should survive the rotation and `RegenerateNew()` when it should not.
- Keep logout flows aligned with `auth.Logout(sess)`.
- Keep cookie security settings environment-aware and framework-owned.
- Avoid storing unnecessary sensitive or bulky data in the session.
- Remember that full-page rendering consumes flash and field-error state for that request.
