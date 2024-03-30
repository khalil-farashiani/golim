package main

import (
	"net/http/httptest"
	"testing"
)

func TestReadUserIP(t *testing.T) {
	//arrange
	testCases := []struct {
		name      string
		ipHeaders map[string]string
		expected  string
	}{
		{
			name: "X-Real-Ip header",
			ipHeaders: map[string]string{
				"X-Real-Ip": "192.168.0.1",
			},
			expected: "192.168.0.1",
		},
		{
			name: "X-Forwarded-For header",
			ipHeaders: map[string]string{
				"X-Forwarded-For": "192.168.0.2",
			},
			expected: "192.168.0.2",
		},
		{
			name:     "RemoteAddr field",
			expected: "192.168.0.4:67890",
		},
		{
			name: "No headers",
			ipHeaders: map[string]string{
				"Other-Header": "value",
			},
			expected: "192.168.0.4:67890",
		},
	}

	for _, tc := range testCases {
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "192.168.0.4:67890" // Default RemoteAddr

		for k, v := range tc.ipHeaders {
			req.Header.Set(k, v)
		}
		//act
		got := readUserIP(req)

		//assert
		if got != tc.expected {
			t.Errorf("%s: readUserIP() = %s; want %s", tc.name, got, tc.expected)
		}
	}
}
