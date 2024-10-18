package test_test

import (
	"com668-backend/middleware"
	"com668-backend/utility"
	"net/http"
	"testing"
)

func TestGetProviders(t *testing.T) {
	engine := setup()
	jwtString, err := getJWT(engine, TestUserEmail, TestUserPassword)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("LogProviders", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/providers?provider_type=log", nil)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		if code := writer.Code; code != http.StatusOK {
			errorResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(errorResp.Error)
			t.Fatalf("Status code %d != %d", code, http.StatusOK)
		}
		resp, err := utility.ReadJSONStruct[utility.ProvidersGetResponseSchema](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Providers) == 0 {
			t.Fatal("no data was returned")
		}
	})

	t.Run("AlertProviders", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/providers?provider_type=alert", nil)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		if code := writer.Code; code != http.StatusOK {
			errorResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(errorResp.Error)
			t.Fatalf("Status code %d != %d", code, http.StatusOK)
		}
		resp, err := utility.ReadJSONStruct[utility.ProvidersGetResponseSchema](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Providers) == 0 {
			t.Fatal("no data was returned")
		}
	})

	t.Run("InvalidProviderType", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/providers?provider_type=invalid", nil)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		if code := writer.Code; code != http.StatusBadRequest {
			switch code {
			case http.StatusUnauthorized | http.StatusInternalServerError:
				resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
				if err != nil {
					t.Fatal(err)
				}
				t.Log(resp.Error)
			default:
				break
			}
			t.Fatalf("Status code %d != %d", code, http.StatusBadRequest)
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
}
