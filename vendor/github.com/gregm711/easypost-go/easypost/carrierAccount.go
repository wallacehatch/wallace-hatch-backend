package easypost

/*
Official API documentation available at:
https://www.easypost.com/docs/api.html#carrier-accounts
*/

//CarrierAccount is an easypost carrier account
type CarrierAccount struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	Fields          Fields `json:"fields"`
	Clone           bool   `json:"clone"`
	Description     string `json:"description"`
	Reference       string `json:"reference"`
	Readable        string `json:"readable"`
	Credentials     Field  `json:"credentials"`
	TestCredentials Field  `json:"test_credentials"`
}

type Fields struct {
	Credentials     Field `json:"credentials"`
	TestCredentials Field `json:"test_credentials"`
	AutoLink        bool  `json:"auto_link"`
	CustomWorkflow  bool  `json:"custom_workflow"`
}

type Field struct {
	Key        string `json:"key"`
	Visibility string `json:"visibility"`
	Label      string `json:"label"`
	Value      string `json:"value"`
}
