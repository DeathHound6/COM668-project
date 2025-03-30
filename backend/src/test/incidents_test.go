package test_test

import (
	"com668-backend/middleware"
	"com668-backend/utility"
	"crypto/sha1"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"testing"
)

func TestGetIncidents(t *testing.T) {
	engine := setup()
	jwtString, err := getJWT(engine, TestAdminEmail, TestAdminPassword)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("GetIncidents", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/incidents", nil)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusOK
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("status code %d != %d", code, expected)
		}
		res, err := utility.ReadJSONStruct[utility.GetManyResponseSchema[utility.IncidentGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Data) == 0 {
			t.Fatal("no data")
		}
	})

	t.Run("GetIncidents ResolvedQuery", func(t *testing.T) {
		// resolved
		req, _ := http.NewRequest(http.MethodGet, "/incidents?resolved=true", nil)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusOK
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("status code %d != %d", code, expected)
		}
		res, err := utility.ReadJSONStruct[utility.GetManyResponseSchema[utility.IncidentGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Data) == 0 {
			t.Fatal("no data")
		}

		for _, incident := range res.Data {
			if incident.ResolvedAt == nil || incident.ResolvedBy == nil {
				t.Fatal("not resolved")
			}
		}

		// unresolved
		req, _ = http.NewRequest(http.MethodGet, "/incidents?resolved=false", nil)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer = makeRequest(engine, req)

		expected = http.StatusOK
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("status code %d != %d", code, expected)
		}
		res, err = utility.ReadJSONStruct[utility.GetManyResponseSchema[utility.IncidentGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Data) == 0 {
			t.Fatal("no data")
		}

		for _, incident := range res.Data {
			if incident.ResolvedAt != nil || incident.ResolvedBy != nil {
				t.Fatal("resolved")
			}
		}
	})

	t.Run("GetIncidents MyTeamsQuery", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/incidents?myTeams=true", nil)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusOK
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("status code %d != %d", code, expected)
		}
		res, err := utility.ReadJSONStruct[utility.GetManyResponseSchema[utility.IncidentGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Data) == 0 {
			t.Fatal("no data")
		}

		for _, incident := range res.Data {
			if len(incident.ResolutionTeams) == 0 {
				t.Fatal("no teams")
			}
			users := make([]string, 0)
			for _, team := range incident.ResolutionTeams {
				if len(team.Users) == 0 {
					t.Fatal("no users")
				}
				for _, user := range team.Users {
					users = append(users, user.Email)
				}
			}
			if len(users) == 0 {
				t.Fatal("no users")
			}
			if !slices.Contains(users, TestAdminEmail) {
				t.Fatalf("user %s not found", TestAdminEmail)
			}
		}
	})

	t.Run("GetIncidents HashQuery", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/incidents", nil)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusOK
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("status code %d != %d", code, expected)
		}
		res, err := utility.ReadJSONStruct[utility.GetManyResponseSchema[utility.IncidentGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Data) == 0 {
			t.Fatal("no data")
		}

		hash := res.Data[0].Hash
		req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/incidents?hash=%s", hash), nil)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer = makeRequest(engine, req)

		expected = http.StatusOK
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("status code %d != %d", code, expected)
		}
		res, err = utility.ReadJSONStruct[utility.GetManyResponseSchema[utility.IncidentGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Data) != 1 {
			t.Fatal("data length mismatch")
		}
		if res.Data[0].Hash != hash {
			t.Fatal("hash mismatch")
		}
	})
}

func TestCreateIncident(t *testing.T) {
	engine := setup()
	jwtString, err := getJWT(engine, TestAdminEmail, TestAdminPassword)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("CreateIncident", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/hosts", nil)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusOK
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("status code %d != %d", code, expected)
		}
		hostsRes, err := utility.ReadJSONStruct[utility.GetManyResponseSchema[utility.HostMachineGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(hostsRes.Data) == 0 {
			t.Fatal("no data")
		}

		req, _ = http.NewRequest(http.MethodGet, "/teams", nil)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer = makeRequest(engine, req)

		expected = http.StatusOK
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("status code %d != %d", code, expected)
		}
		teamsRes, err := utility.ReadJSONStruct[utility.GetManyResponseSchema[utility.TeamGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(teamsRes.Data) == 0 {
			t.Fatal("no data")
		}

		hasher := sha1.New()
		hasher.Write([]byte("Test Incident"))
		body, err := getJSONBodyAsReader(map[string]any{
			"summary":         "Test Incident",
			"description":     "Test Incident Details",
			"resolutionTeams": []string{teamsRes.Data[0].UUID},
			"hostsAffected":   []string{hostsRes.Data[0].UUID},
			"hash":            fmt.Sprintf("%x", hasher.Sum(nil)),
		})
		if err != nil {
			t.Fatal(err)
		}
		req, _ = http.NewRequest(http.MethodPost, "/incidents", body)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer = makeRequest(engine, req)

		expected = http.StatusCreated
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("status code %d != %d", code, expected)
		}
	})

	t.Run("CreateIncident InvalidBody", func(t *testing.T) {
		body, err := getJSONBodyAsReader(map[string]any{
			"invalidField": "invalidValue",
		})
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest(http.MethodPost, "/incidents", body)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusBadRequest
		if code := writer.Code; code != expected {
			if !strings.Contains(fmt.Sprint(code), "2") {
				resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
				if err != nil {
					t.Fatal(err)
				}
				t.Log(resp.Error)
			}
			t.Fatalf("status code %d != %d", code, expected)
		}
	})
}

func TestUpdateIncident(t *testing.T) {
	engine := setup()
	jwtString, err := getJWT(engine, TestAdminEmail, TestAdminPassword)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("UpdateIncident", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/incidents", nil)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusOK
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("status code %d != %d", code, expected)
		}
		res, err := utility.ReadJSONStruct[utility.GetManyResponseSchema[utility.IncidentGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Data) == 0 {
			t.Fatal("no data")
		}

		body, err := getJSONBodyAsReader(map[string]any{
			"summary":         "Updated Incident",
			"description":     "Updated Incident Details",
			"resolutionTeams": []string{res.Data[0].ResolutionTeams[0].UUID},
			"hostsAffected":   []string{res.Data[0].HostsAffected[0].UUID},
			"resolved":        true,
		})
		if err != nil {
			t.Fatal(err)
		}
		req, _ = http.NewRequest(http.MethodPut, fmt.Sprintf("/incidents/%s", res.Data[0].UUID), body)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer = makeRequest(engine, req)

		expected = http.StatusNoContent
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("status code %d != %d", code, expected)
		}
	})

	t.Run("UpdateIncident InvalidBody", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/incidents", nil)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusOK
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("status code %d != %d", code, expected)
		}
		res, err := utility.ReadJSONStruct[utility.GetManyResponseSchema[utility.IncidentGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Data) == 0 {
			t.Fatal("no data")
		}

		body, err := getJSONBodyAsReader(map[string]any{
			"invalidField": "invalidValue",
		})
		if err != nil {
			t.Fatal(err)
		}
		req, _ = http.NewRequest(http.MethodPut, fmt.Sprintf("/incidents/%s", res.Data[0].UUID), body)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer = makeRequest(engine, req)

		expected = http.StatusBadRequest
		if code := writer.Code; code != expected {
			if !strings.Contains(fmt.Sprint(code), "2") {
				resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
				if err != nil {
					t.Fatal(err)
				}
				t.Log(resp.Error)
			}
			t.Fatalf("status code %d != %d", code, expected)
		}
		errRes, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(errRes.Error, "' is required") {
			t.Fatal("error message mismatch")
		}
	})

	t.Run("UpdateIncident InvalidUUID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/incidents", nil)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusOK
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("status code %d != %d", code, expected)
		}
		res, err := utility.ReadJSONStruct[utility.GetManyResponseSchema[utility.IncidentGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Data) == 0 {
			t.Fatal("no data")
		}

		body, err := getJSONBodyAsReader(map[string]any{
			"summary":         "Updated Incident",
			"description":     "Updated Incident Details",
			"resolutionTeams": []string{res.Data[0].ResolutionTeams[0].UUID},
			"hostsAffected":   []string{res.Data[0].HostsAffected[0].UUID},
			"resolved":        true,
		})
		if err != nil {
			t.Fatal(err)
		}
		uuid, err := utility.GenerateRandomUUID()
		if err != nil {
			t.Fatal(err)
		}
		req, _ = http.NewRequest(http.MethodPut, fmt.Sprintf("/incidents/%s", uuid), body)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer = makeRequest(engine, req)

		expected = http.StatusNotFound
		if code := writer.Code; code != expected {
			if !strings.Contains(fmt.Sprint(code), "2") {
				resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
				if err != nil {
					t.Fatal(err)
				}
				t.Log(resp.Error)
			}
			t.Fatalf("status code %d != %d", code, expected)
		}
	})
}

func TestGetIncident(t *testing.T) {
	engine := setup()
	jwtString, err := getJWT(engine, TestAdminEmail, TestAdminPassword)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("GetIncident", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/incidents", nil)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusOK
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("status code %d != %d", code, expected)
		}
		res, err := utility.ReadJSONStruct[utility.GetManyResponseSchema[utility.IncidentGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Data) == 0 {
			t.Fatal("no data")
		}

		req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/incidents/%s", res.Data[0].UUID), nil)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer = makeRequest(engine, req)

		expected = http.StatusOK
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("status code %d != %d", code, expected)
		}
		incidentRes, err := utility.ReadJSONStruct[utility.IncidentGetResponseBodySchema](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if incidentRes.UUID != res.Data[0].UUID {
			t.Fatal("uuid mismatch")
		}
	})

	t.Run("GetIncident InvalidUUID", func(t *testing.T) {
		uuid, err := utility.GenerateRandomUUID()
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/incidents/%s", uuid), nil)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusNotFound
		if code := writer.Code; code != expected {
			if !strings.Contains(fmt.Sprint(code), "2") {
				resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
				if err != nil {
					t.Fatal(err)
				}
				t.Log(resp.Error)
			}
			t.Fatalf("status code %d != %d", code, expected)
		}
	})
}

func TestCreateIncidentComment(t *testing.T) {
	engine := setup()
	jwtString, err := getJWT(engine, TestAdminEmail, TestAdminPassword)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("CreateIncidentComment", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/incidents", nil)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusOK
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("status code %d != %d", code, expected)
		}
		res, err := utility.ReadJSONStruct[utility.GetManyResponseSchema[utility.IncidentGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Data) == 0 {
			t.Fatal("no data")
		}

		body, err := getJSONBodyAsReader(map[string]any{
			"comment": "Test Comment",
		})
		if err != nil {
			t.Fatal(err)
		}
		req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("/incidents/%s/comments", res.Data[0].UUID), body)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer = makeRequest(engine, req)

		expected = http.StatusCreated
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("status code %d != %d", code, expected)
		}
	})

	t.Run("CreateIncidentComment InvalidBody", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/incidents", nil)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusOK
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("status code %d != %d", code, expected)
		}
		res, err := utility.ReadJSONStruct[utility.GetManyResponseSchema[utility.IncidentGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Data) == 0 {
			t.Fatal("no data")
		}

		body, err := getJSONBodyAsReader(map[string]any{
			"invalidField": "invalidValue",
		})
		if err != nil {
			t.Fatal(err)
		}
		req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("/incidents/%s/comments", res.Data[0].UUID), body)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer = makeRequest(engine, req)

		expected = http.StatusBadRequest
		if code := writer.Code; code != expected {
			if !strings.Contains(fmt.Sprint(code), "2") {
				resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
				if err != nil {
					t.Fatal(err)
				}
				t.Log(resp.Error)
			}
			t.Fatalf("status code %d != %d", code, expected)
		}
		if !strings.Contains(writer.Body.String(), "'comment' is required") {
			t.Fatal("error message mismatch")
		}
	})

	t.Run("CreateIncidentComment InvalidUUID", func(t *testing.T) {
		body, err := getJSONBodyAsReader(map[string]any{
			"comment": "Test Comment",
		})
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest(http.MethodPost, "/incidents/invalidUUID/comments", body)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusBadRequest
		if code := writer.Code; code != expected {
			if !strings.Contains(fmt.Sprint(code), "2") {
				resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
				if err != nil {
					t.Fatal(err)
				}
				t.Log(resp.Error)
			}
			t.Fatalf("status code %d != %d", code, expected)
		}
		res, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(res.Error, "invalid incident UUID") {
			t.Fatal("error message mismatch")
		}

		uuid, err := utility.GenerateRandomUUID()
		if err != nil {
			t.Fatal(err)
		}
		req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("/incidents/%s/comments", uuid), body)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer = makeRequest(engine, req)

		expected = http.StatusNotFound
		if code := writer.Code; code != expected {
			if !strings.Contains(fmt.Sprint(code), "2") {
				resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
				if err != nil {
					t.Fatal(err)
				}
				t.Log(resp.Error)
			}
			t.Fatalf("status code %d != %d", code, expected)
		}
	})
}

func TestDeleteIncidentComment(t *testing.T) {
	engine := setup()
	jwtString, err := getJWT(engine, TestAdminEmail, TestAdminPassword)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("DeleteIncidentComment InvalidIncidentUUID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/incidents", nil)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusOK
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("status code %d != %d", code, expected)
		}
		resp, err := utility.ReadJSONStruct[utility.GetManyResponseSchema[utility.IncidentGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Data) == 0 {
			t.Fatal("no data")
		}
		if len(resp.Data[0].Comments) == 0 {
			t.Fatal("no comments")
		}

		uuid, err := utility.GenerateRandomUUID()
		if err != nil {
			t.Fatal(err)
		}
		req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("/incidents/%s/comments/%s", uuid, resp.Data[0].Comments[0].UUID), nil)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer = makeRequest(engine, req)

		expected = http.StatusNotFound
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("status code %d != %d", code, expected)
		}

		req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("/incidents/invalidUUID/comments/%s", resp.Data[0].Comments[0].UUID), nil)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer = makeRequest(engine, req)

		expected = http.StatusBadRequest
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("status code %d != %d", code, expected)
		}
		res, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(res.Error, "invalid incident UUID") {
			t.Fatal("error message mismatch")
		}
	})

	t.Run("DeleteIncidentComment InvalidCommentUUID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/incidents", nil)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusOK
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("status code %d != %d", code, expected)
		}
		resp, err := utility.ReadJSONStruct[utility.GetManyResponseSchema[utility.IncidentGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Data) == 0 {
			t.Fatal("no data")
		}
		if len(resp.Data[0].Comments) == 0 {
			t.Fatal("no comments")
		}

		uuid, err := utility.GenerateRandomUUID()
		if err != nil {
			t.Fatal(err)
		}
		req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("/incidents/%s/comments/%s", resp.Data[0].UUID, uuid), nil)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer = makeRequest(engine, req)

		expected = http.StatusNotFound
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("status code %d != %d", code, expected)
		}

		req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("/incidents/%s/comments/invalidUUID", resp.Data[0].UUID), nil)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer = makeRequest(engine, req)

		expected = http.StatusBadRequest
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("status code %d != %d", code, expected)
		}
		res, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(res.Error, "invalid comment UUID") {
			t.Fatal("error message mismatch")
		}
	})

	t.Run("DeleteIncidentComment Forbidden", func(t *testing.T) {
		jwtString, err := getJWT(engine, TestAdminEmail, TestAdminPassword)
		if err != nil {
			t.Fatal(err)
		}

		// fetch incidents to get UUID for comment create
		req, _ := http.NewRequest(http.MethodGet, "/incidents", nil)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusOK
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("status code %d != %d", code, expected)
		}
		resp, err := utility.ReadJSONStruct[utility.GetManyResponseSchema[utility.IncidentGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Data) == 0 {
			t.Fatal("no data")
		}

		// create a comment as admin user
		body, err := getJSONBodyAsReader(map[string]any{
			"comment": "Test Comment",
		})
		if err != nil {
			t.Fatal(err)
		}
		req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("/incidents/%s/comments", resp.Data[0].UUID), body)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer = makeRequest(engine, req)

		expected = http.StatusCreated
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
		}

		// fetch incidents to get newly created comment UUID
		req, _ = http.NewRequest(http.MethodGet, "/incidents", nil)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer = makeRequest(engine, req)

		expected = http.StatusOK
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("status code %d != %d", code, expected)
		}
		resp, err = utility.ReadJSONStruct[utility.GetManyResponseSchema[utility.IncidentGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Data) == 0 {
			t.Fatal("no data")
		}

		// delete the comment as non-admin user
		jwtString, err = getJWT(engine, TestUserEmail, TestUserPassword)
		if err != nil {
			t.Fatal(err)
		}
		req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("/incidents/%s/comments/%s", resp.Data[0].UUID, resp.Data[0].Comments[0].UUID), nil)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer = makeRequest(engine, req)

		expected = http.StatusForbidden
		if code := writer.Code; code != expected {
			if !strings.Contains(fmt.Sprint(code), "2") {
				resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
				if err != nil {
					t.Fatal(err)
				}
				t.Log(resp.Error)
			}
			t.Fatalf("status code %d != %d", code, expected)
		}
		res, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(res.Error, "you are not allowed to delete this comment") {
			t.Fatal("error message mismatch")
		}
	})

	t.Run("DeleteIncidentComment", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/incidents", nil)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusOK
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("status code %d != %d", code, expected)
		}
		res, err := utility.ReadJSONStruct[utility.GetManyResponseSchema[utility.IncidentGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Data) == 0 {
			t.Fatal("no data")
		}
		if len(res.Data[0].Comments) == 0 {
			t.Fatal("no comments")
		}

		req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("/incidents/%s/comments/%s", res.Data[0].UUID, res.Data[0].Comments[0].UUID), nil)
		req.Header.Set(middleware.AuthHeaderNameString, jwtString)
		writer = makeRequest(engine, req)

		expected = http.StatusNoContent
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("status code %d != %d", code, expected)
		}
	})
}
