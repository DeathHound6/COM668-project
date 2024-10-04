package utility

type ErrorSchema struct {
	Message string `json:"message"`
}

type ErrorResponseSchema struct {
	Errors []ErrorSchema `json:"errors"`
}

type UserPostRequestBodySchema struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Teams    []string `json:"teams"`
}

type TeamPostRequestBodySchema struct {
	Name string `json:"name"`
}
