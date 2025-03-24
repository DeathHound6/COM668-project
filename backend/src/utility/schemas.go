package utility

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

type BodySchema interface {
	Validate() (int, error)
}
type ResponseSchema interface {
	String() string
	JSON() map[string]any
}

type ErrorResponseSchema struct {
	ResponseSchema
	Error string `json:"error"`
}

func (e ErrorResponseSchema) JSON() map[string]any {
	return map[string]any{"error": e.Error}
}
func (e ErrorResponseSchema) String() string {
	return fmt.Sprintf("{'error': '%s'}", e.Error)
}

type UserPostRequestBodySchema struct {
	BodySchema
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Teams    []string `json:"teams"`
}

func (u UserPostRequestBodySchema) Validate() (int, error) {
	if len(u.Name) == 0 {
		return 400, errors.New("'name' is required")
	}
	if len(u.Name) > 30 {
		return 400, errors.New("'name' cannot be longer than 30 characters")
	}
	if len(u.Email) == 0 {
		return 400, errors.New("'email' is required")
	}
	if len(u.Email) > 30 {
		return 400, errors.New("'email' cannot be longer than 30 characters")
	}
	matched, err := regexp.Match("[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+.[A-Za-z]{2,}", []byte(u.Email))
	if err != nil {
		return 500, err
	}
	if !matched {
		return 400, errors.New("'email' is not valid")
	}
	if len(u.Password) == 0 {
		return 400, errors.New("'password' is required")
	}
	if len(u.Password) > 72 {
		return 400, errors.New("'password' cannot be greater than 72 characters")
	}
	return -1, nil
}

type UserLoginRequestBodySchema struct {
	BodySchema
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u UserLoginRequestBodySchema) Validate() (int, error) {
	if len(u.Email) == 0 {
		return 400, errors.New("'email' is required")
	}
	if len(u.Password) == 0 {
		return 400, errors.New("'password' is required")
	}
	return -1, nil
}

type KeyValueSchema struct {
	ResponseSchema
	BodySchema
	Key      string `json:"key"`
	Value    string `json:"value"`
	Type     string `json:"type"`
	Required *bool  `json:"required"`
}

func (k KeyValueSchema) JSON() map[string]any {
	return map[string]any{"key": k.Key, "value": k.Value, "type": k.Type, "required": k.Required}
}
func (k KeyValueSchema) String() string {
	required := "nil"
	if k.Required != nil {
		required = fmt.Sprintf("%t", *k.Required)
	}
	return fmt.Sprintf("{'key': '%s', 'value': '%s', 'type': '%s', 'required': %s}", k.Key, k.Value, k.Type, required)
}
func (k KeyValueSchema) Validate() (int, error) {
	if len(k.Key) == 0 {
		return 400, errors.New("'key' is required")
	}
	if len(k.Key) > 20 {
		return 400, errors.New("'key' cannot be longer than 20 characters")
	}
	if len(k.Value) == 0 {
		return 400, errors.New("'value' is required")
	}
	if len(k.Value) > 30 {
		return 400, errors.New("'value' cannot be longer than 30 characters")
	}
	if !SliceHasElement([]string{"string", "number", "bool"}, k.Type) {
		return 400, errors.New("'type' must be one of 'string', 'number', or 'bool'")
	}
	if k.Required == nil {
		return 400, errors.New("'required' is required")
	}
	return -1, nil
}

type ProviderGetResponseSchema struct {
	ResponseSchema
	UUID   string           `json:"uuid"`
	Name   string           `json:"name"`
	Fields []KeyValueSchema `json:"fields"`
	Type   string           `json:"type"`
}

func (p ProviderGetResponseSchema) JSON() map[string]any {
	fields := make([]map[string]any, 0)
	for _, f := range p.Fields {
		fields = append(fields, f.JSON())
	}
	return map[string]any{"uuid": p.UUID, "name": p.Name, "fields": fields, "type": p.Type}
}
func (p ProviderGetResponseSchema) String() string {
	fields := make([]string, 0)
	for _, f := range p.Fields {
		fields = append(fields, f.String())
	}
	return fmt.Sprintf("{'uuid': '%s', 'name': '%s', 'fields': [%s], 'type': '%s'}", p.UUID, p.Name, strings.Join(fields, " "), p.Type)
}

type MetaSchema struct {
	ResponseSchema
	TotalItems int64 `json:"total"`
	Pages      int   `json:"pages"`
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
}

func (m MetaSchema) JSON() map[string]any {
	return map[string]any{"total": m.TotalItems, "pages": m.Pages, "page": m.Page, "pageSize": m.PageSize}
}
func (m MetaSchema) String() string {
	return fmt.Sprintf("{'total': %d, 'pages': %d, 'page': %d, 'pageSize': %d}", m.TotalItems, m.Pages, m.Page, m.PageSize)
}

type GetManyResponseSchema[T ResponseSchema] struct {
	ResponseSchema
	Data []T        `json:"data"`
	Meta MetaSchema `json:"meta"`
}

func (g GetManyResponseSchema[T]) JSON() map[string]any {
	data := make([]map[string]any, 0)
	for _, d := range g.Data {
		data = append(data, d.JSON())
	}
	return map[string]any{"data": data, "meta": g.Meta.JSON()}
}
func (g GetManyResponseSchema[T]) String() string {
	data := make([]string, 0)
	for _, d := range g.Data {
		data = append(data, d.String())
	}
	return fmt.Sprintf("{'data': [%s], 'meta': %s}", strings.Join(data, " "), g.Meta.String())
}

type IncidentPostRequestBodySchema struct {
	BodySchema
	Summary         string   `json:"summary"`
	Description     string   `json:"description"`
	ResolutionTeams []string `json:"resolutionTeams"`
	HostsAffected   []string `json:"hostsAffected"`
	Hash            string   `json:"hash"`
}

func (i IncidentPostRequestBodySchema) Validate() (int, error) {
	if len(i.Summary) == 0 {
		return 400, errors.New("'summary' is required")
	}
	if len(i.Summary) > 100 {
		return 400, errors.New("'summary' cannot be longer than 100 characters")
	}
	if len(i.Description) == 0 {
		return 400, errors.New("'description' is required")
	}
	if len(i.Description) > 500 {
		return 400, errors.New("'description' cannot be longer than 500 characters")
	}
	for _, r := range i.ResolutionTeams {
		if _, err := uuid.Parse(r); err != nil {
			return 400, errors.New("'resolutionTeams' must be a list of valid UUIDs")
		}
	}
	for _, h := range i.HostsAffected {
		if _, err := uuid.Parse(h); err != nil {
			return 400, errors.New("'hostsAffected' must be a list of valid UUIDs")
		}
	}
	if len(i.Hash) == 0 {
		return 400, errors.New("'hash' is required")
	}
	if len(i.Hash) > 40 {
		return 400, errors.New("'hash' cannot be longer than 40 characters")
	}
	return -1, nil
}

type UserGetResponseBodySchema struct {
	ResponseSchema
	UUID    string                      `json:"uuid"`
	Name    string                      `json:"name"`
	Email   string                      `json:"email"`
	Teams   []TeamGetResponseBodySchema `json:"teams"`
	SlackID string                      `json:"slackID"`
	Admin   *bool                       `json:"admin"`
}

func (u UserGetResponseBodySchema) JSON() map[string]any {
	teams := make([]map[string]any, 0)
	for _, t := range u.Teams {
		teams = append(teams, t.JSON())
	}
	return map[string]any{"uuid": u.UUID, "name": u.Name, "email": u.Email, "teams": teams, "slackID": u.SlackID, "admin": u.Admin}
}
func (u UserGetResponseBodySchema) String() string {
	teams := make([]string, 0)
	for _, t := range u.Teams {
		teams = append(teams, t.String())
	}
	admin := "nil"
	if u.Admin != nil {
		admin = fmt.Sprintf("%t", *u.Admin)
	}
	return fmt.Sprintf("{'uuid': '%s', 'name': '%s', 'email': '%s', 'teams': [%s], 'slackID': '%s', 'admin': %s}", u.UUID, u.Name, u.Email, strings.Join(teams, " "), u.SlackID, admin)
}

type TeamGetResponseBodySchema struct {
	ResponseSchema
	UUID  string                      `json:"uuid"`
	Name  string                      `json:"name"`
	Users []UserGetResponseBodySchema `json:"users"`
}

func (t TeamGetResponseBodySchema) JSON() map[string]any {
	users := make([]map[string]any, 0)
	for _, u := range t.Users {
		users = append(users, u.JSON())
	}
	return map[string]any{"uuid": t.UUID, "name": t.Name, "users": users}
}
func (t TeamGetResponseBodySchema) String() string {
	users := make([]string, 0)
	for _, u := range t.Users {
		users = append(users, u.String())
	}
	return fmt.Sprintf("{'uuid': '%s', 'name': '%s', 'users': ['%s']}", t.UUID, t.Name, strings.Join(users, " "))
}

type IncidentCommentGetResponseBodySchema struct {
	ResponseSchema
	UUID        string                    `json:"uuid"`
	Comment     string                    `json:"comment"`
	CommentedBy UserGetResponseBodySchema `json:"commentedBy"`
	CommentedAt time.Time                 `json:"commentedAt"`
}

func (i IncidentCommentGetResponseBodySchema) JSON() map[string]any {
	return map[string]any{"uuid": i.UUID, "comment": i.Comment, "commentedBy": i.CommentedBy.JSON(), "commentedAt": i.CommentedAt}
}
func (i IncidentCommentGetResponseBodySchema) String() string {
	return fmt.Sprintf("{'uuid': '%s', 'comment': '%s', 'commentedBy': %s, 'commentedAt': '%s'}", i.UUID, i.Comment, i.CommentedBy.String(), i.CommentedAt)
}

type IncidentGetResponseBodySchema struct {
	ResponseSchema
	UUID            string                                 `json:"uuid"`
	Comments        []IncidentCommentGetResponseBodySchema `json:"comments"`
	HostsAffected   []HostMachineGetResponseBodySchema     `json:"hostsAffected"`
	Summary         string                                 `json:"summary"`
	Description     string                                 `json:"description"`
	CreatedAt       time.Time                              `json:"createdAt"`
	ResolvedAt      *time.Time                             `json:"resolvedAt"`
	ResolvedBy      *UserGetResponseBodySchema             `json:"resolvedBy"`
	ResolutionTeams []TeamGetResponseBodySchema            `json:"resolutionTeams"`
	Hash            string                                 `json:"hash"`
}

func (i IncidentGetResponseBodySchema) JSON() map[string]any {
	comments := make([]map[string]any, 0)
	for _, c := range i.Comments {
		comments = append(comments, c.JSON())
	}
	hosts := make([]map[string]any, 0)
	for _, h := range i.HostsAffected {
		hosts = append(hosts, h.JSON())
	}
	resolutionTeams := make([]map[string]any, 0)
	for _, r := range i.ResolutionTeams {
		resolutionTeams = append(resolutionTeams, r.JSON())
	}
	var resolvedBy *map[string]any = nil
	if i.ResolvedBy != nil {
		resolvedBy = Pointer(i.ResolvedBy.JSON())
	}
	return map[string]any{"uuid": i.UUID, "comments": comments, "hostsAffected": hosts, "summary": i.Summary, "description": i.Description, "createdAt": i.CreatedAt, "resolvedAt": i.ResolvedAt, "resolvedBy": resolvedBy, "resolutionTeams": resolutionTeams, "hash": i.Hash}
}
func (i IncidentGetResponseBodySchema) String() string {
	comments := make([]string, 0)
	for _, c := range i.Comments {
		comments = append(comments, c.String())
	}
	hosts := make([]string, 0)
	for _, h := range i.HostsAffected {
		hosts = append(hosts, h.String())
	}
	resolutionTeams := make([]string, 0)
	for _, r := range i.ResolutionTeams {
		resolutionTeams = append(resolutionTeams, r.String())
	}
	resolvedAt := "nil"
	if i.ResolvedAt != nil {
		resolvedAt = fmt.Sprintf("'%s'", *i.ResolvedAt)
	}
	resolvedBy := "nil"
	if i.ResolvedBy != nil {
		resolvedBy = i.ResolvedBy.String()
	}
	return fmt.Sprintf("{'uuid': '%s', 'comments': [%s], 'hostsAffected': [%s], 'summary': '%s', 'description': '%s', 'createdAt': '%s', 'resolvedAt': '%s', 'resolvedBy': %s, 'resolutionTeams': [%s], 'hash': '%s'}", i.UUID, strings.Join(comments, " "), strings.Join(hosts, " "), i.Summary, i.Description, i.CreatedAt, resolvedAt, resolvedBy, strings.Join(resolutionTeams, " "), i.Hash)
}

type HostMachineGetResponseBodySchema struct {
	ResponseSchema
	UUID     string                    `json:"uuid"`
	OS       string                    `json:"os"`
	Hostname string                    `json:"hostname"`
	IP4      *string                   `json:"ip4"`
	IP6      *string                   `json:"ip6"`
	Team     TeamGetResponseBodySchema `json:"team"`
}

func (h HostMachineGetResponseBodySchema) JSON() map[string]any {
	return map[string]any{"uuid": h.UUID, "os": h.OS, "hostname": h.Hostname, "ip4": h.IP4, "ip6": h.IP6, "team": h.Team.JSON()}
}
func (h HostMachineGetResponseBodySchema) String() string {
	ip4 := "nil"
	if h.IP4 != nil {
		ip4 = fmt.Sprintf("'%s'", *h.IP4)
	}
	ip6 := "nil"
	if h.IP6 != nil {
		ip6 = fmt.Sprintf("'%s'", *h.IP6)
	}
	return fmt.Sprintf("{'uuid': '%s', 'os': '%s', 'hostname': '%s', 'ip4': %s, 'ip6': %s, 'team': %s}", h.UUID, h.OS, h.Hostname, ip4, ip6, h.Team.String())
}

type ProviderPostRequestBodySchema struct {
	BodySchema
	Name string `json:"name"`
}

func (p ProviderPostRequestBodySchema) Validate() (int, error) {
	if len(p.Name) == 0 {
		return 400, errors.New("'name' is required")
	}
	if len(p.Name) > 30 {
		return 400, errors.New("'name' cannot be longer than 30 characters")
	}
	return -1, nil
}

type ProviderPutRequestBodySchema struct {
	BodySchema
	ProviderPostRequestBodySchema
	Fields []KeyValueSchema `json:"fields"`
}

func (p ProviderPutRequestBodySchema) Validate() (int, error) {
	if len(p.Name) == 0 {
		return 400, errors.New("'name' is required")
	}
	if len(p.Name) > 30 {
		return 400, errors.New("'name' cannot be longer than 30 characters")
	}
	for _, f := range p.Fields {
		if status, err := f.Validate(); err != nil {
			return status, err
		}
	}
	return -1, nil
}

type HostMachinePostPutRequestBodySchema struct {
	BodySchema
	OS       string  `json:"os"`
	Hostname string  `json:"hostname"`
	IP4      *string `json:"ip4"`
	IP6      *string `json:"ip6"`
	TeamID   string  `json:"teamID"`
}

func (h HostMachinePostPutRequestBodySchema) Validate() (int, error) {
	if len(h.OS) == 0 {
		return 400, errors.New("'os' is required")
	}
	if !SliceHasElement([]string{"Windows", "Linux", "MacOS"}, h.OS) {
		return 400, errors.New("'os' must be either 'Windows', 'Linux', or 'MacOS'")
	}
	if len(h.Hostname) == 0 {
		return 400, errors.New("'hostname' is required")
	}
	if len(h.Hostname) > 30 {
		return 400, errors.New("'hostname' cannot be longer than 30 characters")
	}
	if h.IP4 == nil && h.IP6 == nil {
		return 400, errors.New("either 'ip4' or 'ip6' is required")
	}
	if h.IP4 != nil {
		if len(*h.IP4) == 0 {
			return 400, errors.New("'ip4' is required")
		}
		ip := net.ParseIP(*h.IP4)
		if ip == nil || ip.To4() == nil {
			return 400, errors.New("'ip4' is not a valid IPv4 address")
		}
	}
	if h.IP6 != nil {
		if len(*h.IP6) == 0 {
			return 400, errors.New("'ip6' is required")
		}
		ip := net.ParseIP(*h.IP6)
		if ip == nil || ip.To16() == nil {
			return 400, errors.New("'ip6' is not a valid IPv6 address")
		}
	}
	if _, err := uuid.Parse(h.TeamID); err != nil {
		return 400, errors.New("'teamID' must be a valid UUID")
	}
	return -1, nil
}

type IncidentCommentPostRequestBodySchema struct {
	BodySchema
	Comment string `json:"comment"`
}

func (i IncidentCommentPostRequestBodySchema) Validate() (int, error) {
	if len(i.Comment) == 0 {
		return 400, errors.New("'comment' is required")
	}
	if len(i.Comment) > 200 {
		return 400, errors.New("'comment' cannot be longer than 200 characters")
	}
	return -1, nil
}

type IncidentPutRequestBodySchema struct {
	BodySchema
	Summary         string   `json:"summary"`
	Description     string   `json:"description"`
	HostsAffected   []string `json:"hostsAffected"`
	ResolutionTeams []string `json:"resolutionTeams"`
	Resolved        *bool    `json:"resolved"`
}

func (i IncidentPutRequestBodySchema) Validate() (int, error) {
	if len(i.Summary) == 0 {
		return 400, errors.New("'summary' is required")
	}
	if len(i.Summary) > 100 {
		return 400, errors.New("'summary' cannot be longer than 100 characters")
	}
	if len(i.Description) == 0 {
		return 400, errors.New("'description' is required")
	}
	if len(i.Description) > 500 {
		return 400, errors.New("'description' cannot be longer than 500 characters")
	}
	for _, r := range i.ResolutionTeams {
		if _, err := uuid.Parse(r); err != nil {
			return 400, errors.New("'resolutionTeams' must be a list of valid UUIDs")
		}
	}
	for _, h := range i.HostsAffected {
		if _, err := uuid.Parse(h); err != nil {
			return 400, errors.New("'hostsAffected' must be a list of valid UUIDs")
		}
	}
	if i.Resolved == nil {
		return 400, errors.New("'resolved' is required")
	}
	return -1, nil
}
