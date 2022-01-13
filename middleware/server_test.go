package middleware

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlob(t *testing.T) {
	url := "/biz/foo/bar"

	pattern := "/biz/*/*"

	matched, err := filepath.Match(pattern, url)

	assert.NoError(t, err)

	assert.NotNil(t, matched)
}
