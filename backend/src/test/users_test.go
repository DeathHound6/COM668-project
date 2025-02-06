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
		req, _ := http.NewRequest(http.MethodGet, "/me", nil)
		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
		writer := makeRequest(engine, req)

		expected := http.StatusOK
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

		resp, err := utility.ReadJSONStruct[utility.UserGetResponseBodySchema](writer.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}

		token, err := jwt.Parse(
			strings.Split(jwtString, " ")[1],
			func(t *jwt.Token) (interface{}, error) {
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
		if resp.Email != string(subBytes) {
			t.Fatal("email mismatch")
		}
	})

	t.Run("GetMe Unauthorized", func(t *testing.T) {
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

// func TestCreateUser(t *testing.T) {
// 	engine := setup()
// 	jwtString, err := getJWT(engine, TestAdminEmail, TestAdminPassword)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Run("CreateUser", func(t *testing.T) {
// 		body, err := getJSONBodyAsReader(map[string]any{
// 			"name":     "A",
// 			"email":    "a@example.com",
// 			"password": "password",
// 			"teams":    []string{},
// 		})
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		req, _ := http.NewRequest(http.MethodPost, "/users", body)
// 		req.Header.Add(middleware.AuthHeaderNameString, jwtString)
// 		writer := makeRequest(engine, req)

// 		expected := http.StatusCreated
// 		if code := writer.Code; code != expected {
// 			if strings.HasPrefix(fmt.Sprint(writer.Code), "4") || strings.HasPrefix(fmt.Sprint(writer.Code), "5") {
// 				resp, err := utility.ReadJSONStruct[utility.ErrorResponseSchema](writer.Body.Bytes())
// 				if err != nil {
// 					t.Fatal(err)
// 				}
// 				t.Log(resp.Error)
// 			}
// 			t.Fatalf("status code %d != %d", code, expected)
// 		}
// 	})
// }
