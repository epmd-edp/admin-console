package context

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertToBoolMethod_ShouldBeExecutedSuccessfully(t *testing.T) {
	assert.False(t, convertToBool())
}
