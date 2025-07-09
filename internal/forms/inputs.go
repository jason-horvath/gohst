package forms

// Input is a generic interface for all form input types
type Input interface{}

// Text represents an HTML text input element
type Text struct {
    Input       Input
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
    Input   Input
    Name    string
    ID      string
    Value   string
    Type    string
    Errors  []string
    Label   Label
    Checked bool
}

// Radio represents an HTML radio input element
type Radio struct {
    Input   Input
    Name    string
    ID      string
    Value   string
    Type    string
    Errors  []string
    Label   Label
    Checked bool
}

// Select represents an HTML select element
type Select struct {
    Input   Input
    Name    string
    ID      string
    Value   string
    Type    string
    Errors  []string
    Label   Label
    Options []Option
}
