package render

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var defaultTmplExt string = "*.html"

var templates *template.Template

func LoadTemplateDir(dir string, ext ...string) *template.Template{
    templateExt := "*.html"
    if len(ext) > 0 && ext[0] != "" {
        templateExt = ext[0]
    }

	pattern := filepath.Join(dir, templateExt)
	log.Println("Loading templates from directory", pattern)
    templates = template.Must(template.ParseGlob(pattern))

	var err error
    templates, err = template.ParseGlob(pattern)
    if err != nil {
        log.Fatalf("Error loading templates from directory %s: %v", dir, err)
    }

	return templates
}

func LoadAllTemplates(dir string) {
    var allFiles []string

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".html" {
			allFiles = append(allFiles, path)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking through templates directory: %v", err)
	}

	templates, err = template.ParseFiles(allFiles...)
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	log.Println("Templates loaded successfully:", templates.DefinedTemplates())
}

func Text(w http.ResponseWriter, text string) {
	w.Write([]byte(text))
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
    log.Println("Available templates:", templates.DefinedTemplates())
	tmpl := templates.Lookup(name)

    if tmpl == nil {
        http.Error(w, "The template does not exist.", http.StatusInternalServerError)
        return
    }

	// tmpl.Funcs(templateFuncs)
    err := tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
