package render

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

// IconData represents the structure of an icon JSON file
type IconData struct {
    Path    string `json:"path"`
    ViewBox string `json:"viewBox,omitempty"`
}

// Cache for loaded icons to improve performance
var iconCache = make(map[string]IconData)

// Base directory for icon files
var iconPath = "assets/icons"

// Icon renders an SVG icon by name with consistent wrapper
func Icon(name string, classes ...string) template.HTML {
    // Try to get icon data from cache
    iconData, ok := iconCache[name]
    if !ok {
        // Load icon data from file
        filePath := filepath.Join(iconPath, name+".json")
        content, err := os.ReadFile(filePath)
        if err != nil {
            // Try with .svg extension if .json not found
            filePath = filepath.Join(iconPath, name+".svg")
            content, err = os.ReadFile(filePath)
            if err != nil {
                return template.HTML(fmt.Sprintf("<!-- Icon not found: %s -->", name))
            }

            // For .svg files, assume they only contain the path data
            iconData = IconData{
                Path: strings.TrimSpace(string(content)),
            }
        } else {
            // Parse JSON data
            if err := json.Unmarshal(content, &iconData); err != nil {
                return template.HTML(fmt.Sprintf("<!-- Error parsing icon: %s -->", name))
            }
        }

        // Cache for future use
        iconCache[name] = iconData
    }

    // Default classes if none provided
    className := "w-6 h-6"
    if len(classes) > 0 {
        className = strings.Join(classes, " ")
    }

    // Default viewBox if not specified
    viewBox := "0 0 24 24"
    if iconData.ViewBox != "" {
        viewBox = iconData.ViewBox
    }

    // Create SVG with consistent attributes and the path data
    svg := fmt.Sprintf(`<svg class="%s" xmlns="http://www.w3.org/2000/svg" viewBox="%s" fill="currentColor">
        <path d="%s" />
    </svg>`, className, viewBox, iconData.Path)

    return template.HTML(svg)
}

// InitIcons initializes the icon system with a custom path
func InitIcons(path string) {
    if path != "" {
        iconPath = path
    }

    // Clear cache when initializing with a new path
    iconCache = make(map[string]IconData)
}
