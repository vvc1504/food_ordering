package spec

// StatusCode represents the status of an operation.
type StatusCode int

// Status defines the interface for operation status reporting.
type Status interface {
	Code() StatusCode
	Message() string
}

// statusObj is the internal implementation of the Status interface.
type statusObj struct {
	code    StatusCode
	message string
}

func (s *statusObj) Code() StatusCode { return s.code }
func (s *statusObj) Message() string  { return s.message }

// NewStatus creates a new Status instance.
func NewStatus(code StatusCode, message string) Status {
	return &statusObj{
		code:    code,
		message: message,
	}
}
