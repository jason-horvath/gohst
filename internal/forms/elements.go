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
    Value string
    Text  string
}
