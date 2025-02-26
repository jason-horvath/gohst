package render

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type ViewProps struct {
    ViewsDir      	string
    ComponentsDir   string
	LayoutsDir      string
	ContentDir		string
	TemplatesDir	string
}

type View struct {
	componentsDir 	string
	layoutsDir		string
	Template 		*template.Template
	templatesDir 	string
	viewsDir		string
	contentDir		string
}

func NewView(props ViewProps) *View {

	if props.ViewsDir == "" {
		props.ViewsDir = "views"
	}

	if props.ComponentsDir == "" {
		props.ComponentsDir = "components"
	}

	if props.LayoutsDir == "" {
		props.LayoutsDir = "layouts"
	}

	if props.TemplatesDir == "" {
		props.TemplatesDir = "templates"
	}

	return &View{
		componentsDir: 	props.ComponentsDir,
		layoutsDir: 	props.LayoutsDir,
		Template:		template.New(""),
		templatesDir: 	props.TemplatesDir,
		viewsDir: 		props.ViewsDir,
		contentDir: 	props.ContentDir,
	}
}

func (v *View) Init() {
	v.LoadAllTemplates()
	v.LoadLayouts()
	v.LoadContentTemplates()
	log.Println("VIEW TEMPLATES:" , v.Template.DefinedTemplates())
}

func (v *View) LoadComponentsDir() {
	v.LoadTemplateDir(v.GetComponentsPath())
}

func (v *View) LoadLayoutsDir() {
	v.LoadTemplateDir(v.GetLayoutsPath())
}

func (v *View) LoadContentTemplates() {
	v.LoadTemplateDir(v.GetContentPath())
}

func (v * View) LoadAllTemplates() {
    var allFiles []string

	err := filepath.WalkDir(v.GetTemplatesPath(), func(path string, d os.DirEntry, err error) error {
		log.Println("Walking through", path)
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".html" {
			allFiles = append(allFiles, path)
		}
		return nil
	})
	log.Println("All files:", allFiles)
	if err != nil {
		log.Fatalf("Error walking through templates directory: %v", err)
	}

	templates, err = v.Template.ParseFiles(allFiles...)
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}
}

func (v * View) LoadTemplateDir(dir string, ext ...string) *template.Template{
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

func (v *View) LoadLayouts() {
	layoutsPath := v.GetLayoutsPath()
    err := filepath.Walk(layoutsPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            content, err := os.ReadFile(path)

            if err != nil {
                return err
            }
            name := filepath.ToSlash(path)
			log.Println("Loading layout:", name)
			fullName := filepath.Join(v.layoutsDir, info.Name())
            _, err = v.Template.New(fullName).Parse(string(content))
            if err != nil {
                return err
            }
        }
        return nil
    })

	if err != nil {
        log.Fatalf("Error loading layouts from directory %s: %v", layoutsPath, err)
    }
}

func (v *View) RenderTemplate(w http.ResponseWriter, name string, data ...interface{}) {
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
	log.Println("Available templates:", v.Template.DefinedTemplates())
	tmpl := v.Template.Lookup(name)
    if tmpl == nil {
        http.Error(w, "The template does not exist.", http.StatusInternalServerError)
        return
    }
    err := tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func (v *View) GetComponentsPath() string {
	relComponents := filepath.Join(v.templatesDir, v.componentsDir)

	return v.GetAbsPath(relComponents)
}

func (v *View) GetLayoutsPath() string {
	relLayouts := filepath.Join(v.templatesDir, v.layoutsDir)

	return v.GetAbsPath(relLayouts)
}

func (v *View) GetTemplatesPath() string {
	return v.GetAbsPath(v.templatesDir)
}

func (v *View) GetViewsPath() string {
	relViews := filepath.Join(v.templatesDir, v.viewsDir)

	return v.GetAbsPath(relViews)
}

func (v *View) GetContentPath() string {
	relContent := filepath.Join(v.GetViewsPath(), v.contentDir)

	return v.GetAbsPath(relContent)
}

func (v *View) GetAbsPath(relPath string) string {
	absPath, err := filepath.Abs(relPath)

	if err != nil {
		log.Fatal(err)
	}

	return absPath
}
func (v *View) GetViewsDir() string {
	return v.viewsDir
}
