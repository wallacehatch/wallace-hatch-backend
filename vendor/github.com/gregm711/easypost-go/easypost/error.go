package easypost

/*
Official API documentation available at:
https://www.easypost.com/docs/api.html#errors
*/

//Error is an EasyPost object representing an error
type Error struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Errors  []FieldError `json:"errors"`
}

//FieldError is an EasyPost object that defines an error in the verification
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
