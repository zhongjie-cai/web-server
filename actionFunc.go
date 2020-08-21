package webserver

// ActionFunc defines the action function to be called for route processing logic
type ActionFunc func(
	session Session,
) (
	responseObject interface{},
	responseError error,
)
