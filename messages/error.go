package messages

import "fmt"

const (
	ErrorInvalidRequest = "invalid_request"
	ErrorServerError    = "server_error"
	ErrorUnknownCommand = "unknown_command"
)

type ErrorMessage struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
}

func NewInvalidRequestMalformed() *ErrorMessage {
	return &ErrorMessage{
		Error:            ErrorInvalidRequest,
		ErrorDescription: "Malformed json message",
	}
}

func NewInvalidRequestMissingParameter(param string) *ErrorMessage {
	return &ErrorMessage{
		Error:            ErrorInvalidRequest,
		ErrorDescription: fmt.Sprintf("Missing required parameter: %s", param),
	}
}

func NewServerError() *ErrorMessage {
	return &ErrorMessage{
		Error:            ErrorServerError,
		ErrorDescription: "Internal server error",
	}
}

func NewUnknownCommandError(c string) *ErrorMessage {
	return &ErrorMessage{
		Error:            ErrorUnknownCommand,
		ErrorDescription: fmt.Sprintf("Unknown command: %s", c),
	}
}
