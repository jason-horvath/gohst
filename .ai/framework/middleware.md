# Middleware

Related: `overview.md`, `routes.md`, `sessions.md`, `rendering.md`

## Purpose

Middleware in `internal/middleware/` provides reusable request wrapping for the framework. Controller route groups opt into middleware explicitly from `RegisterRoutes()`, while the outer app router applies global wrappers.

## Signature

Framework middleware uses the standard shape:

```go
func MiddlewareName(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        next.ServeHTTP(w, r)
    })
}
```

Use middleware only for cross-cutting request concerns.

## Intent

This file should be treated as implementation guidance for how middleware is composed, ordered, and applied in the framework. It is not just a catalog of functions.

## Chaining Model

`middleware.Chain(handler, middlewares...)` wraps in reverse order, so the last middleware passed becomes the outermost request wrapper.

Use an order like this for browser route groups:

```go
middleware.Chain(
    mux,
    session.SM.SessionMiddleware,
    middleware.CSRF,
    middleware.Logger,
    middleware.Guest, // or middleware.Auth when needed
)
```

That means the wrapped handler becomes:

```go
session.SM.SessionMiddleware(
    middleware.CSRF(
        middleware.Logger(
            middleware.Guest(mux),
        ),
    ),
)
```

So request flow enters `SessionMiddleware` first, then `CSRF`, then `Logger`, then `Guest` or `Auth`, then the handler.

## Why Order Matters

Middleware order is part of the behavior contract.

- short-circuiting middleware can prevent all inner middleware and handlers from running
- session-dependent middleware must have session access available when they execute
- outer wrappers define what gets logged, recovered, denied, or cached

When changing middleware order, treat it as a behavioral change rather than a cosmetic refactor.

## Correct And Incorrect Ordering

Use examples like these when reasoning about middleware order.

Correct when auth or guest middleware needs session state:

```go
middleware.Chain(
    mux,
    session.SM.SessionMiddleware,
    middleware.CSRF,
    middleware.Logger,
    middleware.Auth,
)
```

In that chain, session is available before CSRF and auth logic runs.

Incorrect if session-dependent middleware would run before session is attached:

```go
middleware.Chain(
    mux,
    middleware.Auth,
    middleware.CSRF,
    middleware.Logger,
)
```

That shape leaves `Auth` and `CSRF` without session middleware in front of them.

## Framework Middleware Surface

### Recover

- Registered at the framework router level through `internal/routes.RegisterRouter`.
- Logs stack traces server-side.
- Shows detailed error output only in development.
- Returns a generic error message in production.

`Recover` is also paired with `RecoverGoroutine(name, fn)` for background goroutines that are outside the request chain. Use that pattern for any framework or app goroutine that must not crash the process on panic.

### SessionMiddleware

- Loads or starts the session.
- Wraps raw session data in `*session.Session`.
- Injects the session into request context.
- Must run before handlers that expect `session.FromContext`.

Any middleware that reads auth state, flash state, or CSRF session data depends on this behavior.

### CSRF

- Ensures a session token exists.
- Validates state-changing requests using form value first, then `X-CSRF-Token` header.
- Uses `crypto/subtle.ConstantTimeCompare`.
- Depends on session context already being present when executed.

This still applies even though page rendering is templ-based now. The templ change affects how forms render the token, not whether CSRF middleware is needed.

### Logger

- Logs request metadata and timing.
- Useful as a per-route-group diagnostic layer.

The current implementation is intentionally simple: method, path, and elapsed time.

### Auth

- Requires authenticated session state.
- Sets a flash error and redirects to `/auth/login` when access is denied.

This middleware is for authenticated access control, not role authorization.

### Guest

- Prevents authenticated users from guest-only flows.
- Redirects authenticated users away from login and registration routes.

If a downstream project uses a different post-login home, document that in project guidance.

### Role

- Lives in `internal/middleware/role.go`.
- Requires an authenticated user and checks for allowed roles through the auth role interface.
- Redirects unauthenticated users to login.
- Returns `403 Forbidden` when the user lacks an allowed role.

Use this for authorization gates that are stricter than simple authentication.

Example:

```go
adminRoutes := middleware.Chain(
    mux,
    session.SM.SessionMiddleware,
    middleware.CSRF,
    middleware.Logger,
    middleware.Role("admin"),
)
```

### NotFound

- Lives in `internal/middleware/not_found.go`.
- Wraps the response writer to intercept `404 Not Found` responses.
- Suppresses the default Go 404 text.
- Renders the framework's templ-based not-found page through `pages.NotFoundPage()`.

### SecurityHeaders

- Lives in `internal/middleware/security.go`.
- Adds standard response security headers such as CSP, frame blocking, referrer policy, permissions policy, and HSTS in production.
- Adjusts CSP in development to allow the Vite dev server and websocket/HMR flow.

Use this for response hardening at the outer router layer or on specific route groups.

### NoCacheHeaders

- Lives in `internal/middleware/security.go`.
- Adds no-store and related cache-control headers.
- Intended for authenticated or sensitive pages.

Use it on dashboard, account, billing, admin, or other user-specific pages where browser caching is undesirable.

### MaxBodySize

- Lives in `internal/middleware/body_limit.go`.
- Wraps the request body with `http.MaxBytesReader`.
- Returns `413` behavior when the request exceeds the configured limit.

Use it where body size needs to be constrained explicitly, especially auth forms, uploads, and endpoints exposed to untrusted clients.

Examples:

- `64 * 1024` for small auth forms
- `1 << 20` for general form submissions
- larger explicit limits for upload endpoints that are intentionally allowed to receive files

### Template

- Lives in `internal/middleware/template.go`.
- Currently acts as a stub and is not carrying meaningful framework behavior.

Do not build new conventions around it until it has an actual framework responsibility.

### Framework Rate Limiter

- Lives in `internal/ratelimit/` rather than `internal/middleware/`.
- Exposes middleware-compatible handlers through `Limiter.Middleware`.
- Supports preset policies such as public browsing, API defaults, auth-sensitive endpoints, and exports.
- Can use memory or Redis storage depending on config.

The rate limiter is part of the framework middleware surface and should be applied where abuse protection matters.

## Middleware Rules

- Keep middleware reusable and framework-oriented when it lives in `internal/middleware/`.
- Apply middleware from route registration, not from handler methods.
- Be explicit about middleware order for every route group.
- If a middleware depends on session data, make sure the session middleware is deeper in the chain so the session exists before the dependent middleware executes.
- If a concern is app-specific rather than reusable, prefer a project-layer convention or a new app package instead of expanding framework middleware unnecessarily.

## Middleware Surface Summary

Useful framework middleware currently available across `internal/middleware/` and `internal/ratelimit/` includes:

- `Recover`
- `RecoverGoroutine`
- `SessionMiddleware`
- `CSRF`
- `Logger`
- `Auth`
- `Guest`
- `Role(...)`
- `NotFound()`
- `SecurityHeaders`
- `NoCacheHeaders`
- `MaxBodySize(...)`
- rate-limiter middleware via `Limiter.Middleware`

## Rate Limiting Guidance

Use rate limiting at the route-group level, for example on:

- login and registration flows
- password reset flows
- public browse endpoints if scraping protection is needed
- API endpoints
- heavy export or report generation endpoints

## Common Middleware Patterns

### Global outer-router concerns

These are commonly applied at the outer router layer:

- `Recover`
- `NotFound()`
- `SecurityHeaders`

### Authenticated route groups

These commonly appear on protected route groups:

- `SessionMiddleware`
- `CSRF`
- `Logger`
- `NoCacheHeaders`
- `Auth`
- `Role(...)` when authorization is required
- rate limiting when abuse or costly operations are involved

### Guest route groups

These commonly appear on login, register, and reset routes:

- `SessionMiddleware`
- `CSRF`
- `Logger`
- `Guest`
- auth-sensitive rate limiting
- `MaxBodySize(...)` when you want stricter body limits for small forms

## Example Route-Group Patterns

These examples are useful because they show how middleware groups are usually structured.

### Controller-owned `RegisterRoutes()` style

```go
func (c *AuthController) RegisterRoutes() http.Handler {
    mux := http.NewServeMux()

    mux.HandleFunc("GET /login", c.Login)
    mux.HandleFunc("POST /login", c.HandleLogin)
    mux.HandleFunc("GET /register", c.Register)
    mux.HandleFunc("POST /register", c.HandleRegister)

    return middleware.Chain(
        mux,
        session.SM.SessionMiddleware,
        middleware.CSRF,
        middleware.Logger,
        middleware.Guest,
    )
}
```

### Mixed guest and authenticated subgroups

This is useful when one controller owns both guest-only and authenticated routes:

```go
func (c *AuthController) RegisterRoutes() http.Handler {
    guestMux := http.NewServeMux()
    guestMux.HandleFunc("GET /login", c.Login)
    guestMux.HandleFunc("POST /login", c.HandleLogin)

    guestRoutes := middleware.Chain(
        guestMux,
        session.SM.SessionMiddleware,
        middleware.CSRF,
        middleware.Logger,
        middleware.Guest,
    )

    authMux := http.NewServeMux()
    authMux.HandleFunc("POST /logout", c.HandleLogout)

    authRoutes := middleware.Chain(
        authMux,
        session.SM.SessionMiddleware,
        middleware.CSRF,
        middleware.Logger,
        middleware.Auth,
    )

    parentMux := http.NewServeMux()
    parentMux.Handle("/", guestRoutes)
    parentMux.Handle("/logout", authRoutes)

    return parentMux
}
```

### Outer-router global wrapper style

Use the app router for outer concerns such as recovery, security headers, or not-found handling:

```go
func (r *AppRouter) SetupRoutes() http.Handler {
    mainMux := http.NewServeMux()

    auth := controllers.NewAuthController()
    pages := controllers.NewPagesController()

    mainMux.Handle("/auth/", http.StripPrefix("/auth", auth.RegisterRoutes()))
    mainMux.Handle("/", pages.RegisterRoutes())

    return middleware.Chain(
        mainMux,
        middleware.SecurityHeaders,
        middleware.NotFound(),
    )
}
```

These examples are framework patterns, not application business rules. They are appropriate to keep because they help agents compose middleware predictably.

## Short-Circuiting Rule

Some middleware may stop the request chain and write a response immediately.

Examples in the framework:

- `Auth` redirects unauthenticated users
- `Guest` redirects authenticated users
- `CSRF` rejects invalid or missing tokens
- `Role(...)` rejects forbidden users
- rate limiting returns `429 Too Many Requests`

Agents should treat these as control-flow boundaries when reasoning about request behavior.

## Practical Middleware Rules

Use these as enforcement guidance for agents:

- Keep middleware concerns explicit rather than implicit.
- Prefer one well-defined route-group chain over scattered per-handler wrapping.
- Reuse framework middleware when the behavior already exists instead of rewriting it in controllers.
- Add new framework middleware only when it is reusable across projects.
- If templ rendering changed but the security or request lifecycle concern still exists, keep the middleware guidance and only update the rendering-related parts.

## Middleware Checklist

Use this when defining or reviewing a middleware chain:

- Start with the route group’s purpose and threat profile.
- Ensure any middleware that needs session state runs with session middleware available beneath it in the chain.
- Apply CSRF to all state-changing browser flows.
- Apply `Auth` to protected routes and `Guest` to guest-only flows.
- Add rate limiting for auth-sensitive, public-browse, API, or heavy-operation endpoints as needed.
- Keep middleware order explicit and consistent.
- Keep framework-wide reusable middleware in `internal/middleware/` or `internal/ratelimit/`.
- Keep app-specific middleware out of the framework layer unless it is becoming reusable framework behavior.
