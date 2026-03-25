package render

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/a-h/templ"

	"gohst/internal/auth"
	"gohst/internal/config"
	"gohst/internal/session"
)

// LayoutFunc creates a layout component; content is injected via templ.WithChildren.
type LayoutFunc func(title string, meta *PageMeta) templ.Component

// layoutRegistry is the package-level registry of named layout functions.
var layoutRegistry = map[string]LayoutFunc{}

// RegisterLayout registers a named layout function for use by View.Render.
// Call once per layout during application startup (e.g. in main.go).
func RegisterLayout(name string, fn LayoutFunc) {
	layoutRegistry[name] = fn
}

// Page bundles a view's title and rendered content component.
type Page struct {
	Title   string
	Content templ.Component
	Meta    *PageMeta
}

// RequestProps holds request metadata available to any templ component via context.
type RequestProps struct {
	Path   string
	Method string
	URL    string
}

type PageMeta struct {
	Title       string
	Description string
	Canonical   string
	OGImage     string
	OGType      string
	TwitterCard string
	NoIndex     bool
	Schema      any
}

// View manages per-controller layout and title state.
type View struct {
	Layout string
	title  string
	meta   *PageMeta
}

const defaultLayout = "layouts/default"

const errorTemplate string = `<!DOCTYPE html>
<html>
<head>
    <title>%s</title>
    <style>
        body { background-color: #1a1a1a; color: #e0e0e0; font-family: system-ui, -apple-system, sans-serif; padding: 2rem; }
        .container { max-width: 800px; margin: 0 auto; }
        .error-card { background-color: #2d2d2d; border: 1px solid #404040; border-radius: 8px; padding: 1.5rem; box-shadow: 0 4px 6px rgba(0,0,0,0.3); }
        h1 { color: #ff6b6b; margin-top: 0; font-size: 1.5rem; border-bottom: 1px solid #404040; padding-bottom: 1rem; }
        pre { background-color: #1a1a1a; padding: 1rem; border-radius: 4px; overflow-x: auto; color: #ffb86c; font-family: monospace; font-size: 0.9rem; border: 1px solid #404040; }
        .footer { margin-top: 1rem; color: #888; font-size: 0.8rem; }
    </style>
</head>
<body>
    <div class="container">
        <div class="error-card">
            <h1>👻 %s</h1>
            <pre>%s</pre>
            <div class="footer">
                Gohst Framework • Development Mode
            </div>
        </div>
    </div>
</body>
</html>`

func NewView() *View {
	return &View{
		Layout: defaultLayout,
	}
}

func (v *View) SetLayout(layout string) {
	v.Layout = layout
}

// Render assembles request-scoped page context and renders a Page through its layout.
func (v *View) Render(w http.ResponseWriter, r *http.Request, page Page) error {
	sess := session.FromContext(r.Context())
	csrf := GetCSRF(r)

	var authData any
	flash := make(map[string]any)
	fieldErrors := make(map[string][]string)

	if sess != nil {
		authData = auth.GetAuthData(sess)
		if f := sess.GetAllFlash(); f != nil {
			flash = f
		}
		if fe := sess.GetAllFieldErrors(); fe != nil {
			fieldErrors = fe
		}
	}

	baseURL := strings.TrimRight(config.GetEnv("APP_URL", "http://localhost:3030").(string), "/")
	req := &RequestProps{
		Path:   r.URL.Path,
		Method: r.Method,
		URL:    baseURL + r.URL.RequestURI(),
	}

	ctx := SetPageContext(r.Context(), csrf, authData, flash, fieldErrors, req)

	title := page.Title
	if title == "" {
		title = v.title
	}
	v.title = ""
	v.meta = nil

	layoutFn, ok := layoutRegistry[v.Layout]
	if !ok {
		layoutFn, ok = layoutRegistry[defaultLayout]
		if !ok {
			return v.handleError(w, fmt.Errorf("no layout registered (wanted %q)", v.Layout))
		}
	}

	component := templ.ComponentFunc(func(ctx context.Context, wr io.Writer) error {
		return layoutFn(title, page.Meta).Render(templ.WithChildren(ctx, page.Content), wr)
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(ctx, w); err != nil {
		return v.handleError(w, err)
	}
	return nil
}

// RenderPartial renders a templ component directly without a layout wrapper.
// Use for HTMX partial responses or bare HTML fragments.
func (v *View) RenderPartial(w http.ResponseWriter, r *http.Request, component templ.Component) error {
	sess := session.FromContext(r.Context())
	csrf := GetCSRF(r)

	var authData any
	flash := make(map[string]any)
	fieldErrors := make(map[string][]string)

	if sess != nil {
		authData = auth.GetAuthData(sess)
		if f := sess.GetAllFlash(); f != nil {
			flash = f
		}
		if fe := sess.GetAllFieldErrors(); fe != nil {
			fieldErrors = fe
		}
	}

	req := &RequestProps{
		Path:   r.URL.Path,
		Method: r.Method,
		URL:    r.URL.String(),
	}

	ctx := SetPageContext(r.Context(), csrf, authData, flash, fieldErrors, req)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return component.Render(ctx, w)
}


// RenderError renders the error template directly to the response writer.
// This is useful for middleware or other parts of the framework that need to show an error page.
func RenderError(w http.ResponseWriter, title string, err interface{}) {
	if config.GetAppConfig().IsDevelopment() {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, errorTemplate, title, title, err)
	} else {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (v *View) handleError(w http.ResponseWriter, err error) error {
	log.Println("Error executing template:", err)
	RenderError(w, "Template Rendering Error", err)
	return nil // Error handled
}

// SetTitle sets the page title for the next render.
func (v *View) SetTitle(title string) {
	v.title = title
}

// GetTitle gets the page title for the current request.
func (v *View) GetTitle() string {
	return v.title
}

// SetMeta sets the page meta for the next render.
func (v *View) SetMeta(meta *PageMeta) {
	v.meta = meta
}

// GetMeta gets the page meta for the current request.
func (v *View) GetMeta() *PageMeta {
	return v.meta
}
