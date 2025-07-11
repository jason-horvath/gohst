package render

import (
	"fmt"
	"html/template"
	"log"

	"gohst/internal/config"
	"gohst/internal/utils"
)

// Asset represents a single asset file with its type, file path, and optional metadata.
type Asset struct{
	Type string
	File string
	Meta map[string]string
}

type Assets 	[]Asset

// ManifestAsset converts a ManifestEntry to an Asset.
func ManifestAsset(entry config.ManifestEntry) *Asset {
	if entry.IsEntry {
		return &Asset {
			Type: entry.GetType(),
			File: entry.File,
		}
	}

	return &Asset{}
}

// AssetsHead returns the appropriate HTML for including assets in the head section of a page.
func AssetsHead() template.HTML {
	app := config.GetAppConfig()
    if app.IsProduction() {
        return AssetsHeadProd()
    }
    return AssetsHeadDev()
}

// AssetsHeadDev returns the HTML for development mode assets, including Vite client and entry files.
func AssetsHeadDev() template.HTML {
	vitePort := config.Vite.Port
    html := fmt.Sprintf(`
        <script type="module" src="http://localhost:%d/@vite/client"></script>
        <script type="module" src="http://localhost:%d/assets/js/entry.ts"></script>
        <link rel="stylesheet" href="http://localhost:%d/assets/css/entry.css" />
  `, vitePort, vitePort, vitePort)

	return template.HTML(html)
}

// AssetsHeadProd returns the HTML for production mode assets, using the Vite manifest.
func AssetsHeadProd() template.HTML {
	var html string
	log.Println("MANIFEST:", config.Vite.Manifest)
	for _, entry := range config.Vite.Manifest {
		asset := ManifestAsset(entry)
		log.Println("Asset:", asset)
		if len(asset.File) > 0 {
			url := utils.BuildDistURL(asset.File)
            switch asset.Type {
            case "javascript":
                html += fmt.Sprintf(`<script type="module" src="%s"></script>`, url)
            case "css":
                html += fmt.Sprintf(`<link rel="stylesheet" href="%s" />`, url)
            }
		}
	}

	return template.HTML(html)
}

// StaticAssetURL returns the URL for a static asset file, using the configured distribution path.
func StaticAssetURL(file string) string {
	return utils.BuildDistURL(file)
}
