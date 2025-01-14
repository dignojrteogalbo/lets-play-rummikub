package constants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexOutOfBounds(t *testing.T) {
	t.Run("ShouldReturnWithMinAndMax", func(t *testing.T) {
		message := IndexOutOfBounds(0, 2)
		assert.Equal(t, message, "index must be > 0 and < 2")
	})
	t.Run("ShouldReturnWithLabel", func(t *testing.T) {
		message := IndexOutOfBounds(0, 2, "custom")
		assert.Equal(t, message, "custom must be > 0 and < 2")
	})
}
