package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gohst/internal/auth"
	"gohst/internal/config"
	"gohst/internal/session"
	"gohst/internal/utils"
)

type TemplateData struct {
	CSRF        *CSRF               // CSRF token for form protection
	Auth        any                 // Pointer to the authenticated user (if any)
	Flash       map[string]any      // Slice for any flash messages (success/error)
	OldData     map[string]any      // Map for old input values (for form repopulation)
	FieldErrors map[string][]string // Per-field validation errors
	Data        any                 // Additional dynamic data specific to each page
	Request     *RequestProps       // Request metadata (path, method, etc.)
}

type RequestProps struct {
	Path   string // Request URI path
	Method string // HTTP method (GET, POST, etc.)
	URL    string // Full URL
}

type ViewData struct {
	CSRF        *CSRF
	Auth        any // Pointer to the authenticated user (if any)
	Data        any // Same data passed to the view, available in layout blocks
	Flash       map[string]any
	FieldErrors map[string][]string // Per-field validation errors
	Request     *RequestProps       // Request metadata available in layout blocks
	Props       any
	Content     template.HTML
	Title       string    // Page title set via SetTitle
	Meta        *PageMeta // Page metadata for SEO/social
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

type View struct {
	Template  *template.Template            // base template: layouts, partials, components
	viewFiles map[string]string             // view name → raw file content (parsed per-request via Clone)
	viewCache map[string]*template.Template // cached per-view template sets
	mu        sync.RWMutex                  // protects viewCache
	Layout    string
	Dirs      ViewDirs
	title     string    // holds per-request page title
	meta      *PageMeta // holds per-request page meta
}

type ViewDirs struct {
	Layouts    string
	Templates  string
	Views      string
	Partials   string
	Components string
}

const defaultLayout string = "layouts/default"

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
	templateFuncs := TemplateFuncs()
	view := &View{
		Template: template.New("").Funcs(templateFuncs),
		Layout:   defaultLayout,
		Dirs: ViewDirs{
			Layouts:    "layouts",
			Templates:  "templates",
			Views:      "views",
			Partials:   "partials",
			Components: "components",
		},
	}

	view.Init()

	return view
}

func (v *View) Init() {
	v.loadAll()
}

func (v *View) LoadTemplates() {
	// Go’s ParseGlob does not support brace expansion.
	// Load HTML and tmpl files separately if needed.
	_, err := v.Template.ParseGlob("templates/**/*.html")
	if err != nil {
		log.Fatalf("Error parsing HTML templates: %v", err)
	}
	_, err = v.Template.ParseGlob("templates/**/*.tmpl")
	if err != nil {
		log.Fatalf("Error parsing Tmpl templates: %v", err)
	}
}

// loadAll loads all templates from the specified directory.
// View files (under views/) are stored separately so they can be cloned per-request,
// isolating each view's {{ define }} blocks (e.g., "title", "admin-header") from other views.
// Shared templates (layouts, partials, components) are parsed into the base template.
func (v *View) loadAll() {
	v.viewFiles = make(map[string]string)
	v.viewCache = make(map[string]*template.Template)

	dirPath := v.Dirs.Templates
	viewsPrefix := v.Dirs.Views + "/"

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			name, err := filepath.Rel(dirPath, path)
			if err != nil {
				log.Fatal("No rel base path found.", err)
			}

			// Standardize path and remove extension
			name = filepath.ToSlash(name)
			name = name[:len(name)-len(filepath.Ext(name))]

			// View files are stored separately for per-view scoping
			if strings.HasPrefix(name, viewsPrefix) {
				v.viewFiles[name] = string(content)
			} else {
				// Shared templates (layouts, partials, components) go into the base set
				_, err = v.Template.New(name).Parse(string(content))
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error loading templates: %v", err)
	}
}

// getViewTemplate returns a template set with base templates + the specific view.
// Results are cached so each view is only cloned and parsed once.
func (v *View) getViewTemplate(viewName string) (*template.Template, error) {
	v.mu.RLock()
	tmpl, ok := v.viewCache[viewName]
	v.mu.RUnlock()
	if ok {
		return tmpl, nil
	}

	content, ok := v.viewFiles[viewName]
	if !ok {
		return nil, fmt.Errorf("view template %q not found", viewName)
	}

	tmpl, err := v.Template.Clone()
	if err != nil {
		return nil, err
	}
	if _, err = tmpl.New(viewName).Parse(content); err != nil {
		return nil, err
	}

	v.mu.Lock()
	v.viewCache[viewName] = tmpl
	v.mu.Unlock()

	return tmpl, nil
}

// Render renders a view with the given name and data. Render the content then the whole view with the layout.
func (v *View) Render(w http.ResponseWriter, r *http.Request, viewName string, data ...interface{}) error {
	var viewContent bytes.Buffer
	useData := utils.StructSafe(data)
	sess := session.FromContext(r.Context())
	csrf := GetCSRF(r)

	var authData any
	flash := make(map[string]any)
	oldData := make(map[string]any)
	fieldErrors := make(map[string][]string)

	if sess != nil {
		authData = auth.GetAuthData(sess)
		if f := sess.GetAllFlash(); f != nil {
			flash = f
		}
		if o := sess.GetAllOld(); o != nil {
			oldData = o
		}
		if fe := sess.GetAllFieldErrors(); fe != nil {
			fieldErrors = fe
		}
	}

	// Define template data for globlal use in templates
	baseURL := strings.TrimRight(config.GetEnv("APP_URL", "http://localhost:3030").(string), "/")
	fullURL := baseURL + r.URL.RequestURI()
	templateData := TemplateData{
		CSRF:        csrf,
		Auth:        authData,
		Data:        useData,
		Flash:       flash,
		OldData:     oldData,
		FieldErrors: fieldErrors,
		Request: &RequestProps{
			Path:   r.URL.Path,
			Method: r.Method,
			URL:    fullURL,
		},
	}

	useViewName := v.Dirs.Views + "/" + viewName

	// Get a template set scoped to this specific view
	tmpl, err := v.getViewTemplate(useViewName)
	if err != nil {
		return v.handleError(w, err)
	}

	err = tmpl.ExecuteTemplate(&viewContent, useViewName, templateData)
	if err != nil {
		return v.handleError(w, err)
	}

	// Capture and clear the title so it doesn't persist across requests
	title := v.title
	v.title = ""
	meta := v.meta
	v.meta = nil

	meta = applyMetaDefaults(meta, baseURL, fullURL)
	td := ViewData{
		CSRF:        csrf,
		Auth:        authData,
		Data:        useData,
		Flash:       flash,
		FieldErrors: fieldErrors,
		Request:     templateData.Request,
		Title:       title,
		Meta:        meta,
		Props:       struct{}{},
		Content:     template.HTML(viewContent.String()),
	}

	// Buffer the final output too, to catch layout errors
	var finalOutput bytes.Buffer
	err = tmpl.ExecuteTemplate(&finalOutput, v.Layout, td)
	if err != nil {
		return v.handleError(w, err)
	}

	_, err = finalOutput.WriteTo(w)
	return err
}

// RenderPartial renders a named template without wrapping it in a layout.
// Use this for HTMX partial responses or any request that needs a bare HTML fragment.
// The templateName should be the full template name (e.g. "partials/profile-cards").
func (v *View) RenderPartial(w http.ResponseWriter, r *http.Request, templateName string, data ...interface{}) error {
	useData := utils.StructSafe(data)
	sess := session.FromContext(r.Context())
	csrf := GetCSRF(r)

	var authData any
	flash := make(map[string]any)
	oldData := make(map[string]any)
	fieldErrors := make(map[string][]string)

	if sess != nil {
		authData = auth.GetAuthData(sess)
		if f := sess.GetAllFlash(); f != nil {
			flash = f
		}
		if o := sess.GetAllOld(); o != nil {
			oldData = o
		}
		if fe := sess.GetAllFieldErrors(); fe != nil {
			fieldErrors = fe
		}
	}

	templateData := TemplateData{
		CSRF:        csrf,
		Auth:        authData,
		Data:        useData,
		Flash:       flash,
		OldData:     oldData,
		FieldErrors: fieldErrors,
		Request: &RequestProps{
			Path:   r.URL.Path,
			Method: r.Method,
			URL:    r.URL.String(),
		},
	}

	var output bytes.Buffer
	err := v.Template.ExecuteTemplate(&output, templateName, templateData)
	if err != nil {
		return v.handleError(w, err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = output.WriteTo(w)
	return err
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

func applyMetaDefaults(meta *PageMeta, baseURL, fullURL string) *PageMeta {
	if meta == nil {
		meta = &PageMeta{}
	}

	if meta.Title == "" {
		meta.Title = "Humaculture – Cultivating Better Workplaces"
	}
	if meta.Description == "" {
		meta.Description = "A marketplace connecting organizations with expert HR practitioners for strategic consulting, assessments, and workforce development."
	}
	if meta.Canonical == "" {
		meta.Canonical = fullURL
	}
	if meta.OGImage == "" {
		meta.OGImage = baseURL + "/static/images/social/og-default.jpg"
	}
	if meta.OGType == "" {
		meta.OGType = "website"
	}
	if meta.TwitterCard == "" {
		meta.TwitterCard = "summary_large_image"
	}
	if meta.Schema == nil {
		meta.Schema = map[string]any{
			"@context":    "https://schema.org",
			"@type":       "Organization",
			"name":        "Humaculture",
			"url":         baseURL,
			"logo":        baseURL + "/static/images/logo.png",
			"description": "A marketplace connecting organizations with expert HR practitioners.",
			"sameAs":      []string{},
		}
	}

	return meta
}

// Set the layout to be used for rendering
func (v *View) SetLayout(layout string) {
	if layout != "" {
		v.Layout = layout
	} else {
		v.Layout = defaultLayout
	}
}

// Get the lauout to be used for rendering
func (v *View) GetLayout() string {
	if v.Layout == "" {
		return defaultLayout
	}
	return v.Layout
}
