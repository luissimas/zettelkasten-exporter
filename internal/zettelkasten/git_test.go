package zettelkasten

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeAuthenticatedURL(t *testing.T) {
	data := []struct {
		name     string
		rawURL   string
		token    string
		expected string
	}{
		{
			name:     "empty token",
			rawURL:   "https://github.com/user-name/repository-name",
			token:    "",
			expected: "https://github.com/user-name/repository-name",
		},
		{
			name:     "token",
			rawURL:   "https://github.com/user-name/repository-name",
			token:    "any-token",
			expected: "https://git:any-token@github.com/user-name/repository-name",
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			result, err := makeAuthenticatedUrl(d.rawURL, d.token)
			assert.Nil(t, err)
			assert.Equal(t, d.expected, result)
		})
	}
}
