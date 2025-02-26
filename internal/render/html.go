package render

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Html struct {
	Template 	*template.Template
	devMode 	bool
	dirs 		Dirs
	layout 		string
}

type Dirs struct {
	layouts		string
	templates 	string
	views 		string
	partials 	string
}

func NewHtml() *Html {
	templateFuncs := TemplateFuncs()
	html := &Html{
		Template: template.New("").Funcs(templateFuncs),
		devMode: true,
		dirs: Dirs{
			layouts: "layouts",
			templates: "templates",
			views: "views",
			partials: "partials",
		},
		layout: "default",
	}

	html.Init()

	return html
}

func (h *Html) Init() {
	h.LoadPartials()
}

func (h *Html) LoadPartials() {
	partialsPath := h.baseRelPath(h.dirs.partials)
	err := filepath.Walk(partialsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			name, err := filepath.Rel(partialsPath, path)
			if err != nil {
				log.Fatal("No rel base path found.", err)
			}

			// Standardize path and remove extension
			name = filepath.ToSlash(name)
			name = name[:len(name)-len(filepath.Ext(name))]

			// Parse into base template
			_, err = h.Template.New(name).Parse(string(content))
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error loading partials: %v", err)
	}

	log.Println("Loaded partials:", h.Template.DefinedTemplates())
}

func (h *Html) RenderTemplate(w http.ResponseWriter, name string, data ...interface{}) {
	useData := h.useData(data)
	log.Println("REL:", name)
	log.Println("REL:", h.baseRelPath(name))

	tmplPath, err := filepath.Abs(h.baseRelPath(name))
	if err != nil {
		log.Println("Error loading Abs path:", err)
	}

	t, err := h.Template.Clone()
	if err != nil {
		log.Println("Error cloning template:", err)
	}
	layoutPath := h.getLayoutPath()
	tmpl, err := t.ParseFiles(layoutPath, tmplPath)

	if err != nil {
		log.Println("Error loading template:", err)
	}

	_ = tmpl.ExecuteTemplate(w, h.TmplLayoutName(), useData)
    // if err != nil {
    //     http.Error(w, err.Error(), http.StatusInternalServerError)
    // }
}

func (h *Html) RenderView(w http.ResponseWriter, name string, data ...interface{}) {
	view := h.dirs.views + "/" + name

	h.RenderTemplate(w, view, data)
}

func (h *Html) RenderPartial(w http.ResponseWriter, name string, data ...interface{}) {
	view := h.dirs.partials + "/" + name

	h.RenderTemplate(w, view, data)
}

func (h *Html) useData(data ...interface{}) interface{} {
	var useData interface{}
    if len(data) > 0 {
        useData = data[0]
		if _, ok := useData.(interface{}); !ok {
            log.Println("Warning: data[0] is not of type interface{}")
            useData = struct{}{} // Provide a default empty struct if data is not of the expected type
        }
    } else {
        useData = struct{}{} // Provide a default empty struct if data is nil
    }

	return useData
}

func (h *Html) SetLayout(name string) *Html {
	h.layout = name;

	return h
}

func (h *Html) GetLayoutName() string {
	if h.layout == "" {
		return "default"
	}

	return h.layout
}

func (h *Html) TmplLayoutName() string {
	return "layout/" + h.GetLayoutName()
}

func (h *Html) getLayoutPath() string {
	name := h.GetLayoutName()
	rel := h.dirs.layouts + "/" + name + ".html"

	return h.absPath(rel)
}

func (h *Html) baseRelPath(rel string) string {
	return h.dirs.templates + "/" + rel
}

func (h *Html) absPath(name string) string {
	absPath, err := filepath.Abs(h.baseRelPath(name))

	if err != nil {
		log.Fatal("Error getting absPath", err)
	}

	return absPath
}
