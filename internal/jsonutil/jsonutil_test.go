package jsonutil_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rubensdev/inventoryflow-backend/internal/jsonutil"
)

func TestReadJSON(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		err := jsonutil.ReadJSON(w, r, &input)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte(input.Username + "," + input.Password))
	}

	tests := []struct {
		Name        string
		PostData    []byte
		ExpectedRes string
	}{
		{
			Name:        "Empty body should return an error",
			PostData:    nil,
			ExpectedRes: jsonutil.EmptyBodyJSONErrorMsg,
		},
		{
			Name:        "Badly-formed JSON should return an error",
			PostData:    []byte("{\"hello:1}"),
			ExpectedRes: jsonutil.BadlyFormedJSONErrorMsg,
		},
		{
			Name:        "Badly-formed JSON at character 12 should return an error",
			PostData:    []byte("{\"hello\":1,}"),
			ExpectedRes: jsonutil.BadlyFormedJSONError(12).Error(),
		},
		{
			Name:        "Incorrect JSON type for the field \"username\" should return an error",
			PostData:    []byte("{\"username\":2}"),
			ExpectedRes: jsonutil.FieldTypeError("username").Error(),
		},
		{
			Name:        "Incorrect JSON type for the field \"username\" should return an error",
			PostData:    []byte("{\"username\":2}"),
			ExpectedRes: jsonutil.FieldTypeError("username").Error(),
		},
		{
			Name:        "A well formed JSON should return an OK as response",
			PostData:    []byte("{\"username\":\"foo\", \"password\": \"foo123\"}"),
			ExpectedRes: "foo,foo123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(tt.PostData))
			w := httptest.NewRecorder()

			handler(w, req)

			res := w.Result()
			defer res.Body.Close()

			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("error: %v", err)
			}
			result := string(data)

			if result != tt.ExpectedRes {
				t.Errorf("got \"%s\", expected \"%s\"", result, tt.ExpectedRes)
			}
		})
	}

}

func TestWriteJSON(t *testing.T) {
	expectedStatusCode := http.StatusOK
	cacheControlValue := "max-age=604800"

	handler := func(w http.ResponseWriter, r *http.Request) {
		headers := http.Header{}
		headers.Add("Cache-Control", cacheControlValue)

		err := jsonutil.WriteJSON(w, expectedStatusCode, jsonutil.H{
			"message": "Hello World",
		}, headers)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("an error has occurred"))
			t.Error(err)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	result := string(data)
	expected := "{\n\t\"message\": \"Hello World\"\n}\n"

	if result != expected {
		t.Errorf("got \"%s\", expected \"%s\"", result, expected)
	}

	if res.StatusCode != expectedStatusCode {
		t.Errorf("got status code %d, expected %d", res.StatusCode, expectedStatusCode)
	}

	if content := res.Header.Get("Content-Type"); content != "application/json" {
		t.Errorf("Content-Type header value must be application/json")
	}

	if content := res.Header.Get("Cache-Control"); content != cacheControlValue {
		t.Errorf("Cache-Control header value must be present and the value maxage=604800")
	}
}
