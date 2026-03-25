package forms

// Text represents an HTML text input element
type Text struct {
    Element                  // HTML attributes (Alpine, HTMX, aria-*, data-*)
    Name        string
    ID          string
    Value       string
    Type        string
    Errors      []string
    Label       Label
    Placeholder string
}

// CheckBox represents an HTML checkbox input element
type CheckBox struct {
    Element                  // HTML attributes (Alpine, HTMX, aria-*, data-*)
    Name     string
    ID       string
    Type     string
    Errors   []string
    Label    Label
    Options  []CheckBoxOption
}

// Radio represents an HTML radio input element
type Radio struct {
    Element                  // HTML attributes (Alpine, HTMX, aria-*, data-*)
    Name     string
    ID       string
    Type     string
    Errors   []string
    Label    Label
    Options  []RadioOption
}

// Select represents an HTML select element
type Select struct {
    Element                  // HTML attributes (Alpine, HTMX, aria-*, data-*)
    Name    string
    ID      string
    Value   string
    Type    string
    Errors  []string
    Label   Label
    Options []SelectOption
}

// File represents an HTML file input element
type File struct {
    Element                  // HTML attributes (Alpine, HTMX, aria-*, data-*)
    Name        string
    ID          string
    Type        string
    Errors      []string
    Label       Label
    Accept      string   // MIME types or file extensions (e.g., "image/*", ".jpg,.png")
    Multiple    bool     // Allow multiple file selection
    MaxSize     int64    // Maximum file size in bytes
    Placeholder string   // Helper text for file selection
}
