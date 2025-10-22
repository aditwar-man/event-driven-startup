package domain

// Auth-specific domain logic can stay here if needed
// For now, we'll use the shared domain package

// Domain errors specific to auth service
var (
	ErrInvalidEmail = NewDomainError("invalid email address")
	ErrUserNotFound = NewDomainError("user not found")
)

type DomainError struct {
	message string
}

func NewDomainError(message string) *DomainError {
	return &DomainError{message: message}
}

func (d *DomainError) Error() string {
	return d.message
}
