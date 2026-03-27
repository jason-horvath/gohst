# Routes

Related: `overview.md`, `controllers.md`, `middleware.md`, `sessions.md`

## Purpose

Routing in this framework is controller-owned. Feature controllers define their own route trees through `RegisterRoutes()`, and the application router composes those controllers into the final mux.

## Architecture

There are two routing layers:

1. `app/routes/routes.go` builds the top-level application mux and mounts feature controllers.
2. `internal/routes/router.go` defines the framework router contract and wraps the application router with `middleware.Recover`.

The expected `AppRouter.SetupRoutes()` flow is:

1. Create a top-level `http.ServeMux`.
2. Mount static files under `/static/`.
3. Instantiate feature controllers.
4. Mount controller route trees at their prefixes.
5. Return the mux to `internal/routes.RegisterRouter`, which adds panic recovery.

Example:

```go
func (r *AppRouter) SetupRoutes() http.Handler {
	mainMux := http.NewServeMux()

	auth := controllers.NewAuthController()
	pages := controllers.NewPagesController()

	mainMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mainMux.Handle("/auth/", http.StripPrefix("/auth", auth.RegisterRoutes()))
	mainMux.Handle("/", pages.RegisterRoutes())

	return mainMux
}
```

## Route Group Pattern

Each controller-owned route group should:

- instantiate the controller it needs
- register Go 1.22 `http.ServeMux` patterns directly
- apply the exact middleware chain for that group
- return a self-contained `http.Handler`

In practice, those responsibilities live inside the controller’s `RegisterRoutes()` implementation, while `AppRouter` stays focused on mounting and top-level composition.

## Conventions

- Public pages live in the controller that owns them.
- Auth routes may be split into guest-only and auth-only subgroups inside the auth controller.
- Static files are served once from the top-level mux.
- Prefix stripping is only used when the mounted path differs from the inner mux paths.

## Path Pattern Rules

Use Go 1.22 method-aware patterns consistently:

- `GET /{$}` for the exact root of a group
- `GET /post/{id}` for path parameters
- `POST /login` or `POST /logout` for mutating actions

Inside handlers, read dynamic values with `r.PathValue("name")`.

## Router Responsibilities

The router layer is responsible for:

- path ownership
- mounting controller route trees at the correct prefixes
- mounting static files and prefixed route trees
- applying outer-router concerns that span the whole application

Controller route registration is responsible for deciding which middleware wraps which route group.

The router is not the place for business logic or templ composition.

## RegisterRoutes Convention

`RegisterRoutes()` is the framework convention.

Use it to keep:

- handler method registration close to the controller
- controller-scoped middleware explicit
- mount-point ownership in `app/routes/routes.go`
- shared outer-router concerns at the app router layer when needed

## Adding New Routes

When adding a new feature area:

1. Add `RegisterRoutes()` to the controller for the feature area.
2. Register explicit method-aware mux patterns inside that controller.
3. Apply controller-scoped middleware such as auth, guest, CSRF, and rate limiting where appropriate.
4. Mount the controller from `app/routes/routes.go`.
5. Keep top-level global composition in `SetupRoutes()`.

If a new route group introduces a reusable framework concern, document that concern in `.ai/framework/`. If it is only an app policy for the cloned project, document it in `.ai/project/`.

## Rate Limiting

Rate limiting is a framework capability under `internal/ratelimit/`.

What exists:

- reusable limiter middleware via `Limiter.Middleware`
- in-memory and Redis-backed stores
- preset policies such as public browse, auth-sensitive, API default, and exports
- key strategies based on IP, user, token, or composite identifiers
- standard rate-limit response headers and optional deny logging

Attach it where the route group is defined, usually inside `RegisterRoutes()`:

```go
func (c *AuthController) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /login", c.HandleLogin)

	store := ratelimit.NewStore()
	authLimiter := ratelimit.NewAuthSensitiveLimiter(store, "email")

	return middleware.Chain(
		mux,
		session.SM.SessionMiddleware,
		middleware.CSRF,
		middleware.Logger,
		authLimiter.Middleware,
		middleware.Guest,
	)
}
```

## Routing Checklist

Use this as an execution checklist when adding or reviewing routes:

- Keep mount-point composition in `app/routes/routes.go`.
- Put route registration in controller-owned `RegisterRoutes()` methods.
- Specify HTTP methods explicitly for every route pattern.
- Use path parameters for resource identifiers and query parameters for filtering, search, or pagination.
- Serve static files from the top-level router only once.
- Apply session, CSRF, auth, guest, logger, and rate-limit middleware intentionally for each route group.
- Attach rate limiting where abuse risk exists.
- Keep global middleware and cross-controller composition at the app router layer.
