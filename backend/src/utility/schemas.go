package utility

import (
	"time"
)

type ErrorResponseSchema struct {
	Error string `json:"error" binding:"required"`
}

type UserPostRequestBodySchema struct {
	Name     string   `json:"name" binding:"required"`
	Email    string   `json:"email" binding:"required"`
	Password string   `json:"password" binding:"required"`
	Teams    []string `json:"teams" binding:"required"`
}

type UserLoginRequestBodySchema struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TeamPostRequestBodySchema struct {
	Name string `json:"name" binding:"required"`
}

type KeyValueSchema struct {
	Key      string `json:"key" binding:"required"`
	Value    string `json:"value" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Required bool   `json:"required" binding:"required"`
}

type ProviderGetResponseSchema struct {
	UUID   string           `json:"uuid" binding:"required"`
	Name   string           `json:"name" binding:"required"`
	Fields []KeyValueSchema `json:"fields" binding:"required"`
	Type   string           `json:"type" binding:"required"`
}

type MetaSchema struct {
	TotalItems int64 `json:"total" binding:"required"`
	Pages      int   `json:"pages" binding:"required"`
	Page       int   `json:"page" binding:"required"`
	PageSize   int   `json:"pageSize" binding:"required"`
}

type GetManyResponseSchema[T any] struct {
	Data []T        `json:"data" binding:"required"`
	Meta MetaSchema `json:"meta" binding:"required"`
}

type IncidentPostRequestBodySchema struct {
	Summary       string   `json:"summary" binding:"required"`
	HostsAffected []string `json:"hostsAffected" binding:"required"`
}

type UserGetResponseBodySchema struct {
	UUID    string                      `json:"uuid" binding:"required"`
	Name    string                      `json:"name" binding:"required"`
	Email   string                      `json:"email" binding:"required"`
	Teams   []TeamGetResponseBodySchema `json:"teams" binding:"required"`
	SlackID string                      `json:"slackID" binding:"required"`
	Admin   bool                        `json:"admin" binding:"required"`
}

type TeamGetResponseBodySchema struct {
	UUID string `json:"uuid" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type IncidentCommentGetResponseBodySchema struct {
	UUID        string                    `json:"uuid" binding:"required"`
	Comment     string                    `json:"comment" binding:"required"`
	CommentedBy UserGetResponseBodySchema `json:"commentedBy" binding:"required"`
	CommentedAt time.Time                 `json:"commentedAt" binding:"required"`
}

type IncidentGetResponseBodySchema struct {
	UUID            string                                 `json:"uuid" binding:"required"`
	Comments        []IncidentCommentGetResponseBodySchema `json:"comments" binding:"required"`
	HostsAffected   []HostMachineGetResponseBodySchema     `json:"hostsAffected" binding:"required"`
	Summary         string                                 `json:"summary" binding:"required"`
	Description     string                                 `json:"description" binding:"required"`
	CreatedAt       time.Time                              `json:"createdAt" binding:"required"`
	ResolvedAt      *time.Time                             `json:"resolvedAt" binding:"required"`
	ResolvedBy      *UserGetResponseBodySchema             `json:"resolvedBy" binding:"required"`
	ResolutionTeams []TeamGetResponseBodySchema            `json:"resolutionTeams" binding:"required"`
}

type HostMachineGetResponseBodySchema struct {
	UUID     string                    `json:"uuid" binding:"required"`
	OS       string                    `json:"os" binding:"required"`
	Hostname string                    `json:"hostname" binding:"required"`
	IP4      *string                   `json:"ip4" binding:"required"`
	IP6      *string                   `json:"ip6" binding:"required"`
	Team     TeamGetResponseBodySchema `json:"team" binding:"required"`
}

type ProviderPostRequestBodySchema struct {
	Name string `json:"name" binding:"required"`
}

type ProviderPutRequestBodySchema struct {
	ProviderPostRequestBodySchema
	Fields []KeyValueSchema `json:"fields" binding:"required"`
}

type HostMachinePostPutRequestBodySchema struct {
	OS       string  `json:"os" binding:"required"`
	Hostname string  `json:"hostname" binding:"required"`
	IP4      *string `json:"ip4" binding:"required"`
	IP6      *string `json:"ip6" binding:"required"`
	TeamID   string  `json:"teamID" binding:"required"`
}
