package utility

import (
	"time"
)

type ErrorResponseSchema struct {
	Error string `json:"error" binding:"required"`
}

type UserPostRequestBodySchema struct {
	Name     string   `json:"name" binding:"required,max=30,min=1"`
	Email    string   `json:"email" binding:"required,email,max=30,min=1"`
	Password string   `json:"password" binding:"required,max=72,min=1"`
	Teams    []string `json:"teams" binding:"required"`
}

type UserLoginRequestBodySchema struct {
	Email    string `json:"email" binding:"required,email,max=30,min=1"`
	Password string `json:"password" binding:"required,max=72,min=1"`
}

type TeamPostRequestBodySchema struct {
	Name string `json:"name" binding:"required,max=30"`
}

type KeyValueSchema struct {
	Key      string `json:"key" binding:"required,max=20"`
	Value    string `json:"value" binding:"required,max=30"`
	Type     string `json:"type" binding:"required,oneOf=string number boolean"`
	Required *bool  `json:"required" binding:"required"`
}

type ProviderGetResponseSchema struct {
	UUID   string           `json:"uuid" binding:"required,uuid4"`
	Name   string           `json:"name" binding:"required,max=30"`
	Fields []KeyValueSchema `json:"fields" binding:"required"`
	Type   string           `json:"type" binding:"required,oneOf=alert log"`
}

type MetaSchema struct {
	TotalItems int64 `json:"total" binding:"required,number"`
	Pages      int   `json:"pages" binding:"required,number"`
	Page       int   `json:"page" binding:"required,number"`
	PageSize   int   `json:"pageSize" binding:"required,number"`
}

type GetManyResponseSchema[T any] struct {
	Data []T        `json:"data" binding:"required"`
	Meta MetaSchema `json:"meta" binding:"required"`
}

type IncidentPostRequestBodySchema struct {
	Summary         string   `json:"summary" binding:"required,max=100,min=1"`
	Description     string   `json:"description" binding:"required,max=500,min=1"`
	ResolutionTeams []string `json:"resolutionTeams" binding:"required"`
	HostsAffected   []string `json:"hostsAffected" binding:"required"`
}

type UserGetResponseBodySchema struct {
	UUID    string                      `json:"uuid" binding:"required,uuid4"`
	Name    string                      `json:"name" binding:"required,max=30"`
	Email   string                      `json:"email" binding:"required,email,max=30"`
	Teams   []TeamGetResponseBodySchema `json:"teams" binding:"required"`
	SlackID string                      `json:"slackID" binding:"required,max=30"`
	Admin   *bool                       `json:"admin" binding:"required"`
}

type TeamGetResponseBodySchema struct {
	UUID string `json:"uuid" binding:"required,uuid4"`
	Name string `json:"name" binding:"required,max=30"`
}

type IncidentCommentGetResponseBodySchema struct {
	UUID        string                    `json:"uuid" binding:"required,uuid4"`
	Comment     string                    `json:"comment" binding:"required,max=200"`
	CommentedBy UserGetResponseBodySchema `json:"commentedBy" binding:"required"`
	CommentedAt time.Time                 `json:"commentedAt" binding:"required"`
}

type IncidentGetResponseBodySchema struct {
	UUID            string                                 `json:"uuid" binding:"required,uuid4"`
	Comments        []IncidentCommentGetResponseBodySchema `json:"comments" binding:"required"`
	HostsAffected   []HostMachineGetResponseBodySchema     `json:"hostsAffected" binding:"required"`
	Summary         string                                 `json:"summary" binding:"required,max=100"`
	Description     string                                 `json:"description" binding:"required,max=500"`
	CreatedAt       time.Time                              `json:"createdAt" binding:"required"`
	ResolvedAt      *time.Time                             `json:"resolvedAt" binding:"required"`
	ResolvedBy      *UserGetResponseBodySchema             `json:"resolvedBy" binding:"required"`
	ResolutionTeams []TeamGetResponseBodySchema            `json:"resolutionTeams" binding:"required"`
}

type HostMachineGetResponseBodySchema struct {
	UUID     string                    `json:"uuid" binding:"required,uuid4"`
	OS       string                    `json:"os" binding:"required,oneof=Windows Linux MacOS"`
	Hostname string                    `json:"hostname" binding:"required,hostname,max=255,min=1"`
	IP4      *string                   `json:"ip4" binding:"required_if=IP6 nil,ipv4"`
	IP6      *string                   `json:"ip6" binding:"required_if=IP4 nil,ipv6"`
	Team     TeamGetResponseBodySchema `json:"team" binding:"required"`
}

type ProviderPostRequestBodySchema struct {
	Name string `json:"name" binding:"required,max=30,min=1"`
}

type ProviderPutRequestBodySchema struct {
	ProviderPostRequestBodySchema
	Fields []KeyValueSchema `json:"fields" binding:"required"`
}

type HostMachinePostPutRequestBodySchema struct {
	OS       string  `json:"os" binding:"required,oneof=Windows Linux MacOS"`
	Hostname string  `json:"hostname" binding:"required,hostname,max=255,min=1"`
	IP4      *string `json:"ip4" binding:"required_if=IP6 nil,ipv4"`
	IP6      *string `json:"ip6" binding:"required_if=IP4 nil,ipv6"`
	TeamID   string  `json:"teamID" binding:"required,uuid4"`
}

type IncidentCommentPostRequestBodySchema struct {
	Comment string `json:"comment" binding:"required,max=200,min=1"`
}

type IncidentPutRequestBodySchema struct {
	Summary         string   `json:"summary" binding:"required,max=100,min=1"`
	Description     string   `json:"description" binding:"required,max=500,min=1"`
	HostsAffected   []string `json:"hostsAffected" binding:"required"`
	ResolutionTeams []string `json:"resolutionTeams" binding:"required"`
	Resolved        *bool    `json:"resolved" binding:"required"`
}
