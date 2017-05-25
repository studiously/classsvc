package errors

import "errors"

var (
	ErrDatabaseUnresponsive = errors.New("could not connect to the database")
	// ErrUnauthenticated is returned when the authentication mechanism was not engaged; that is, the server didn't have a chance to check whether an action was permitted because no token was passed or the provided token was malformed.
	ErrUnauthenticated = errors.New("either no token was passed or the provided token was malformed")
	// ErrForbidden indicates that the client failed the authorization check.
	ErrForbidden = errors.New("client does not have correct scope for request")
	// ErrNotFound is returned when a resource was not found or, in sensitive areas, if the client is not authorized to perform an action.
	ErrNotFound     = errors.New("not found")
	ErrUserEnrolled = errors.New("user is already enrolled")
	ErrMustSetOwner = errors.New("new owner must be chosen before user can leave class")
)
