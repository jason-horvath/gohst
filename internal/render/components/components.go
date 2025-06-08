package components

import (
	"bytes"
	"html/template"
	"io"
)

// Shared writer reference from the parent package
var TemplateWriter io.Writer

// CaptureOutput executes a function and returns its output as a string
func CaptureOutput(fn func()) string {
    var buf bytes.Buffer

    oldWriter := TemplateWriter
    TemplateWriter = &buf
    fn()
    TemplateWriter = oldWriter

    return buf.String()
}

// Helper to merge function maps
func mergeFuncs(target, source template.FuncMap) {
    for name, fn := range source {
        target[name] = fn
    }
}
