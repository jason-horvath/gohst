# Controllers

Related: `overview.md`, `routes.md`, `middleware.md`, `sessions.md`, `rendering.md`

## Purpose

Controllers are request handlers for the application layer. They live in `app/controllers/` and coordinate HTTP input, session state, services, and templ rendering. The framework controller code in `internal/controllers/` provides the reusable rendering and response primitives that every project should build on.

This file should be treated as implementation guidance for building controllers consistently, not just as a description of what happens to exist today.

## Repository Boundary

- `internal/controllers/base_controller.go` is framework code.
- `app/controllers/app_controller.go` is the application base controller.
- Feature handlers such as `AuthController` and `PagesController` live in `app/controllers/`.
- Feature controllers should own their route registration through `RegisterRoutes()`.

## Required Hierarchy

Use this embedding chain consistently:

```go
internal/controllers.BaseController
        ↓
app/controllers.AppController
        ↓
app/controllers.FeatureController
```

### BaseController contract

`BaseController` owns request response helpers and the render view instance. The current framework API is:

- `Render(w, r, page render.Page)` for full-page templ responses
- `RenderPartial(w, r, component templ.Component)` for fragment responses
- `Redirect(w, r, url, statusCode)` for HTTP redirects
- `SetError(r, message)` for flash errors
- `SetTitle(title)` and `SetMeta(meta)` for per-request view metadata
- `JSON(w, status, data)` for JSON responses

Do not put app-specific behavior into `BaseController`.

### AppController contract

`AppController` embeds `BaseController` and is the place for shared application behavior. Every feature controller should be constructed through `NewAppController()` and can apply layout or app-specific defaults in its constructor.

Current example:

```go
type AuthController struct {
    *AppController
}

func NewAuthController() *AuthController {
    auth := &AuthController{
        AppController: NewAppController(),
    }
    auth.View.SetLayout("layouts/auth")
    return auth
}
```

## Route Registration Convention

Use controller-owned route registration as the standard pattern:

```go
func (c *AuthController) RegisterRoutes() http.Handler {
    mux := http.NewServeMux()

    mux.HandleFunc("GET /login", c.Login)
    mux.HandleFunc("POST /login", c.HandleLogin)
    mux.HandleFunc("GET /register", c.Register)
    mux.HandleFunc("POST /register", c.HandleRegister)

    store := ratelimit.NewStore()
    limiter := ratelimit.NewAuthSensitiveLimiter(store, "email")

    return middleware.Chain(
        mux,
        session.SM.SessionMiddleware,
        middleware.CSRF,
        middleware.Logger,
        limiter.Middleware,
        middleware.Guest,
    )
}
```

This pattern keeps route definitions close to the handlers and makes controller-level middleware explicit.

## Reference Controller Example

A full controller example is useful for agents, but only when it reflects the current templ render flow and the framework's target routing direction. Do not preserve older examples that still assume the pre-templ template system or string-based view rendering.

Use an example like this as the reference pattern:

```go
package controllers

import (
    "net/http"

    "gohst/views/pages"
)

type PagesController struct {
    *AppController
}

func NewPagesController() *PagesController {
    pagesController := &PagesController{
        AppController: NewAppController(),
    }
    return pagesController
}

func (c *PagesController) Index(w http.ResponseWriter, r *http.Request) {
    c.Render(w, r, pages.IndexPage())
}

func (c *PagesController) Post(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")

    response := struct {
        ID      string `json:"id"`
        Message string `json:"message"`
    }{
        ID:      id,
        Message: "This is post " + id,
    }

    c.JSON(w, http.StatusOK, response)
}

func (c *PagesController) NotFound(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNotFound)
    c.Render(w, r, pages.NotFoundPage())
}

func (c *PagesController) RegisterRoutes() http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /{$}", c.Index)
    mux.HandleFunc("GET /post/{id}", c.Post)
    mux.HandleFunc("GET /", c.NotFound)

    return mux
}
```

This example is intentionally aligned with the templ page functions, `Render` and `JSON` helpers, path-parameter handling, and controller-owned `RegisterRoutes()` architecture. More complex flows such as auth forms, session-backed validation, and rate limiting should follow the same controller structure but can build on additional patterns documented elsewhere in this file.

## Handler Pattern

Handlers are plain `net/http` methods registered from the router. Keep them thin and predictable:

1. Pull request-scoped dependencies from context.
2. Parse and validate request input.
3. Call services or models for business work.
4. Write flash or field error state when redirecting.
5. Render a templ page, redirect, or return JSON.

Current examples in this repository:

- `PagesController.Index` renders a templ page directly.
- `PagesController.Post` reads a path param with `r.PathValue("id")` and returns JSON.
- `AuthController.HandleLogin` and `HandleRegister` parse form input, persist old values in session, and redirect on failure.

## Controller Rules

- Keep business rules in `app/services/` when logic grows beyond request orchestration.
- Always access session state through `session.FromContext(r.Context())`.
- Use flash state for redirect-driven user feedback.
- Use field errors and old-input helpers for forms instead of ad hoc query parameters.
- Set layout changes in the controller constructor, not inline in every handler.
- Return JSON only for endpoints intended to be data responses. Standard page handlers should render templ pages.

If two controllers solve the same kind of problem, they should follow the same handler, session, rendering, and route-registration structure unless the framework docs are updated to permit a different pattern.

## Rendering Rules

- Full pages must call `Render` with a `render.Page` returned from a templ page function.
- Use `RenderPartial` only for partial HTML fragments such as HTMX responses.
- View components should rely on render context helpers for shared auth, CSRF, flash, and request metadata rather than passing those values down manually through every component.

## Form Error Pattern

The earlier `buildFormWithErrors` example was removed because it was tied to an `AccountController` and route pattern that do not exist in this repository. The underlying pattern is still valid and should stay documented.

When a form fails validation:

1. Persist submitted values with `sess.SetOld(...)` when you are redirecting.
2. Persist field-specific errors with `sess.SetFieldErrors(...)` or `sess.AddFieldError(...)`.
3. On the next render, rebuild the form using the submitted values and field errors rather than the original defaults.

That keeps the user’s input intact and lets templ components render field-level feedback consistently.

Minimal example:

```go
sess := session.FromContext(r.Context())

email := r.FormValue("email")
sess.SetOld("email", email)

if !validation.IsEmail(email) {
    sess.SetFieldErrors("email", []string{"Invalid email format"})
    c.Redirect(w, r, "/auth/login", http.StatusSeeOther)
    return
}
```

When rendering the form, pull the old value and field errors back into the form definition or templ component props.

## Controller Checklist

When creating or updating a controller in this framework:

- Keep it in `app/controllers/`.
- Embed `*AppController`.
- Construct it through `NewAppController()`.
- Set layout defaults in the constructor when needed.
- Add `RegisterRoutes() http.Handler`.
- Keep handlers focused on HTTP orchestration, not heavy business logic.
- Read session state from `session.FromContext(r.Context())`.
- Use flash, old-input, and field-error helpers for redirect-based forms.
- Return full pages through `Render`, fragments through `RenderPartial`, and data responses through `JSON`.
- If the controller needs abuse protection, attach the appropriate rate limiter where the route group is defined.

## Enforcement Note

Agents should default to this checklist when creating controllers. Deviations should be intentional and documented, not accidental or stylistic.

## What Not To Do

- Do not place application controllers in `internal/`.
- Do not bypass `AppController` and embed `BaseController` directly in feature controllers.
- Do not mix layout registration logic into controllers; layout functions are registered during web boot.
