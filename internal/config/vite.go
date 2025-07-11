package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const DEFAULT_VITE_PORT = 5173

type ViteConfig struct {
	Port			int
	ManifestPath 	string
	Manifest		Manifest
}

var Vite *ViteConfig

type Manifest map[string]ManifestEntry

type ManifestEntry struct {
    File    string 	`json:"file"`
	IsEntry bool   	`json:"isEntry"`
	Name 	*string  `json:"name,omitempty"`
    Src     string 	`json:"src"`
}

func (me *ManifestEntry) GetName() string {
	if me.Name != nil {
		return *me.Name
	}
	return ""
}

func (me *ManifestEntry) GetType() string {
	ext := filepath.Ext(me.File)
    switch strings.ToLower(ext) {
    case ".css":
        return "css"
    case ".js":
        return "javascript"
    case ".jpg", ".jpeg", ".png", ".gif", ".svg":
        return "image"
    default:
        return "unknown"
    }
}

func LoadManifest(path string) (Manifest, error) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Println("Error getting current working directory in LoadManifest:", err)
	}

	manifestPath := cwd + "/" + path
	log.Println("Manifest Path:", manifestPath)
	jsonFile, err := os.Open(manifestPath)
	if err != nil {
        return Manifest{}, err
    }
	defer jsonFile.Close()

	byteValue, _ := os.ReadFile(manifestPath)

    var manifest Manifest
    json.Unmarshal(byteValue, &manifest)
	log.Printf("Manifest after change: %+v", manifest) // Add this line
    return manifest, nil
}

func initVite() {
	app := GetAppConfig()
	vitePort := GetEnv("VITE_PORT", DEFAULT_VITE_PORT).(int)

	manifestPath := GetEnv("VITE_MANIFEST_PATH").(string)
	var err error
	var manifest Manifest

    if app.IsProduction() {
        // Load manifest only if not in development mode
        manifest, err = LoadManifest(manifestPath)

        if err != nil {
            log.Printf("Error loading manifest: %v", err)
        }
    }

	Vite = &ViteConfig{
		Port: vitePort,
		ManifestPath: manifestPath,
		Manifest: manifest,
	}
}

func (vc *ViteConfig) AbsManifestPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Println("Error getting current working directory in ManifestAbsPath:", err)
	}
	return cwd + "/" + vc.ManifestPath
}


