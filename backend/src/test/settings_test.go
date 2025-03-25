package test_test

import (
	"com668-backend/middleware"
	"com668-backend/utility"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestGetProviders(t *testing.T) {
	engine := setup()
	jwtString, err := getJWT(engine, TestAdminEmail, TestAdminPassword)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("GetProviders ValidProviderType", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest(http.MethodGet, "/providers?provider_type=alert", nil)
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
		resp, err := utility.ReadJSONStruct[utility.GetManyResponseSchema[*utility.ProviderGetResponseSchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Data) == 0 {
			t.Fatal("no data was returned")
		}
	})

	t.Run("GetProviders InvalidProviderType", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest(http.MethodGet, "/providers?provider_type=invalid", nil)
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
		if resp.Error != "'provider_type' query parameter must be either 'log' or 'alert'" {
			t.Log(resp.Error)
			t.Fatal("error response message was not expected message")
		}
	})

	t.Run("GetProviders InvalidCommonParams", func(t *testing.T) {
		t.Parallel()
		page := "invalid"
		pageSize := "10"

		// invalid page
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/providers?provider_type=log&page=%s&pageSize=%s", page, pageSize), nil)
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

		req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/providers?provider_type=log&page=%s&pageSize=%s", page, pageSize), nil)
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

	t.Run("GetProviders Forbidden", func(t *testing.T) {
		t.Parallel()
		jwtString, err := getJWT(engine, TestUserEmail, TestUserPassword)
		if err != nil {
			t.Fatal(err)
		}

		req, _ := http.NewRequest(http.MethodGet, "/providers?provider_type=log", nil)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusForbidden
		if code := writer.Code; code != expected {
			errorResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(errorResp.Error)
			t.Fatalf("Status code %d != %d", code, expected)
		}
	})
}

func TestCreateProvider(t *testing.T) {
	engine := setup()
	jwtString, err := getJWT(engine, TestAdminEmail, TestAdminPassword)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("CreateProvider", func(t *testing.T) {
		t.Parallel()
		body := map[string]any{
			"name": "Test Provider",
		}
		bodyReader, err := getJSONBodyAsReader(body)
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest(http.MethodPost, "/providers?provider_type=log", bodyReader)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusCreated
		if code := writer.Code; code != expected {
			errorResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(errorResp.Error)
			t.Fatalf("Status code %d != %d", code, expected)
		}
	})

	t.Run("CreateProvider InvalidProviderType", func(t *testing.T) {
		t.Parallel()
		body := map[string]any{
			"name": "Test Provider",
		}
		bodyReader, err := getJSONBodyAsReader(body)
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest(http.MethodPost, "/providers?provider_type=invalid", bodyReader)
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
		if resp.Error != "'provider_type' query parameter must be either 'log' or 'alert'" {
			t.Log(resp.Error)
			t.Fatal("error response message was not expected message")
		}
	})

	t.Run("CreateProvider InvalidBody", func(t *testing.T) {
		t.Parallel()
		body := map[string]any{
			"a": 1,
		}
		bodyReader, err := getJSONBodyAsReader(body)
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest(http.MethodPost, "/providers?provider_type=log", bodyReader)
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

		if !strings.Contains(resp.Error, "' is required") {
			t.Log(resp.Error)
			t.Fatal("error response message was not expected message")
		}
	})
}

func TestUpdateProvider(t *testing.T) {
	engine := setup()
	jwtString, err := getJWT(engine, TestAdminEmail, TestAdminPassword)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("UpdateProvider", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest(http.MethodGet, "/providers?provider_type=log", nil)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusOK
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("Status code %d != %d", code, expected)
		}

		resp, err := utility.ReadJSONStruct[utility.GetManyResponseSchema[*utility.ProviderGetResponseSchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Data) == 0 {
			t.Fatal("no providers found")
		}

		body := map[string]any{
			"name":   "Test 2",
			"fields": []map[string]any{{"key": "test", "value": "test", "type": "string", "required": false}},
		}
		bodyReader, err := getJSONBodyAsReader(body)
		if err != nil {
			t.Fatal(err)
		}
		req, _ = http.NewRequest(http.MethodPut, fmt.Sprintf("/providers/%s", resp.Data[0].UUID), bodyReader)
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

	t.Run("UpdateProvider InvalidUUID", func(t *testing.T) {
		t.Parallel()
		body := map[string]any{
			"name":   "Test 2",
			"fields": []map[string]any{{"key": "test", "value": "test", "type": "string", "required": false}},
		}
		bodyReader, err := getJSONBodyAsReader(body)
		if err != nil {
			t.Fatal(err)
		}
		uuid, err := utility.GenerateRandomUUID()
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/providers/%s", uuid), bodyReader)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusNotFound
		if code := writer.Code; code != expected {
			errorResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(errorResp.Error)
			t.Fatalf("Status code %d != %d", code, expected)
		}
	})
}

func TestDeleteProvider(t *testing.T) {
	engine := setup()
	jwtString, err := getJWT(engine, TestAdminEmail, TestAdminPassword)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("DeleteProvider", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest(http.MethodGet, "/providers?provider_type=log", nil)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusOK
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("Status code %d != %d", code, expected)
		}

		resp, err := utility.ReadJSONStruct[utility.GetManyResponseSchema[*utility.ProviderGetResponseSchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Data) == 0 {
			t.Fatal("no providers found")
		}

		req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("/providers/%s", resp.Data[len(resp.Data)-1].UUID), nil)
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

	t.Run("DeleteProvider InvalidUUID", func(t *testing.T) {
		t.Parallel()
		uuid, err := utility.GenerateRandomUUID()
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/providers/%s", uuid), nil)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusNotFound
		if code := writer.Code; code != expected {
			errorResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(errorResp.Error)
			t.Fatalf("Status code %d != %d", code, expected)
		}
	})
}
