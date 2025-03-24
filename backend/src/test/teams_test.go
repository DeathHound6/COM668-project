package test_test

import (
	"com668-backend/middleware"
	"com668-backend/utility"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestGetTeams(t *testing.T) {
	engine := setup()
	jwtString, err := getJWT(engine, TestAdminEmail, TestAdminPassword)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("GetTeams", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest(http.MethodGet, "/teams", nil)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
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
		resp, err := utility.ReadJSONStruct[utility.GetManyResponseSchema[utility.TeamGetResponseBodySchema]](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Data) == 0 {
			t.Fatal("no teams")
		}
	})

	t.Run("GetTeams InvalidCommonParams", func(t *testing.T) {
		t.Parallel()
		// invalid page
		page := "invalid"
		pageSize := "10"

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/teams?page=%s&pageSize=%s", page, pageSize), nil)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusBadRequest
		if code := writer.Code; code != expected {
			if !strings.Contains(fmt.Sprint(writer.Code), "2") {
				resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
				if err != nil {
					t.Fatal(err)
				}
				t.Log(resp.Error)
			}
			t.Fatalf("status code %d != %d", code, expected)
		}

		// invalid page size
		page = "1"
		pageSize = "invalid"

		req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/teams?page=%s&pageSize=%s", page, pageSize), nil)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer = makeRequest(engine, req)

		expected = http.StatusBadRequest
		if code := writer.Code; code != expected {
			if !strings.Contains(fmt.Sprint(writer.Code), "2") {
				resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
				if err != nil {
					t.Fatal(err)
				}
				t.Log(resp.Error)
			}
			t.Fatalf("status code %d != %d", code, expected)
		}

		// valid params
		page = "1"
		pageSize = "10"

		req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/teams?page=%s&pageSize=%s", page, pageSize), nil)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
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
	})
}
