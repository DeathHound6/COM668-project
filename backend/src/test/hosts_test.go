package test_test

import (
	"com668-backend/middleware"
	"com668-backend/utility"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"testing"
)

func TestGetHosts(t *testing.T) {
	engine := setup()
	jwtString, err := getJWT(engine, TestAdminEmail, TestAdminPassword)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("GetHosts", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest(http.MethodGet, "/hosts", nil)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusOK
		if code := writer.Code; code != expected {
			errorResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(errorResp.Error)
			t.Fatalf("Status code %d != %d", code, expected)
		}
		resp, err := utility.ReadJSONStruct[utility.GetManyResponseSchema[*utility.HostMachineGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Data) == 0 {
			t.Fatal("no data was returned")
		}
	})

	t.Run("GetHosts InvalidCommonParams", func(t *testing.T) {
		t.Parallel()
		page := "invalid"
		pageSize := "10"

		// invalid page
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/hosts?page=%s&pageSize=%s", page, pageSize), nil)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusBadRequest
		if code := writer.Code; code != expected {
			if !strings.HasPrefix(fmt.Sprint(code), "2") {
				resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
				if err != nil {
					t.Fatal(err)
				}
				t.Log(resp.Error)
			}
			t.Fatalf("Status code %d != %d", code, expected)
		}
		resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if resp.Error != "page query parameter must be an integer" {
			t.Log(resp.Error)
			t.Fatal("error response message was not expected message")
		}

		// invalid pageSize
		page = "1"
		pageSize = "invalid"

		req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/hosts?page=%s&pageSize=%s", page, pageSize), nil)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer = makeRequest(engine, req)

		expected = http.StatusBadRequest
		if code := writer.Code; code != expected {
			if !strings.HasPrefix(fmt.Sprint(code), "2") {
				resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
				if err != nil {
					t.Fatal(err)
				}
				t.Log(resp.Error)
			}
			t.Fatalf("Status code %d != %d", code, expected)
		}

		resp, err = utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if resp.Error != "pageSize query parameter must be an integer" {
			t.Log(resp.Error)
			t.Fatal("error response message was not expected message")
		}
	})

	t.Run("GetHosts ValidHostnames", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest(http.MethodGet, "/hosts", nil)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusOK
		if code := writer.Code; code != expected {
			errorResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(errorResp.Error)
			t.Fatalf("Status code %d != %d", code, expected)
		}
		resp, err := utility.ReadJSONStruct[utility.GetManyResponseSchema[*utility.HostMachineGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Data) == 0 {
			t.Fatal("no data was returned")
		}

		hosts := []string{resp.Data[0].Hostname}
		req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/hosts?hostnames=%s", strings.Join(hosts, ",")), nil)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer = makeRequest(engine, req)

		expected = http.StatusOK
		if code := writer.Code; code != expected {
			errorResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(errorResp.Error)
			t.Fatalf("Status code %d != %d", code, expected)
		}
		resp, err = utility.ReadJSONStruct[utility.GetManyResponseSchema[*utility.HostMachineGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Data) == 0 {
			t.Fatal("no data was returned")
		}
		responseHostNames := make([]string, 0)
		for _, host := range resp.Data {
			responseHostNames = append(responseHostNames, host.Hostname)
		}
		for _, host := range hosts {
			if !slices.Contains(responseHostNames, host) {
				t.Fatalf("host %s was not returned", host)
			}
		}
	})

	t.Run("GetHosts InvalidHostnames", func(t *testing.T) {
		t.Parallel()
		hosts := []string{"invalid"}
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/hosts?hostnames=%s", strings.Join(hosts, ",")), nil)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusOK
		if code := writer.Code; code != expected {
			errorResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(errorResp.Error)
			t.Fatalf("Status code %d != %d", code, expected)
		}
		resp, err := utility.ReadJSONStruct[utility.GetManyResponseSchema[*utility.HostMachineGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Data) != 0 {
			t.Fatal("data was returned")
		}
	})
}

func TestGetHost(t *testing.T) {
	engine := setup()
	jwtString, err := getJWT(engine, TestAdminEmail, TestAdminPassword)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("GetHost", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest(http.MethodGet, "/hosts", nil)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusOK
		if code := writer.Code; code != expected {
			errorResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(errorResp.Error)
			t.Fatalf("Status code %d != %d", code, expected)
		}
		resp, err := utility.ReadJSONStruct[utility.GetManyResponseSchema[*utility.HostMachineGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Data) == 0 {
			t.Fatal("no data was returned")
		}

		req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/hosts/%s", resp.Data[0].UUID), nil)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer = makeRequest(engine, req)

		expected = http.StatusOK
		if code := writer.Code; code != expected {
			errorResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(errorResp.Error)
			t.Fatalf("Status code %d != %d", code, expected)
		}
		hostResp, err := utility.ReadJSONStruct[utility.HostMachineGetResponseBodySchema](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if hostResp.UUID != resp.Data[0].UUID {
			t.Fatalf("UUID %s != %s", hostResp.UUID, resp.Data[0].UUID)
		}
	})

	t.Run("GetHost InvalidUUID", func(t *testing.T) {
		t.Parallel()
		uuid, err := utility.GenerateRandomUUID()
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/hosts/%s", uuid), nil)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusNotFound
		if code := writer.Code; code != expected {
			if !strings.HasPrefix(fmt.Sprint(code), "2") {
				errorResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
				if err != nil {
					t.Fatal(err)
				}
				t.Log(errorResp.Error)
			}
			t.Fatalf("Status code %d != %d", code, expected)
		}
	})
}

func TestCreateHost(t *testing.T) {
	engine := setup()
	jwtString, err := getJWT(engine, TestAdminEmail, TestAdminPassword)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("CreateHost", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest(http.MethodGet, "/teams", nil)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusOK
		if code := writer.Code; code != expected {
			errorResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(errorResp.Error)
			t.Fatalf("Status code %d != %d", code, expected)
		}
		teamsResp, err := utility.ReadJSONStruct[utility.GetManyResponseSchema[*utility.TeamGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(teamsResp.Data) == 0 {
			t.Fatal("no data was returned")
		}

		body, err := getJSONBodyAsReader(map[string]any{
			"hostname": "test-host",
			"ip4":      "192.168.0.4",
			"ip6":      "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			"os":       "Linux",
			"teamID":   teamsResp.Data[0].UUID,
		})
		if err != nil {
			t.Fatal(err)
		}
		req, _ = http.NewRequest(http.MethodPost, "/hosts", body)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer = makeRequest(engine, req)

		expected = http.StatusCreated
		if code := writer.Code; code != expected {
			errorResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(errorResp.Error)
			t.Fatalf("Status code %d != %d", code, expected)
		}
	})

	t.Run("CreateHost InvalidBody", func(t *testing.T) {
		t.Parallel()
		body, err := getJSONBodyAsReader(map[string]any{
			"invalidField": "invalidValue",
		})
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest(http.MethodPost, "/hosts", body)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusBadRequest
		if code := writer.Code; code != expected {
			if !strings.HasPrefix(fmt.Sprint(code), "2") {
				errorResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
				if err != nil {
					t.Fatal(err)
				}
				t.Log(errorResp.Error)
			}
			t.Fatalf("Status code %d != %d", code, expected)
		}
		errResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(errResp.Error, "' is required") {
			t.Log(errResp.Error)
			t.Fatal("error response message was not expected message")
		}
	})
}

func TestUpdateHost(t *testing.T) {
	engine := setup()
	jwtString, err := getJWT(engine, TestAdminEmail, TestAdminPassword)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("UpdateHost", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest(http.MethodGet, "/teams", nil)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusOK
		if code := writer.Code; code != expected {
			errorResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(errorResp.Error)
			t.Fatalf("Status code %d != %d", code, expected)
		}
		teamsResp, err := utility.ReadJSONStruct[utility.GetManyResponseSchema[*utility.TeamGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(teamsResp.Data) == 0 {
			t.Fatal("no data was returned")
		}

		req, _ = http.NewRequest(http.MethodGet, "/hosts", nil)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer = makeRequest(engine, req)

		expected = http.StatusOK
		if code := writer.Code; code != expected {
			errorResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(errorResp.Error)
			t.Fatalf("Status code %d != %d", code, expected)
		}
		resp, err := utility.ReadJSONStruct[utility.GetManyResponseSchema[*utility.HostMachineGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Data) == 0 {
			t.Fatal("no data was returned")
		}

		body, err := getJSONBodyAsReader(map[string]any{
			"hostname": "test-host2",
			"ip4":      "192.168.0.2",
			"ip6":      "2001:0db8:85a3:0000:0000:8a2e:4500:7334",
			"os":       "Linux",
			"teamID":   teamsResp.Data[0].UUID,
		})
		if err != nil {
			t.Fatal(err)
		}
		req, _ = http.NewRequest(http.MethodPut, fmt.Sprintf("/hosts/%s", resp.Data[0].UUID), body)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer = makeRequest(engine, req)

		expected = http.StatusNoContent
		if code := writer.Code; code != expected {
			errorResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(errorResp.Error)
			t.Fatalf("Status code %d != %d", code, expected)
		}
	})

	t.Run("UpdateHost InvalidUUID", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest(http.MethodGet, "/teams", nil)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusOK
		if code := writer.Code; code != expected {
			errorResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(errorResp.Error)
			t.Fatalf("Status code %d != %d", code, expected)
		}
		teamsResp, err := utility.ReadJSONStruct[utility.GetManyResponseSchema[*utility.TeamGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(teamsResp.Data) == 0 {
			t.Fatal("no data was returned")
		}

		body, err := getJSONBodyAsReader(map[string]any{
			"hostname": "test-host",
			"ip4":      "192.168.0.6",
			"ip6":      "2001:0db8:85a3:0000:0000:8a2e:0370:7454",
			"os":       "Windows",
			"teamID":   teamsResp.Data[0].UUID,
		})
		if err != nil {
			t.Fatal(err)
		}
		uuid, err := utility.GenerateRandomUUID()
		if err != nil {
			t.Fatal(err)
		}
		req, _ = http.NewRequest(http.MethodPut, fmt.Sprintf("/hosts/%s", uuid), body)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer = makeRequest(engine, req)

		expected = http.StatusNotFound
		if code := writer.Code; code != expected {
			if !strings.Contains(fmt.Sprint(code), "2") {
				errorResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
				if err != nil {
					t.Fatal(err)
				}
				t.Log(errorResp.Error)
			}
			t.Fatalf("Status code %d != %d", code, expected)
		}
	})

	t.Run("UpdateHost InvalidBody", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest(http.MethodGet, "/hosts", nil)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusOK
		if code := writer.Code; code != expected {
			errorResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(errorResp.Error)
			t.Fatalf("Status code %d != %d", code, expected)
		}
		resp, err := utility.ReadJSONStruct[utility.GetManyResponseSchema[*utility.HostMachineGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Data) == 0 {
			t.Fatal("no data was returned")
		}

		body, err := getJSONBodyAsReader(map[string]any{
			"invalidField": "invalidValue",
		})
		if err != nil {
			t.Fatal(err)
		}
		req, _ = http.NewRequest(http.MethodPut, fmt.Sprintf("/hosts/%s", resp.Data[0].UUID), body)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer = makeRequest(engine, req)

		expected = http.StatusBadRequest
		if code := writer.Code; code != expected {
			if !strings.HasPrefix(fmt.Sprint(code), "2") {
				errorResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
				if err != nil {
					t.Fatal(err)
				}
				t.Log(errorResp.Error)
			}
			t.Fatalf("Status code %d != %d", code, expected)
		}
		errResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(errResp.Error, "' is required") {
			t.Log(errResp.Error)
			t.Fatal("error response message was not expected message")
		}
	})
}

func TestDeleteHost(t *testing.T) {
	engine := setup()
	jwtString, err := getJWT(engine, TestAdminEmail, TestAdminPassword)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("DeleteHost", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest(http.MethodGet, "/hosts", nil)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusOK
		if code := writer.Code; code != expected {
			errorResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(errorResp.Error)
			t.Fatalf("Status code %d != %d", code, expected)
		}
		resp, err := utility.ReadJSONStruct[utility.GetManyResponseSchema[*utility.HostMachineGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Data) == 0 {
			t.Fatal("no data was returned")
		}

		req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("/hosts/%s", resp.Data[len(resp.Data)-1].UUID), nil)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer = makeRequest(engine, req)

		expected = http.StatusNoContent
		if code := writer.Code; code != expected {
			errorResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(errorResp.Error)
			t.Fatalf("Status code %d != %d", code, expected)
		}
	})

	t.Run("DeleteHost InvalidUUID", func(t *testing.T) {
		t.Parallel()
		uuid, err := utility.GenerateRandomUUID()
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/hosts/%s", uuid), nil)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusNotFound
		if code := writer.Code; code != expected {
			if !strings.HasPrefix(fmt.Sprint(code), "2") {
				errorResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
				if err != nil {
					t.Fatal(err)
				}
				t.Log(errorResp.Error)
			}
			t.Fatalf("Status code %d != %d", code, expected)
		}
	})
}
