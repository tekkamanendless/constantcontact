package constantcontact

// Error is an error returned by the Constant Contact API.
type Error struct {
	ErrorKey     string `json:"error_key"`
	ErrorMessage string `json:"error_message"`
}

// ErrorOutput is a list of errors returned by the Constant Contact API.
type ErrorOutput []Error
