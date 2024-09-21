package mongo_migrate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractVersionDescription(t *testing.T) {
	version, description, err := extractVersionDescription("00001_test_hello_abc.go")
	assert.NoError(t, err)
	assert.Falsef(t, version != 1 || description != "test_hello_abc", "Bad version/description: %v %v", version, description)

	_, _, err = extractVersionDescription("test.go")
	assert.Error(t, err, "Unexpected nil error")

	_, _, err = extractVersionDescription("test")
	assert.Error(t, err, "Unexpected nil error")

	_, _, err = extractVersionDescription("test_test.go")
	assert.Error(t, err, "Unexpected nil error")
}
