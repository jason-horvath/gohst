package forms

// Element holds common HTML attributes for form elements
type Element struct {
		Attrs map[string]string // HTML attributes (e.g. Alpine or HTMX directives)
}

// ensureAttrs initializes the Attrs map if nil
func (e *Element) ensureAttrs() {
	if e.Attrs == nil {
		 e.Attrs = make(map[string]string)
	}
}

// SetAttr sets or updates a single HTML attribute
func (e *Element) SetAttr(key, value string) {
	e.ensureAttrs()
	e.Attrs[key] = value
}

// GetAttr retrieves the value for a given attribute key
// Returns the value and a bool indicating presence
func (e *Element) GetAttr(key string) (string, bool) {
	if e.Attrs == nil {
		 return "", false
	}
	val, ok := e.Attrs[key]
	return val, ok
}

// AddAttrs merges multiple attributes into the existing map
func (e *Element) AddAttrs(attrs map[string]string) {
	e.ensureAttrs()
	for k, v := range attrs {
		 e.Attrs[k] = v
	}
}

// ClearAttrs removes all attributes
func (e *Element) ClearAttrs() {
	e.Attrs = make(map[string]string)
}

// Button represents an HTML button element
type Button struct {
	Element               // embed common attributes
	ID      string        // optional element ID
	Type    string        // button type (submit, button, etc.)
	Text    string        // button label text
}

// Label represents an HTML label element
type Label struct {
	Element               // embed common attributes
	For    string         // 'for' attribute
	Text   string         // label text
}

// Option represents an option in a select element
type Option struct {
    Value string `json:"value"`
	Label string `json:"label"`
}
type CheckBoxOption struct {
	Option
	Name	 string // group name
	ID		 string // unique per option, e.g. groupName + "-" + value
	Checked  bool `json:"checked"`
}

type RadioOption struct {
    Option
	Name	 string // group name
	ID		 string // unique per option, e.g. groupName + "-" + value
	Selected bool `json:"selected"`
}

type SelectOption struct {
	Option
	Selected bool `json:"selected"`
}
