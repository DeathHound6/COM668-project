package test_test

import (
	"com668-backend/middleware"
	"com668-backend/utility"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestGetMe(t *testing.T) {
	engine := setup()
	jwtString, err := getJWT(engine, TestAdminEmail, TestAdminPassword)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("GetMe", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest(http.MethodGet, "/me", nil)
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

		resp, err := utility.ReadJSONStruct[utility.UserGetResponseBodySchema](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}

		token, err := jwt.Parse(
			strings.Split(jwtString, " ")[1],
			func(t *jwt.Token) (any, error) {
				return []byte(os.Getenv("JWT_SIGNING_KEY")), nil
			},
			jwt.WithValidMethods([]string{
				middleware.JWTSigningMethod.Name,
			}),
		)
		if err != nil {
			t.Fatal(err)
		}
		sub, err := token.Claims.GetSubject()
		if err != nil {
			t.Fatal(err)
		}
		subBytes, err := base64.StdEncoding.DecodeString(sub)
		if err != nil {
			t.Fatal(err)
		}
		if resp.UUID != string(subBytes) {
			t.Fatal("uuid mismatch")
		}
	})

	t.Run("GetMe Unauthorized", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest(http.MethodGet, "/me", nil)
		writer := makeRequest(engine, req)

		expected := http.StatusUnauthorized
		if code := writer.Code; code != expected {
			if strings.HasPrefix(fmt.Sprint(writer.Code), "4") || strings.HasPrefix(fmt.Sprint(writer.Code), "5") {
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

func TestCreateUser(t *testing.T) {
	engine := setup()
	jwtString, err := getJWT(engine, TestAdminEmail, TestAdminPassword)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("CreateUser", func(t *testing.T) {
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

		body, err := getJSONBodyAsReader(map[string]any{
			"name":     "A",
			"email":    "a@example.com",
			"password": "password",
			"teams":    []string{resp.Data[0].UUID},
		})
		if err != nil {
			t.Fatal(err)
		}
		req, _ = http.NewRequest(http.MethodPost, "/users", body)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
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

	t.Run("CreateUser InvalidBody", func(t *testing.T) {
		t.Parallel()
		body, err := getJSONBodyAsReader(map[string]any{
			"invalidField": "invalidValue",
		})
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest(http.MethodPost, "/users", body)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusBadRequest
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("status code %d != %d", code, expected)
		}
		resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(resp.Error, "' is required") {
			t.Fatal("error message mismatch")
		}
	})
}

func TestUserLogin(t *testing.T) {
	engine := setup()

	t.Run("UserLogin", func(t *testing.T) {
		body, err := getJSONBodyAsReader(map[string]any{
			"email":    TestAdminEmail,
			"password": TestAdminPassword,
		})
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest(http.MethodPost, "/users/login", body)
		writer := makeRequest(engine, req)

		expected := http.StatusNoContent
		if code := writer.Code; code != expected {
			resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			t.Fatalf("status code %d != %d", code, expected)
		}
	})

	t.Run("UserLogin Unauthorized", func(t *testing.T) {
		body, err := getJSONBodyAsReader(map[string]any{
			"email":    TestAdminEmail,
			"password": "wrongpassword",
		})
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest(http.MethodPost, "/users/login", body)
		writer := makeRequest(engine, req)

		expected := http.StatusBadRequest
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
