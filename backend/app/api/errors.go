package api

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

type AuthenticationError struct {
	Message string
}

func (e AuthenticationError) Error() string {
	if e.Message == "" {
		return "Authentication error"
	}
	return e.Message
}

type AuthorizationError struct {
	Message string
}

func (e AuthorizationError) Error() string {
	if e.Message == "" {
		return "Authorization error"
	}
	return e.Message
}

// const (
// 	HTTP_ERROR_CODE_RENEW_TOKEN = 701
// 	HTTP_ERROR_CODE_VALID       = 702
// 	HTTP_ERROR_CODE_MONGO       = 703
// )
