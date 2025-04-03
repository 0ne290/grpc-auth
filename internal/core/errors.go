package core

type InvariantViolationError struct {
	Message string
}

func (e *InvariantViolationError) Error() string {
	return e.Message
}
