package domain

// User service specific domain logic can stay here
// For quota management, we'll use the shared domain

// Domain errors specific to user service
var (
	ErrUserNotFound  = NewDomainError("user not found")
	ErrQuotaExceeded = NewDomainError("quota exceeded")
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
