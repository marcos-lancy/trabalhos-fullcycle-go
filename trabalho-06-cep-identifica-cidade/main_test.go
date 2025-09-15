package main

import (
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestValidateCEP(t *testing.T) {
	tests := []struct {
		cep      string
		expected bool
	}{
		{"12345678", true},
		{"01234567", true},
		{"1234567", false},
		{"123456789", false},
		{"1234567a", false},
		{"", false},
		{"1234-5678", false},
	}

	for _, test := range tests {
		result := validateCEP(test.cep)
		if result != test.expected {
			t.Errorf("validateCEP(%s) = %v, expected %v", test.cep, result, test.expected)
		}
	}
}

func TestCelsiusToFahrenheit(t *testing.T) {
	tests := []struct {
		celsius    float64
		fahrenheit float64
	}{
		{0, 32},
		{100, 212},
		{28.5, 83.3},
		{-40, -40},
	}

	for _, test := range tests {
		result := celsiusToFahrenheit(test.celsius)
		if math.Abs(result-test.fahrenheit) > 0.01 {
			t.Errorf("celsiusToFahrenheit(%f) = %f, expected %f", test.celsius, result, test.fahrenheit)
		}
	}
}

func TestCelsiusToKelvin(t *testing.T) {
	tests := []struct {
		celsius float64
		kelvin  float64
	}{
		{0, 273},
		{100, 373},
		{28.5, 301.5},
		{-273, 0},
	}

	for _, test := range tests {
		result := celsiusToKelvin(test.celsius)
		if math.Abs(result-test.kelvin) > 0.01 {
			t.Errorf("celsiusToKelvin(%f) = %f, expected %f", test.celsius, result, test.kelvin)
		}
	}
}

func TestWeatherHandlerInvalidCEP(t *testing.T) {
	req, err := http.NewRequest("GET", "/weather/1234567", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := map[string]string{"cep": "1234567"}
		r = mux.SetURLVars(r, vars)
		weatherHandler(w, r)
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnprocessableEntity)
	}

	expected := "invalid zipcode"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
