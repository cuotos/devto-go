package devto

const (
	errEmptyCredentials        = "invalid credentials: key must not be empty"
	errArticleBodyAlreadyTaken = "body markdown has already been taken"
	errUnmarshalError          = "error unmarshalling the JSON response"
)

type Error struct {
	ErrorMsg string `json:"error"`
	Status   int    `json:"status"`
}

func (e Error) Error() string {
	return e.ErrorMsg
}
