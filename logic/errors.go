package logic

// ValidationError representa un error de validación de reglas de negocio.
type ValidationError struct {
	Msg string
}

func (e *ValidationError) Error() string {
	return e.Msg
}

func NewValidationError(msg string) error {
	return &ValidationError{Msg: msg}
}
