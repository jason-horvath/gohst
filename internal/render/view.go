package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"gohst/internal/auth"
	"gohst/internal/session"
	"gohst/internal/utils"
)

type TemplateData struct {
    CSRF    		*CSRF      			// CSRF token for form protection
    Auth	  		*auth.AuthData      // Pointer to the authenticated user (if any)
    Flash 	map[string]any 		// Slice for any flash messages (success/error)
	OldData    	map[string]any 			// Map for old input values (for form repopulation)
    Data         	any  				// Additional dynamic data specific to each page
}

type ViewData struct {
    Props    any
    Content  template.HTML
	CSRF	 *CSRF
}

type View struct {
	Template 	*template.Template
	Layout  	string
	Dirs 		ViewDirs
}

type ViewDirs struct {
	Layouts 	string
	Templates 	string
	Views 		string
	Partials 	string
	Components 	string
}


const defaultLayout string = "layout/default"

func NewView() *View {
	templateFuncs := TemplateFuncs()
	view := &View{
		Template: template.New("").Funcs(templateFuncs),
		Layout: defaultLayout,
		Dirs: ViewDirs{
			Layouts: "layouts",
			Templates: "templates",
			Views: "views",
			Partials: "partials",
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

// loadAll loads all templates from the specified directory and parses them into the template engine.
func (v *View) loadAll() {
	dirPath := v.Dirs.Templates
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

			// Parse into base template
			_, err = v.Template.New(name).Parse(string(content))
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error loading tempaltes: %v", err)
	}

	log.Println("Loaded Dir:", dirPath, v.Template.DefinedTemplates())
}

// Render renders a view with the given name and data. Render the content then the whole view with the layout.
func (v *View) Render(w http.ResponseWriter, r *http.Request, viewName string, data ...interface{}) error {
	var viewContent bytes.Buffer
	useData := utils.StructSafe(data)
    sess := session.FromContext(r.Context())
	csrf := GetCSRF(r)
	authData := auth.GetAuthData(sess)

	// Define template data for globlal use in templates
	templateData := TemplateData{
		CSRF: csrf,
		Auth: authData,
		Data: useData,
		Flash: sess.GetAllFlash(),
        OldData:   sess.GetAllOld(),
	}

	useViewName := v.Dirs.Views + "/" + viewName
	err := v.Template.ExecuteTemplate(&viewContent, useViewName, templateData)
	if err != nil {
		log.Println("Error executing template:", err)
	}

	td := ViewData{
        Props:   struct{}{}, // Possibly to load other view data
        Content: template.HTML(viewContent.String()),
    }

	log.Println("DEFINED TEMPLATES:", v.Template.DefinedTemplates())
	return v.Template.ExecuteTemplate(w, v.Layout, td)
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

