package test_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateIncident(t *testing.T) {
	engine := setup()
	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/incidents", nil)
	engine.ServeHTTP(writer, req)

	if code := writer.Code; code != 201 {
		t.Logf("writer code: '%d'", code)
		t.FailNow()
	}
	if header := writer.Header().Get("location"); header == "" {
		t.Logf("location header: '%s'", header)
		t.FailNow()
	}
}
