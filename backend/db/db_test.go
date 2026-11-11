package db

import "testing"

func TestCleanURL(t *testing.T) {
	cases := map[string]string{
		"postgresql://u:p@ep-x-pooler.neon.tech/neondb?sslmode=require&channel_binding=require":   "postgresql://u:p@ep-x-pooler.neon.tech/neondb?sslmode=require",
		"psql 'postgresql://u:p@ep-x.neon.tech/neondb?sslmode=require&channel_binding=require'": "postgresql://u:p@ep-x.neon.tech/neondb?sslmode=require",
		" postgres://u:p@h/db ": "postgres://u:p@h/db?sslmode=require",
	}
	for in, want := range cases {
		if got := cleanURL(in); got != want {
			t.Errorf("cleanURL(%q) = %q, want %q", in, got, want)
		}
	}
}
