package forms

// Button represents an HTML button element
type Button struct {
    ID   string
    Type string
    Text string
}

// Label represents an HTML label element
type Label struct {
    For  string
    Text string
}

// Option represents an option in a select element
type Option struct {
    Value string `json:"value"`
	Label string `json:"label"`
}
type CheckBoxOption struct {
	Option
	Checked bool `json:"checked"`
}

type RadioOption struct {
    Option
    Selected bool `json:"selected"`
}

type SelectOption struct {
	Option
	Selected bool `json:"selected"`
}
