package test_test

import (
	"com668-backend/utility"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetProviders(t *testing.T) {
	engine := setup()

	t.Run("LogProviders", func(t *testing.T) {
		writer := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/providers?provider_type=log", nil)
		engine.ServeHTTP(writer, req)

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
		writer := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/providers?provider_type=alert", nil)
		engine.ServeHTTP(writer, req)

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
		writer := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/providers?provider_type=invalid", nil)
		engine.ServeHTTP(writer, req)

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

func TestGetSettings(t *testing.T) {
	engine := setup()

	t.Run("LogSettings", func(t *testing.T) {
		writer := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/providers?provider_type=log", nil)
		engine.ServeHTTP(writer, req)

		if code := writer.Code; code != http.StatusOK {
			errorResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(errorResp.Error)
			t.Fatalf("Status code %d != %d", code, http.StatusOK)
		}
		providerResp, err := utility.ReadJSONStruct[utility.ProvidersGetResponseSchema](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(providerResp.Providers) == 0 {
			t.Fatal("no data was returned")
		}
		uuid := providerResp.Providers[0].ID

		writer = httptest.NewRecorder()
		req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/providers/%s/settings?provider_type=log", uuid), nil)
		engine.ServeHTTP(writer, req)

		if code := writer.Code; code != http.StatusOK {
			errorResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(errorResp.Error)
			t.Fatalf("Status code %d != %d", code, http.StatusOK)
		}
		settingsResp, err := utility.ReadJSONStruct[utility.SettingsGetResponseSchema](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(settingsResp.Settings) == 0 {
			t.Fatal("no data was returned")
		}
	})

	t.Run("AlertSettings", func(t *testing.T) {
		writer := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/providers?provider_type=alert", nil)
		engine.ServeHTTP(writer, req)

		if code := writer.Code; code != http.StatusOK {
			errorResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(errorResp.Error)
			t.Fatalf("Status code %d != %d", code, http.StatusOK)
		}
		providerResp, err := utility.ReadJSONStruct[utility.ProvidersGetResponseSchema](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(providerResp.Providers) == 0 {
			t.Fatal("no data was returned")
		}
		uuid := providerResp.Providers[0].ID

		writer = httptest.NewRecorder()
		req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/providers/%s/settings?provider_type=alert", uuid), nil)
		engine.ServeHTTP(writer, req)

		if code := writer.Code; code != http.StatusOK {
			errorResp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(errorResp.Error)
			t.Fatalf("Status code %d != %d", code, http.StatusOK)
		}
		settingsResp, err := utility.ReadJSONStruct[utility.SettingsGetResponseSchema](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(settingsResp.Settings) == 0 {
			t.Fatal("no data was returned")
		}
	})

	t.Run("InvalidProviderUUID", func(t *testing.T) {
		uuid, err := utility.GenerateRandomUUID()
		if err != nil {
			t.Fatal(err)
		}

		writer := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/providers/%s/settings?provider_type=alert", uuid), nil)
		engine.ServeHTTP(writer, req)

		if code := writer.Code; code != http.StatusNotFound {
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
		if resp.Error != fmt.Sprintf("provider with uuid '%s' not found", uuid) {
			t.Log(resp.Error)
			t.Fatal("error response message was not expected message")
		}
	})

	t.Run("InvalidProviderType", func(t *testing.T) {
		uuid, err := utility.GenerateRandomUUID()
		if err != nil {
			t.Fatal(err)
		}

		writer := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/providers/%s/settings?provider_type=invalid", uuid), nil)
		engine.ServeHTTP(writer, req)

		if code := writer.Code; code != http.StatusBadRequest {
			switch code {
			case http.StatusUnauthorized | http.StatusNotFound | http.StatusInternalServerError:
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
