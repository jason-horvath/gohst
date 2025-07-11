package forms

// Form represents an HTML form with fields and buttons
type Form struct {
    Method   string
    Action   string
    Fieldset Fieldset
    Buttons  map[string]Button
}

// Fieldset is a collection of form fields indexed by field name
type Fieldset map[string]Field

// Field represents a form field with label, input, and error
type Field struct {
    Label Label
    Input Input
    Error string
}
