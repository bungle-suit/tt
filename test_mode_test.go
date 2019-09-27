package tt_test

import (
	"testing"

	"github.com/bungle-suit/tt"
	"github.com/stretchr/testify/assert"
)

func TestTestMode(t *testing.T) {
	assert.True(t, tt.TestMode())
}
