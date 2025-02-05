package utility

import "time"

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
	Key      string `json:"key"`
	Value    string `json:"value"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

type ProviderGetResponseSchema struct {
	UUID   string           `json:"uuid"`
	Name   string           `json:"name"`
	Fields []KeyValueSchema `json:"fields"`
	Type   string           `json:"type"`
}

type MetaSchema struct {
	TotalItems int64 `json:"total"`
	Pages      int   `json:"pages"`
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
}

type GetManyResponseSchema struct {
	Data []any      `json:"data"`
	Meta MetaSchema `json:"meta"`
}

type IncidentPostRequestBodySchema struct {
	Summary       string   `json:"summary"`
	HostsAffected []string `json:"hostsAffected"`
}

type UserGetResponseBodySchema struct {
	UUID    string                      `json:"uuid"`
	Name    string                      `json:"name"`
	Email   string                      `json:"email"`
	Teams   []TeamGetResponseBodySchema `json:"teams"`
	SlackID string                      `json:"slackID"`
	Admin   bool                        `json:"admin"`
}

type TeamGetResponseBodySchema struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type IncidentGetResponseBodySchema struct {
	UUID          string                             `json:"uuid"`
	HostsAffected []HostMachineGetResponseBodySchema `json:"hostsAffected"`
	Summary       string                             `json:"summary"`
	CreatedAt     time.Time                          `json:"createdAt"`
	ResolvedAt    *time.Time                         `json:"resolvedAt"`
	ResolvedBy    *UserGetResponseBodySchema         `json:"resolvedBy"`
}

type HostMachineGetResponseBodySchema struct {
	UUID     string                    `json:"uuid"`
	OS       string                    `json:"os"`
	Hostname string                    `json:"hostname"`
	IP4      *string                   `json:"ip4"`
	IP6      *string                   `json:"ip6"`
	Team     TeamGetResponseBodySchema `json:"team"`
}

type ProviderPostRequestBodySchema struct {
	Name string `json:"name"`
}

type ProviderPutRequestBodySchema struct {
	Fields []KeyValueSchema `json:"fields"`
}

type HostMachinePostPutRequestBodySchema struct {
	OS       string `json:"os"`
	Hostname string `json:"hostname"`
	IP4      string `json:"ip4"`
	IP6      string `json:"ip6"`
	TeamID   string `json:"teamID"`
}
