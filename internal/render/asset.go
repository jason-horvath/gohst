package render

import (
	"fmt"
	"html/template"
	"log"

	"gohst/internal/config"
)

type Asset struct{
	Type string
	File string
	Meta map[string]string
}

type Assets 	[]Asset

func ManifestAsset(entry config.ManifestEntry) *Asset {
	if entry.IsEntry {
		return &Asset {
			Type: entry.GetType(),
			File: entry.File,
		}
	}

	return &Asset{}
}

func AssetsHead() template.HTML {
    if config.App.IsProduction() {
        return AssetsHeadProd()
    }
    return AssetsHeadDev()
}

func AssetsHeadDev() template.HTML {
	vitePort := config.Vite.Port
    html := fmt.Sprintf(`
        <script type="module" src="http://localhost:%d/@vite/client"></script>
        <script type="module" src="http://localhost:%d/assets/js/entry.ts"></script>
        <link rel="stylesheet" href="http://localhost:%d/assets/css/entry.css" />
  `, vitePort, vitePort, vitePort)

	return template.HTML(html)
}

func AssetsHeadProd() template.HTML {
	var html string
	log.Println("MANIFEST:", config.Vite.Manifest)
	for _, entry := range config.Vite.Manifest {
		asset := ManifestAsset(entry)
		log.Println("Asset:", asset)
		if len(asset.File) > 0 {
			url := config.App.DistURL(asset.File)
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

func StaticAssetURL(file string) string {
	return config.App.FullURL(file)
}
