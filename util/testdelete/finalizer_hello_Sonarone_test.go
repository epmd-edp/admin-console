package finalizer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContainsStringHelloSonarone(t *testing.T) {
	e := ContainsStringHelloSonarone([]string{"a", "b"}, "b")
	assert.True(t, e)
}

