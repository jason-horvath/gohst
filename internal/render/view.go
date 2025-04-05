package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"gohst/internal/utils"
)

type ViewData struct {
    Props    interface{}
    Content  template.HTML
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

func NewView() *View {
	templateFuncs := TemplateFuncs()
	view := &View{
		Template: template.New("").Funcs(templateFuncs),
		Layout: "layout/default",
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

func (v *View) Render(r http.ResponseWriter, viewName string, data ...interface{}) error {
	var viewContent bytes.Buffer
	useData := utils.StructEmpty(data)
	log.Printf("useData type: %T", useData)
	useViewName := v.Dirs.Views + "/" + viewName
	err := v.Template.ExecuteTemplate(&viewContent, useViewName, useData)
	if err != nil {
		log.Println("Error executing template:", err)
	}

	td := ViewData{
        Props:   struct{}{}, // Possibly to load other view data
        Content: template.HTML(viewContent.String()),
    }

	log.Println("DEFINED TEMPLATES:", v.Template.DefinedTemplates())
	return v.Template.ExecuteTemplate(r, v.Layout, td)
}

