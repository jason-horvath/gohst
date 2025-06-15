package types

type Button struct {
	ID	   string
	Type   string
	Text   string
}

type Label struct {
    For  string
    Text string
}

type Input interface{}

type Text struct {
	Input 		Input
    Name        string
	ID		  	string
    Value       string
    Type        string
    Errors      []string
    Label       Label
    Placeholder string
}

type CheckBox struct {
	Input 		Input
    Name    	string
	ID		  	string
    Value   	string
    Type    	string
    Errors  	[]string
    Label   	Label
    Checked 	bool
}

type Radio struct {
	Input 		Input
    Name    	string
	ID		  	string
    Value   	string
    Type    	string
    Errors  	[]string
    Label   	Label
    Checked 	bool
}

type Option struct {
    Value string
    Text  string
}

type Select struct {
	Input 		Input
    Name    	string
	ID		  	string
    Value   	string
    Type    	string
    Errors  	[]string
    Label   	Label
    Options 	[]Option
}

type Form struct {
    Method 	 string
    Action 	 string
    Fieldset Fieldset
	Buttons  map[string]Button
}

type Fieldset map[string]Field

type Field struct {
	Label Label
	Input Input
	Error  string
}
