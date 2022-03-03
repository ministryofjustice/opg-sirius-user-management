package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(Get("PATH", "123"), os.Getenv("PATH"))
	assert.Equal(Get("IFTHISEXISTSTHENTHETESTWILLFAILBUTITSHOULDNOT", "123"), "123")
}
