package messages

const (
	ErrorInvalidRequest = "invalid_request"
	ErrorServerError    = "server_error"
)

type ErrorMessage struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
}
