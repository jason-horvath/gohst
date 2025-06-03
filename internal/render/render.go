package render

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var defaultTmplExt string = "*.html"

var templates *template.Template

func LoadTemplateDir(dir string, ext ...string) *template.Template{
    templateExt := defaultTmplExt
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

func JSON(w http.ResponseWriter, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    err := json.NewEncoder(w).Encode(data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}
