// Copyright 2026

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNormalizeTheme(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "auto", input: "auto", want: themeAuto},
		{name: "light", input: "light", want: themeLight},
		{name: "dark", input: "dark", want: themeDark},
		{name: "trimmed", input: " dark ", want: themeDark},
		{name: "upper case", input: "DARK", want: themeDark},
		{name: "invalid", input: "night", want: themeAuto},
		{name: "empty", input: "", want: themeAuto},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeTheme(tt.input); got != tt.want {
				t.Fatalf("normalizeTheme(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestCurrentTheme(t *testing.T) {
	tests := []struct {
		name      string
		cookieVal string
		want      string
	}{
		{name: "no cookie", cookieVal: "", want: themeAuto},
		{name: "valid cookie", cookieVal: "light", want: themeLight},
		{name: "invalid cookie", cookieVal: "sunset", want: themeAuto},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.cookieVal != "" {
				req.AddCookie(&http.Cookie{Name: cookieTheme, Value: tt.cookieVal})
			}

			if got := currentTheme(req); got != tt.want {
				t.Fatalf("currentTheme() = %q, want %q", got, tt.want)
			}
		})
	}
}
