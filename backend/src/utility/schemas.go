package utility

type ErrorResponseSchema struct {
	Error string `json:"error"`
}

type UserPostRequestBodySchema struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Teams    []string `json:"teams"`
}

type UserLoginRequestBodySchema struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TeamPostRequestBodySchema struct {
	Name string `json:"name"`
}

type KeyValueSchema struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type ProviderGetResponseSchema struct {
	ID     string           `json:"id"`
	Name   string           `json:"name"`
	Fields []KeyValueSchema `json:"fields"`
	Type   string           `json:"type"`
}

type ProvidersGetResponseSchema struct {
	Providers []ProviderGetResponseSchema `json:"providers"`
}

type IncidentPostRequestBodySchema struct {
	Team    uint   `json:"team"`
	Summary string `json:"summary"`
}
