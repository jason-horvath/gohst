package render

import "context"

type ctxKey string

const (
	csrfCtxKey        ctxKey = "gohst_csrf"
	authCtxKey        ctxKey = "gohst_auth"
	flashCtxKey       ctxKey = "gohst_flash"
	fieldErrorsCtxKey ctxKey = "gohst_field_errors"
	requestCtxKey     ctxKey = "gohst_request"
)

// AuthUser is the minimal interface views need to display auth data.
// The application's auth/session data type must implement this.
type AuthUser interface {
	GetEmail() string
	GetName() string
}

// SetPageContext stores all shared page data in the request context so any
// templ component in the tree can access it without being passed explicitly.
func SetPageContext(
	ctx context.Context,
	csrf *CSRF,
	auth any,
	flash map[string]any,
	fieldErrors map[string][]string,
	req *RequestProps,
) context.Context {
	ctx = context.WithValue(ctx, csrfCtxKey, csrf)
	ctx = context.WithValue(ctx, authCtxKey, auth)
	ctx = context.WithValue(ctx, flashCtxKey, flash)
	ctx = context.WithValue(ctx, fieldErrorsCtxKey, fieldErrors)
	ctx = context.WithValue(ctx, requestCtxKey, req)
	return ctx
}

func GetCSRFFromCtx(ctx context.Context) *CSRF {
	if v, ok := ctx.Value(csrfCtxKey).(*CSRF); ok {
		return v
	}
	return &CSRF{}
}

func GetAuthFromCtx(ctx context.Context) any {
	return ctx.Value(authCtxKey)
}

func GetFlashFromCtx(ctx context.Context) map[string]any {
	if v, ok := ctx.Value(flashCtxKey).(map[string]any); ok {
		return v
	}
	return map[string]any{}
}

func GetFieldErrorsFromCtx(ctx context.Context) map[string][]string {
	if v, ok := ctx.Value(fieldErrorsCtxKey).(map[string][]string); ok {
		return v
	}
	return map[string][]string{}
}

func GetRequestFromCtx(ctx context.Context) *RequestProps {
	if v, ok := ctx.Value(requestCtxKey).(*RequestProps); ok {
		return v
	}
	return &RequestProps{}
}
