package render

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Cache for loaded icons to improve performance
var iconCache = make(map[string]template.HTML)

// Base directory for icon files
var iconPath = "assets/icons"

// Icon renders an SVG icon by name with consistent wrapper
func Icon(name string, classes ...string) template.HTML {
    // Generate a cache key that includes the classes
    cacheKey := name
    if len(classes) > 0 {
        cacheKey = fmt.Sprintf("%s|%s", name, strings.Join(classes, ","))
    }

    // Try to get icon from cache
    if cachedSvg, ok := iconCache[cacheKey]; ok {
        return cachedSvg
    }

    // Find SVG file
    filePath := filepath.Join(iconPath, name+".svg")
    content, err := os.ReadFile(filePath)
    if err != nil {
        return template.HTML(fmt.Sprintf("<!-- Icon not found: %s -->", name))
    }

    // Parse SVG to extract and modify attributes
    svg := string(content)

    // Apply custom classes
    className := "w-6 h-6" // Default classes
    if len(classes) > 0 {
        className = strings.Join(classes, " ")
    }

    // Replace class attribute or add if missing
    classRegex := regexp.MustCompile(`class="[^"]*"`)
    if classRegex.MatchString(svg) {
        svg = classRegex.ReplaceAllString(svg, fmt.Sprintf(`class="%s"`, className))
    } else {
        svg = strings.Replace(svg, "<svg", fmt.Sprintf(`<svg class="%s"`, className), 1)
    }

    // Store in cache
    result := template.HTML(svg)
    iconCache[cacheKey] = result

    return result
}

// InitIcons initializes the icon system with a custom path
func InitIcons(path string) {
    if path != "" {
        iconPath = path
    }
    iconCache = make(map[string]template.HTML)
}
